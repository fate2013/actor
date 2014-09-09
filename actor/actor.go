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

	mutex *sync.Mutex
	queue *pqueue.PriorityQueue
}

func NewActor(server *server.Server) *Actor {
	this := new(Actor)
	this.server = server
	this.queue = pqueue.New()
	heap.Init(this.queue)
	this.mutex = new(sync.Mutex)

	return this

}

func (this *Actor) add() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

}

func (this *Actor) PeekOrPop() interface{} {
	return this.queue.Peek()

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

			go this.handleUpstreamRequest(conn)
		}
	}()

}

func (this *Actor) handleUpstreamRequest(client net.Conn) {
	defer client.Close()

	buf := make([]byte, 1024) // TODO reusable mem pool
	bytesRead, err := client.Read(buf)
	if err != nil {

	}
	client.Write([]byte("ok"))

	this.parseUpstreamRequest(buf[:bytesRead])

}

func (this *Actor) parseUpstreamRequest(body []byte) {
	var req map[string]interface{}
	err := json.Unmarshal(body, req)
	if err != nil {

	}

	cmd := req["cmd"].(string)
	if cmd != "" {

	}

}
