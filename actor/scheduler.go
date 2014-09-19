package actor

import (
	log "github.com/funkygao/log4go"
	"time"
)

type Scheduler struct {
	interval time.Duration

	// Poller -> jobChan -> Callbacker
	jobChan    chan Job
	pollers    map[string]Poller
	callbacker Callbacker
}

func NewScheduler(interval time.Duration, callbackUrl string,
	conf *ConfigMysql) *Scheduler {
	this := new(Scheduler)
	this.interval = interval
	this.jobChan = make(chan Job, 1000) // TODO
	this.pollers = make(map[string]Poller)
	this.callbacker = NewPhpCallbacker(callbackUrl)
	return this
}

func (this *Scheduler) Len() int {
	return len(this.jobChan)
}

func (this *Scheduler) Run(myconf *ConfigMysql) {
	go this.runCallback()

	for pool, my := range myconf.Servers {
		this.pollers[pool] = NewMysqlPoller(this.interval, my, &myconf.Breaker)
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
