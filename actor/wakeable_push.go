package actor

import (
	"github.com/kr/beanstalk"
	"time"
    "strings"
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

func (this *Push) SplitMsgAndChannels(body string) (msg string, channels []string) {
    endChannelPos := int64(strings.Index(body, "|"))
    channels = strings.Split(string(body[:endChannelPos]), ",")
    msg = string(body[endChannelPos+1:])
    return
}
