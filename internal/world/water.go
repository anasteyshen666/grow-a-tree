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
	sourceCount     = 5
	SourceMaxAmount = 120.0
	mineRatePerSec  = 8.0
)

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

func (g *Grid) spawnSources() {
	for len(g.sources) < sourceCount {
		c, r := rand.Intn(Cols), rand.Intn(Rows)
		if g.cells[r][c] != Empty {
			continue
		}
		amt := sourceCapacity()
		g.cells[r][c] = Water
		g.sources = append(g.sources, &waterSource{col: c, row: r, amount: amt, max: amt})
	}
}

// MineWater draws water from every source touched by the live network, depletes
// those sources, removes any that ran dry, and returns the total mined this tick.
// Each source that dries up is replaced by a fresh one elsewhere, so the map
// always holds sourceCount puddles.
func (g *Grid) MineWater(dt float64) float64 {
	total := 0.0
	kept := g.sources[:0]
	for _, s := range g.sources {
		if g.touchesNetwork(s.col, s.row) {
			take := mineRatePerSec * dt
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
