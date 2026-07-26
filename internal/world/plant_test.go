package world

import "testing"

func TestCanPlantNearCoreOnly(t *testing.T) {
	g := bareGrid()
	col0, row0 := coreCorner()

	if !g.CanPlant(col0-1, row0) {
		t.Fatal("should be able to plant on an empty cell next to the Core")
	}
	if g.CanPlant(0, 0) {
		t.Fatal("planted far from any Core")
	}

	g.Grow(col0-1, row0) // now occupied by a root
	if g.CanPlant(col0-1, row0) {
		t.Fatal("planted on an occupied cell")
	}
}

func TestSetPlantMarksCell(t *testing.T) {
	g := bareGrid()
	col0, row0 := coreCorner()
	g.SetPlant(col0-1, row0)
	if g.Kind(col0-1, row0) != Plant {
		t.Fatal("SetPlant did not mark the cell")
	}
	if g.CanGrow(col0-1, row0) {
		t.Fatal("a root should not grow on a plant cell")
	}
}
