package actor

import (
	"encoding/json"
	"fmt"
	"time"
)

type Job struct {
	Uid       int64     `json:"uid"`
	JobId     int64     `json:"job_id"`
	CityId    int64     `json:"city_id"`
	Type      uint16    `json:"event_type"`
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	Trace     string    `json:"trace"`
}

func (this *Job) Marshal() []byte {
	b, _ := json.Marshal(*this)
	return b
}

func (this *Job) DueTime() time.Time {
	return this.TimeEnd
}

func (this *Job) Ignored() bool {
	return false
}

func (this *Job) FlightKey() interface{} {
	return *this
}

func (this Job) String() string {
	return fmt.Sprintf("Job{uid:%d, jid:%d, type:%d, due:%s}",
		this.Uid, this.JobId, this.Type, this.TimeEnd)
}
