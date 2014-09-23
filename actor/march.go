package actor

import (
	"time"
)

type March struct {
	Uid       int64
	MarchId   int64
	CityId    int64
	Type      string
	State     string
	X0        int16
	Y0        int16
	X1        int16
	Y1        int16
	StartTime time.Time
	EndTime   time.Time
}
