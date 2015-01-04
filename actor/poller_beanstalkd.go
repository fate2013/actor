package actor

import (
	log "github.com/funkygao/log4go"
	"github.com/kr/beanstalk"
	"time"
)

type BeanstalkdPoller struct {
	conn *beanstalk.Conn
}

func NewBeanstalkdPoller(addr string) (this *BeanstalkdPoller, err error) {
	this = new(BeanstalkdPoller)
	this.conn, err = beanstalk.Dial("tcp", addr)
	return
}

func (this *BeanstalkdPoller) Poll(ch chan<- Wakeable) {
	var (
		id   uint64
		body []byte
		err  error
		push *Push
	)
	for {
		id, body, err = this.conn.Reserve(time.Hour * 100) // TODO
		if err != nil {
			log.Error("beanstalk.reserve: %v", err)
			continue
		}

		push = new(Push) // TODO mem pool
		push.Uid = int64(id)
		push.Body = body

		ch <- push
	}
}

func (this *BeanstalkdPoller) Stop() {
	this.conn.Close()
}
