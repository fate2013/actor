package actor

import (
	rest "github.com/funkygao/dragon/http"
	"github.com/funkygao/dragon/server"
	log "github.com/funkygao/log4go"
	"github.com/gorilla/mux"
	"net/http"
	"runtime"
	"time"
)

func (this *Actor) showConsoleStats() {
	log.Info("ver: %s, elapsed:%s, sess:%d, req:%d, jobs:%d, goroutine:%d",
		server.BuildID,
		time.Since(this.server.StartedAt),
		this.totalSessionN,
		this.totalReqN,
		this.jobs.Len(),
		runtime.NumGoroutine())
}

func (this *Actor) launchHttpServ() {
	var restListenAddr string = this.server.String("rest_listen_addr", "")
	if restListenAddr == "" {
		return
	}

	rest.LaunchHttpServ(restListenAddr, this.server.String("prof_listen_addr", ""))
	rest.RegisterHttpApi("/s/{cmd}",
		func(w http.ResponseWriter, req *http.Request,
			params map[string]interface{}) (interface{}, error) {
			return this.handleHttpQuery(w, req, params)
		}).Methods("GET")
}

func (this *Actor) stopHttpServ() {
	rest.StopHttpServ()
}

func (this *Actor) handleHttpQuery(w http.ResponseWriter, req *http.Request,
	params map[string]interface{}) (interface{}, error) {
	var (
		vars   = mux.Vars(req)
		cmd    = vars["cmd"]
		output = make(map[string]interface{})
	)

	switch cmd {
	case "ver":
		output["ver"] = server.BuildID

	case "stat":
		output["job_size"] = this.jobs.Len()

	case "conf":
		output["conf"] = *this.server.Conf

	case "guide", "help", "h":
		output["uris"] = []string{
			"/s/stat",
			"/s/conf",
		}
		pprofAddr := this.server.String("prof_listen_addr", "")
		if pprofAddr != "" {
			output["pprof"] = "http://" + pprofAddr + "/debug/pprof/"
		}

	default:
		return nil, rest.ErrHttp404
	}

	return output, nil
}
