package world

import "testing"

// clearSources removes the random water so it can't block test paths.
func clearSources(g *Grid) {
	for _, s := range g.sources {
		g.cells[s.row][s.col] = Empty
	}
	g.sources = nil
}

func TestAdjacentTargetFindsRoot(t *testing.T) {
	g := NewGrid()
	clearSources(g)
	col0, row0 := coreCorner()
	g.Grow(col0-1, row0)

	c, r, ok := g.AdjacentTarget(col0-2, row0)
	if !ok || g.Kind(c, r) != Root {
		t.Fatalf("did not find the adjacent root, got (%d,%d) ok=%v", c, r, ok)
	}
}

func TestNextStepMovesTowardNetwork(t *testing.T) {
	g := NewGrid()
	clearSources(g)
	col0, row0 := coreCorner()
	sc, sr := col0-5, row0

	c, r, ok := g.NextStep(sc, sr)
	if !ok {
		t.Fatal("no path to the network")
	}
	if c != sc+1 || r != sr {
		t.Fatalf("expected first step (%d,%d), got (%d,%d)", sc+1, sr, c, r)
	}
}

func TestCoreTakesDamage(t *testing.T) {
	g := NewGrid()
	if g.CoreHP() != CoreMaxHP {
		t.Fatalf("core should start at %d HP, got %d", CoreMaxHP, g.CoreHP())
	}
	col0, row0 := coreCorner()
	g.Damage(col0, row0, 7)
	if g.CoreHP() != CoreMaxHP-7 {
		t.Fatalf("gnawing the core did not chip its HP: %d", g.CoreHP())
	}
}

func TestDamageDestroysRootAndBreaksNetwork(t *testing.T) {
	g := NewGrid()
	clearSources(g)
	row, cols := growChainLeft(g, 2) // cols[0] by the Core, cols[1] beyond it

	g.Damage(cols[0], row, 0)
	if g.Kind(cols[0], row) != Empty {
		t.Fatal("gnawed root was not destroyed")
	}
	if g.IsConnected(cols[1], row) {
		t.Fatal("root beyond the destroyed one should be disconnected")
	}
}
