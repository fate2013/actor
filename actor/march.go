package actor

/*
actor schedules march event: arrive(speedup), recall(homeBack)

critical sections:
. encamp same tile at roughly the same time
. attack same city at roughly the same time
. arrive at gathering tile at roughly the same time
. wonder? mini wonder?
. TODO see march state diagram

all the critical sections share the same attribute:
same destination(geohash) at roughly the same time
*/
type march struct {
	Uid     int64 `json:"uid"`
	MarchId int64 `json:"march_id"`
	At      int   `json:"at,omitempty"`

	// php.EventConst EVT_PVP_MARCH_ARRIVE | EVT_PVP_GATHER_DONE | EVT_PVP_HOME_BACK
	Evt int `json:"evt,omitempty"`

	X  int   `json:"x"`
	Y  int   `json:"y"`
	T0 int64 `json:"t0,omitempty"` // sent timestamp from php-fpm
}

// given a hash, x=h>>10, y=h&((1<<10)-1)
func (this *march) geoHash() int {
	return (this.X << GEO_HASH_SHIFT) | this.Y
}

func (this *march) ident() marchIdent {
	return marchIdent{this.Uid, this.MarchId}
}

type marchIdent struct {
	Uid, MarchId int64
}
