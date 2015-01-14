package actor

import (
	log "github.com/funkygao/log4go"
	"github.com/kr/beanstalk"
	"time"
)

type BeanstalkdPoller struct {
	conn *beanstalk.Conn
}

func NewBeanstalkdPoller(addr string, watchTubes ...string) (this *BeanstalkdPoller, err error) {
	this = new(BeanstalkdPoller)
	this.conn, err = beanstalk.Dial("tcp", addr)
    this.conn.TubeSet = *beanstalk.NewTubeSet(this.conn, watchTubes...)
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
        this.conn.Delete(id) // FIXME
		if err != nil {
			log.Error("beanstalk.reserve: %v", err)
			continue
		}

		push = new(Push) // TODO mem pool
		push.Uid = int64(id)
		push.Body = body
		push.conn = this.conn
		push.id = id

		ch <- push
	}
}

func (this *BeanstalkdPoller) Stop() {
	this.conn.Close()
}
