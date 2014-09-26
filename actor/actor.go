/*
   statsRunner
     |
   actor --- schduler --- pollers --- mysql farm
               |            |
               |            V channel of Wakeable's
               |            |
             worker --------+
               |
          Wake | Wakeable
               V
    -------------------------------------------------------- http | mq | ?
               |        |       |
              php      php     php

*/
package actor

import (
	"github.com/funkygao/golib/server"
)

type Actor struct {
	server *server.Server
	config *ActorConfig

	statsRunner *StatsRunner
	scheduler   *Scheduler
}

func New(server *server.Server) (this *Actor) {
	this = new(Actor)
	this.server = server

	this.config = new(ActorConfig)
	if err := this.config.LoadConfig(server.Conf); err != nil {
		panic(err)
	}

	this.scheduler = NewScheduler(this.config.ScheduleInterval,
		this.config.SchedulerBacklog, this.config.WorkerConfig)
	this.statsRunner = NewStatsRunner(this, this.scheduler)

	return
}

func (this *Actor) ServeForever() {
	go this.scheduler.Run(this.config.MysqlConfig)

	this.statsRunner.Run()
}
