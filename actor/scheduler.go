package actor

import (
	log "github.com/funkygao/log4go"
	"time"
)

type Scheduler struct {
	interval time.Duration

	// Poller -> ch -> Callbacker
	ch chan Schedulable

	pollers    map[string]Poller // key is db pool name
	callbacker Callbacker
}

func NewScheduler(interval time.Duration, callbackConf *ConfigCallback,
	conf *ConfigMysql) *Scheduler {
	this := new(Scheduler)
	this.interval = interval
	this.ch = make(chan Schedulable, 10000)
	this.pollers = make(map[string]Poller)
	this.callbacker = NewPhpCallbacker(callbackConf)
	return this
}

func (this *Scheduler) Outstandings() int {
	return len(this.ch)
}

func (this *Scheduler) Run(myconf *ConfigMysql) {
	go this.scheduleCallback()

	for pool, my := range myconf.Servers {
		this.pollers[pool] = NewMysqlPoller(this.interval, my, &myconf.Breaker)
		if this.pollers[pool] != nil {
			log.Debug("started %s poller", pool)

			go this.pollers[pool].Run(this.ch)
		}
	}

	log.Info("scheduler started")
}

// TODO do we need finish t jobs callback before callback t+1?
func (this *Scheduler) scheduleCallback() {
	for {
		select {
		case s, ok := <-this.ch:
			if !ok {
				log.Critical("scheduler chan closed")
				return
			}

			if s.Ignored() {
				log.Debug("ignored: %+v", s)
				continue
			}

			if time.Since(s.DueTime()).Seconds() > this.interval.Seconds() {
				log.Debug("late %+v", s)
			}

			go this.callbacker.Call(s)
		}
	}
}
