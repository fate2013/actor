package actor

import (
	"github.com/funkygao/actor/config"
	log "github.com/funkygao/log4go"
    proxy "github.com/fate2013/go-rtm/proxy"
    "math/rand"
    "strconv"
)

const (
    msg_type_general = 100
    rtm_user_system = 0
)

type WorkerRtm struct {
	config *config.ConfigWorkerRtm
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
    log.Debug("rtm: +v", push)
    this.backlog <- push
}

func (this *WorkerRtm) doPublish(push *Push) {
    log.Debug("push.body in rtm: ID[%d], body[%s]", push.id, push.Body)
    rtmClient, err := proxy.NewRtmClient(this.config.PrimaryHosts[0])
    if err != nil {
        //TODO
    }
    msg, channels := push.SplitMsgAndChannels(string(push.Body))
    log.Debug("message : %s", msg)
    log.Debug("push to : %s", channels)
    
    for _, channel := range channels {
        go func(channel string) {
            intChannel, _ := strconv.ParseInt(channel, 0, 0)
            rtmClient.SendMsg(int32(this.config.ProjectId), this.config.SecretKey, msg_type_general, rtm_user_system, intChannel, this.midgen(), msg)
        }(channel)
    }
}

func (this *WorkerRtm) midgen() int64 {
    return rand.Int63()
}
