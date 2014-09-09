package actor

import (
	"time"
)

func (this *Actor) ServeForever() {
	this.waitForUpstreamRequests()

	ticker := time.NewTicker(time.Duration(this.server.Int("upstream_tick", 1)) * time.Second)
	defer ticker.Stop()

	for _ = range ticker.C {

	}

}
