package actor

import (
	"encoding/json"
	"fmt"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"io/ioutil"
	"net/http"
	"time"
)

type PhpCallbacker struct {
	url     string
	latency metrics.Histogram
	flight  *Flight
}

func NewPhpCallbacker(url string) *PhpCallbacker {
	this := new(PhpCallbacker)
	this.url = url
	this.flight = NewFlight(10000) // FIXME
	this.latency = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("latency.php", this.latency)
	return this
}

func (this *PhpCallbacker) Call(j Job) (retry bool) {
	if !this.flight.Go(j) { // lock failed
		log.Warn("locked %+v", j)
		return
	}

	params := j.marshal()
	url := fmt.Sprintf(this.url, string(params))
	log.Debug("callback: %s", url)

	t0 := time.Now()
	res, err := http.Get(url)
	this.latency.Update(time.Since(t0).Nanoseconds() / 1e6)
	if err != nil {
		log.Error("http err: %s", err.Error())
		this.flight.Return(j)
		return
	}

	defer func() {
		res.Body.Close()
		this.flight.Return(j)
	}()

	payload, err := ioutil.ReadAll(res.Body)
	log.Debug("payload: %s, elapsed: %v", string(payload), time.Since(t0))
	if err != nil {
		log.Error("payload: %s", err.Error())
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Error("callback err: %+v", res)
		return
	}

	return

	// parse php payload to check if to retry
	var (
		objmap map[string]*json.RawMessage
		ok     int
	)
	err = json.Unmarshal(payload, &objmap)
	if err != nil {
		log.Error("payload err: %s", err.Error())
		return
	}

	json.Unmarshal(*objmap["ok"], &ok)
	log.Debug("payload ok: %d", ok)

	switch ok {
	case RESPONSE_OK:
		break

	case RESPONSE_RETRY:
		// php tells me to retry
		retry = true
	}

	return
}
