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
	defer statsTicker.Stop()

	inChan := syslogng.Subscribe()
L:
	for {
		select {
		case req, ok := <-inChan:
			if !ok {
				// inChan closed
				log.Info("inChan closed")
				break L
			}

			proxy.totalReqN++ // no other goroutine will update it, so it's safe

			log.Debug("got event: %s", string(req))
			proxy.dispatch(req)

		case <-statsTicker.C:
			proxy.showStats()
		}
	}

	log.Info("stopping the world")
	proxy.stop()
	proxy.wg.Wait()
}
