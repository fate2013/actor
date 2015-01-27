package actor

import (
	"github.com/kr/beanstalk"
	"strconv"
	"strings"
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

func (this *Push) Unmarshal(body string) (msg string, from int64, channels []string) {
	arrBody := strings.SplitN(body, "|", 3)
	strChannels := arrBody[0]
	from, _ = strconv.ParseInt(arrBody[1], 0, 0)
	channels = strings.Split(strChannels, ",")
	msg = string(arrBody[2])
	return
}
