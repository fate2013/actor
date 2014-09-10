package actor

import (
	"encoding/json"
	log "github.com/funkygao/log4go"
	"net"
)

func (this *Actor) runInboundSession(conn net.Conn) {
	buf := make([]byte, 1024) // TODO reusable mem pool
	var (
		ever      = false
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
			continue
		}

		this.jobs.set(req)

		select {
		case <-this.stopChan:
			ever = false

		default:
			break
		}

	}

}
