package world

// recomputeConnectivity flood-fills from every Core cell across adjacent roots,
// marking which roots are still fed by a Core. Roots the fill never reaches are
// "cut off": they conduct nothing until the gap is patched (GDD §3).
func (g *Grid) recomputeConnectivity() {
	var queue [][2]int
	for r := 0; r < Rows; r++ {
		for c := 0; c < Cols; c++ {
			seed := g.cells[r][c] == Core
			g.connected[r][c] = seed
			if seed {
				queue = append(queue, [2]int{c, r})
			}
		}
	}

	for len(queue) > 0 {
		c0, r0 := queue[0][0], queue[0][1]
		queue = queue[1:]
		for _, d := range neighbors {
			c, r := c0+d[0], r0+d[1]
			if !g.InBounds(c, r) || g.connected[r][c] {
				continue
			}
			if g.cells[r][c] == Root {
				g.connected[r][c] = true
				queue = append(queue, [2]int{c, r})
			}
		}
	}
}

// IsConnected reports whether the cell is part of the live network (a Core, or
// a Root still reachable from one).
func (g *Grid) IsConnected(col, row int) bool {
	return g.InBounds(col, row) && g.connected[row][col]
}
