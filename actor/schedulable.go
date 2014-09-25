package actor

import (
	"time"
)

type Schedulable interface {
	DueTime() time.Time
	Marshal() []byte
	FlightKey() interface{}
	Ignored() bool
}
