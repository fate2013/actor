package actor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/cookieo9/go-misc/slice"
	"time"
)

type March struct {
	Uid     int64          `json:"uid"`
	MarchId int64          `json:"march_id"`
	Type    sql.NullString `json:"type"`
	State   string         `json:"-"`
	X1      int16          `json:"-"`
	Y1      int16          `json:"-"`
	EndTime time.Time      `json:"-"`
}

func (this *March) GeoHash() int {
	return int(this.X1)<<11 | int(this.Y1)
}

func (this *March) DecodeGeoHash(hash int) (x, y int16) {
	x = int16(hash >> 11)
	y = int16(((1 << 11) - 1) & hash)
	return
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
	return fmt.Sprintf("March{uid:%d, mid:%d, type:%v, state:%s, (%d, %d), due:%s}",
		this.Uid, this.MarchId, this.Type, this.State, this.X1, this.Y1, this.EndTime)
}

// FIXME not used yet
type MarchGroup []March

func (this *MarchGroup) sortByDestination() {
	byDestination := func(m1, m2 interface{}) bool {
		return m1.(March).X1 < m2.(March).X1
	}
	slice.Sort(this, byDestination)
}
