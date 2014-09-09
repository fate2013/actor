package actor

import (
	"testing"
)

func TestMarchesStruct(t *testing.T) {
	marches := newMarches()
	marches.set(march{Uid: 1, MarchId: 1, Optime: 3, Op: "startMarch"})
	marches.set(march{Uid: 2, MarchId: 2, Optime: 2, Op: "startMarch"})
	marches.set(march{Uid: 2, MarchId: 2, Optime: 8, Op: "recall"})
	t.Logf("%#v %#v\n", marches.m, marches.sortedKeys())
	for i := range marches.sortedKeys() {
		t.Logf("%#v\n", marches.m[marches.k[i]])
	}

}

func TestSmoking(t *testing.T) {
	marches := newMarches()
	marches.set(march{Uid: 1, MarchId: 1, Optime: 3, Op: "startMarch"})
	marches.set(march{Uid: 2, MarchId: 2, Optime: 2, Op: "startMarch"})
	marches.set(march{Uid: 2, MarchId: 2, Optime: 8, Op: "recall"})
	t.Logf("%#v %#v\n", marches.m, marches.sortedKeys())
	actor := NewActor(nil)
	for i := range marches.sortedKeys() {
		t.Logf("%#v\n", marches.m[marches.k[i]])
		actor.callback(marches.m[marches.k[i]])
	}

}
