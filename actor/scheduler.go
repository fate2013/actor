package actor

import (
	"github.com/funkygao/actor/config"
	log "github.com/funkygao/log4go"
	"time"
)

// serial scheduler
type Scheduler struct {
	config *config.ConfigActor

	interval time.Duration
	stopCh   chan bool

	// Poller -> WakeableChannel -> Worker
	backlog chan Wakeable

	pollers map[string]Poller // key is db pool name
	worker  Worker
}

func NewScheduler(cf *config.ConfigActor) *Scheduler {
	this := new(Scheduler)
	this.config = cf
	this.stopCh = make(chan bool)
	this.backlog = make(chan Wakeable, cf.SchedulerBacklog)
	this.pollers = make(map[string]Poller)
	this.worker = NewPhpWorker(&cf.Worker.Php)
	return this
}

func (this *Scheduler) Outstandings() int {
	return len(this.backlog)
}

func (this *Scheduler) Stat() map[string]interface{} {
	return map[string]interface{}{
		"backlog": this.Outstandings(),
	}
}

func (this *Scheduler) Stop() {
	for _, p := range this.pollers {
		p.Stop()
	}

	close(this.stopCh)
}

func (this *Scheduler) Run() {
	this.worker.Start()
	go this.runWorker()

	var err error
	var pollersN int
	for pool, my := range this.config.Poller.Mysql.Servers {
		this.pollers[pool], err = NewMysqlPoller(
			this.config.ScheduleInterval,
			this.config.Poller.Mysql.SlowThreshold,
			this.config.Poller.Mysql.ManyWakeupsThreshold,
			my,
			&this.config.Poller.Mysql.Query,
			&this.config.Poller.Mysql.Breaker)
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
			if elapsed.Seconds() > this.config.ScheduleInterval.Seconds()+1 {
				log.Debug("late %s for %+v", elapsed, w)
			}

			go this.worker.Wake(w)

		case <-this.stopCh:
			log.Info("scheduler stopped")
			return

		}
	}
}
