package actor

import (
	"sort"
	"sync"
)

type march struct {
	Uid     int64  `json:"uid"`
	MarchId int64  `json:"march_id"`
	Optime  int    `json:"optime,omitempty"`
	Op      string `json:"op,omitempty"`
}

// sorted map
type marches struct {
	lock sync.Mutex

	m map[int]march
	k []int
}

func newMarches() *marches {
	this := new(marches)
	this.m = make(map[int]march)
	return this
}

func (this *marches) Len() int {
	return len(this.m)
}

func (this *marches) Less(i, j int) bool {
	return this.m[this.k[i]].Optime < this.m[this.k[j]].Optime

}

func (this *marches) Swap(i, j int) {
	this.k[i], this.k[j] = this.k[j], this.k[i]
}

func (this *marches) sortedKeys() []int {
	this.k = make([]int, len(this.m))

	i := 0
	for k, _ := range this.m {
		this.k[i] = k
		i++
	}

	sort.Sort(this)
	return this.k

}

func (this *marches) set(march march) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.m[march.Optime] = march
}

func (this *marches) del(march march) {
	delete(this.m, march.Optime)
}
