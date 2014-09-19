/*
   statsRunner
     |
   actor --- schduler --- pollers --- mysql farm
               |            |
               |            | callbackCh
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

	statsRunner *statsRunner
	scheduler   *scheduler
}

func New(server *server.Server) (this *Actor) {
	this = new(Actor)
	this.server = server

	this.config = new(ActorConfig)
	if err := this.config.LoadConfig(server.Conf); err != nil {
		panic(err)
	}

	this.statsRunner = newStatsRunner(this)
	this.statsRunner.init()

	this.scheduler = newScheduler(this.config.ScheduleInterval,
		this.config.CallbackUrl, this.config.MysqlConfig)

	return
}

func (this *Actor) ServeForever() {
	go this.scheduler.run()

	this.statsRunner.run()
}
