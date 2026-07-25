package world

import "testing"

func TestRotDecaysAway(t *testing.T) {
	g := bareGrid()
	col0, row0 := coreCorner()
	g.Grow(col0-1, row0)
	g.Cut(col0-1, row0) // -> rot with a lifetime

	if g.Kind(col0-1, row0) != Rot || !g.hasRot() {
		t.Fatal("cut did not create a tracked rot tile")
	}

	g.decayRot(rotLifetime + 1)
	if g.Kind(col0-1, row0) != Empty {
		t.Fatal("rot did not crumble away after its lifetime")
	}
	if g.hasRot() {
		t.Fatal("decayed rot was not removed from the list")
	}
}

func TestGrowingOverRotUntracksIt(t *testing.T) {
	g := bareGrid()
	col0, row0 := coreCorner()
	g.Grow(col0-1, row0)
	g.Cut(col0-1, row0)  // rot
	g.Grow(col0-1, row0) // reclaim it

	if g.hasRot() {
		t.Fatal("rot reclaimed by a root should no longer be tracked")
	}
}
