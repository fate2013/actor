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
	body := bytes.NewBuffer(params)
	url := fmt.Sprintf(this.server.String("callback_url", ""), string(params))
	log.Debug("callback: %s", url)

	res, err := http.Post(url, CONTENT_TYPE_JSON, body)
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
