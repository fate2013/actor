package actor

import (
	log "github.com/funkygao/log4go"
	"time"
)

type Scheduler struct {
	interval time.Duration

	// poller -> job -> callback
	jobChan    chan Job
	pollers    map[string]Poller
	callbacker Callbacker
}

func newScheduler(interval time.Duration, callbackUrl string,
	conf *ConfigMysql) *Scheduler {
	this := new(Scheduler)
	this.interval = interval
	this.jobChan = make(chan Job, 1000) // TODO
	this.pollers = make(map[string]Poller)
	this.callbacker = newPhpCallbacker(callbackUrl)
	return this
}

func (this *Scheduler) Run(myconf *ConfigMysql) {
	go this.runCallback()

	for pool, my := range myconf.Servers {
		this.pollers[pool] = newMysqlPoller(this.interval, my, &myconf.Breaker)
		if this.pollers[pool] != nil {
			log.Debug("started %s poller", pool)

			go this.pollers[pool].Run(this.jobChan)
		}
	}

	log.Info("scheduler started")
}

func (this *Scheduler) runCallback() {
	for {
		select {
		case job, ok := <-this.jobChan:
			if !ok {
				log.Critical("job chan closed")
				return
			}

			go this.callbacker.Call(job)
		}
	}
}
