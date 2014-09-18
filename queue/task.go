package queue

/*
job:
	uid, city_id, job_id,   evt, trace,   due_time, t0
march:
    uid, city_id, march_id, evt, k, x, y, due_time, t0

1. found same marchId, then delete old
2. sort by due_time
list(when, March)
list(jobId, March)

list(when, Job)
list(jobId, Job)
*/
type Task struct {
	Typ     string
	DueTime int64
	Event   int64 // event type, php(MarchConst, JobConst)
	Uid     int64
	CityId  int64
	JobId   int64 // march_id or job_id
	T0      int64 // for debugging latency between actord and php
	Payload []byte

	done chan bool
}

func newTask() (this *Task) {
	this = new(Task)
	this.done = make(chan bool)
	return
}

func (this *Task) Terminate() {
	close(this.done)
}
