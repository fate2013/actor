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

}

func BenchmarkJsonUnmarshal(b *testing.B) {
	var v map[string]interface{}
	for i := 0; i < b.N; i++ {
		json.Unmarshal([]byte(`{"uid":189845,"march_id":1444}`), v)
	}
}
