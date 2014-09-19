package actor

import (
	"github.com/funkygao/dragon/server"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"github.com/gorilla/mux"
	logger "log"
	"net/http"
	"os"
	"runtime"
	"time"
)

type StatsRunner struct {
	actor *Actor
}

func NewStatsRunner(actor *Actor) *StatsRunner {
	this := new(StatsRunner)
	this.actor = actor

	return this
}

func (this *StatsRunner) Run() {
	this.launchHttpServ()
	defer this.stopHttpServ()

	go metrics.Log(metrics.DefaultRegistry, this.actor.config.ConsoleStatsInterval,
		logger.New(os.Stderr, "", logger.LstdFlags))

	ticker := time.NewTicker(this.actor.config.ConsoleStatsInterval)
	defer ticker.Stop()

	for _ = range ticker.C {
		this.showStats()
	}
}

func (this *StatsRunner) showStats() {
	log.Info("ver: %s, elapsed:%s, goroutine:%d",
		server.BuildID,
		time.Since(this.actor.server.StartedAt),
		runtime.NumGoroutine())
}

func (this *StatsRunner) launchHttpServ() {
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

func (this *StatsRunner) stopHttpServ() {
	server.StopHttpServ()
}

func (this *StatsRunner) handleHttpQuery(w http.ResponseWriter, req *http.Request,
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
