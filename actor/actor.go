/*
   StatsRunner
     |
   Actor --- Schduler --- Pollers --- mysql farm
               |            |
               |            V backlog(channel of Wakeables)
               |            |
             Worker --------+
               |
               | go Wake(Wakeable) with retries in case of lock conflict
               V
    -------------------------------------------------------------------- callback
               |        |       |
              php      php     php

*/
package actor

import (
	"github.com/funkygao/etclib"
	"github.com/funkygao/golib/server"
	log "github.com/funkygao/log4go"
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
		this.config.SchedulerBacklog, this.config.WorkerConfig,
		this.config.HttpApiListenAddr)
	this.statsRunner = NewStatsRunner(this, this.scheduler)

	return
}

func (this *Actor) ServeForever() {
	go this.scheduler.Run(this.config.MysqlConfig)

	if this.config.EtcdSelfAddr != "" {
		etclib.Init(this.config.EtcdServers, "dw")
		etclib.BootActor(this.config.EtcdSelfAddr)
	}

	this.statsRunner.Run()
}

func (this *Actor) Stop() {
	if this.config.EtcdSelfAddr != "" {
		etclib.Init(this.config.EtcdServers, "dw")

		log.Info("shutdown actor node[%s] in etcd", this.config.EtcdSelfAddr)
		etclib.ShutdownActor(this.config.EtcdSelfAddr)
	}

	this.statsRunner.stopHttpServ()

	this.scheduler.Stop()
}
