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

	for {
		select {
		case <-schedTicker.C:
			dueTasks := this.queue.Wakeup(time.Now().Unix())
			if len(dueTasks) > 0 {
				log.Debug("%d marches waked up: %+v", len(dueTasks), dueTasks)
			}

		case <-statsTicker.C:
			this.showConsoleStats()

		}
	}

}
