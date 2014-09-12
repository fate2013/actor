package actor

import (
	"github.com/funkygao/assert"
	"math/rand"
	"testing"
	"time"
)

func TestMarchesChunks(t *testing.T) {
	m := marches{}
	m = append(m, march{Uid: 1, MarchId: 2, At: 1, X: 10, Y: 10})
	m = append(m, march{Uid: 2, MarchId: 3, At: 1, X: 10, Y: 10})
	m = append(m, march{Uid: 3, MarchId: 4, At: 2, X: 20, Y: 10})
	m = append(m, march{Uid: 4, MarchId: 5, At: 1, X: 10, Y: 20})
	m = append(m, march{Uid: 8, MarchId: 3, At: 1, X: 10, Y: 10})
	chunks := m.chunks()
	assert.Equal(t, 3, len(chunks))

	maxChunkLen := -1
	for _, c := range chunks {
		if len(c) > maxChunkLen {
			maxChunkLen = len(c)
		}
	}
	assert.Equal(t, 3, maxChunkLen)
}

func BenchmarkWakeup(b *testing.B) {
	jobN := 100000
	jobs := newJobs()
	now := time.Now().Unix()
	var marchId int64 = 0
	events := []int{
		12,  // arrive
		15,  // gather done
		20,  // back home
		501, // speedup
		502, // recall
	}
	for i := 0; i < jobN; i++ {
		marchId++
		uid := rand.Int63n(10000) + 1
		x := rand.Intn(1024)
		y := rand.Intn(512)
		event := events[rand.Intn(len(events))]
		at := time.Now().Add(time.Duration(rand.Intn(1000)) * time.Second)
		jobs.sched(march{Uid: uid, MarchId: marchId, At: int(at.Unix()), Evt: event, X: x, Y: y})
	}
	for i := 0; i < b.N; i++ {
		jobs.wakeup(now)
	}

}

func BenchmarkSched(b *testing.B) {
	var marchId int64 = 0
	jobs := newJobs()
	for i := 0; i < b.N; i++ {
		marchId++
		uid := rand.Int63n(10000) + 1
		x := rand.Intn(1024)
		y := rand.Intn(512)
		event := 12
		at := time.Now().Add(time.Duration(rand.Intn(1000)) * time.Second)
		jobs.sched(march{Uid: uid, MarchId: marchId, At: int(at.Unix()), Evt: event, X: x, Y: y})
	}

}
