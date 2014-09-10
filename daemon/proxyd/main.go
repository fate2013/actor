package main

import (
	"github.com/funkygao/dragon/server"
	"github.com/funkygao/golib/syslogng"
	log "github.com/funkygao/log4go"
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
	inChan := syslogng.Subscribe()
	for {
		select {
		case req := <-inChan:
			proxy.totalReqN++
			log.Debug("got event: %s", req)
			proxy.dispatch([]byte(req))

		case <-statsTicker.C:
			proxy.showStats()
		}
	}

}
