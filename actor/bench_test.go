package actor

import (
	"encoding/json"
	"github.com/funkygao/assert"
	"testing"
	"time"
)

func BenchmarkJobJsonEncode(b *testing.B) {
	job := Job{Uid: 534343, JobId: 5677, dueTime: time.Now()}
	for i := 0; i < b.N; i++ {
		json.Marshal(job)
	}
}

func BenchmarkPhpPayloadPartialDecode(b *testing.B) {
	payload := []byte(`{"ok":0,"msg":"0:Unknown event: . -- #0 \/sgn\/htdocs\/dev-dev\/dragon-server-code\/v2\/classes\/Services\/ActorService.php(20): Event\\EventEngine::fire(NULL, Array)\n#1 \/sgn\/htdocs\/dev-dev\/dragon-server-code\/v2\/classes\/System\/Application.php(160): Services\\ActorService->play(Array)\n#2 \/sgn\/htdocs\/dev-dev\/dragon-server-code\/docroot\/api\/index.php(12): System\\Application->execute()\n#3 {main}"}`)
	var (
		objmap map[string]*json.RawMessage
		ok     int
	)
	for i := 0; i < b.N; i++ {
		json.Unmarshal(payload, &objmap)
		json.Unmarshal(*objmap["ok"], &ok)
	}
}

func BenchmarkPhpPayloadFullDecode(b *testing.B) {
	payload := []byte(`{"ok":0,"msg":"0:Unknown event: . -- #0 \/sgn\/htdocs\/dev-dev\/dragon-server-code\/v2\/classes\/Services\/ActorService.php(20): Event\\EventEngine::fire(NULL, Array)\n#1 \/sgn\/htdocs\/dev-dev\/dragon-server-code\/v2\/classes\/System\/Application.php(160): Services\\ActorService->play(Array)\n#2 \/sgn\/htdocs\/dev-dev\/dragon-server-code\/docroot\/api\/index.php(12): System\\Application->execute()\n#3 {main}"}`)
	var decoded map[string]interface{}
	for i := 0; i < b.N; i++ {
		json.Unmarshal(payload, &decoded)
	}
}

func TestJobEncode(t *testing.T) {
	job := Job{Uid: 534343, JobId: 5677, dueTime: time.Now()}
	body, _ := json.Marshal(job)
	assert.Equal(t, `{"uid":534343,"job_id":5677}`, string(body))
}
