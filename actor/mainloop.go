package actor

import (
	log "github.com/funkygao/golib/log4go"
	"net"
)

func (this *Actor) Mainloop() {
	// listen for incoming directed request and put into chan
	this.run()

	for {
		// loop through the requests and dispatch to tile bucket handlers

	}

}

func (this *Actor) run() {
	listener, err := net.Listen("tcp4", this.server.String("listen_addr", ":9002"))
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Error(err)
				continue
			}

			go this.handleRequest(conn)

		}
	}()

}

func (this *Actor) handleRequest(client net.Conn) {

}
