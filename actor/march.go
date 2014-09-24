package actor

import (
	"encoding/json"
	"github.com/cookieo9/go-misc/slice"
	//"time"
)

type March struct {
	Uid     int64 `json:"uid"`
	MarchId int64 `json:"march_id"`
	//CityId    int64
	//Type      string
	//State     string
	//X0        int16
	//Y0        int16
	X1 int16
	Y1 int16
	//StartTime time.Time
	//EndTime   time.Time
}

func (this *March) GeoHash() int {
	return int(this.X1<<11 | this.Y1)
}

func (this *March) Marshal() []byte {
	b, _ := json.Marshal(*this)
	return b
}

type MarchGroup []March

func (this *MarchGroup) sortByDestination() {
	byDestination := func(m1, m2 interface{}) bool {
		return m1.(March).X1 < m2.(March).X1
	}
	slice.Sort(this, byDestination)
}
