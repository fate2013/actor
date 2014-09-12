package actor

import (
	"github.com/funkygao/assert"
	"testing"
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
