package actor

import (
	"github.com/funkygao/golib/cache"
	log "github.com/funkygao/log4go"
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

// FIXME for March, key is (x, y), if 100 march to the same tile, max retries will
// be reached early, can refuse to serve the remaining march
func (this *Flight) canPass(key cache.Key) (ok, firstTimeFail bool) {
	ok, firstTimeFail = true, false
	if this.maxRetries == 0 {
		return
	}
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
	ok, firstTimeFail := this.canPass(key)
	if !ok {
		if firstTimeFail {
			log.Warn("max retries[%d] reached: %+v", this.maxRetries, key)
		}

		return false
	}

	// FIXME Get and Set is not atomic
	if _, present := this.lock.Get(key); !present {
		this.lock.Set(key, true)
		return true
	}

	log.Debug("locked %#v", key)
	return false
}

func (this *Flight) Land(key cache.Key, ok bool) {
	this.lock.Del(key)
	if this.maxRetries > 0 && ok {
		this.retry.Set(key, 0) // reset the counter
	}

}

func (this *Flight) Flying(key cache.Key) bool {
	if _, present := this.lock.Get(key); present {
		return true
	}
	return false
}

func (this *Flight) Len() int {
	return this.lock.Len()
}
