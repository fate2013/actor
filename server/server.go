package server

import (
	conf "github.com/funkygao/jsconf"
	"time"
)

type Server struct {
	*conf.Conf

	configFile string
	StartedAt  time.Time
	pid        int
	hostname   string
}

func NewServer() (this *Server) {
	this = new(Server)

	return
}
