package actor

import (
	"time"
)

type Wakeable interface {
	DueTime() time.Time
	Marshal() []byte
	FlightKey() interface{}
	Ignored() bool
}
