package world

import "testing"

func coreCorner() (col, row int) { return Cols/2 - 1, Rows/2 - 1 }

func TestNewGridHasCore(t *testing.T) {
	g := NewGrid()
	col, row := coreCorner()
	for dr := 0; dr < 2; dr++ {
		for dc := 0; dc < 2; dc++ {
			if g.Kind(col+dc, row+dr) != Core {
				t.Fatalf("cell (%d,%d) is not Core", col+dc, row+dr)
			}
		}
	}
}

func TestGrowRequiresAdjacency(t *testing.T) {
	g := NewGrid()
	if g.Grow(0, 0) {
		t.Fatal("grew a root with no adjacent network")
	}
	col, row := coreCorner()
	if !g.Grow(col-1, row) {
		t.Fatal("could not grow next to the Core")
	}
	if g.Kind(col-1, row) != Root {
		t.Fatal("cell did not become Root")
	}
}

func TestGrowChainsAndRejectsOccupied(t *testing.T) {
	g := NewGrid()
	col, row := coreCorner()
	if !g.Grow(col-1, row) {
		t.Fatal("first root failed")
	}
	if !g.Grow(col-2, row) {
		t.Fatal("chained root off an existing root failed")
	}
	if g.Grow(col-2, row) {
		t.Fatal("grew onto an occupied cell")
	}
}
