package actor

import (
	"github.com/funkygao/actor/config"
	log "github.com/funkygao/log4go"
    "github.com/fate2013/pubnub-go/messaging"
)

type WorkerPnb struct {
	config  *config.ConfigWorkerPnb
	backlog chan *Push
}

func NewPnbWorker(config *config.ConfigWorkerPnb) *WorkerPnb {
	this := new(WorkerPnb)
	this.config = config
	this.backlog = make(chan *Push, config.Backlog)
	return this
}

func (this *WorkerPnb) Start() {
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

func (this *WorkerPnb) Wake(w Wakeable) {
	push := w.(*Push)
	log.Debug("pnb: +v", *push)
	this.backlog <- push // TODO timeout
}

// TODO how to get channel and msg body from beanstalk msg
func (this *WorkerPnb) doPublish(push *Push) {
    log.Debug("push.body in pnb: ID[%d], body[%s]", push.id, push.Body)
	pnb := messaging.NewPubnub(this.config.PublishKey,
		this.config.SubscribeKey, this.config.SecretKey,
		this.config.CipherKey, this.config.UseSsl, "")
    msg, channels := push.SplitMsgAndChannels(string(push.Body))
    log.Debug("message : %s", msg)
    log.Debug("push to : %s", channels)
    for _, channel := range channels {
        go func(channel string) {
            successChannel := make(chan []byte)
            errorChannel := make(chan []byte)
            log.Debug("channel : %s", channel)
            go pnb.Publish(channel, msg, successChannel, errorChannel)
            select {
            case msg := <-successChannel:
                //push.conn.Delete(push.id) // ack success
                log.Debug("pnb: %s", string(msg))

            case err := <-errorChannel:
                //push.conn.Bury(push.id, 1) // ack fail
                log.Error("pnb: %s", string(err))
            }
        }(channel)
    }

}

