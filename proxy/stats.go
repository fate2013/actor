package proxy

import (
	"github.com/funkygao/dragon/server"
	log "github.com/funkygao/log4go"
	"runtime"
	"sync/atomic"
	"time"
)

func (this *Proxy) showStats() {
	spareSessionN := atomic.LoadInt32(&this.spareSessionN)
	totalReqN := atomic.LoadInt64(&this.totalReqN)
	sessionN := atomic.LoadInt32(&this.activeSessionN)
	log.Info("ver: %s, elapsed:%s, sess:%d, spare:%d, req:%d, goroutine:%d",
		server.BuildID,
		time.Since(this.server.StartedAt),
		sessionN,
		spareSessionN,
		totalReqN,
		runtime.NumGoroutine())
}
