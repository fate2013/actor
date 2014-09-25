package actor

import (
	"encoding/json"
	"fmt"
	"time"
)

type Pve struct {
	Uid     int64 `json:"uid"`
	MarchId int64 `json:"march_id"`

	State string `json:"state"`

	EndTime time.Time
}

func (this *Pve) Marshal() []byte {
	b, _ := json.Marshal(*this)
	return b
}

func (this *Pve) Ignored() bool {
	return this.State == "done"
}

func (this *Pve) DueTime() time.Time {
	return this.EndTime
}

func (this Pve) String() string {
	return fmt.Sprintf("Pve{uid:%d, mid:%d, due:%s, state:%s}",
		this.Uid, this.MarchId, this.EndTime, this.State)
}

func (this *Pve) FlightKey() interface{} {
	return *this
}
