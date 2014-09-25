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
	pveChan   chan Pve

	pollers    map[string]Poller // key is db pool name
	callbacker Callbacker
}

func NewScheduler(interval time.Duration, callbackConf *ConfigCallback,
	conf *ConfigMysql) *Scheduler {
	this := new(Scheduler)
	this.interval = interval
	this.jobChan = make(chan Job, 1000)     // TODO
	this.marchChan = make(chan March, 1000) // TODO
	this.pveChan = make(chan Pve, 1000)
	this.pollers = make(map[string]Poller)
	this.callbacker = NewPhpCallbacker(callbackConf)
	return this
}

func (this *Scheduler) Len() int {
	return len(this.jobChan) + len(this.marchChan) + len(this.pveChan)
}

func (this *Scheduler) Run(myconf *ConfigMysql) {
	go this.scheduleCallback()

	for pool, my := range myconf.Servers {
		this.pollers[pool] = NewMysqlPoller(this.interval, my, &myconf.Breaker)
		if this.pollers[pool] != nil {
			log.Debug("started %s poller", pool)

			go this.pollers[pool].Run(this.jobChan, this.marchChan, this.pveChan)
		}
	}

	log.Info("scheduler started")
}

// TODO do we need finish t jobs callback before callback t+1?
func (this *Scheduler) scheduleCallback() {
	var now time.Time
	for {
		now = time.Now().UTC()

		select {
		case job, ok := <-this.jobChan:
			if !ok {
				log.Critical("job chan closed")
				return
			}

			if now.Sub(job.TimeEnd).Seconds() > 1 {
				log.Warn("late job: %+v, now: %s", job, now)
			}

			go this.callbacker.Wakeup(job)

		case march, ok := <-this.marchChan:
			if !ok {
				log.Critical("march chan closed")
				return
			}

			if march.Ignored() {
				log.Warn("ignored: %+v", march)
				continue
			}

			if now.Sub(march.EndTime).Seconds() > 1 {
				log.Warn("late march: %+v now: %s", march, now)
			}

			go this.callbacker.Play(march)

		case pve, ok := <-this.pveChan:
			if !ok {
				log.Critical("pve chan closed")
				return
			}

			if pve.Ignored() {
				log.Warn("ignored: %+v", pve)
				continue
			}

			if now.Sub(pve.EndTime).Seconds() > 1 {
				log.Warn("late pve: %+v, now: %s", pve, now)
			}

			go this.callbacker.Pve(pve)
		}

	}
}

// TODO kill this
func (this *Scheduler) callbackMarch(m March) {
	rtos := []time.Duration{1, 2, 3, 4, 5}
	for i := 0; i < 5; i++ {
		if retry := this.callbacker.Play(m); !retry {
			return
		}

		return
		time.Sleep(time.Second * rtos[i])
		log.Warn("recallback %+v", m)
	}
}

// TODO kill this
func (this *Scheduler) callbackJob(j Job) {
	rtos := []time.Duration{1, 2, 3, 4, 5}
	for i := 0; i < 5; i++ {
		if retry := this.callbacker.Wakeup(j); !retry {
			log.Debug("ok job callback: %+v", j)
			return
		}

		return

		time.Sleep(time.Second * rtos[i])
		log.Warn("recallback %+v", j)
	}
}
