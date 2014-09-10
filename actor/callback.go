package actor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (this *Actor) callback(m march) {
	m.Evt = "" // omitempty
	m.At = 0
	buf, _ := json.Marshal(m)
	fmt.Println(string(buf), m)
	body := bytes.NewBuffer(buf)
	url := fmt.Sprintf("http://localhost/api/?class=r&method=%s", m.Evt)
	res, err := http.Post(url, "application/json", body)
	if err != nil {

	}

	fmt.Println(res)

}
