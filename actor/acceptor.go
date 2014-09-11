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
		go this.runReceiverSession(conn, atomic.AddInt32(&this.totalSessionN, 1))
	}
}

// a single tcp conn that will recv march job, then put to central scheduler
func (this *Actor) runReceiverSession(conn net.Conn, sessionNo int32) {
	defer func() {
		log.Info("session[%d] closed", sessionNo)

		conn.Close()
		atomic.AddInt32(&this.totalSessionN, -1)
	}()

	log.Info("session[%d] started", sessionNo)

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
				log.Error("session[%d] err:%s", sessionNo, err.Error())
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

		log.Debug("session[%d] req: %#v, elapsed:%dus",
			sessionNo,
			req, (time.Now().UnixNano()-req.T0)/1000)
		atomic.AddInt64(&this.totalReqN, 1)
		this.jobs.sched(req)
	}

}
