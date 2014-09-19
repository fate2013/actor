package actor

import (
	"encoding/json"
	"sync"
	"time"
)

type Job struct {
	Uid     int64 `json:"uid"`
	JobId   int64 `json:"job_id"`
	dueTime time.Time
}

func (this *Job) marshal() []byte {
	b, _ := json.Marshal(*this)
	return b
}

type outstandingJobs struct {
	mutex sync.Mutex
	jobs  map[Job]bool
}

func newOutstandingJobs() *outstandingJobs {
	this := new(outstandingJobs)
	this.jobs = make(map[Job]bool)
	return this
}

func (this *outstandingJobs) enter(j Job) {
	this.mutex.Lock()
	this.jobs[j] = true
	this.mutex.Unlock()
}

func (this *outstandingJobs) leave(j Job) {
	this.mutex.Lock()
	delete(this.jobs, j)
	this.mutex.Unlock()
}
