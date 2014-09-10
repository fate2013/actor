package main

import (
	"github.com/funkygao/dragon/server"
	"github.com/funkygao/golib/syslogng"
	log "github.com/funkygao/log4go"
	"os"
	"time"
)

func main() {
	server := server.NewServer("proxyd")
	server.LoadConfig("etc/proxyd.cf")
	server.Launch()

	proxy := newProxy()
	proxy.loadConfig(server.Conf)
	proxy.start(server)

	statsTicker := time.NewTicker(time.Second * time.Duration(proxy.config.statsInterval))
	defer statsTicker.Stop()

	inChan := syslogng.Subscribe()
	for {
		select {
		case req := <-inChan:
			proxy.totalReqN++
			switch req.(type) {
			case error:
				log.Error("subscriber err: %v", req)
				os.Exit(1)

			case []byte:
				reqBody := req.([]byte)
				log.Debug("got event: %s", string(reqBody))

				proxy.dispatch(reqBody)
			}

		case <-statsTicker.C:
			proxy.showStats()
		}
	}

}
