package actor

import (
	log "github.com/funkygao/log4go"
	"net"
	"sync/atomic"
)

func (this *Actor) runAcceptor(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("accept error: %s", err.Error())
			continue
		}

		defer conn.Close()

		// each conn is persitent conn
		atomic.AddInt32(&this.totalSessionN, 1)
		go this.runInboundSession(conn)
	}
}
