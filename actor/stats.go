package actor

import (
	"github.com/funkygao/dragon/server"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"github.com/gorilla/mux"
	"net/http"
	"runtime"
	"time"
)

type actorStats struct {
	actor             *Actor
	dbLatencies       metrics.Histogram
	callbackLatencies metrics.Histogram
}

func newActorStats(actor *Actor) *actorStats {
	this := new(actorStats)
	this.actor = actor
	this.dbLatencies = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))

	this.callbackLatencies = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))

	return this
}

func (this *actorStats) init() {
	metrics.Register("latency.db", this.dbLatencies)
	metrics.Register("latency.callback", this.callbackLatencies)
}

func (this *actorStats) showConsoleStats() {
	log.Info("ver: %s, elapsed:%s, goroutine:%d",
		server.BuildID,
		time.Since(this.actor.server.StartedAt),
		runtime.NumGoroutine())
}

func (this *actorStats) launchHttpServ() {
	if this.actor.config.RestListenAddr == "" {
		return
	}

	server.LaunchHttpServ(this.actor.config.RestListenAddr, this.actor.config.ProfListenAddr)
	server.RegisterHttpApi("/s/{cmd}",
		func(w http.ResponseWriter, req *http.Request,
			params map[string]interface{}) (interface{}, error) {
			return this.handleHttpQuery(w, req, params)
		}).Methods("GET")
}

func (this *actorStats) stopHttpServ() {
	server.StopHttpServ()
}

func (this *actorStats) handleHttpQuery(w http.ResponseWriter, req *http.Request,
	params map[string]interface{}) (interface{}, error) {
	var (
		vars   = mux.Vars(req)
		cmd    = vars["cmd"]
		output = make(map[string]interface{})
	)

	switch cmd {
	case "ver":
		output["ver"] = server.BuildID

	case "conf":
		output["conf"] = *this.actor.server.Conf

	case "guide", "help", "h":
		output["uris"] = []string{
			"/s/ver",
			"/s/conf",
		}
		if this.actor.config.ProfListenAddr != "" {
			output["pprof"] = "http://" + this.actor.config.ProfListenAddr + "/debug/pprof/"
		}

	default:
		return nil, server.ErrHttp404
	}

	return output, nil
}
