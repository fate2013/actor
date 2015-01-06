package actor

import (
	"time"
)

type Wakeable interface {
	String() string
	DueTime() time.Time
	Marshal() []byte
	Ignored() bool
	GetUid() int64 // each Wakeable has a uid
}

type WakeablePhp interface {
}

type Lockable interface {
	LockKey() string
}
