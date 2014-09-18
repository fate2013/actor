/*
                           +-- mysql
                           |
          persistent conn  |-- mysql
   actord -----------------|
     |           job poll  |-- mysql
     | http                |
     | callback            +-- mysql
     |
    php ---+
     |     | lock
     +-----+
*/
package actor

import (
	"github.com/funkygao/dragon/server"
)

type Actor struct {
	server *server.Server
	config *ActorConfig

	stats     *statsRunner
	scheduler *scheduler
}

func New(server *server.Server) (this *Actor) {
	this = new(Actor)
	this.server = server

	this.config = new(ActorConfig)
	if err := this.config.LoadConfig(server.Conf); err != nil {
		panic(err)
	}

	this.stats = newStatsRunner(this)
	this.stats.init()

	this.scheduler = newScheduler(this.config.MysqlConfig)

	return
}
