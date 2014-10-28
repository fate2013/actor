package actor

import (
	"fmt"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"io/ioutil"
	"net/http"
	"time"
)

type PhpWorker struct {
	config *ConfigWorker

	latency metrics.Histogram

	userFlight *Flight
	tileFlight *Flight
}

func NewPhpWorker(config *ConfigWorker) *PhpWorker {
	this := new(PhpWorker)
	this.config = config
	this.userFlight = NewFlight(config.MaxFlightEntries,
		this.config.MaxRetryEntries, this.config.MaxRetries)
	this.tileFlight = NewFlight(config.MaxFlightEntries,
		this.config.MaxRetryEntries, this.config.MaxRetries)
	this.latency = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("latency.php", this.latency)

	api := ApiRunner{userFlight: this.userFlight}
	api.Run()
	return this
}

func (this PhpWorker) Flights() map[string]interface{} {
	return map[string]interface{}{
		"user.items": this.userFlight.items.Len(),
		"user.retry": this.userFlight.retries.Len(),
		"tile.items": this.tileFlight.items.Len(),
		"tile.retry": this.tileFlight.retries.Len(),
	}
}

func (this *PhpWorker) FlightCount() int {
	return this.userFlight.Len() + this.tileFlight.Len()
}

// FIXME retry is not used now
func (this *PhpWorker) Wake(w Wakeable) (retry bool) {
	var (
		params = string(w.Marshal())
		url    string
	)
	switch w := w.(type) {
	case *Pve:
		url = fmt.Sprintf(this.config.Pve, params)

	case *March:
		url = fmt.Sprintf(this.config.March, params)
		if !this.tileFlight.Takeoff(w.GeoHash()) {
			log.Debug("tile locked (%d, %d)", w.X1, w.Y1)
			return
		}

		// attackee lock
		// FIXME not atomic
		if w.OppUid.Valid &&
			w.OppUid.Int64 > 0 &&
			!this.userFlight.Takeoff(w.OppUid.Int64) {
			log.Debug("attackee[%d] locked", w.OppUid.Int64)
			return
		}

	case *Job:
		url = fmt.Sprintf(this.config.Job, params)
	}

	// FIXME atomic for both tile and user flight
	if !this.userFlight.Takeoff(w.GetUid()) {
		log.Debug("user[%d] locked", w.GetUid())
		return
	}

	if this.config.DryRun {
		log.Debug("dry run: %s", url)
		return
	}

	log.Debug("%s", url)

	t0 := time.Now()
	res, err := http.Get(url)
	if err != nil {
		log.Error("http: %s", err.Error())

		this.userFlight.Land(w.GetUid(), false)
		if m, ok := w.(*March); ok {
			this.tileFlight.Land(m.GeoHash(), false)
			if m.OppUid.Valid && m.OppUid.Int64 > 0 {
				this.userFlight.Land(m.OppUid.Int64, false)
			}
		}

		return
	}

	this.latency.Update(time.Since(t0).Nanoseconds() / 1e6)

	defer func() {
		res.Body.Close()
	}()

	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("http: %s", err.Error())
		this.userFlight.Land(w.GetUid(), false)
		if m, ok := w.(*March); ok {
			this.tileFlight.Land(m.GeoHash(), false)
			if m.OppUid.Valid && m.OppUid.Int64 > 0 {
				this.userFlight.Land(m.OppUid.Int64, false)
			}
		}
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Error("unexpected php status: %s", res.Status)
		this.userFlight.Land(w.GetUid(), false)
		if m, ok := w.(*March); ok {
			this.tileFlight.Land(m.GeoHash(), false)
			if m.OppUid.Valid && m.OppUid.Int64 > 0 {
				this.userFlight.Land(m.OppUid.Int64, false)
			}
		}
		return
	}

	if payload[0] == '{' {
		// php.Application json payload means Exception thrown
		log.Error("payload: %s, elapsed: %v, %+v", string(payload), time.Since(t0), w)
		this.userFlight.Land(w.GetUid(), false)
		if m, ok := w.(*March); ok {
			this.tileFlight.Land(m.GeoHash(), false)
			if m.OppUid.Valid && m.OppUid.Int64 > 0 {
				this.userFlight.Land(m.OppUid.Int64, false)
			}
		}
	} else {
		log.Debug("payload: %s, elapsed: %v, %+v", string(payload), time.Since(t0), w)
		this.userFlight.Land(w.GetUid(), true)
		if m, ok := w.(*March); ok {
			this.tileFlight.Land(m.GeoHash(), true)
			if m.OppUid.Valid && m.OppUid.Int64 > 0 {
				this.userFlight.Land(m.OppUid.Int64, true)
			}
		}
	}

	return
}
