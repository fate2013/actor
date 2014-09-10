/*
	  (upstream)
      php server -----------------+
        |                         |
        | request                 |
        V                         |
   +------------------------+     |
   |        |         actor |     |
   |        |               |     |
   |        | inject        |     |
   |        V               |     |
   |      queue             |     |
   |        ^               |     |
   |        | tick          |     |
   |        | peekOrPop     |     |
   |        V               |     |
   |      dispatcher        |     |
   |        |               |     |
   |        |               |     |
   +------------------------+     |
            |                     |
            +-------->------------+
            	downstream

*/
package actor

import (
	"github.com/funkygao/dragon/server"
	log "github.com/funkygao/log4go"
	"sync"
)

type Actor struct {
	server *server.Server

	stopChan chan bool
	mutex    *sync.Mutex

	totalReqN     int64
	totalSessionN int

	jobs *jobs
}

func NewActor(server *server.Server) *Actor {
	this := new(Actor)
	this.server = server
	this.mutex = new(sync.Mutex)
	this.stopChan = make(chan bool)
	this.jobs = newJobs()

	return this
}

func (this *Actor) stop() {
	log.Info("actor stopped")
	close(this.stopChan)
}
