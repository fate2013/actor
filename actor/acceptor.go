package actor

import (
	"encoding/json"
	log "github.com/funkygao/log4go"
	"io"
	"net"
	"sync/atomic"
	"time"
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
		go this.runAcceptorSession(conn)
	}
}

func (this *Actor) runAcceptorSession(conn net.Conn) {
	defer func() {
		log.Info("session[%+v] closed", conn)

		conn.Close()
		atomic.AddInt32(&this.totalSessionN, -1)
	}()

	buf := make([]byte, 1024) // TODO reusable mem pool
	var (
		ever      = true
		err       error
		bytesRead int
		req       march
	)

	for ever {
		//conn.SetDeadline(time.Now().Add(
		//		time.Duration(this.server.Int("tcp_io_timeout", 5)) * time.Second))
		bytesRead, err = conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Error("session[%+v] err:%s", conn, err.Error())
			}

			ever = false
			continue
		}

		_, err = conn.Write([]byte(RESPONSE_OK))
		if err == io.EOF {
			ever = false

			continue
		}

		err = json.Unmarshal(buf[:bytesRead], &req)
		if err != nil {
			log.Error(err.Error())

			continue
		}

		log.Debug("elapsed:%dus, req: %#v", (time.Now().UnixNano()-req.T0)/1000, req)
		atomic.AddInt64(&this.totalReqN, 1)
		this.jobs.sched(req)

		select {
		case <-this.stopChan:
			ever = false

		default:
			break
		}

	}

}
