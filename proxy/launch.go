package proxy

import (
	"github.com/funkygao/dragon/server"
	log "github.com/funkygao/log4go"
	"time"
)

func (this *Proxy) ServeForever() {
	statsTicker := time.NewTicker(time.Second * time.Duration(this.config.statsInterval))
	defer statsTicker.Stop()

	inChan := this.input.Reader()
L:
	for {
		select {
		case req, ok := <-inChan: // FIXME race condition with this.Stop
			if !ok {
				// inChan closed
				log.Info("syslog-ng closed")
				break L
			}

			this.totalReqN++ // no other goroutine will update it, so it's safe

			log.Debug("got event: %s", string(req))
			this.dispatch(req)

		case <-statsTicker.C:
			this.showStats()
		}
	}

	log.Info("stopping the world")
	this.Stop()
	this.wg.Wait()

}

func (this *Proxy) Start(server *server.Server) *Proxy {
	this.server = server
	this.loadConfig(server.Conf)
	this.reqChan = make(chan []byte, this.config.pm.maxServerN)
	this.spawnOutputSessions(this.config.pm.startServerN)
	return this
}

func (this *Proxy) Stop() {
	close(this.stopChan)
}
