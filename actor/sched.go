package actor

import (
	log "github.com/funkygao/log4go"
	"time"
)

func (this *Actor) runScheduler() {
	log.Info("scheduler start")

	schedTicker := time.NewTicker(
		time.Duration(this.server.Int("sched_interval", 1)) * time.Second)
	defer schedTicker.Stop()

	statsTicker := time.NewTicker(
		time.Duration(this.server.Int("stats_interval", 5)) * time.Second)
	defer statsTicker.Stop()

	for {
		select {
		case <-schedTicker.C:

		case <-statsTicker.C:
			this.stats.showConsoleStats()

		}
	}

}
