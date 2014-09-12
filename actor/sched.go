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

	go this.runAcceptor(listener)

	this.runScheduler()
}

func (this *Actor) runScheduler() {
	schedTicker := time.NewTicker(
		time.Duration(this.server.Int("sched_interval", 1)) * time.Second)
	defer schedTicker.Stop()

	statsTicker := time.NewTicker(
		time.Duration(this.server.Int("stats_interval", 5)) * time.Second)
	defer statsTicker.Stop()

	dumpTicker := time.NewTicker(
		time.Duration(this.server.Int("dump_interval", 500)) * time.Second)
	defer dumpTicker.Stop()

	for {
		select {
		case <-schedTicker.C:
			dueMarches := this.jobs.wakeup(time.Now().Unix())
			if len(dueMarches) != 0 {
				log.Debug("%d marches waked up: %+v", len(dueMarches), dueMarches)

				chunks := marches(dueMarches)
				for _, chunk := range chunks.chunks() {
					log.Debug("chunk: %#v", chunk)

					//go this.callback(march)
				}
			}

		case <-statsTicker.C:
			this.showStats()

		case <-dumpTicker.C:
			go this.dump()
		}
	}

}
