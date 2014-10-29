package actor

import (
	"github.com/funkygao/golib/cache"
	log "github.com/funkygao/log4go"
	"time"
)

// used as locking Wakeable's
// TODO add auto expire for a lock
type Flight struct {
	debug bool

	items *cache.LruCache
}

func NewFlight(maxLockEntries int, maxRetryEntries int, maxRetries int, debug bool) *Flight {
	this := new(Flight)
	this.debug = debug
	this.items = cache.NewLruCache(maxLockEntries)

	return this
}

// return true if accquired the lock
func (this *Flight) Takeoff(key cache.Key) (success bool) {
	// FIXME Get and Set is not atomic
	if _, present := this.items.Get(key); !present {
		this.items.Set(key, time.Now())
		if this.debug {
			log.Debug("locking[%#v]", key)
		}
		return true
	}

	log.Warn("already locked: %#v", key)
	return false
}

func (this *Flight) Land(key cache.Key, ok bool) {
	this.items.Del(key)
	if this.debug {
		log.Debug("unlock[%#v]", key)
	}
}

func (this *Flight) Len() int {
	return this.items.Len()
}
