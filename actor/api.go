package actor

import (
	"github.com/funkygao/golib/cache"
	"github.com/funkygao/golib/server"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ApiRunner struct {
	listenAddr string

	userFlight *Flight
	tileFlight *Flight
}

func NewApiRunner(listenAddr string, userFlight, tileFlight *Flight) *ApiRunner {
	this := new(ApiRunner)
	this.listenAddr = listenAddr
	this.userFlight = userFlight
	this.tileFlight = tileFlight
	return this
}

func (this *ApiRunner) Run() {
	server.LaunchHttpServ(this.listenAddr, "")
	server.RegisterHttpApi("/{op}/{type}/{id}",
		func(w http.ResponseWriter, req *http.Request,
			params map[string]interface{}) (interface{}, error) {
			return this.handleHttpQuery(w, req, params)
		}).Methods("GET")
}

func (this *ApiRunner) handleHttpQuery(w http.ResponseWriter, req *http.Request,
	params map[string]interface{}) (interface{}, error) {
	var (
		vars = mux.Vars(req)
		op   = vars["op"]   // lock | unlock
		typ  = vars["type"] // user | tile

		output = make(map[string]interface{})
		key    cache.Key
		flight *Flight
	)

	switch typ {
	case API_TYPE_USER:
		uid, err := strconv.Atoi(vars["id"])
		if err != nil {
			return nil, err
		}

		flight = this.userFlight
		key = User{Uid: int64(uid)}

	case API_TYPE_TILE:
		geohash, err := strconv.Atoi(vars["id"])
		if err != nil {
			return nil, err
		}

		flight = this.tileFlight
		key = Tile{Geohash: geohash}

	default:
		output["ok"] = false
		output["msg"] = "invalid type"
		return output, nil
	}

	switch op {
	case API_OP_LOCK:
		output["ok"] = flight.Takeoff(key)

	case API_OP_UNLOCK:
		flight.Land(key)
		output["ok"] = true

	default:
		output["ok"] = false
		output["msg"] = "invalid operation"
	}

	return output, nil
}
