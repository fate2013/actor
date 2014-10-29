package actor

import (
	"github.com/funkygao/golib/cache"
	log "github.com/funkygao/log4go"
)

// used as locking Wakeable's
// TODO add auto expire for a lock
type Flight struct {
	items *cache.LruCache

	maxRetries int
	retries    *cache.LruCache
}

func NewFlight(maxLockEntries int, maxRetryEntries int, maxRetries int) *Flight {
	this := new(Flight)
	this.items = cache.NewLruCache(maxLockEntries)
	this.maxRetries = maxRetries
	this.retries = cache.NewLruCache(maxRetryEntries)
	return this
}

// FIXME for March, key is (x, y), if 100 march to the same tile, max retries will
// be reached early, can refuse to serve the remaining march
func (this *Flight) canPass(key cache.Key) (ok, firstTimeFail bool) {
	ok, firstTimeFail = true, false
	if this.maxRetries == 0 {
		return
	}
	retried := this.retries.Inc(key)
	if retried >= this.maxRetries {
		ok = false

		if retried == this.maxRetries {
			firstTimeFail = true
		}
	}

	return
}

// return true if accquired the lock
func (this *Flight) Takeoff(key cache.Key) (success bool) {
	ok, firstTimeFail := this.canPass(key)
	if !ok {
		if firstTimeFail {
			log.Warn("max retries[%d] reached: %+v", this.maxRetries, key)
		}

		return false
	}

	// FIXME Get and Set is not atomic
	if _, present := this.items.Get(key); !present {
		this.items.Set(key, true)
		return true
	}

	log.Debug("already locked: %#v", key)
	return false
}

func (this *Flight) Land(key cache.Key, ok bool) {
	this.items.Del(key)
	if this.maxRetries > 0 && ok {
		this.retries.Set(key, 0) // reset the counter
	}
}

func (this *Flight) Flying(key cache.Key) bool {
	if _, present := this.items.Get(key); present {
		return true
	}
	return false
}

func (this *Flight) Len() int {
	return this.items.Len()
}
