package actor

import (
	"fmt"
	"sort"
	"sync"
)

type march struct {
	uid     int64
	marchId int64
	opTime  int
	op      string
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
	return this.m[this.k[i]].opTime < this.m[this.k[j]].opTime

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

	fmt.Println(this.k)

	sort.Sort(this)
	return this.k

}

func (this *marches) set(march march) {
	this.m[march.opTime] = march

}

func (this *marches) del(march march) {
	delete(this.m, march.opTime)
}
