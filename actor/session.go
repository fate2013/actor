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
		bytesRead, err = conn.Read(buf)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		conn.Write([]byte(RESPONSE_OK))

		err = json.Unmarshal(buf[:bytesRead], req)
		if err != nil {
			log.Error(err.Error())

			if err == io.EOF {
				ever = false
			}

			continue
		}

		log.Debug("req: %#v", req)
		this.jobs.enque(req)

		select {
		case <-this.stopChan:
			ever = false

		default:
			break
		}

	}

}
