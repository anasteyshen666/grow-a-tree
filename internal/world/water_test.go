package world

import "testing"

func TestMineWaterNeedsLiveNetwork(t *testing.T) {
	g := NewGrid()
	g.sources = nil // drop the random spawn, place a known source

	col0, row0 := coreCorner()
	sc, sr := col0-2, row0 // one gap away from the Core
	g.cells[sr][sc] = Water
	g.sources = append(g.sources, &waterSource{col: sc, row: sr, amount: 100})

	if got := g.MineWater(1); got != 0 {
		t.Fatalf("mined water with no adjacent network: %v", got)
	}

	g.Grow(col0-1, row0) // connected root now sits beside the source
	if got := g.MineWater(1); got <= 0 {
		t.Fatal("no water mined despite an adjacent connected root")
	}
}

func TestDepletedSourceIsReplaced(t *testing.T) {
	g := NewGrid()
	for _, s := range g.sources { // clear the random spawn
		g.cells[s.row][s.col] = Empty
	}
	g.sources = nil

	col0, row0 := coreCorner()
	sc, sr := col0-1, row0 // directly beside the Core
	g.cells[sr][sc] = Water
	old := &waterSource{col: sc, row: sr, amount: 1}
	g.sources = append(g.sources, old)

	g.MineWater(1) // rate outpaces the tiny amount, so it empties and respawns
	if len(g.sources) != sourceCount {
		t.Fatalf("map should refill to %d sources, got %d", sourceCount, len(g.sources))
	}
	for _, s := range g.sources {
		if s == old {
			t.Fatal("depleted source should have been removed")
		}
	}
}
