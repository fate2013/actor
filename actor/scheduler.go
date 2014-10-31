package actor

import (
	log "github.com/funkygao/log4go"
	"time"
)

// serial scheduler
type Scheduler struct {
	interval time.Duration
	stopCh   chan bool

	// Poller -> WakeableChannel -> Worker
	backlog chan Wakeable

	pollers map[string]Poller // key is db pool name
	worker  Worker
}

func NewScheduler(interval time.Duration, backlogSize int,
	workerConf *ConfigWorker, apiListenAddr string) *Scheduler {
	this := new(Scheduler)
	this.interval = interval
	this.stopCh = make(chan bool)
	this.backlog = make(chan Wakeable, backlogSize)
	this.pollers = make(map[string]Poller)
	this.worker = NewPhpWorker(apiListenAddr, workerConf)
	return this
}

func (this *Scheduler) Outstandings() int {
	return len(this.backlog)
}

func (this *Scheduler) Stat() map[string]interface{} {
	return map[string]interface{}{
		"worker_flights": this.worker.Flights(),
		"backlog":        len(this.backlog),
	}
}

func (this *Scheduler) FlightCount() int {
	return this.worker.FlightCount()
}

func (this *Scheduler) Stop() {
	for _, p := range this.pollers {
		p.Stop()
	}

	close(this.stopCh)
}

func (this *Scheduler) Run(myconf *ConfigMysql) {
	go this.runWorker()

	var err error
	var pollersN int
	for pool, my := range myconf.Servers {
		this.pollers[pool], err = NewMysqlPoller(this.interval, myconf.SlowThreshold,
			my, &myconf.Query, &myconf.Breaker)
		if err != nil {
			log.Error("poller[%s]: %s", pool, err)
			continue
		}

		log.Debug("started poller[%s]", pool)

		pollersN++
		go this.pollers[pool].Poll(this.backlog)
	}

	if pollersN == 0 {
		panic("Zero poller")
	}

	log.Info("scheduler started with %d pollers", pollersN)
}

func (this *Scheduler) runWorker() {
	for {
		select {
		case w, open := <-this.backlog:
			if !open {
				log.Critical("scheduler chan closed")
				return
			}

			if w.Ignored() {
				log.Debug("ignored: %+v", w)
				continue
			}

			elapsed := time.Since(w.DueTime())
			if elapsed.Seconds() > this.interval.Seconds() {
				log.Warn("late %s for %+v", elapsed, w)
			}

			go this.worker.Wake(w)

		case <-this.stopCh:
			log.Info("scheduler stopped")
			return

		}
	}
}
