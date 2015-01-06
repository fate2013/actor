package actor

import (
	"github.com/funkygao/lockkey"
	log "github.com/funkygao/log4go"
)

const (
	LOCKER_REASON = "actor.sched"
	LOCKER_LOCK   = "lock"
	LOCKER_UNLOCK = "unlock"
)

type Locker []string

func NewLocker() Locker {
	this := make([]string, 0)
	return this
}

func (this *Locker) LockUser(uid int64) bool {
	return this.Lock(lockkey.User(uid))
}

func (this *Locker) Lock(key string) (success bool) {
	svt, err := fae.proxy.ServantByKey(key)
	if err != nil {
		log.Error("fae.lock[%s]: %s", key, err.Error())
		return false
	}

	log.Debug("fae.lock[%s]: %s", key, svt.Addr())
	if success, _ = svt.GmLock(fae.Context(LOCKER_REASON), LOCKER_LOCK, key); success {
		*this = append(*this, key)
	}

	svt.Recycle()

	return
}

func (this *Locker) ReleaseAll() {
	for _, key := range *this {
		svt, err := fae.proxy.ServantByKey(key)
		if err != nil {
			log.Error("fae.unlock[%s]: %s", key, err.Error())
			continue
		}

		if err = svt.GmUnlock(fae.Context(LOCKER_REASON),
			LOCKER_UNLOCK, key); err != nil {
			log.Error("fae.unlock[%s]: %s", key, err.Error())
		} else {
			log.Debug("fae.unlock[%s]: %s", key, svt.Addr())
		}

		svt.Recycle()
	}
}
