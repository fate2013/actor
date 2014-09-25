package actor

import (
	"fmt"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"io/ioutil"
	"net/http"
	"time"
)

type PhpCallbacker struct {
	config *ConfigCallback

	latency metrics.Histogram

	jobFlight   *Flight
	marchFlight *Flight
	pveFlight   *Flight
}

func NewPhpCallbacker(config *ConfigCallback) *PhpCallbacker {
	this := new(PhpCallbacker)
	this.config = config
	this.jobFlight = NewFlight(10000) // FIXME
	this.marchFlight = NewFlight(10000)
	this.pveFlight = NewFlight(10000)
	this.latency = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("latency.php", this.latency)
	return this
}

// FIXME retry is not used now
func (this *PhpCallbacker) Call(s Schedulable) (retry bool) {
	var (
		params          = string(s.Marshal())
		url             string
		flightContainer *Flight
		flightKey       = s.FlightKey()
	)
	switch s.(type) {
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
		log.Warn("locked %+v", s)
		return
	}

	log.Debug("%s", url)

	t0 := time.Now()
	res, err := http.Get(url)
	this.latency.Update(time.Since(t0).Nanoseconds() / 1e6)
	if err != nil {
		log.Error("http err: %s", err.Error())
		flightContainer.Land(flightKey)
		return
	}

	defer func() {
		res.Body.Close()
		flightContainer.Land(flightKey)
	}()

	payload, err := ioutil.ReadAll(res.Body)
	log.Debug("payload: %s, elapsed: %v, %+v", string(payload), time.Since(t0), s)
	if err != nil {
		log.Error("payload: %s", err.Error())
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Error("callback err: %+v", res)
		return
	}

	return
}
