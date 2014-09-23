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
	go this.scheduleCallback()

	for pool, my := range myconf.Servers {
		this.pollers[pool] = NewMysqlPoller(this.interval, my, &myconf.Breaker)
		if this.pollers[pool] != nil {
			log.Debug("started %s poller", pool)

			go this.pollers[pool].Run(this.jobChan)
		}
	}

	log.Info("scheduler started")
}

// TODO do we need finish t jobs callback before callback t+1?
func (this *Scheduler) scheduleCallback() {
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

func (this *Scheduler) callback(j Job) {
	rtos := []time.Duration{1, 2, 3, 4, 5}
	for i := 0; i < 5; i++ {
		if retry := this.callbacker.Call(j); !retry {
			return
		}

		time.Sleep(time.Second * rtos[i])
		log.Warn("recallback %v", j)
	}
}
