package actor

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/funkygao/log4go"
	"net/http"
)

func (this *Actor) callback(m march) {
	buf, _ := json.Marshal(m)

	body := bytes.NewBuffer(buf)
	url := fmt.Sprintf(this.server.String("callback_url", ""), m.Evt)
	log.Debug("%+v %s %s", m, string(buf), url)

	res, err := http.Post(url, "application/json", body)
	if err != nil {
		log.Error(err)
	} else {
		if res.StatusCode != http.StatusOK {
			log.Error("callback error: %+v", res)
		}
	}

}
