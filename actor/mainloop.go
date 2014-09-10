package actor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (this *Actor) ServeForever() {
	this.waitForUpstreamRequests()

	ticker := time.NewTicker(
		time.Duration(this.server.Int("upstream_tick", 1)) * time.Second)
	defer ticker.Stop()

	var now time.Time
	for {
		select {
		case <-ticker.C:
			now = time.Now()
			for _, m := range this.marches.pullInBatch(now) {
				go this.callback(m)
			}

		}
	}

}

func (this *Actor) callback(m march) {
	m.Op = "" // omitempty
	m.At = 0
	buf, _ := json.Marshal(m)
	fmt.Println(string(buf), m)
	body := bytes.NewBuffer(buf)
	url := fmt.Sprintf("http://localhost/api/?class=r&method=%s", m.Op)
	res, err := http.Post(url, "application/json", body)
	if err != nil {

	}

	fmt.Println(res)

}
