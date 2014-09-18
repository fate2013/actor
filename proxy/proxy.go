package proxy

import (
	"github.com/funkygao/dragon/server"
	"sync"
)

type Proxy struct {
	config proxyConfig
	server *server.Server
	input  Input

	stopChan chan interface{}

	wg *sync.WaitGroup

	sessionNo      int64 // keep track of most recent session number
	activeSessionN int32 // active sessions
	spareSessionN  int32 // maintains persistent tcp conn pool with upstream
	totalReqN      int64 // how many requests served since startup

	reqChan chan []byte // max outstanding session throttle
}

func New() *Proxy {
	this := new(Proxy)
	this.wg = new(sync.WaitGroup)
	this.stopChan = make(chan interface{})
	this.input = new(syslogngInput)
	return this
}

// FIXME do we need dispatch? input can be shared across output sessions
func (this *Proxy) dispatch(req []byte) {
	this.reqChan <- req
}
