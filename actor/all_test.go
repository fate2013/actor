package actor

import (
	"encoding/json"
	"github.com/funkygao/assert"
	"sync"
	"testing"
	"time"
)

func BenchmarkDefer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		defer func() {

		}()
	}
}

func BenchmarkSwitch(b *testing.B) {
	var x = 10
	for i := 0; i < b.N; i++ {
		switch x {
		case 1:
		case 2:
		case 4:
		case 10:
		default:
		}
	}
}

func BenchmarkAdd(b *testing.B) {
	b.ReportAllocs()
	var x int64
	for i := 0; i < b.N; i++ {
		x = x + 1
	}
}

func BenchmarkJobMarshal(b *testing.B) {
	b.ReportAllocs()
	job := Job{Uid: 534343, JobId: 5677, TimeEnd: time.Now()}
	for i := 0; i < b.N; i++ {
		job.Marshal()
	}
	b.SetBytes(int64(len(job.Marshal())))
}

func BenchmarkMarchMarshal(b *testing.B) {
	b.ReportAllocs()
	march := March{Uid: 232323, MarchId: 23223232, State: "marching", X1: 12, Y1: 122, EndTime: time.Now()}
	for i := 0; i < b.N; i++ {
		march.Marshal()
	}
	b.SetBytes(int64(len(march.Marshal())))
}

func BenchmarkPveMarshal(b *testing.B) {
	b.ReportAllocs()
	pve := Pve{Uid: 3434343, MarchId: 343433434333, State: "marching", EndTime: time.Now()}
	for i := 0; i < b.N; i++ {
		pve.Marshal()
	}
	b.SetBytes(int64(len(pve.Marshal())))
}

func BenchmarkMutex(b *testing.B) {
	b.ReportAllocs()
	var mutex sync.Mutex
	for i := 0; i < b.N; i++ {
		mutex.Lock()
		mutex.Unlock()
	}
}

func BenchmarkFlightTakeoff(b *testing.B) {
	b.ReportAllocs()
	f := NewFlight(100000)
	job := Job{Uid: 534343, JobId: 5677, TimeEnd: time.Now()}
	for i := 0; i < b.N; i++ {
		f.Takeoff(job.FlightKey())
	}
}

func BenchmarkFlightLand(b *testing.B) {
	b.ReportAllocs()
	f := NewFlight(100000)
	job := Job{Uid: 534343, JobId: 5677, TimeEnd: time.Now()}
	for i := 0; i < b.N; i++ {
		f.Land(job.FlightKey())
	}
}

func BenchmarkPhpPayloadPartialDecode(b *testing.B) {
	b.ReportAllocs()
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
	b.ReportAllocs()
	payload := []byte(`{"ok":0,"msg":"0:Unknown event: . -- #0 \/sgn\/htdocs\/dev-dev\/dragon-server-code\/v2\/classes\/Services\/ActorService.php(20): Event\\EventEngine::fire(NULL, Array)\n#1 \/sgn\/htdocs\/dev-dev\/dragon-server-code\/v2\/classes\/System\/Application.php(160): Services\\ActorService->play(Array)\n#2 \/sgn\/htdocs\/dev-dev\/dragon-server-code\/docroot\/api\/index.php(12): System\\Application->execute()\n#3 {main}"}`)
	var decoded map[string]interface{}
	for i := 0; i < b.N; i++ {
		json.Unmarshal(payload, &decoded)
	}
}

func TestJobEncode(t *testing.T) {
	job := Job{Uid: 534343, JobId: 5677, TimeEnd: time.Now()}
	body, _ := json.Marshal(job)
	assert.Equal(t, `{"uid":534343,"job_id":5677}`, string(body))
}

func TestMarchGeoHash(t *testing.T) {
	march := March{Uid: 22323, MarchId: 34343434343, X1: 2323, Y1: 343}
	assert.Equal(t, 12, march.GeoHash())
}
