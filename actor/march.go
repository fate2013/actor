package actor

// actor schedules march event: arrive(speedup), recall
type march struct {
	Uid     int64  `json:"uid"`
	MarchId int64  `json:"march_id"`
	At      int    `json:"at,omitempty"`
	Evt     string `json:"evt,omitempty"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
}
