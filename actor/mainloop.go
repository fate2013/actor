package actor

import (
	log "github.com/funkygao/log4go"
	"net"
	"time"
)

func (this *Actor) ServeForever() {
	listener, err := net.Listen("tcp4", this.server.String("listen_addr", ":9002"))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Error(err)
				continue
			}

			defer conn.Close()

			// each conn is persitent conn
			go this.runInboundSession(conn)
		}
	}()

	schedTicker := time.NewTicker(
		time.Duration(this.server.Int("sched_interval", 1)) * time.Second)
	defer schedTicker.Stop()

	statsTicker := time.NewTicker(
		time.Duration(this.server.Int("stats_interval", 5)) * time.Second)
	defer statsTicker.Stop()

	var now time.Time
	for {
		select {
		case <-schedTicker.C:
			now = time.Now()
			for _, m := range this.jobs.pullInBatch(now) {
				go this.callback(m)
			}

		case <-statsTicker.C:
			this.showStats()
		}
	}

}
