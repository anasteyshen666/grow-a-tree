package world

import "math/rand"

// walkable reports whether a bug can crawl across the cell.
func (g *Grid) walkable(col, row int) bool {
	k := g.Kind(col, row)
	return k == Empty || k == Rot
}

// hasRootOrCore reports whether any orthogonal neighbor is a Root or Core,
// regardless of connectivity — bugs gnaw physical roots, alive or not.
func (g *Grid) hasRootOrCore(col, row int) bool {
	for _, d := range neighbors {
		switch g.Kind(col+d[0], row+d[1]) {
		case Root, Core:
			return true
		}
	}
	return false
}

// AdjacentTarget returns an orthogonally adjacent Root (preferred) or Core for a
// bug to gnaw.
func (g *Grid) AdjacentTarget(col, row int) (int, int, bool) {
	coreC, coreR, haveCore := 0, 0, false
	for _, d := range neighbors {
		c, r := col+d[0], row+d[1]
		switch g.Kind(c, r) {
		case Root:
			return c, r, true
		case Core:
			coreC, coreR, haveCore = c, r, true
		}
	}
	return coreC, coreR, haveCore
}

// IsRot reports whether the cell is rot (bug bait).
func (g *Grid) IsRot(col, row int) bool { return g.Kind(col, row) == Rot }

// EatRot removes a rot tile a bug has consumed.
func (g *Grid) EatRot(col, row int) {
	if g.Kind(col, row) == Rot {
		g.cells[row][col] = Empty
		g.rotCount--
	}
}

// NextRotStep runs a BFS toward the nearest rot tile and returns the first step
// of the shortest path. ok is false if there is no rot or no path to it.
func (g *Grid) NextRotStep(sc, sr int) (int, int, bool) {
	if g.rotCount == 0 {
		return 0, 0, false
	}
	var visited [Rows][Cols]bool
	var firstC, firstR [Rows][Cols]int
	visited[sr][sc] = true
	queue := [][2]int{{sc, sr}}

	for len(queue) > 0 {
		c0, r0 := queue[0][0], queue[0][1]
		queue = queue[1:]
		for _, d := range neighbors {
			c, r := c0+d[0], r0+d[1]
			if !g.InBounds(c, r) || visited[r][c] || !g.walkable(c, r) {
				continue
			}
			visited[r][c] = true
			if c0 == sc && r0 == sr {
				firstC[r][c], firstR[r][c] = c, r
			} else {
				firstC[r][c], firstR[r][c] = firstC[r0][c0], firstR[r0][c0]
			}
			if g.cells[r][c] == Rot {
				return firstC[r][c], firstR[r][c], true
			}
			queue = append(queue, [2]int{c, r})
		}
	}
	return 0, 0, false
}

// NextStep runs a BFS from (sc,sr) across walkable cells and returns the first
// step of the shortest path toward the network — a cell adjacent to a Root or
// Core. ok is false if no path exists.
func (g *Grid) NextStep(sc, sr int) (int, int, bool) {
	var visited [Rows][Cols]bool
	var firstC, firstR [Rows][Cols]int
	visited[sr][sc] = true
	queue := [][2]int{{sc, sr}}

	for len(queue) > 0 {
		c0, r0 := queue[0][0], queue[0][1]
		queue = queue[1:]
		for _, d := range neighbors {
			c, r := c0+d[0], r0+d[1]
			if !g.InBounds(c, r) || visited[r][c] || !g.walkable(c, r) {
				continue
			}
			visited[r][c] = true
			if c0 == sc && r0 == sr {
				firstC[r][c], firstR[r][c] = c, r
			} else {
				firstC[r][c], firstR[r][c] = firstC[r0][c0], firstR[r0][c0]
			}
			if g.hasRootOrCore(c, r) {
				return firstC[r][c], firstR[r][c], true
			}
			queue = append(queue, [2]int{c, r})
		}
	}
	return 0, 0, false
}

const coreBiteDamage = 5

// Damage destroys a gnawed root (the network beyond it falls off) or chips the
// Core's health. Reaching 0 HP triggers game over in a later stage.
func (g *Grid) Damage(col, row int) {
	switch g.Kind(col, row) {
	case Root:
		g.cells[row][col] = Empty
		g.recomputeConnectivity()
	case Core:
		g.coreHP -= coreBiteDamage
		if g.coreHP < 0 {
			g.coreHP = 0
		}
	}
}

// RandomEdgeSpawn returns a random empty cell on the map border, where a bug
// can crawl in from.
func (g *Grid) RandomEdgeSpawn() (int, int, bool) {
	for try := 0; try < 40; try++ {
		var c, r int
		switch rand.Intn(4) {
		case 0:
			c, r = rand.Intn(Cols), 0
		case 1:
			c, r = rand.Intn(Cols), Rows-1
		case 2:
			c, r = 0, rand.Intn(Rows)
		default:
			c, r = Cols-1, rand.Intn(Rows)
		}
		if g.cells[r][c] == Empty {
			return c, r, true
		}
	}
	return 0, 0, false
}
