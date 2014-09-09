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

	ticker := time.NewTicker(time.Duration(this.server.Int("upstream_tick", 1)) * time.Second)
	defer ticker.Stop()

	var march march
	var now time.Time
	for _ = range ticker.C {
		now = time.Now()
		for optime := range this.marches.sortedKeys() {
			march = this.marches.m[this.marches.k[optime]]
			if now.Unix() >= int64(march.opTime) {
				// time to do op
				go this.callback(march)
			}
		}

	}

}

func (this *Actor) callback(m march) {
	buf, _ := json.Marshal(m)
	fmt.Println(string(buf), m)
	body := bytes.NewBuffer(buf)
	url := fmt.Sprintf("http://localhost/api/?class=r&method=%s", m.op)
	res, err := http.Post(url, "application/json", body)
	if err != nil {

	}

	fmt.Println(res)

}
