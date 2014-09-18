package actor

import (
	log "github.com/funkygao/log4go"
)

type scheduler struct {
}

func newScheduler(mysql *ConfigMysql) *scheduler {
	this := new(scheduler)
	return this
}

func (this *scheduler) run() {
	log.Info("scheduler start")

}
