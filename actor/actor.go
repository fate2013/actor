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
	"github.com/funkygao/dragon/queue"
	"github.com/funkygao/dragon/server"
)

type Actor struct {
	server *server.Server

	replicator *replicator

	totalReqN      int64
	totalSessionN  int64
	activeSessionN int32

	queue *queue.Queue
}

func New(server *server.Server) (this *Actor) {
	this = new(Actor)
	this.server = server
	this.replicator = NewReplicator()
	replicatorConf, err := server.Section("replicator")
	if err != nil {
		panic(err)
	}
	this.replicator.LoadConfig(replicatorConf)
	this.queue = queue.New()

	return
}
