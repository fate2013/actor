package actor

import (
	"testing"
)

func TestMarchesStruct(t *testing.T) {
	marches := newMarches()
	marches.set(march{uid: 1, marchId: 1, opTime: 3, op: "startMarch"})
	marches.set(march{uid: 2, marchId: 2, opTime: 2, op: "startMarch"})
	marches.set(march{uid: 2, marchId: 2, opTime: 8, op: "recall"})
	t.Logf("%#v %#v\n", marches.m, marches.sortedKeys())
	for i := range marches.sortedKeys() {
		t.Logf("%#v\n", marches.m[marches.k[i]])
	}

}
