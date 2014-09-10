package actor

import (
	"encoding/json"
	log "github.com/funkygao/log4go"
	"io"
	"net"
	"sync/atomic"
)

func (this *Actor) runInboundSession(conn net.Conn) {
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
			log.Error(err.Error())
			if err == io.EOF {
				ever = false
			}

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

		log.Debug("req: %#v", req)
		atomic.AddInt64(&this.totalReqN, 1)
		this.jobs.submit(req)

		select {
		case <-this.stopChan:
			ever = false

		default:
			break
		}

	}

}
