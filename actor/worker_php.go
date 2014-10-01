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

	jobFlight   *Flight
	marchFlight *Flight
	pveFlight   *Flight
}

func NewPhpWorker(config *ConfigWorker) *PhpWorker {
	this := new(PhpWorker)
	this.config = config
	this.jobFlight = NewFlight(config.MaxFlightEntries,
		this.config.MaxRetryEntries, this.config.MaxRetries)
	this.marchFlight = NewFlight(config.MaxFlightEntries,
		this.config.MaxRetryEntries, this.config.MaxRetries)
	this.pveFlight = NewFlight(config.MaxFlightEntries,
		this.config.MaxRetryEntries, this.config.MaxRetries)
	this.latency = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("latency.php", this.latency)
	return this
}

func (this PhpWorker) Flights() map[string]interface{} {
	return map[string]interface{}{
		"job.lock":    this.jobFlight.lock.Len(),
		"job.retry":   this.jobFlight.retry.Len(),
		"march.lock":  this.marchFlight.lock.Len(),
		"march.retry": this.marchFlight.retry.Len(),
		"pve.lock":    this.pveFlight.lock.Len(),
		"pve.retry":   this.pveFlight.retry.Len(),
	}
}

func (this *PhpWorker) InFlight() int {
	return this.jobFlight.Len() + this.marchFlight.Len() + this.pveFlight.Len()
}

// FIXME retry is not used now
func (this *PhpWorker) Wake(w Wakeable) (retry bool) {
	var (
		params          = string(w.Marshal())
		url             string
		flightContainer *Flight
		flightKey       = w.FlightKey()
	)
	switch w.(type) {
	case *Pve:
		url = fmt.Sprintf(this.config.Pve, params)
		flightContainer = this.pveFlight

	case *March:
		url = fmt.Sprintf(this.config.March, params)
		flightContainer = this.marchFlight

	case *Job:
		url = fmt.Sprintf(this.config.Job, params)
		flightContainer = this.jobFlight
	}

	if !flightContainer.Takeoff(flightKey) {
		return
	}

	log.Debug("%s", url)
	if this.config.DryRun {
		return
	}

	t0 := time.Now()
	res, err := http.Get(url)
	if err != nil {
		log.Error("http: %s", err.Error())

		flightContainer.Land(flightKey)
		return
	}

	this.latency.Update(time.Since(t0).Nanoseconds() / 1e6)

	defer func() {
		res.Body.Close()
		flightContainer.Land(flightKey)
	}()

	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("http: %s", err.Error())
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Error("unexpected php status: %s", res.Status)
		return
	}

	if payload[0] == '{' {
		// php.Application json payload means Exception thrown
		log.Error("payload: %s, elapsed: %v, %+v", string(payload), time.Since(t0), w)
	} else {
		log.Debug("payload: %s, elapsed: %v, %+v", string(payload), time.Since(t0), w)
	}

	return
}
