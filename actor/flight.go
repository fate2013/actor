package actor

import (
	"github.com/funkygao/golib/cache"
)

type Flight struct {
	entries *cache.LruCache
}

func NewFlight(maxEntries int) *Flight {
	this := new(Flight)
	this.entries = cache.NewLruCache(maxEntries)
	return this
}

// return true if accquired the lock
func (this *Flight) Takeoff(key cache.Key) bool {
	// FIXME Get and Set is not atomic
	if _, present := this.entries.Get(key); !present {
		this.entries.Set(key, true)
		return true
	}

	return false
}

func (this *Flight) Land(key cache.Key) {
	this.entries.Del(key)
}
