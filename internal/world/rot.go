package world

// rotTile is a piece of rot with a lifetime. Rot is bug bait, but it is not
// permanent — fresh rot is what lures them (GDD §4).
type rotTile struct {
	col, row int
	ttl      float64
}

const rotLifetime = 12.0 // seconds a rot tile lasts before crumbling away

func (g *Grid) addRot(col, row int) {
	g.cells[row][col] = Rot
	g.rots = append(g.rots, &rotTile{col: col, row: row, ttl: rotLifetime})
}

func (g *Grid) removeRot(col, row int) {
	for i, rt := range g.rots {
		if rt.col == col && rt.row == row {
			g.rots = append(g.rots[:i], g.rots[i+1:]...)
			return
		}
	}
}

func (g *Grid) hasRot() bool { return len(g.rots) > 0 }

// decayRot ages every rot tile, clearing any that crumble away or that were
// already consumed or grown over.
func (g *Grid) decayRot(dt float64) {
	kept := g.rots[:0]
	for _, rt := range g.rots {
		if g.cells[rt.row][rt.col] != Rot {
			continue // consumed by a bug or reclaimed by a root
		}
		rt.ttl -= dt / g.season.RotLifetimeMul()
		if rt.ttl <= 0 {
			g.cells[rt.row][rt.col] = Empty
			continue
		}
		kept = append(kept, rt)
	}
	g.rots = kept
}
