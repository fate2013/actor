package actor

import (
	"encoding/json"
	"time"
)

type Pve struct {
	Uid     int64 `json:"uid"`
	MarchId int64 `json:"march_id"`

	State string

	EndTime time.Time
}

func (this *Pve) Marshal() []byte {
	b, _ := json.Marshal(*this)
	return b
}

func (this *Pve) Ignored() bool {
	return this.State == "done"
}
