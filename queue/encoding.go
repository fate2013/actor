package queue

import (
	"bytes"
	log "github.com/funkygao/log4go"
	"strconv"
)

func (this *Task) toInt64(data []byte) int64 {
	i, _ := strconv.Atoi(string(data))
	return int64(i)
}

func (this *Task) unmarshal(data []byte) error {
	parts := bytes.SplitN(data, []byte(","), 8)
	if len(parts) != 8 {
		return ERR_INVALID_REQ
	}

	this.Typ = string(parts[0])
	this.DueTime = this.toInt64(parts[1])
	this.Event = this.toInt64(parts[2])
	this.Uid = this.toInt64(parts[3])
	this.CityId = this.toInt64(parts[4])
	this.JobId = this.toInt64(parts[5])
	this.T0 = this.toInt64(parts[6])
	this.Payload = parts[7]

	log.Debug("task unmarshalled: %+v", *this)

	return nil
}

func (this *Task) Marshal() []byte {
	return nil
}
