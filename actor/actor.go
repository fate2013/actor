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
)

type Actor struct {
	server *server.Server

	totalReqN      int64
	totalSessionN  int64
	activeSessionN int32
}

func New(server *server.Server) (this *Actor) {
	this = new(Actor)
	this.server = server

	return
}
