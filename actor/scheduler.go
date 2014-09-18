package actor

import (
	log "github.com/funkygao/log4go"
)

type scheduler struct {
	actor      *Actor
	conf       *ConfigMysql
	callbackCh chan string
	pollers    map[string]*poller
}

func newScheduler(actor *Actor, conf *ConfigMysql) *scheduler {
	this := new(scheduler)
	this.actor = actor
	this.conf = conf
	this.pollers = make(map[string]*poller)
	this.callbackCh = make(chan string, 1000)
	return this
}

func (this *scheduler) run() {
	go this.runCallback()

	for pool, my := range this.conf.Servers {
		mysql := newMysql(my.DSN(), &this.conf.Breaker)
		this.pollers[pool] = newPoller(this.actor, mysql)
		go this.pollers[pool].run(this.callbackCh)
	}

	log.Info("scheduler started")

}

func (this *scheduler) runCallback() {
	for {
		select {
		case params := <-this.callbackCh:
			log.Debug("got callback: %s", params)
		}
	}
}
