package actor

// actor schedules march event: arrive(speedup), recall
type march struct {
	Uid     int64 `json:"uid"`
	MarchId int64 `json:"march_id"`
	At      int   `json:"at,omitempty"`
	Evt     int   `json:"evt,omitempty"`
	X       int   `json:"x"`
	Y       int   `json:"y"`
	T0      int64 `json:"t0"`
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
