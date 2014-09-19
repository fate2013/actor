package actor

import (
	"time"
)

type job struct {
	Uid     int64 `json:"uid"`
	JobId   int64 `json:"job_id"`
	dueTime time.Time
}
