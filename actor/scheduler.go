package actor

import (
	log "github.com/funkygao/log4go"
	"time"
)

type Scheduler struct {
	interval time.Duration

	// Poller -> WakeableChannel -> Worker
	backlog chan Wakeable

	pollers map[string]Poller // key is db pool name
	worker  Worker
}

func NewScheduler(interval time.Duration, backlog int,
	workerConf *ConfigWorker) *Scheduler {
	this := new(Scheduler)
	this.interval = interval
	this.backlog = make(chan Wakeable, backlog)
	this.pollers = make(map[string]Poller)
	this.worker = NewPhpWorker(workerConf)
	return this
}

func (this *Scheduler) Outstandings() int {
	return len(this.backlog)
}

func (this *Scheduler) InFlight() int {
	return this.worker.InFlight()
}

func (this *Scheduler) Stop() {
	for _, p := range this.pollers {
		p.Stop()
	}
}

func (this *Scheduler) Run(myconf *ConfigMysql) {
	go this.runWorker()

	for pool, my := range myconf.Servers {
		this.pollers[pool] = NewMysqlPoller(this.interval, my, &myconf.Breaker)
		if this.pollers[pool] != nil {
			log.Debug("started poller[%s]", pool)

			go this.pollers[pool].Poll(this.backlog)
		}
	}

	log.Info("scheduler started")
}

func (this *Scheduler) runWorker() {
	for {
		select {
		case w, ok := <-this.backlog:
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
