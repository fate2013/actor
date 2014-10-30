package actor

import (
	"fmt"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"io/ioutil"
	"math/rand"
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

	api := NewApiRunner(apiListenAddr, this.userFlight, this.tileFlight)
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
	maxRetries := 3
	base := 50
	for i := 0; i < maxRetries; i++ {
		ok := this.tryWake(w)
		if ok {
			return
		}

		waitMs := (maxRetries-i)*base + rand.Intn(base) + 10
		log.Warn("retry after %dms: %+v", waitMs, w)
		time.Sleep(time.Millisecond * time.Duration(waitMs))
	}

	log.Warn("give up waking[%+v] after %d retries, await being rescheduled...", w, maxRetries)
}

func (this *PhpWorker) tryWake(w Wakeable) (ok bool) {
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

	// FIXME dry run didn't release lock correctly
	if this.config.DryRun {
		log.Debug("dry run: %s", url)
		ok = true
		return
	}

	log.Debug("%s", url)

	t0 := time.Now()
	res, err := http.Get(url)
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
		log.Error("payload: %s, elapsed: %v, %+v", string(payload), time.Since(t0), w)
	} else {
		log.Debug("payload: %s, elapsed: %v, %+v", string(payload), time.Since(t0), w)
		ok = true
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
