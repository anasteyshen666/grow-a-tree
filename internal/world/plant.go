package world

const plantCoreDist = 4 // a plant must sit within this many cells of a Core

func (g *Grid) nearCore(col, row, dist int) bool {
	for dr := -dist; dr <= dist; dr++ {
		for dc := -dist; dc <= dist; dc++ {
			if g.Kind(col+dc, row+dr) == Core {
				return true
			}
		}
	}
	return false
}

// CanPlant reports whether a companion plant may be placed on the cell: it must
// be empty and near a Core (GDD §5).
func (g *Grid) CanPlant(col, row int) bool {
	return g.InBounds(col, row) && g.cells[row][col] == Empty && g.nearCore(col, row, plantCoreDist)
}

// SetPlant marks a cell as occupied by a companion plant. The plant's kind and
// aura live in the plants package; the grid only tracks that the cell is taken.
func (g *Grid) SetPlant(col, row int) {
	if g.CanPlant(col, row) {
		g.cells[row][col] = Plant
	}
}
