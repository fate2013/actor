package actor

import (
	"sync"
	"time"
)

type job struct {
	Uid     int64 `json:"uid"`
	JobId   int64 `json:"job_id"`
	dueTime time.Time
}

type outstandingJobs struct {
	mutex sync.Mutex
	jobs  map[job]bool
}

func newOutstandingJobs() *outstandingJobs {
	this := new(outstandingJobs)
	this.jobs = make(map[job]bool)
	return this
}

func (this *outstandingJobs) enter(j job) {
	this.mutex.Lock()
	this.jobs[j] = true
	this.mutex.Unlock()
}

func (this *outstandingJobs) exit(j job) {
	this.mutex.Lock()
	delete(this.jobs, j)
	this.mutex.Unlock()
}
