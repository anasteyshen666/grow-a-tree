package world

import "testing"

func coreCorner() (col, row int) { return Cols/2 - 1, Rows/2 - 1 }

// bareGrid is a fresh grid with the random water sources cleared, so growth
// tests near the Core aren't disturbed by a puddle landing on a target cell.
func bareGrid() *Grid {
	g := NewGrid()
	clearSources(g)
	return g
}

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
	g := bareGrid()
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

func TestCutTurnsRootToRot(t *testing.T) {
	g := bareGrid()
	col, row := coreCorner()
	g.Grow(col-1, row)

	if !g.Cut(col-1, row) || g.Kind(col-1, row) != Rot {
		t.Fatal("cutting a root did not leave Rot")
	}
	if g.Cut(col-1, row) {
		t.Fatal("cut something that was not a root")
	}
	if g.Cut(col, row) {
		t.Fatal("cut the Core")
	}
}

func TestGrowReclaimsRot(t *testing.T) {
	g := bareGrid()
	col, row := coreCorner()
	g.Grow(col-1, row)
	g.Cut(col-1, row)

	if !g.Grow(col-1, row) || g.Kind(col-1, row) != Root {
		t.Fatal("could not regrow over Rot next to the Core")
	}
}

// growChainLeft grows n roots straight left from the Core's left edge and
// returns the columns used (row stays at the Core's top row).
func growChainLeft(g *Grid, n int) (row int, cols []int) {
	col0, row0 := coreCorner()
	for i := 1; i <= n; i++ {
		c := col0 - i
		if !g.Grow(c, row0) {
			panic("chain growth failed")
		}
		cols = append(cols, c)
	}
	return row0, cols
}

func TestCutInMiddleDisconnectsDownstream(t *testing.T) {
	g := bareGrid()
	row, cols := growChainLeft(g, 3) // cols[0] nearest Core ... cols[2] farthest

	if !g.Cut(cols[1], row) {
		t.Fatal("cut of middle root failed")
	}
	if !g.IsConnected(cols[0], row) {
		t.Fatal("root between Core and the gap should stay connected")
	}
	if g.IsConnected(cols[2], row) {
		t.Fatal("root beyond the gap should be disconnected")
	}
}

func TestCannotGrowFromDeadRoot(t *testing.T) {
	g := bareGrid()
	row, cols := growChainLeft(g, 3)
	g.Cut(cols[1], row) // isolates cols[2]

	if g.CanGrow(cols[2], row-1) {
		t.Fatal("grew from a disconnected root")
	}
	if !g.CanGrow(cols[0], row-1) {
		t.Fatal("could not grow from a still-connected root")
	}
}

func TestGrowChainsAndRejectsOccupied(t *testing.T) {
	g := bareGrid()
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
