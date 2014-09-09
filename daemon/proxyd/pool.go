package main

import (
	"net"
	"sync/atomic"
)

type dragonPool struct {
	config       proxyConfig
	spareServerN int32

	reqChan chan []byte // max outstanding session throttle
}

func newDragonPool() *dragonPool {
	this := new(dragonPool)
	return this
}

func (this *dragonPool) start() {
	this.reqChan = make(chan []byte, this.config.pm.maxOutstandingSessionN)

}

func (this *dragonPool) spawnSessions(batchSize int) {
	for i := 0; i < batchSize; i++ {
		go this.foward()
		atomic.AddInt32(&this.spareServerN, 1)
	}
}

func (this *dragonPool) dispatch(req []byte) {
	this.reqChan <- req
}

func (this *dragonPool) foward() {
	conn, err := net.Dial("tcp", this.config.dragonServer)
	if err != nil {
		panic(err)
	}
	if this.config.tcpNoDelay {
		conn.(*net.TCPConn).SetNoDelay(true)
	}

	var (
		response         = make([]byte, 1024)
		bytesRead        int
		expectedResponse = "ok"
	)

	for {
		req := <-this.reqChan

		atomic.AddInt32(&this.spareServerN, -1)
		leftN := atomic.LoadInt32(&this.spareServerN)
		if leftN < int32(this.config.pm.minSpareServerN) {
			go this.spawnSessions(this.config.pm.spawnBatchSize)
		}

		conn.Write(req)

		bytesRead, err = conn.Read(response)
		if err != nil || string(response[:bytesRead]) != expectedResponse {
			panic(err)
		}

		// this req forwarded, I'm spare again, able to handle new req
		atomic.AddInt32(&this.spareServerN, 1)

	}

}
