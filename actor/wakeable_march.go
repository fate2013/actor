package actor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// TODO add K1
type March struct {
	Uid     int64          `json:"uid"`
	MarchId int64          `json:"march_id"`
	Type    sql.NullString `json:"type"`
	OppUid  sql.NullInt64  `json:"-"`
	State   string         `json:"-"`
	X1      int16          `json:"-"`
	Y1      int16          `json:"-"`
	EndTime time.Time      `json:"-"`
}

func (this *March) GeoHash() int {
	return int(this.X1)<<GEOHASH_SHIFT | int(this.Y1)
}

func (this *March) DecodeGeoHash(hash int) (x, y int16) {
	x = int16(hash >> GEOHASH_SHIFT)
	y = int16(((1 << GEOHASH_SHIFT) - 1) & hash)
	return
}

func (this *March) TileKey() Tile {
	return Tile{Geohash: this.GeoHash()}
}

func (this *March) Marshal() []byte {
	b, _ := json.Marshal(*this)
	return b
}

func (this *March) GetUid() int64 {
	return this.Uid
}

func (this *March) Ignored() bool {
	return this.State == MARCH_DONE || this.State == MARCH_ENCAMP
}

// rally里有自己向自己行军的情况，这时像自己打自己，只锁该uid一次，否则锁失败
func (this *March) IsRally() bool {
	return this.Type.Valid &&
		this.Type.String == "rally" &&
		this.OppUid.Valid &&
		this.Uid == this.OppUid.Int64
}

func (this *March) DueTime() time.Time {
	return this.EndTime
}

func (this March) String() string {
	return fmt.Sprintf("March{uid:%d, opp:%d, mid:%d, type:%s, state:%s, (%d, %d)}",
		this.Uid, this.OppUid.Int64,
		this.MarchId, this.Type.String, this.State, this.X1, this.Y1)
}
