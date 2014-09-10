package actor

import (
	"sort"
	"sync"
	"time"
)

// sorted map
type jobs struct {
	lock sync.Mutex

	m map[int]march
	k []int
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
	for k, _ := range this.m {
		this.k[i] = k
		i++
	}

	sort.Sort(this)
	return this.k

}

func (this *jobs) set(march march) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.m[march.At] = march
}

func (this *jobs) del(march march) {
	delete(this.m, march.At)
}

func (this *jobs) pullInBatch(t time.Time) []march {
	r := make([]march, 0)
	for optime := range this.sortedKeys() {
		march := this.m[this.k[optime]]
		if t.Unix() >= int64(march.At) {
			r = append(r, march)

		}
	}

	return r
}
