package actor

import (
	"encoding/json"
	"github.com/funkygao/assert"
	"testing"
	"time"
)

func BenchmarkJobJsonEncode(b *testing.B) {
	job := job{Uid: 534343, JobId: 5677, dueTime: time.Now()}
	for i := 0; i < b.N; i++ {
		json.Marshal(job)
	}
}

func TestJobEncode(t *testing.T) {
	job := job{Uid: 534343, JobId: 5677, dueTime: time.Now()}
	body, _ := json.Marshal(job)
	assert.Equal(t, `{"uid":534343,"job_id":5677}`, string(body))
}
