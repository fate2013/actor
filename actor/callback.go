package actor

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/funkygao/log4go"
	"net/http"
)

func (this *Actor) callback(m march) {
	params, _ := json.Marshal(m)
	url := fmt.Sprintf(this.server.String("callback_url", ""), string(params))
	log.Debug("callback: %s", url)

	// may fail, because php will throw LockException
	// in that case, will reschedule the job after 1s
	res, err := http.Post(url, CONTENT_TYPE_JSON, bytes.NewBuffer(params))
	defer func() {
		res.Body.Close()
	}()

	if err != nil {
		log.Error("post error: %s", err.Error())
	} else {
		if res.StatusCode != http.StatusOK {
			log.Error("callback error: %+v", res)
		}
	}

}

// coordinate marches with same destination at the same time
func (this *Actor) coordinate(chunk []march) {

}
