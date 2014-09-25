package actor

import (
	"encoding/json"
	"fmt"
	"github.com/cookieo9/go-misc/slice"
	"time"
)

type March struct {
	Uid     int64  `json:"uid"`
	MarchId int64  `json:"march_id"`
	State   string `json:"state"`
	X1      int16
	Y1      int16
	EndTime time.Time
}

func (this *March) GeoHash() int {
	return int(this.X1<<11 | this.Y1)
}

func (this *March) Marshal() []byte {
	b, _ := json.Marshal(*this)
	return b
}

func (this *March) FlightKey() interface{} {
	return this.GeoHash()
}

func (this *March) Ignored() bool {
	return this.State == "done" || this.State == "encamping"
}

func (this *March) DueTime() time.Time {
	return this.EndTime
}

func (this March) String() string {
	return fmt.Sprintf("March{uid:%d, mid:%d, due:%s, state:%s}",
		this.Uid, this.MarchId, this.EndTime, this.State)
}

type MarchGroup []March

func (this *MarchGroup) sortByDestination() {
	byDestination := func(m1, m2 interface{}) bool {
		return m1.(March).X1 < m2.(March).X1
	}
	slice.Sort(this, byDestination)
}
