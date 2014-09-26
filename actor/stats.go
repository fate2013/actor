package actor

import (
	"github.com/funkygao/golib/gofmt"
	"github.com/funkygao/golib/server"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"github.com/gorilla/mux"
	"io"
	logger "log"
	"net/http"
	"os"
	"runtime"
	"syscall"
	"time"
)

type StatsRunner struct {
	actor     *Actor
	scheduler *Scheduler
}

func NewStatsRunner(actor *Actor, scheduler *Scheduler) *StatsRunner {
	this := new(StatsRunner)
	this.actor = actor
	this.scheduler = scheduler
	return this
}

func (this *StatsRunner) Run() {
	this.launchHttpServ()
	defer this.stopHttpServ()

	var (
		metricsWriter io.Writer
		err           error
	)
	if this.actor.config.MetricsLogfile == "" ||
		this.actor.config.MetricsLogfile == "stdout" {
		metricsWriter = os.Stdout
	} else {
		metricsWriter, err = os.OpenFile(this.actor.config.MetricsLogfile,
			os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
	}
	go metrics.Log(metrics.DefaultRegistry, this.actor.config.ConsoleStatsInterval,
		logger.New(metricsWriter, "", logger.LstdFlags))

	ticker := time.NewTicker(this.actor.config.ConsoleStatsInterval)
	defer ticker.Stop()

	var (
		ms           = new(runtime.MemStats)
		rusage       = &syscall.Rusage{}
		lastUserTime int64
		lastSysTime  int64
		userTime     int64
		sysTime      int64
		userCpuUtil  float64
		sysCpuUtil   float64
	)
	for _ = range ticker.C {
		runtime.ReadMemStats(ms)

		syscall.Getrusage(syscall.RUSAGE_SELF, rusage)
		syscall.Getrusage(syscall.RUSAGE_SELF, rusage)
		userTime = rusage.Utime.Sec*1000000000 + int64(rusage.Utime.Usec)
		sysTime = rusage.Stime.Sec*1000000000 + int64(rusage.Stime.Usec)
		userCpuUtil = float64(userTime-lastUserTime) * 100 / float64(this.actor.config.ConsoleStatsInterval)
		sysCpuUtil = float64(sysTime-lastSysTime) * 100 / float64(this.actor.config.ConsoleStatsInterval)

		lastUserTime = userTime
		lastSysTime = sysTime

		log.Info("ver:%s, elapsed:%s, backlog:%d, flight:%d, goroutine:%d, mem:%s, cpu:%3.2f%%us,%3.2f%%sy",
			server.BuildID,
			time.Since(this.actor.server.StartedAt),
			this.scheduler.Outstandings(),
			this.scheduler.InFlight(),
			runtime.NumGoroutine(),
			gofmt.ByteSize(ms.Alloc),
			userCpuUtil,
			sysCpuUtil)
	}
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
	case "ping":
		output["status"] = "ok"

	case "ver":
		output["ver"] = server.BuildID

	case "debug":
		stack := make([]byte, 1<<20)
		stackSize := runtime.Stack(stack, true)
		output["callstack"] = string(stack[:stackSize])

	case "stat":
		output["goroutines"] = runtime.NumGoroutine()

		memStats := new(runtime.MemStats)
		runtime.ReadMemStats(memStats)
		output["memory"] = *memStats

		rusage := syscall.Rusage{}
		syscall.Getrusage(0, &rusage)
		output["rusage"] = rusage

	case "conf":
		output["conf"] = *this.actor.server.Conf

	default:
		return nil, server.ErrHttp404
	}

	output["links"] = []string{
		"/s/ping",
		"/s/ver",
		"/s/conf",
		"/s/stat",
		"/s/debug",
	}
	if this.actor.config.ProfListenAddr != "" {
		output["pprof"] = "http://" +
			this.actor.config.ProfListenAddr + "/debug/pprof/"
	}

	return output, nil
}
