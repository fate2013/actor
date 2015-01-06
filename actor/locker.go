package actor

import (
	log "github.com/funkygao/log4go"
)

type Locker []string

func NewLocker() Locker {
	this := make([]string, 0)
	return this
}

func (this *Locker) Lock(key string) (success bool) {
	svt, err := fae.proxy.Servant("localhost:9001")
	if err != nil {
		return false
	}

	// FIXME
	if success, _ = svt.GmLock(fae.Context("actor"), "lock", key); success {
		*this = append(*this, key)
	}

	svt.Recycle()

	return
}

func (this *Locker) ReleaseAll() {
	for _, key := range *this {
		svt, err := fae.proxy.Servant("localhost:9001")
		if err != nil {
			log.Error("fae.servant[%s]: %s", key, err.Error())
			continue
		}

		if err = svt.GmUnlock(fae.Context("actor"), "unlock", key); err != nil {
			log.Error("fae.unlock[%s]: %s", key, err.Error())
		}

		svt.Recycle()
	}
}
