package server

import (
	"github.com/funkygao/golib/signal"
	log "github.com/funkygao/log4go"
	"os"
	"runtime"
	"syscall"
	"time"
)

func (this *Server) Launch() {
	this.StartedAt = time.Now()
	this.hostname, _ = os.Hostname()
	this.pid = os.Getpid()
	signal.IgnoreSignal(syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGSTOP)

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Info("Server %s ready", this.name)
}
