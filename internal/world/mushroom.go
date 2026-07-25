package world

import "math/rand"

const (
	maxSpores          = 4
	sporeSpawnInterval = 7.0 // seconds between spore appearances
	sporeStartWave     = 5   // spores only begin appearing from this wave on
	mushroomSlowRadius = 3   // cells around a mushroom root where bugs are slowed
)

// Update ticks the mushroom system: from wave sporeStartWave on, spores appear
// over time and infect any root they end up next to, turning it into a mushroom
// root (GDD §4).
func (g *Grid) Update(dt float64, wave int) {
	if wave >= sporeStartWave {
		g.sporeTimer += dt
		if g.sporeTimer >= sporeSpawnInterval {
			g.sporeTimer = 0
			g.spawnSpore()
		}
	}
	g.infectFromSpores()
}

func (g *Grid) spawnSpore() {
	if len(g.spores) >= maxSpores {
		return
	}
	for try := 0; try < 40; try++ {
		c, r := rand.Intn(Cols), rand.Intn(Rows)
		if g.cells[r][c] == Empty {
			g.cells[r][c] = Spore
			g.spores = append(g.spores, [2]int{c, r})
			return
		}
	}
}

// infectFromSpores turns roots adjacent to a spore into mushroom roots and
// consumes the spore.
func (g *Grid) infectFromSpores() {
	kept := g.spores[:0]
	for _, s := range g.spores {
		c, r := s[0], s[1]
		if g.cells[r][c] != Spore {
			continue // already cleared some other way
		}
		infected := false
		for _, d := range neighbors {
			nc, nr := c+d[0], r+d[1]
			if g.InBounds(nc, nr) && g.cells[nr][nc] == Root {
				g.mushroom[nr][nc] = true
				infected = true
			}
		}
		if infected {
			g.cells[r][c] = Empty
		} else {
			kept = append(kept, s)
		}
	}
	g.spores = kept
}

// IsMushroom reports whether the cell is a mushroom-infected root.
func (g *Grid) IsMushroom(col, row int) bool {
	return g.InBounds(col, row) && g.cells[row][col] == Root && g.mushroom[row][col]
}

// SlowsBug reports whether a mushroom root sits within mushroomSlowRadius cells
// of the cell, halving the speed of a bug crawling there.
func (g *Grid) SlowsBug(col, row int) bool {
	for dr := -mushroomSlowRadius; dr <= mushroomSlowRadius; dr++ {
		for dc := -mushroomSlowRadius; dc <= mushroomSlowRadius; dc++ {
			c, r := col+dc, row+dr
			if g.InBounds(c, r) && g.cells[r][c] == Root && g.mushroom[r][c] {
				return true
			}
		}
	}
	return false
}
