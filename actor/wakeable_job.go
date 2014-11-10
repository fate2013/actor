package actor

import (
	"encoding/json"
	"fmt"
	"time"
)

type Job struct {
	Uid   int64 `json:"uid"`
	JobId int64 `json:"job_id"`

	CityId    int64     `json:"-"` // ignored json field, json:"myname,omitempty"
	Type      uint16    `json:"-"`
	TimeStart time.Time `json:"-"`
	TimeEnd   time.Time `json:"-"`
	Trace     string    `json:"-"`
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

func (this *Job) GetUid() int64 {
	return this.Uid
}

func (this Job) String() string {
	return fmt.Sprintf("Job{uid:%d, jid:%d, type:%d}",
		this.Uid, this.JobId, this.Type)
}
