package actor

import (
	mq "github.com/funkygao/lentil"
	log "github.com/funkygao/log4go"
)

type MqCallbacker struct {
	beanstalkd *mq.Beanstalkd
}

func NewMqCallbacker(addr string) (this *MqCallbacker) {
	this = new(MqCallbacker)
	var err error
	this.beanstalkd, err = mq.Dial(addr)
	if err != nil {
		log.Error("mq: %s", err.Error())
		return nil
	}

	return
}

func (this *MqCallbacker) Call(j Job) (retry bool) {
	jobId, err := this.beanstalkd.Put(0, 0, 60, j.Marshal())
	if err != nil {
		log.Error("mq put: %s", err.Error())

		retry = true
		return
	}

	log.Debug("mq jobId: %d", jobId)
	return
}
