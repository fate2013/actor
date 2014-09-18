package actor

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"net"
)

type replicator struct {
	addr       string
	tcpNoDelay bool
}

func NewReplicator() *replicator {
	this := new(replicator)
	return this
}

func (this *replicator) LoadConfig(cf *conf.Conf) *replicator {
	this.addr = cf.String("addr", "")
	this.tcpNoDelay = cf.Bool("tcp_nodelay", true)
	return this
}

func (this *replicator) Replay() {

}

func (this *replicator) enabled() bool {
	return this.addr != ""
}

func (this *replicator) Start() {
	if !this.enabled() {
		log.Info("replication disabled")
		return
	}

	go this.runAcceptor()
	go this.runSender()

	log.Info("replication started at %s", this.addr)
}

func (this *replicator) runAcceptor() {
	listener, err := net.Listen("tcp4", this.addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("accept error: %s", err.Error())
			continue
		}

		defer conn.Close()

	}

}

// todo add a chan in Actor, acceptor -> chan -> replicatorSender
func (this *replicator) runSender() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", this.addr)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}

	defer func() {
		log.Info("replicate sender died")

		conn.Close()
	}()

	if this.tcpNoDelay {
		conn.SetNoDelay(true)
	}

}
