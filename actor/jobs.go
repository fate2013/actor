package actor

import (
	log "github.com/funkygao/log4go"
	"sort"
	"sync"
)

// sorted map
type jobs struct {
	sort.Interface // jobs is sortable

	lock sync.Mutex

	m map[int]*march        // key is timestamp('at')
	n map[marchIdent]*march // key is hash(uid, march_id)
	k []int                 // array of timestamp('at')
}

func newJobs() *jobs {
	this := new(jobs)
	this.m = make(map[int]*march)
	this.n = make(map[marchIdent]*march)
	return this
}

func (this *jobs) Len() int {
	return len(this.m)
}

func (this *jobs) Less(i, j int) bool {
	return this.m[this.k[i]].At < this.m[this.k[j]].At
}

func (this *jobs) Swap(i, j int) {
	this.k[i], this.k[j] = this.k[j], this.k[i]
}

func (this *jobs) sortedKeys() []int {
	this.k = make([]int, len(this.m))

	i := 0
	for at, _ := range this.m {
		this.k[i] = at
		i++
	}

	sort.Sort(this)
	return this.k
}

func (this *jobs) sched(march march) {
	this.lock.Lock()

	this.m[march.At] = &march
	this.n[march.ident()] = &march
	// first lookup this march
	m, present := this.n[march.ident()]
	if present {
		// modify existent march entry
		delete(this.m, m.At)
		delete(this.n, m.ident())
	}

	log.Debug("jobs: %+v", *this)

	this.lock.Unlock()
}

func (this *jobs) wakeup(tillWhen int64) []march {
	this.lock.Lock()
	defer this.lock.Unlock()

	r := make([]march, 0)
	for at := range this.sortedKeys() {
		march := this.m[this.k[at]]
		dueTime := int64(march.At)
		if tillWhen >= dueTime {
			r = append(r, *march)

			delete(this.m, march.At) // this job is done

			if tillWhen > dueTime {
				// scheduler is late to wake it up
				log.Warn("late schedule march[%d > %d]: %+v", tillWhen, dueTime, *march)
			}

		}

	}

	return r
}
