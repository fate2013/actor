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

	pool := newDragonPool()
	pool.loadConfig(server.Conf)
	pool.start()

	var reqN int64 = 0
	tick := time.NewTicker(time.Second * time.Duration(pool.config.ticker))
	input := syslogng.Subscribe()
	for {
		select {
		case req := <-input:
			reqN++
			log.Debug("got event: %s", req)
			pool.dispatch([]byte(req))

		case <-tick.C:
			log.Info("req: %d", reqN)
		}
	}

}
