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
	proxy.start()

	var reqN int64 = 0
	tick := time.NewTicker(time.Second * time.Duration(proxy.config.ticker))
	input := syslogng.Subscribe()
	for {
		select {
		case req := <-input:
			reqN++
			log.Debug("got event: %s", req)
			proxy.dispatch([]byte(req))

		case <-tick.C:
			log.Info("req: %d", reqN)
		}
	}

}
