package actor

import (
	"fmt"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"time"
)

type PhpWorker struct {
	config *ConfigWorker

	latency metrics.Histogram

	userFlight *Flight // key is uid
	tileFlight *Flight // key is geohash
}

func NewPhpWorker(apiListenAddr string, config *ConfigWorker) *PhpWorker {
	this := new(PhpWorker)
	this.config = config
	this.userFlight = NewFlight(config.MaxFlightEntries,
		config.DebugLocking, config.LockExpires)
	this.tileFlight = NewFlight(config.MaxFlightEntries,
		config.DebugLocking, config.LockExpires)
	this.latency = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("latency.php", this.latency)

	api := NewHttpApiRunner(apiListenAddr, this.userFlight, this.tileFlight)
	api.Run()
	return this
}

func (this PhpWorker) Flights() map[string]interface{} {
	return map[string]interface{}{
		"user.items": this.userFlight.items.Len(),
		"tile.items": this.tileFlight.items.Len(),
	}
}

func (this *PhpWorker) FlightCount() int {
	return this.userFlight.Len() + this.tileFlight.Len()
}

func (this *PhpWorker) Wake(w Wakeable) {
	var (
		maxRetries  = 3
		randbaseMs  = 50
		totalWaitMs = 0
		waitMs      int
	)
	for i := 0; i < maxRetries; i++ {
		ok := this.tryWake(w)
		if ok {
			elapsed := time.Since(w.DueTime())
			if elapsed.Seconds() > 2 {
				log.Info("late after %s ok: %+v", elapsed, w)
			} else if i > 0 {
				// ever retried
				log.Info("retry[%d] after %dms ok: %+v", i+1, waitMs, w)
			}

			return
		}

		waitMs = (maxRetries-i)*randbaseMs + rand.Intn(randbaseMs) + 10
		totalWaitMs += waitMs
		log.Debug("retry[%d] after %dms: %+v", i+1, waitMs, w)
		time.Sleep(time.Millisecond * time.Duration(waitMs))
	}

	log.Warn("Quit after %dms: %+v", totalWaitMs, w)
}

func (this *PhpWorker) dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, this.config.Timeout)
}

// callback with timeout
func (this *PhpWorker) callPhp(url string) (resp *http.Response, err error) {
	client := http.Client{Transport: &http.Transport{Dial: this.dialTimeout}}
	return client.Get(url)
}

func (this *PhpWorker) tryWake(w Wakeable) (success bool) {
	var (
		params = string(w.Marshal())
		url    string
	)
	switch w := w.(type) {
	case *Pve:
		url = fmt.Sprintf(this.config.Pve, params)

	case *Job:
		url = fmt.Sprintf(this.config.Job, params)

	case *March:
		url = fmt.Sprintf(this.config.March, params)
		if !this.tileFlight.Takeoff(w.TileKey()) {
			return
		}

		// FIXME not atomic
		if w.OppUid.Valid &&
			w.OppUid.Int64 > 0 &&
			!this.userFlight.Takeoff(User{Uid: w.OppUid.Int64}) {
			this.tileFlight.Land(w.TileKey())
			return
		}
	}

	// FIXME atomic for both tile and user flight
	if !this.userFlight.Takeoff(User{Uid: w.GetUid()}) {
		if m, ok := w.(*March); ok {
			this.tileFlight.Land(m.TileKey())

			if m.OppUid.Valid && m.OppUid.Int64 > 0 {
				this.userFlight.Land(User{Uid: m.OppUid.Int64})
			}
		}

		return
	}

	if this.config.DryRun {
		log.Debug("dry run: %s", url)

		this.userFlight.Land(User{Uid: w.GetUid()})

		if m, ok := w.(*March); ok {
			this.tileFlight.Land(m.TileKey())
			if m.OppUid.Valid && m.OppUid.Int64 > 0 {
				this.userFlight.Land(User{Uid: m.OppUid.Int64})
			}
		}

		success = true
		return
	}

	log.Debug("%s", url)

	t0 := time.Now()
	res, err := this.callPhp(url)
	if err != nil {
		log.Error("http: %s", err.Error())

		this.userFlight.Land(User{Uid: w.GetUid()})

		if m, ok := w.(*March); ok {
			this.tileFlight.Land(m.TileKey())
			if m.OppUid.Valid && m.OppUid.Int64 > 0 {
				this.userFlight.Land(User{Uid: m.OppUid.Int64})
			}
		}

		return
	}

	defer func() {
		res.Body.Close()
	}()

	this.latency.Update(time.Since(t0).Nanoseconds() / 1e6)

	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("php err: %s", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		log.Error("unexpected php status: %s", res.Status)
	}

	if payload[0] == '{' {
		// php.Application json payload means Exception thrown
		log.Error("payload:%s, %+v %d %s",
			string(payload), w,
			res.StatusCode, time.Since(t0))
	} else {
		log.Debug("payload:%s, %+v %d %s",
			string(payload), w,
			res.StatusCode, time.Since(t0))
		success = true
	}

	this.userFlight.Land(User{Uid: w.GetUid()})

	if m, ok := w.(*March); ok {
		this.tileFlight.Land(m.TileKey())
		if m.OppUid.Valid && m.OppUid.Int64 > 0 {
			this.userFlight.Land(User{Uid: m.OppUid.Int64})
		}
	}

	return
}
