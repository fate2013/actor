/*
      php server -----------------------+
        |                               |
        | job{who, when, where, op}     |
        V                               |
      syslogng                          |
        |                               |
      proxyd                            |
        |                               |
        V                               |
   +------------------------+           |
   |        |         actor |           |
   |        |               |           |
   |        V               |           |
   |      queue             |           |
   |        ^               |           |
   |        | tick          |           |
   |        | peekOrPop     |           |
   |        V               |           |
   |      dispatcher        |           |
   |        |               |           |
   +------------------------+           |
            |                           |
            +-------->------------------+
            	callback

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

func New(server *server.Server) *Actor {
	this := new(Actor)
	this.server = server
	this.mutex = new(sync.Mutex)
	this.jobs = newJobs()

	return this
}
