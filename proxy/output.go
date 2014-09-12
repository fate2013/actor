package proxy

import (
	log "github.com/funkygao/log4go"
	"io"
	"net"
	"sync/atomic"
)

func (this *Proxy) spawnOutputSessions(batchSize int) {
	atomic.AddInt32(&this.spareSessionN, int32(batchSize))

	for i := 0; i < batchSize; i++ {
		this.wg.Add(1)

		sessionNo := atomic.AddInt64(&this.sessionNo, 1)
		go this.runOutputSession(sessionNo)
	}
}

// TODO session killed after N idle seconds
func (this *Proxy) runOutputSession(sessionNo int64) {
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
		log.Info("session[%d] died", sessionNo)

		conn.Close()
		atomic.AddInt32(&this.activeSessionN, -1)
		atomic.AddInt32(&this.spareSessionN, -1)
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
		// a spare session can wait for inbound request
		atomic.AddInt32(&this.spareSessionN, 1)

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
		atomic.AddInt32(&this.spareSessionN, -1)
		leftN := atomic.LoadInt32(&this.spareSessionN)
		if leftN < int32(this.config.pm.minSpareServerN) {
			log.Info("session[%d] server busy, spare left:%d, spawn %d processes",
				sessionNo,
				leftN,
				this.config.pm.spawnBatchSize)

			go this.spawnOutputSessions(this.config.pm.spawnBatchSize)
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

	}

}
