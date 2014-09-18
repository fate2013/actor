package queue

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

func (this *march) geoHash() int {
	return this.where()
}

func (this *march) ident() marchIdent {
	return marchIdent{this.Uid, this.MarchId}
}

// given a hash, x=h>>10, y=h&((1<<10)-1)
func (this *march) where() int {
	// TODO what if negative?
	return (this.X << GEO_HASH_SHIFT) | this.Y
}
func (this *march) who() marchIdent {
	return this.ident()
}

func (this *march) when() int {
	return this.At
}

func (this *march) what() int {
	return this.Evt
}

type marchIdent struct {
	Uid, MarchId int64
}
