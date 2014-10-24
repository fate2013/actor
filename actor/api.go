package actor

import (
	"github.com/funkygao/golib/server"
	"net/http"
)

type ApiRunner struct {
}

func (this *ApiRunner) Run() {
	server.LaunchHttpServ(":9898", "") // TODO
	server.RegisterHttpApi("/q/{uid}",
		func(w http.ResponseWriter, req *http.Request,
			params map[string]interface{}) (interface{}, error) {
			return this.handleHttpQuery(w, req, params)
		}).Methods("GET")
}

func (this *ApiRunner) handleHttpQuery(w http.ResponseWriter, req *http.Request,
	params map[string]interface{}) (output interface{}, err error) {
	return
}
