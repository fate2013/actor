package actor

import (
	"github.com/funkygao/dragon/server"
	log "github.com/funkygao/log4go"
	"runtime"
	"time"
)

func (this *Actor) showStats() {
	log.Info("ver: %s, elapsed:%s, sess:%d, req:%d, goroutine:%d",
		server.BuildID,
		time.Since(this.server.StartedAt),
		this.totalSessionN,
		this.totalReqN,
		runtime.NumGoroutine())
}
