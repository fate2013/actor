package actor

import (
	"github.com/funkygao/dragon/server"
)

type Actor struct {
	server *server.Server
}

func NewActor(server *server.Server) *Actor {
	actor := new(Actor)
	actor.server = server
	return actor

}
