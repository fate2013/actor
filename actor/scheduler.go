package actor

import (
	"github.com/funkygao/actor/config"
	log "github.com/funkygao/log4go"
	"time"
)

type Scheduler struct {
	config *config.ConfigActor

	stopCh chan bool

	jobN, marchN, pnbN int64

	// Poller -> WakeableChannel -> Worker
	backlog chan Wakeable

	mysqlPollers     map[string]Poller // key is db pool name
	beanstalkPollers map[string]Poller // key is tube

	phpWorker Worker
	pnbWorker Worker
	rtmWorker Worker
}

func NewScheduler(cf *config.ConfigActor) *Scheduler {
	this := new(Scheduler)
	this.config = cf
	this.stopCh = make(chan bool)
	this.backlog = make(chan Wakeable, cf.SchedulerBacklog)

	this.mysqlPollers = make(map[string]Poller)
	this.beanstalkPollers = make(map[string]Poller)

	this.phpWorker = NewPhpWorker(&cf.Worker.Php)
	this.pnbWorker = NewPnbWorker(&cf.Worker.Pnb)
	this.rtmWorker = NewRtmWorker(&cf.Worker.Rtm)

	var err error
	for pool, my := range this.config.Poller.Mysql.Servers {
		this.mysqlPollers[pool], err = NewMysqlPoller(
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

		log.Info("Started poller[mysql]: %s", pool)

		go this.mysqlPollers[pool].Poll(this.backlog)
	}

	for tube, beanstalk := range this.config.Poller.Beanstalk.Servers {
		this.beanstalkPollers[tube], err = NewBeanstalkdPoller(beanstalk.ServerAddr, tube)
		if err != nil {
			log.Error("poller[%s]: %s", tube, err)
			continue
		}

		log.Info("Started poller[beanstalk]: %s", tube)

		go this.beanstalkPollers[tube].Poll(this.backlog)
	}

	return this
}

func (this *Scheduler) Outstandings() int {
	return len(this.backlog)
}

func (this *Scheduler) Stat() map[string]interface{} {
	return map[string]interface{}{
		"backlog": this.Outstandings(),
		"jobN":    this.jobN,
		"marchN":  this.marchN,
		"pnbN":    this.pnbN,
	}
}

func (this *Scheduler) Stop() {
	for _, p := range this.beanstalkPollers {
		p.Stop()
	}
	for _, p := range this.mysqlPollers {
		p.Stop()
	}

	close(this.stopCh)
}

func (this *Scheduler) Run() {
	go this.phpWorker.Start()
	go this.pnbWorker.Start()
	go this.rtmWorker.Start()

	for {
		select {
		case w, open := <-this.backlog:
			if !open {
				log.Critical("Scheduler chan closed")
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

			switch w := w.(type) {
			case *Push:
				this.pnbN++

				// TODO -- only use one driver
				wRtm := w
				go this.pnbWorker.Wake(w)
				go this.rtmWorker.Wake(wRtm)

			case *Job:
				this.jobN++

				go this.phpWorker.Wake(w)

			case *March:
				this.marchN++

				go this.phpWorker.Wake(w)
			}

		case <-this.stopCh:
			log.Info("Scheduler stopped")
			return

		}
	}
}
