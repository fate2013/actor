package actor

import (
	"github.com/funkygao/golib/cache"
)

type Outstandings struct {
	cache *cache.LruCache
}

func NewOutstandings(maxEntries int) *Outstandings {
	this := new(Outstandings)
	this.cache = cache.NewLruCache(maxEntries)
	return this
}

// return true if accquired the lock
func (this *Outstandings) Lock(key cache.Key) bool {
	if _, present := this.cache.Get(key); !present {
		this.cache.Set(key, true)
		return true
	}

	return false
}

func (this *Outstandings) Unlock(key cache.Key) {
	this.cache.Del(key)
}
