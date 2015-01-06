package actor

import (
	"fmt"
	"github.com/funkygao/actor/config"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"time"
)

type WorkerPhp struct {
	config *config.ConfigWorkerPhp

	latency metrics.Histogram
}

func NewPhpWorker(config *config.ConfigWorkerPhp) *WorkerPhp {
	this := new(WorkerPhp)
	this.config = config
	this.latency = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("latency.php", this.latency)

	return this
}

func (this *WorkerPhp) Start() {

}

func (this *WorkerPhp) Wake(w Wakeable) {
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

func (this *WorkerPhp) dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, this.config.Timeout)
}

// callback with timeout
func (this *WorkerPhp) callPhp(url string) (resp *http.Response, err error) {
	client := http.Client{Transport: &http.Transport{Dial: this.dialTimeout}}
	return client.Get(url)
}

func (this *WorkerPhp) tryWake(w Wakeable) (success bool) {
	var (
		params  = string(w.Marshal())
		url     string
		lockKey string
		locker  = NewLocker()
	)

	defer locker.ReleaseAll()

	switch w := w.(type) {
	case *Pve:
		url = fmt.Sprintf(this.config.Pve, params)

	case *Job:
		url = fmt.Sprintf(this.config.Job, params)
		lockKey = w.LockKey()

	case *March:
		url = fmt.Sprintf(this.config.March, params)
		lockKey = w.LockKey()

		// lock the attackee
		if !w.IsOpponentSelf() &&
			w.OppUid.Valid &&
			w.OppUid.Int64 > 0 &&
			!locker.Lock(w.OppUidKey()) {
			log.Warn("lock[%s] fails", w.OppUidKey())
			return
		}
	}

	if lockKey != "" && !locker.Lock(lockKey) {
		log.Warn("lock[%s] fails", lockKey)
		return
	}

	log.Debug("calling %s", url)

	t0 := time.Now()
	res, err := this.callPhp(url)
	if err != nil {
		log.Error("http: %s", err.Error())

		return
	}

	defer res.Body.Close()

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
		log.Error("payload:%s, %+v %s",
			string(payload), w, time.Since(t0))
	} else {
		log.Debug("payload:%s, %+v %s",
			string(payload), w, time.Since(t0))
		success = true
	}

	return
}
