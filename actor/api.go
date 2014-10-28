package actor

import (
	"github.com/funkygao/golib/server"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ApiRunner struct {
	userFlight *Flight
}

func (this *ApiRunner) Run() {
	server.LaunchHttpServ(":9898", "") // TODO
	server.RegisterHttpApi("/{op}/{uid}",
		func(w http.ResponseWriter, req *http.Request,
			params map[string]interface{}) (interface{}, error) {
			return this.handleHttpQuery(w, req, params)
		}).Methods("GET")
}

func (this *ApiRunner) handleHttpQuery(w http.ResponseWriter, req *http.Request,
	params map[string]interface{}) (interface{}, error) {
	var (
		vars   = mux.Vars(req)
		op     = vars["op"] // lock | unlock
		output = make(map[string]interface{})
	)

	uid, err := strconv.Atoi(vars["uid"])
	if err != nil {
		return nil, err
	}

	switch op {
	case API_OP_LOCK:
		output["ok"] = this.userFlight.Takeoff(int64(uid))
	case API_OP_UNLOCK:
		this.userFlight.Land(int64(uid), true)
		output["ok"] = true
	}
	return output, nil
}
