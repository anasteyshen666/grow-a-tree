package world

import "testing"

// gridWithSource returns a grid with only one known water source (the random
// spawn cleared), so tests aren't disturbed by random placement.
func gridWithSource(sc, sr int) *Grid {
	g := NewGrid()
	clearSources(g)
	g.cells[sr][sc] = Water
	g.sources = append(g.sources, &waterSource{col: sc, row: sr, amount: 100})
	return g
}

func TestMineWaterNeedsLiveNetwork(t *testing.T) {
	col0, row0 := coreCorner()
	sc, sr := col0-2, row0 // one gap away from the Core

	dry := gridWithSource(sc, sr)
	if got := dry.MineWater(1, nil); got != 0 {
		t.Fatalf("mined water with no adjacent network: %v", got)
	}

	tapped := gridWithSource(sc, sr)
	tapped.Grow(col0-1, row0) // connected root now sits beside the source
	if got := tapped.MineWater(1, nil); got <= 0 {
		t.Fatal("no water mined despite an adjacent connected root")
	}
}

func TestNewUnconnectedCoreMinesAdjacentWater(t *testing.T) {
	g := bareGrid()
	col0, row0 := coreCorner()

	// a second core, far from the first and not linked by any roots
	g.placeCore(col0-8, row0) // cells col0-8 and col0-7
	g.recomputeConnectivity()
	if g.CoreMerges() != 0 {
		t.Fatalf("second core should be unconnected, merges=%d", g.CoreMerges())
	}

	sc, sr := col0-9, row0 // source right next to the new core's left cell
	g.cells[sr][sc] = Water
	g.sources = append(g.sources, &waterSource{col: sc, row: sr, amount: 100, max: 100})

	if got := g.MineWater(1, nil); got <= 0 {
		t.Fatal("water not mined from a source beside a new, unconnected core")
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

	g.MineWater(1, nil) // rate outpaces the tiny amount, so it empties and respawns
	if len(g.sources) != baseSourceCount {
		t.Fatalf("map should refill to %d sources, got %d", baseSourceCount, len(g.sources))
	}
	for _, s := range g.sources {
		if s == old {
			t.Fatal("depleted source should have been removed")
		}
	}
}
