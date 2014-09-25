package actor

import (
	log "github.com/funkygao/log4go"
	"time"
)

type Scheduler struct {
	interval time.Duration

	// Poller -> jobChan -> Callbacker
	jobChan   chan Job
	marchChan chan March

	pollers    map[string]Poller // key is db pool name
	callbacker Callbacker
}

func NewScheduler(interval time.Duration, callbackConf *ConfigCallback,
	conf *ConfigMysql) *Scheduler {
	this := new(Scheduler)
	this.interval = interval
	this.jobChan = make(chan Job, 1000)     // TODO
	this.marchChan = make(chan March, 1000) // TODO
	this.pollers = make(map[string]Poller)
	this.callbacker = NewPhpCallbacker(callbackConf)
	return this
}

func (this *Scheduler) Len() int {
	return len(this.jobChan) + len(this.marchChan)
}

func (this *Scheduler) Run(myconf *ConfigMysql) {
	go this.scheduleCallback()

	for pool, my := range myconf.Servers {
		this.pollers[pool] = NewMysqlPoller(this.interval, my, &myconf.Breaker)
		if this.pollers[pool] != nil {
			log.Debug("started %s poller", pool)

			go this.pollers[pool].Run(this.jobChan, this.marchChan)
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

			go this.callbackJob(job)

		case march, ok := <-this.marchChan:
			if !ok {
				log.Critical("march chan closed")
				return
			}

			if march.Ignored() {
				log.Warn("ignored: %+v", march)
				continue
			}

			go this.callbackMarch(march)
		}

	}
}

func (this *Scheduler) callbackMarch(m March) {
	rtos := []time.Duration{1, 2, 3, 4, 5}
	for i := 0; i < 5; i++ {
		if retry := this.callbacker.Play(m); !retry {
			//log.Debug("ok march callback: %+v", m)
			return
		}

		return
		time.Sleep(time.Second * rtos[i])
		log.Warn("recallback %+v", m)
	}
}

func (this *Scheduler) callbackJob(j Job) {
	rtos := []time.Duration{1, 2, 3, 4, 5}
	for i := 0; i < 5; i++ {
		if retry := this.callbacker.Call(j); !retry {
			log.Debug("ok job callback: %+v", j)
			return
		}

		return

		time.Sleep(time.Second * rtos[i])
		log.Warn("recallback %+v", j)
	}
}
