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
	"sync"
)

type Actor struct {
	server *server.Server

	mutex *sync.Mutex

	totalReqN     int64
	totalSessionN int32

	jobs *jobs
}

func NewActor(server *server.Server) *Actor {
	this := new(Actor)
	this.server = server
	this.mutex = new(sync.Mutex)
	this.jobs = newJobs()

	return this
}
