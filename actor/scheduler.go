package actor

import (
	log "github.com/funkygao/log4go"
	"time"
)

type Scheduler struct {
	interval time.Duration

	// Poller -> ch -> Worker
	ch chan Wakeable

	pollers map[string]Poller // key is db pool name
	worker  Worker
}

func NewScheduler(interval time.Duration, workerConf *ConfigWorker) *Scheduler {
	this := new(Scheduler)
	this.interval = interval
	this.ch = make(chan Wakeable, 10000)
	this.pollers = make(map[string]Poller)
	this.worker = NewPhpWorker(workerConf)
	return this
}

func (this *Scheduler) Outstandings() int {
	return len(this.ch)
}

func (this *Scheduler) InFlight() int {
	return this.worker.InFlight()
}

func (this *Scheduler) Run(myconf *ConfigMysql) {
	go this.runWorker()

	for pool, my := range myconf.Servers {
		this.pollers[pool] = NewMysqlPoller(this.interval, my, &myconf.Breaker)
		if this.pollers[pool] != nil {
			log.Debug("started poller[%s]", pool)

			go this.pollers[pool].Run(this.ch)
		}
	}

	log.Info("scheduler started")
}

func (this *Scheduler) runWorker() {
	for {
		select {
		case w, ok := <-this.ch:
			if !ok {
				log.Critical("scheduler chan closed")
				return
			}

			if w.Ignored() {
				log.Debug("ignored: %+v", w)
				continue
			}

			if time.Since(w.DueTime()).Seconds() > this.interval.Seconds() {
				log.Debug("late for %+v", w)
			}

			go this.worker.Wake(w)
		}
	}
}
