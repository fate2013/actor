package actor

import (
	"bytes"
	"fmt"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"io/ioutil"
	"net/http"
	"time"
)

type PhpCallbacker struct {
	url          string
	latency      metrics.Histogram
	outstandings *outstandingJobs // FIXME
}

func NewPhpCallbacker(url string) *PhpCallbacker {
	this := new(PhpCallbacker)
	this.url = url
	this.latency = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("latency.php", this.latency)
	this.outstandings = newOutstandingJobs()
	return this
}

func (this *PhpCallbacker) Call(j Job) {
	if !this.outstandings.lock(j) { // lock failed
		log.Debug("locked %+v", j)
		return
	}

	params := j.marshal()
	url := fmt.Sprintf(this.url, string(params))
	log.Debug("callback: %s", url)

	// may fail, because php will throw LockException
	// in that case, will reschedule the job after 1s
	t0 := time.Now()
	res, err := http.Post(url, CONTENT_TYPE_JSON, bytes.NewBuffer(params))
	if err != nil {
		log.Error("http post err: %s", err.Error())
		this.outstandings.unlock(j)
		return
	}

	defer func() {
		res.Body.Close()
		this.outstandings.unlock(j)
	}()

	payload, err := ioutil.ReadAll(res.Body)
	this.latency.Update(time.Since(t0).Nanoseconds() / 1e6)
	log.Debug("payload: %s", string(payload))

	if err != nil {
		log.Error("post error: %s", err.Error())
	} else {
		if res.StatusCode != http.StatusOK {
			log.Error("callback error: %+v", res)
		}
	}
}
