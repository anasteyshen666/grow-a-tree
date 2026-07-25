package world

import "testing"

func TestSporeInfectsAdjacentRoot(t *testing.T) {
	g := bareGrid()
	col0, row0 := coreCorner()
	g.Grow(col0-1, row0) // a root beside the Core

	sc, sr := col0-2, row0 // spore right next to that root
	g.cells[sr][sc] = Spore
	g.spores = append(g.spores, [2]int{sc, sr})

	g.infectFromSpores()

	if !g.IsMushroom(col0-1, row0) {
		t.Fatal("root next to the spore was not infected")
	}
	if g.Kind(sc, sr) != Empty {
		t.Fatal("spore was not consumed")
	}
	if !g.SlowsBug(col0-1, row0-1) {
		t.Fatal("a cell next to the mushroom root should slow bugs")
	}
}

func TestCutClearsMushroom(t *testing.T) {
	g := bareGrid()
	col0, row0 := coreCorner()
	g.Grow(col0-1, row0)
	g.mushroom[row0][col0-1] = true

	g.Cut(col0-1, row0)
	if g.mushroom[row0][col0-1] {
		t.Fatal("cutting a mushroom root did not clear its flag")
	}
}
