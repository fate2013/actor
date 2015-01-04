package actor

import (
	"github.com/kr/beanstalk"
	"time"
)

type Push struct {
	conn *beanstalk.Conn
	id   uint64

	Uid  int64
	Body []byte
}

func (this *Push) String() string {
	return ""
}

func (this *Push) DueTime() time.Time {
	return time.Now()
}

func (this *Push) Marshal() []byte {
	return nil
}

func (this *Push) Ignored() bool {
	return false
}

func (this *Push) GetUid() int64 {
	return this.Uid
}
