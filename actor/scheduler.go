package actor

import (
	"bytes"
	"fmt"
	log "github.com/funkygao/log4go"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type scheduler struct {
	interval    time.Duration
	callbackUrl string
	conf        *ConfigMysql

	mutex sync.Mutex

	jobChan         chan job
	outstandingJobs map[job]bool

	pollers map[string]*poller
}

func newScheduler(interval time.Duration, callbackUrl string,
	conf *ConfigMysql) *scheduler {
	this := new(scheduler)
	this.interval = interval
	this.callbackUrl = callbackUrl
	this.conf = conf
	this.pollers = make(map[string]*poller)
	this.outstandingJobs = make(map[job]bool)
	this.jobChan = make(chan job, 1000)
	return this
}

func (this *scheduler) run() {
	go this.runCallback()

	var err error
	for pool, my := range this.conf.Servers {
		mysql := newMysql(my.DSN(), &this.conf.Breaker)
		err = mysql.Open()
		if err != nil {
			log.Critical("open mysql[%+v] failed: %s", *my, err)
			continue
		}

		this.pollers[pool] = newPoller(this.interval, mysql)
		if this.pollers[pool] != nil {
			go this.pollers[pool].run(this.jobChan)
		}
	}

	log.Info("scheduler started")
}

func (this *scheduler) runCallback() {
	for {
		select {
		case job, ok := <-this.jobChan:
			if !ok {
				log.Critical("job chan closed")
				return
			}

			go this.callback(job)
		}
	}
}

func (this *scheduler) callback(j job) {
	params := []byte("")
	url := fmt.Sprintf(this.callbackUrl, string(params))
	log.Debug("callback: %s", url)

	// may fail, because php will throw LockException
	// in that case, will reschedule the job after 1s
	res, err := http.Post(url, CONTENT_TYPE_JSON, bytes.NewBuffer(params))
	defer func() {
		res.Body.Close()
	}()

	payload, err := ioutil.ReadAll(res.Body)
	log.Debug("payload: %s", string(payload))

	if err != nil {
		log.Error("post error: %s", err.Error())
	} else {
		if res.StatusCode != http.StatusOK {
			log.Error("callback error: %+v", res)
		}
	}
}
