package main

import (
	"github.com/funkygao/dragon/server"
	"github.com/funkygao/golib/syslogng"
)

func main() {
	server := server.NewServer("proxyd")
	server.LoadConfig("etc/prorxyd.cf")
	server.Launch()

	mainloop(server)
}

func mainloop(server *server.Server) {
	pool := newDragonPool()
	pool.loadConfig(server.Conf)

	in := syslogng.Subscribe()
	for req := range in {
		pool.dispatch([]byte(req))
	}

}
