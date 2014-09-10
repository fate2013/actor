package actor

import (
	"testing"
)

func TestMarchesStruct(t *testing.T) {
	marches := newMarches()
	marches.set(march{Uid: 1, MarchId: 1, At: 3, Evt: "startMarch"})
	marches.set(march{Uid: 2, MarchId: 2, At: 2, Evt: "startMarch"})
	marches.set(march{Uid: 2, MarchId: 2, At: 8, Evt: "recall"})
	t.Logf("%#v %#v\n", marches.m, marches.sortedKeys())
	for i := range marches.sortedKeys() {
		t.Logf("%#v\n", marches.m[marches.k[i]])
	}

}

func TestSmoking(t *testing.T) {
	marches := newMarches()
	marches.set(march{Uid: 1, MarchId: 1, At: 3, Evt: "startMarch"})
	marches.set(march{Uid: 2, MarchId: 2, At: 2, Evt: "startMarch"})
	marches.set(march{Uid: 2, MarchId: 2, At: 8, Evt: "recall"})
	t.Logf("%#v %#v\n", marches.m, marches.sortedKeys())
	actor := NewActor(nil)
	for i := range marches.sortedKeys() {
		t.Logf("%#v\n", marches.m[marches.k[i]])
		actor.callback(marches.m[marches.k[i]])
	}

}
