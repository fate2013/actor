package actor

import (
	"github.com/funkygao/golib/cache"
)

// used as locking Wakeable's
type Flight struct {
	lock *cache.LruCache

	maxRetries int
	retry      *cache.LruCache
}

func NewFlight(maxLockEntries int, maxRetryEntries int, maxRetries int) *Flight {
	this := new(Flight)
	this.lock = cache.NewLruCache(maxLockEntries)
	this.maxRetries = maxRetries
	this.retry = cache.NewLruCache(maxRetryEntries)
	return this
}

func (this *Flight) CanPass(key cache.Key) (ok, firstTimeFail bool) {
	ok, firstTimeFail = true, true
	retried := this.retry.Inc(key)
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
	// FIXME Get and Set is not atomic
	if _, present := this.lock.Get(key); !present {
		this.lock.Set(key, true)
		return true
	}

	return false
}

func (this *Flight) Land(key cache.Key) {
	this.lock.Del(key)
}

func (this *Flight) Len() int {
	return this.lock.Len()
}
