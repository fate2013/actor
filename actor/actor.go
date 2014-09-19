/*
   statsRunner
     |
   actor --- schduler --- pollers --- mysql farm
               |            |
               |            | scheduler.jobChan
             callback ------+
               |
               |
              php

*/
package actor

import (
	"github.com/funkygao/dragon/server"
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
		this.config.CallbackUrl, this.config.MysqlConfig)
	this.statsRunner = NewStatsRunner(this, this.scheduler.jobChan)

	return
}

func (this *Actor) ServeForever() {
	go this.scheduler.Run(this.config.MysqlConfig)

	this.statsRunner.Run()
}
