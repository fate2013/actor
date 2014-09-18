package queue

import (
	"encoding/json"
	"testing"
)

func TestJson(t *testing.T) {
	data := `{"uid":5,"payload":"{city_id:5}"}`
	var v struct {
		Uid     int    `json:"uid"`
		Payload string `json:"payload"` // []byte
	}
	err := json.Unmarshal([]byte(data), &v)
	if err != nil {
		t.Logf("%s", err)
	}
	t.Logf("%+v", v)
}
