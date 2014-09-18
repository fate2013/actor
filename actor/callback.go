package actor

import (
	"bytes"
	"fmt"
	log "github.com/funkygao/log4go"
	"io/ioutil"
	"net/http"
)

func (this *scheduler) callback() {
	params := []byte("")
	url := fmt.Sprintf(this.actor.config.CallbackUrl, string(params))
	log.Debug("callback: %s", url)

	// may fail, because php will throw LockException
	// in that case, will reschedule the job after 1s
	res, err := http.Post(url, CONTENT_TYPE_JSON, bytes.NewBuffer(params))
	defer func() {
		res.Body.Close()
	}()

	ioutil.ReadAll(res.Body)

	if err != nil {
		log.Error("post error: %s", err.Error())
	} else {
		if res.StatusCode != http.StatusOK {
			log.Error("callback error: %+v", res)
		}
	}

}
