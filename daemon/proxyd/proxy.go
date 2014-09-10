package main

import (
	"github.com/funkygao/dragon/server"
	log "github.com/funkygao/log4go"
	"io"
	"net"
	"sync/atomic"
)

type proxy struct {
	config proxyConfig
	server *server.Server

	totalReqN    int64 // how many requests served since startup
	spareServerN int32 // maintains persistent tcp conn pool with upstream

	reqChan chan []byte // max outstanding session throttle
}

func newProxy() *proxy {
	this := new(proxy)
	return this
}

func (this *proxy) start(server *server.Server) {
	this.server = server
	this.reqChan = make(chan []byte, this.config.pm.maxOutstandingSessionN)
	this.spawnSessions(this.config.pm.startServerN)
}

func (this *proxy) spawnSessions(batchSize int) {
	for i := 0; i < batchSize; i++ {
		go this.runForwardSession()

		atomic.AddInt32(&this.spareServerN, 1)
	}
}

func (this *proxy) dispatch(req []byte) {
	this.reqChan <- req
}

func (this *proxy) runForwardSession() {
	// setup the tcp conn
	tcpAddr, err := net.ResolveTCPAddr("tcp", this.config.proxyPass)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if this.config.tcpNoDelay {
		conn.SetNoDelay(true)
	}

	var (
		response         = make([]byte, 1024)
		bytesRead        int
		expectedResponse = "ok"
	)

	// mainloop
	for {
		req, ok := <-this.reqChan
		if !ok {
			log.Warn("reqChan closed")
			break
		}

		// spawn session on demand
		atomic.AddInt32(&this.spareServerN, -1)
		leftN := atomic.LoadInt32(&this.spareServerN)
		if leftN < int32(this.config.pm.minSpareServerN) {
			//go this.spawnSessions(this.config.pm.spawnBatchSize)
		}

		// proxy pass the req
		//conn.SetDeadline(time.Now().Add(this.config.tcpIoTimeout))
		_, err = conn.Write(req)
		if err != nil {
			log.Error("write error: %s", err.Error())

			if err == io.EOF {
				log.Info("session[%+v] closed", conn)
				return
			}

			continue
		}

		bytesRead, err = conn.Read(response)
		if err != nil {
			if err == io.EOF {
				log.Info("session[%+v] closed", conn)
				return
			}

			log.Error(err.Error())
		} else {
			payload := string(response[:bytesRead])
			if payload != expectedResponse {
				log.Warn("invalid response: %s", payload)
			}
		}

		// this req forwarded, I'm spare again, able to handle new req
		atomic.AddInt32(&this.spareServerN, 1)
	}

}
