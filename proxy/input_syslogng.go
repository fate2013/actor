package proxy

import (
	"github.com/funkygao/golib/syslogng"
)

type syslogngInput struct {
}

func (ths *syslogngInput) Reader() chan []byte {
	return syslogng.Subscribe()
}
