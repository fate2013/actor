package actor

import (
	"github.com/funkygao/actor/config"
)

type WorkerPnb struct {
}

func NewPnbWorker(config *config.ConfigWorkerPnb) *WorkerPnb {
	this := new(WorkerPnb)
	return this
}

func (this *WorkerPnb) Start() {

}

func (this *WorkerPnb) Wake(w Wakeable) {
}
