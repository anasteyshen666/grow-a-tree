package world

import "math/rand"

// MatureSeedCost is the seed price of maturing a new Core (spent by the caller).
const MatureSeedCost = 80

// core is one tree's 2x2 heart, tracked by its top-left cell and its own HP.
type core struct {
	col, row int
	hp       int
}

func (g *Grid) coreAt(col, row int) *core {
	for _, c := range g.cores {
		if col >= c.col && col <= c.col+1 && row >= c.row && row <= c.row+1 {
			return c
		}
	}
	return nil
}

// damageCore chips a Core's HP; at 0 it is destroyed, clearing its cells and
// dropping out of the network.
func (g *Grid) damageCore(col, row, dmg int) {
	c := g.coreAt(col, row)
	if c == nil {
		return
	}
	if c.hp -= dmg; c.hp > 0 {
		return
	}
	for dr := 0; dr < 2; dr++ {
		for dc := 0; dc < 2; dc++ {
			g.cells[c.row+dr][c.col+dc] = Empty
		}
	}
	g.removeCore(c)
	g.recomputeConnectivity()
}

func (g *Grid) removeCore(target *core) {
	for i, c := range g.cores {
		if c == target {
			g.cores = append(g.cores[:i], g.cores[i+1:]...)
			return
		}
	}
}

// Mature spawns a new Core on a free 2x2 patch and reports success. The caller
// pays the seed cost.
func (g *Grid) Mature() bool {
	for try := 0; try < 200; try++ {
		// keep the whole 2x2 at least one cell in from every border
		col0, row0 := rand.Intn(Cols-3)+1, rand.Intn(Rows-3)+1
		if g.empty2x2(col0, row0) {
			g.placeCore(col0, row0)
			g.recomputeConnectivity()
			return true
		}
	}
	return false
}

func (g *Grid) empty2x2(col0, row0 int) bool {
	for dr := 0; dr < 2; dr++ {
		for dc := 0; dc < 2; dc++ {
			if g.cells[row0+dr][col0+dc] != Empty {
				return false
			}
		}
	}
	return true
}

// CoreMerges is how many core-to-core links exist: total Cores minus the number
// of connected groups. Linking different trees raises the water cap (GDD §2).
func (g *Grid) CoreMerges() int { return g.coreMerges }

func (g *Grid) isNetworkCell(col, row int) bool {
	k := g.cells[row][col]
	return k == Core || k == Root
}

// computeCoreMerges labels connected components over Core|Root cells, then
// counts cores that share a component beyond the first of each group.
func (g *Grid) computeCoreMerges() {
	var comp [Rows][Cols]int
	id := 0
	for r := 0; r < Rows; r++ {
		for c := 0; c < Cols; c++ {
			if comp[r][c] == 0 && g.isNetworkCell(c, r) {
				id++
				g.labelComponent(c, r, id, &comp)
			}
		}
	}
	groups := map[int]bool{}
	for _, cr := range g.cores {
		groups[comp[cr.row][cr.col]] = true
	}
	g.coreMerges = len(g.cores) - len(groups)
}

func (g *Grid) labelComponent(sc, sr, id int, comp *[Rows][Cols]int) {
	comp[sr][sc] = id
	queue := [][2]int{{sc, sr}}
	for len(queue) > 0 {
		c0, r0 := queue[0][0], queue[0][1]
		queue = queue[1:]
		for _, d := range neighbors {
			c, r := c0+d[0], r0+d[1]
			if g.InBounds(c, r) && comp[r][c] == 0 && g.isNetworkCell(c, r) {
				comp[r][c] = id
				queue = append(queue, [2]int{c, r})
			}
		}
	}
}
