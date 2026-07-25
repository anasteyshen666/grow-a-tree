package world

// CellKind is what occupies a single grid cell. More kinds (Rot, Water, ...)
// arrive in later stages.
type CellKind uint8

const (
	Empty CellKind = iota
	Core
	Root
	Rot
	Water
)

const (
	Cols     = 40
	Rows     = 30
	CellSize = 24

	CoreMaxHP = 100
)

type Grid struct {
	cells [Rows][Cols]CellKind
	// connected[r][c] is true for a Core or a Root reachable from a Core through
	// other roots. Recomputed on every change to the network (see network.go).
	connected [Rows][Cols]bool
	sources   []*waterSource
	coreHP    int
	rotCount  int
}

func NewGrid() *Grid {
	g := &Grid{coreHP: CoreMaxHP}
	g.placeCore()
	g.spawnSources()
	g.recomputeConnectivity()
	return g
}

// CoreHP is the remaining health of the Core (0..CoreMaxHP).
func (g *Grid) CoreHP() int { return g.coreHP }

func (g *Grid) placeCore() {
	col0, row0 := Cols/2-1, Rows/2-1
	for dr := 0; dr < 2; dr++ {
		for dc := 0; dc < 2; dc++ {
			g.cells[row0+dr][col0+dc] = Core
		}
	}
}

func (g *Grid) InBounds(col, row int) bool {
	return col >= 0 && col < Cols && row >= 0 && row < Rows
}

func (g *Grid) Kind(col, row int) CellKind {
	if !g.InBounds(col, row) {
		return Empty
	}
	return g.cells[row][col]
}

var neighbors = [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

// touchesNetwork reports whether the cell shares an edge with the *live*
// network: the Core or a root still connected to it. A disconnected root is
// dead and cannot be built upon.
func (g *Grid) touchesNetwork(col, row int) bool {
	for _, d := range neighbors {
		c, r := col+d[0], row+d[1]
		if !g.InBounds(c, r) {
			continue
		}
		switch g.cells[r][c] {
		case Core:
			return true
		case Root:
			if g.connected[r][c] {
				return true
			}
		}
	}
	return false
}

// CanGrow reports whether a new root may be placed at the cell. Empty and Rot
// cells are both valid targets (a root can reclaim rotted ground).
func (g *Grid) CanGrow(col, row int) bool {
	if !g.InBounds(col, row) {
		return false
	}
	k := g.cells[row][col]
	return (k == Empty || k == Rot) && g.touchesNetwork(col, row)
}

// Grow places a root and returns whether it happened.
func (g *Grid) Grow(col, row int) bool {
	if !g.CanGrow(col, row) {
		return false
	}
	if g.cells[row][col] == Rot {
		g.rotCount--
	}
	g.cells[row][col] = Root
	g.recomputeConnectivity()
	return true
}

// Cut removes a player's root, leaving Rot behind. The Core cannot be cut.
// Returns whether a root was actually removed.
func (g *Grid) Cut(col, row int) bool {
	if !g.InBounds(col, row) || g.cells[row][col] != Root {
		return false
	}
	g.cells[row][col] = Rot
	g.rotCount++
	g.recomputeConnectivity()
	return true
}

// CellAt converts screen pixels to grid coordinates.
func CellAt(x, y int) (col, row int) {
	return x / CellSize, y / CellSize
}
