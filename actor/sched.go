package actor

import (
	log "github.com/funkygao/log4go"
	"net"
	"sync/atomic"
	"time"
)

func (this *Actor) ServeForever() {
	listener, err := net.Listen("tcp4", this.server.String("listen_addr", ":9002"))
	if err != nil {
		panic(err)
	}

	go this.runAcceptor(listener)

	this.runScheduler()
}

func (this *Actor) runAcceptor(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error(err)
			continue
		}

		defer conn.Close()

		// each conn is persitent conn
		atomic.AddInt32(&this.totalSessionN, 1)
		go this.runInboundSession(conn)
	}
}

func (this *Actor) runScheduler() {
	schedTicker := time.NewTicker(
		time.Duration(this.server.Int("sched_interval", 1)) * time.Second)
	defer schedTicker.Stop()

	statsTicker := time.NewTicker(
		time.Duration(this.server.Int("stats_interval", 5)) * time.Second)
	defer statsTicker.Stop()

	for {
		select {
		case <-schedTicker.C:
			marches := this.jobs.pullInBatch(time.Now())
			if len(marches) != 0 {
				log.Debug("%d events waked up: %+v", len(marches), marches)
			}

			for _, march := range marches {
				go this.callback(march)
			}

		case <-statsTicker.C:
			this.showStats()
		}
	}

}