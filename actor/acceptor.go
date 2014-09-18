package actor

import (
	log "github.com/funkygao/log4go"
	"io"
	"net"
	"sync/atomic"
)

func (this *Actor) runAcceptor() {
	listener, err := net.Listen("tcp4", this.server.String("listen_addr", ":9002"))
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("accept error: %s", err.Error())
			continue
		}

		defer conn.Close()

		// each conn is persitent conn
		atomic.AddInt32(&this.activeSessionN, 1)
		go this.runReceiverSession(conn, atomic.AddInt64(&this.totalSessionN, 1))
	}
}

// a single tcp conn that will recv march job, then put to central scheduler
func (this *Actor) runReceiverSession(conn net.Conn, sessionNo int64) {
	defer func() {
		log.Info("session[%d] closed", sessionNo)

		conn.Close()
		atomic.AddInt32(&this.activeSessionN, -1)
	}()

	log.Info("session[%d] started", sessionNo)

	buf := make([]byte, 1024) // TODO reusable mem pool
	var (
		ever      = true
		err       error
		bytesRead int
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

		atomic.AddInt64(&this.totalReqN, 1)

		log.Debug("session[%d] got req: %s", sessionNo, string(buf[:bytesRead]))

		err = this.queue.Enque(buf)
		if err != nil {
			log.Error(err)
		}
	}

}
