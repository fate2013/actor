package main

import (
	"github.com/funkygao/dragon/server"
	log "github.com/funkygao/log4go"
	"runtime"
	"sync/atomic"
	"time"
)

func (this *proxy) showStats() {
	spareServerN := atomic.LoadInt32(&this.spareServerN)
	totalReqN := atomic.LoadInt64(&this.totalReqN)
	sessionN := atomic.LoadInt32(&this.activeSessionN)
	log.Info("ver: %s, elapsed:%s, sess:%d, spare:%d, req:%d, goroutine:%d",
		server.BuildID,
		time.Since(this.server.StartedAt),
		sessionN,
		spareServerN,
		totalReqN,
		runtime.NumGoroutine())
}
