package world

import "testing"

func TestMatureAddsCore(t *testing.T) {
	g := bareGrid()
	if g.CoreCount() != 1 {
		t.Fatalf("should start with one core, got %d", g.CoreCount())
	}
	if !g.Mature() || g.CoreCount() != 2 {
		t.Fatalf("mature did not add a second core, count=%d", g.CoreCount())
	}
}

func TestCoreDestroyedAtZeroHP(t *testing.T) {
	g := bareGrid()
	col0, row0 := coreCorner()
	g.Damage(col0, row0, CoreMaxHP) // one lethal hit

	if g.CoreCount() != 0 {
		t.Fatal("core was not destroyed at 0 HP")
	}
	if g.Kind(col0, row0) != Empty {
		t.Fatal("destroyed core cells were not cleared")
	}
}

func TestLinkingCoresRaisesMerges(t *testing.T) {
	g := bareGrid()
	col0, row0 := coreCorner()
	g.placeCore(col0-4, row0) // second core: cells col0-4 and col0-3
	g.recomputeConnectivity()
	if g.CoreMerges() != 0 {
		t.Fatalf("separate cores should have 0 merges, got %d", g.CoreMerges())
	}

	g.Grow(col0-1, row0) // bridge the two-cell gap between the cores
	g.Grow(col0-2, row0)
	if g.CoreMerges() != 1 {
		t.Fatalf("linked cores should have 1 merge, got %d", g.CoreMerges())
	}
}
