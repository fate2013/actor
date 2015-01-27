package actor

import (
	proxy "github.com/fate2013/go-rtm/proxy"
	"github.com/funkygao/actor/config"
	"github.com/funkygao/golib/idgen"
	log "github.com/funkygao/log4go"
	"math/rand"
	"strconv"
)

const (
	msg_type_general = 100
)

const (
	ticket_placeholder int64 = iota
	ticket_user
	ticket_alliance
	ticket_chat_room
)

type WorkerRtm struct {
	config  *config.ConfigWorkerRtm
	backlog chan *Push
}

func NewRtmWorker(config *config.ConfigWorkerRtm) *WorkerRtm {
	this := new(WorkerRtm)
	this.config = config
	this.backlog = make(chan *Push, config.Backlog)
	return this
}

func (this *WorkerRtm) Start() {
	for i := 0; i < this.config.MaxProcs; i++ {
		go func() {
			for {
				select {
				case push := <-this.backlog:
					this.doPublish(push)
				}
			}
		}()
	}
}

func (this *WorkerRtm) Wake(w Wakeable) {
	push := w.(*Push)
	log.Debug("rtm: +v", *push)
	this.backlog <- push
}

func (this *WorkerRtm) doPublish(push *Push) {
	log.Debug("push.body in rtm: ID[%d], body[%s]", push.id, push.Body)
	rtmClient, err := proxy.NewRtmClient(this.config.PrimaryHosts[0])
	if err != nil {
		//TODO
	}
	msg, from, channels := push.Unmarshal(string(push.Body))
	log.Debug("message : %s", msg)
	log.Debug("push to : %s", channels)

	// TODO -- can optimise to SendMsgs
	for _, channel := range channels {
		go func(channel string) {
			log.Debug("channel : %s", channel)
			intChannel, _ := strconv.ParseInt(channel, 0, 0)
			switch this.channelType(intChannel) {
			case 0:
				rtmClient.SendMsg(int32(this.config.ProjectId), this.config.SecretKey, msg_type_general, from, intChannel, this.midgen(), msg)
				log.Debug("send")
			case 1:
				rtmClient.SendGroupMsg(int32(this.config.ProjectId), this.config.SecretKey, msg_type_general, from, intChannel, this.midgen(), msg)
				log.Debug("send group")
			}
		}(channel)
	}
}

func (this *WorkerRtm) midgen() int64 {
	return rand.Int63()
}

// 0:send, 1:sendGroup, 2:sendAll
func (this *WorkerRtm) channelType(channel int64) int {
	_, tag, _, _ := idgen.DecodeId(channel)
	if tag == ticket_user {
		return 0
	} else {
		return 1
	}
}
