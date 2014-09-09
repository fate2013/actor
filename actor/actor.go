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
	"container/heap"
	"encoding/json"
	"github.com/funkygao/dragon/server"
	"github.com/funkygao/golib/pqueue"
	log "github.com/funkygao/log4go"
	"net"
	"sync"
)

type Actor struct {
	server *server.Server

	stopChan chan bool

	mutex *sync.Mutex
	queue *pqueue.PriorityQueue

	marches *marches
}

func NewActor(server *server.Server) *Actor {
	this := new(Actor)
	this.server = server
	this.queue = pqueue.New()
	heap.Init(this.queue)
	this.mutex = new(sync.Mutex)
	this.stopChan = make(chan bool)
	this.marches = newMarches()

	return this
}

func (this *Actor) waitForUpstreamRequests() {
	listener, err := net.Listen("tcp4", this.server.String("listen_addr", ":9002"))
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Error(err)
				continue
			}

			defer conn.Close()

			// each conn is persitent conn
			go this.handleUpstreamRequest(conn)
		}
	}()

}

func (this *Actor) stop() {
	close(this.stopChan)
}

func (this *Actor) handleUpstreamRequest(conn net.Conn) {
	buf := make([]byte, 1024) // TODO reusable mem pool
	var (
		ever      = false
		err       error
		bytesRead int
		req       march
	)

	for ever {
		bytesRead, err = conn.Read(buf)
		if err != nil {
			log.Error(err.Error())
			break
		}

		conn.Write([]byte(RESPONSE_OK))

		err = json.Unmarshal(buf[:bytesRead], req)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		this.marches.set(req)

		select {
		case <-this.stopChan:
			ever = false

		default:
			break
		}

	}

}
