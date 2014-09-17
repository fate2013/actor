package encoding

import (
	"encoding/json"
	"strings"
	"testing"
)

func BenchmarkStringsSplit(b *testing.B) {
	const s = `14121212121,{"uid":189845,"march_id":1444}`
	for i := 0; i < b.N; i++ {
		strings.SplitN(s, ",", 1)
	}
	b.SetBytes(int64(len(s)))

}

func BenchmarkJsonUnmarshal(b *testing.B) {
	const s = `{"uid":189845,"march_id":1444}`
	var v map[string]interface{}
	for i := 0; i < b.N; i++ {
		json.Unmarshal([]byte(s), v)
	}
	b.SetBytes(int64(len(s)))
}
