package actor

import (
	"github.com/funkygao/golib/cache"
	log "github.com/funkygao/log4go"
	"sync"
	"time"
)

// used as locking Wakeable's
// TODO add auto expire for a lock
type Flight struct {
	debug   bool
	expires time.Duration

	mutex sync.Mutex
	items *cache.LruCache // key: ctime
}

func NewFlight(maxLockEntries int, debug bool, expires time.Duration) *Flight {
	this := new(Flight)
	this.debug = debug
	this.expires = expires
	this.items = cache.NewLruCache(maxLockEntries)

	return this
}

// return true if accquired the lock
func (this *Flight) Takeoff(key cache.Key) (success bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	ctime, present := this.items.Get(key)
	if !present {
		this.items.Set(key, time.Now())

		if this.debug {
			log.Debug("locking: %+v", key)
		}

		return true
	}

	// present, check expires
	if this.expires > 0 && time.Since(ctime.(time.Time)) > this.expires {
		log.Warn("expires[%+v]: %s", key, time.Since(ctime.(time.Time)))

		// refresh the lock
		this.items.Set(key, time.Now())
		return true
	}

	log.Warn("in flight: %+v", key)
	return false
}

func (this *Flight) Land(key cache.Key) {
	this.items.Del(key)
	if this.debug {
		log.Debug("unlock: %+v", key)
	}
}

func (this *Flight) Len() int {
	return this.items.Len()
}
