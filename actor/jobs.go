package actor

import (
	"sort"
	"sync"
	"time"
)

// sorted map
type jobs struct {
	sort.Interface // jobs is sortable

	lock sync.Mutex

	m map[int]march // key is timestamp('at')
	k []int         // array of timestamp('at')
}

func newJobs() *jobs {
	this := new(jobs)
	this.m = make(map[int]march)
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

func (this *jobs) submit(march march) {
	this.lock.Lock()
	this.m[march.At] = march
	this.lock.Unlock()
}

func (this *jobs) wakeup(tillWhen time.Time) []march {
	this.lock.Lock()
	defer this.lock.Unlock()

	r := make([]march, 0)
	for at := range this.sortedKeys() {
		march := this.m[this.k[at]]
		if tillWhen.Unix() >= int64(march.At) {
			r = append(r, march)

			delete(this.m, march.At) // this job is done
		}
	}

	return r
}
