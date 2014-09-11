package main

import (
	"github.com/funkygao/dragon/server"
	log "github.com/funkygao/log4go"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

type proxy struct {
	config proxyConfig
	server *server.Server

	stopChan chan interface{}

	wg *sync.WaitGroup

	sessionNo      int64
	activeSessionN int32 // active sessions
	totalReqN      int64 // how many requests served since startup
	spareServerN   int32 // maintains persistent tcp conn pool with upstream

	reqChan chan []byte // max outstanding session throttle
}

func newProxy() *proxy {
	this := new(proxy)
	this.wg = new(sync.WaitGroup)
	this.stopChan = make(chan interface{})
	return this
}

func (this *proxy) start(server *server.Server) {
	this.server = server
	this.reqChan = make(chan []byte, this.config.pm.maxOutstandingSessionN)
	this.spawnSessions(this.config.pm.startServerN)
}

func (this *proxy) stop() {
	close(this.stopChan)
}

func (this *proxy) spawnSessions(batchSize int) {
	for i := 0; i < batchSize; i++ {
		this.wg.Add(1)
		sessionNo := atomic.AddInt64(&this.sessionNo, 1)
		go this.runForwardSession(sessionNo)

		atomic.AddInt32(&this.spareServerN, 1) // why can't comment this line?
	}
}

func (this *proxy) dispatch(req []byte) {
	this.reqChan <- req
}

func (this *proxy) runForwardSession(sessionNo int64) {
	// setup the tcp conn
	tcpAddr, err := net.ResolveTCPAddr("tcp", this.config.proxyPass)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}

	defer func() {
		log.Info("session[%d] terminated", sessionNo)

		conn.Close()
		atomic.AddInt32(&this.activeSessionN, -1)
		this.wg.Done()
	}()

	if this.config.tcpNoDelay {
		conn.SetNoDelay(true)
	}

	atomic.AddInt32(&this.activeSessionN, 1)
	log.Info("session[%d] started", sessionNo)

	var (
		response         = make([]byte, 1024)
		bytesRead        int
		expectedResponse = "ok"
		ok               bool
		req              []byte
	)

	// mainloop
L:
	for {
		select {
		case req, ok = <-this.reqChan:
			if !ok {
				log.Warn("session[%d] reqChan closed", sessionNo)
				break L
			}

		case <-this.stopChan:
			log.Info("session[%d] stopped", sessionNo)
			break L
		}

		// spawn session on demand
		atomic.AddInt32(&this.spareServerN, -1)
		leftN := atomic.LoadInt32(&this.spareServerN)
		if leftN < int32(this.config.pm.minSpareServerN) {
			log.Info("session[%d] server busy, spare left:%d, spawn %d processes",
				sessionNo,
				leftN,
				this.config.pm.spawnBatchSize)

			go this.spawnSessions(this.config.pm.spawnBatchSize)
		}

		// proxy pass the req
		//conn.SetDeadline(time.Now().Add(this.config.tcpIoTimeout))
		log.Info("session[%d] writing %s", sessionNo, string(req))
		_, err = conn.Write(req)
		if err != nil {
			log.Error("session[%d] write error: %s", sessionNo, err.Error())

			if err == io.EOF {
				log.Info("session[%d] closed", sessionNo)
				return
			}

			continue
		}

		bytesRead, err = conn.Read(response)
		if err != nil {
			if err == io.EOF {
				log.Info("session[%d] closed", sessionNo)
				return
			}

			log.Error(err.Error())
		} else {
			payload := string(response[:bytesRead])
			if payload != expectedResponse {
				log.Warn("session[%d] invalid response: %s", sessionNo, payload)
			}
		}

		// this req forwarded, I'm spare again, able to handle new req
		atomic.AddInt32(&this.spareServerN, 1)
	}

}
