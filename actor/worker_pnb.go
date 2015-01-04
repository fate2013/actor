package actor

import (
	"github.com/funkygao/actor/config"
	log "github.com/funkygao/log4go"
	"github.com/pubnub/go/messaging"
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

func (this *WorkerPnb) doPublish(push *Push) {
	pnb := messaging.NewPubnub(this.config.PublishKey,
		this.config.SubscribeKey, this.config.SecretKey,
		this.config.CipherKey, this.config.UseSsl, "")
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	channel := "pnb_chan"
	go pnb.Publish(channel, push.Body, successChannel, errorChannel)
	select {
	case msg := <-successChannel:
		log.Debug("pnb: %s", string(msg))
		push.conn.Delete(push.id) // ack success

	case err := <-errorChannel:
		push.conn.Bury(push.id, 1) // ack fail
		log.Error("pnb: %s", string(err))
	}

}
