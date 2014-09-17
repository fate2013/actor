package actor

import (
	log "github.com/funkygao/log4go"
	"time"
)

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
			if len(dueMarches) > 0 {
				log.Debug("%d marches waked up: %+v", len(dueMarches), dueMarches)

				chunks := marches(dueMarches)
				for _, chunk := range chunks.chunks() {
					log.Debug("chunk: %#v", chunk)

					this.coordinate(chunk)

					//go this.callback(march)
				}
			}

		case <-statsTicker.C:
			this.showConsoleStats()

		case <-dumpTicker.C:
			go this.dump()
		}
	}

}
