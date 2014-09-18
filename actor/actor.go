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
	stats  *actorStats
}

func New(server *server.Server) (this *Actor) {
	this = new(Actor)
	this.server = server
	this.stats = newActorStats(this)
	this.stats.init()
	this.config = new(ActorConfig)
	if err := this.config.LoadConfig(server.Conf); err != nil {
		panic(err)
	}

	return
}
