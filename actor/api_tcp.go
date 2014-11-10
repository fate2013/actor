package actor

import (
	"bufio"
	log "github.com/funkygao/log4go"
	"net"
	"strings"
)

// line based text protocol with each line as a cmd expected a response
//
// php
//  |
//  | short conn
//  | thrift
//  |
// fae
//  |
//  | persist conn
//  | tcp api
//  |
// actord
//
type TcpApiRunner struct {
	listenAddr string

	userFlight *Flight
	tileFlight *Flight
}

func NewTcpApiRunner(listenAddr string, userFlight, tileFlight *Flight) *TcpApiRunner {
	this := new(TcpApiRunner)
	this.listenAddr = listenAddr
	this.userFlight = userFlight
	this.tileFlight = tileFlight
	return this
}

func (this *TcpApiRunner) Run() {
	listener, err := net.Listen("tcp", this.listenAddr)
	if err != nil {
		panic(err)
	}

	for {
		sock, err := listener.Accept()
		if err != nil {
			log.Error("tcp api: %s", err)
			continue
		}

		go this.handleClient(sock)
	}

}

func (this *TcpApiRunner) handleClient(client net.Conn) {
	defer client.Close()

	// TODO
	// read each cmd line by line, and sendback response by line
	for {
		line, _, err := bufio.NewReader(client).ReadLine()
		// line: /{reason}/{op}/{type}/{id}
		if err != nil {
			log.Error("readline: %s", err)
			break
		}

		strings.Split(string(line), "/")
		client.Write([]byte("ok"))
	}

}
