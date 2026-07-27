package world

import "math/rand"

// waterSource is a puddle the player mines by reaching it with the live network.
// It depletes as water is drawn and vanishes when empty (GDD §2).
type waterSource struct {
	col, row int
	amount   float64
	max      float64
}

const (
	baseSourceCount = 5
	SourceMaxAmount = 120.0
	mineRatePerSec  = 5.0 // water mined per second per tapped source
)

// sourceTarget is how many sources the map should hold: the base count plus one
// per extra Core, so a bigger network has more water to draw from.
func (g *Grid) sourceTarget() int { return baseSourceCount + len(g.cores) - 1 }

// sourceCapacity rolls a source's size: usually single, sometimes double, and
// rarely triple — bigger ones hold more water and so last longer.
func sourceCapacity() float64 {
	switch r := rand.Float64(); {
	case r < 0.04:
		return SourceMaxAmount * 3
	case r < 0.20:
		return SourceMaxAmount * 2
	default:
		return SourceMaxAmount
	}
}

// randInsetCell returns a random cell at least one cell in from every border,
// so features never hug the window edge.
func randInsetCell() (col, row int) {
	return rand.Intn(Cols-2) + 1, rand.Intn(Rows-2) + 1
}

func (g *Grid) spawnSources() {
	for len(g.sources) < g.sourceTarget() {
		c, r := randInsetCell()
		if g.cells[r][c] != Empty {
			continue
		}
		amt := sourceCapacity()
		g.cells[r][c] = Water
		g.sources = append(g.sources, &waterSource{col: c, row: r, amount: amt, max: amt})
		g.emit(FxPlaceWater, c, r)
	}
}

// MineWater draws water from every source touched by the live network, depletes
// those sources, removes any that ran dry, and returns the total mined this tick.
// Each source that dries up is replaced by a fresh one elsewhere. protect (may
// be nil) reports sources shielded by a Thornbush from the winter drain penalty.
func (g *Grid) MineWater(dt float64, protect func(col, row int) bool) float64 {
	total := 0.0
	kept := g.sources[:0]
	for _, s := range g.sources {
		if g.touchesNetwork(s.col, s.row) {
			mul := g.season.SourceDrainMul()
			if mul > 1.0 && protect != nil && protect(s.col, s.row) {
				mul = 1.0 // Thornbush cancels the winter penalty
			}
			take := mineRatePerSec * dt * mul
			if take > s.amount {
				take = s.amount
			}
			s.amount -= take
			total += take
		}
		if s.amount > 0 {
			kept = append(kept, s)
		} else {
			g.cells[s.row][s.col] = Empty
		}
	}
	g.sources = kept
	g.spawnSources() // top back up to sourceCount, respawning any that dried out
	return total
}
