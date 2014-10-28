package actor

import (
	"github.com/funkygao/golib/server"
	"github.com/gorilla/mux"
	"net/http"
)

type ApiRunner struct {
	jobFlight   *Flight
	marchFlight *Flight
	pveFlight   *Flight
}

func (this *ApiRunner) Run() {
	server.LaunchHttpServ(":9898", "") // TODO
	server.RegisterHttpApi("/lock/{uid}",
		func(w http.ResponseWriter, req *http.Request,
			params map[string]interface{}) (interface{}, error) {
			return this.handleHttpQuery(w, req, params)
		}).Methods("GET")
}

func (this *ApiRunner) handleHttpQuery(w http.ResponseWriter, req *http.Request,
	params map[string]interface{}) (interface{}, error) {
	var (
		vars   = mux.Vars(req)
		uid    = vars["uid"]
		output = make(map[string]interface{})
	)
	output["uid"] = uid
	output["ok"] = 1
	return output, nil
}
