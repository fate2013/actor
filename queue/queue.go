package queue

import (
	log "github.com/funkygao/log4go"
	"github.com/huandu/skiplist"
)

type Queue struct {
	tasks *skiplist.SkipList
}

func New() (this *Queue) {
	this = new(Queue)
	this.tasks = skiplist.New(keyFuncWhen)

	return
}

// 0    1       2     3   4      5     6  7
// type,dueTime,event,uid,cityId,jobId,t0,payload
func (this *Queue) Enque(data []byte) (err error) {
	log.Debug("enque %s", string(data))

	task := newTask()
	err = task.unmarshal(data)
	if err != nil {
		log.Error("unmarshal: %s", err.Error())
		return
	}

	log.Debug("task: %+v", *task)

	e := this.tasks.Get(task.DueTime)
	if e == nil {
		tasks := make([]*Task, 1, 10)
		tasks[0] = task
		this.tasks.Set(task.DueTime, tasks)
	} else {
		e.Value = append(e.Value.([]*Task), task)
		this.tasks.Set(task.DueTime, e)
	}

	return
}

func (this *Queue) Wakeup(tillWhen int64) []*Task {
	tasks := this.tasks.Get(tillWhen)
	if tasks == nil {
		return nil
	}

	return tasks.Value.([]*Task)
}

func (this *Queue) Len() int {
	return this.tasks.Len()
}
