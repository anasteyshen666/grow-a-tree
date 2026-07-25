package enemies

// Field is the slice of the world the bug AI needs: where to enter, how to step
// toward the network, what to gnaw, and how to damage it. world.Grid implements it.
type Field interface {
	RandomEdgeSpawn() (col, row int, ok bool)
	NextStep(col, row int) (col2, row2 int, ok bool)
	AdjacentTarget(col, row int) (col2, row2 int, ok bool)
	Damage(col, row int)
}

const (
	moveInterval = 0.18 // seconds to crawl one cell
	gnawTime     = 0.9  // seconds to chew through a root
)

type Bug struct {
	Col, Row int
	timer    float64
}

// update advances one bug: gnaw an adjacent target, otherwise step toward it.
func (b *Bug) update(dt float64, f Field) {
	b.timer += dt
	if tc, tr, ok := f.AdjacentTarget(b.Col, b.Row); ok {
		if b.timer >= gnawTime {
			b.timer = 0
			f.Damage(tc, tr)
		}
		return
	}
	if b.timer >= moveInterval {
		b.timer = 0
		if c, r, ok := f.NextStep(b.Col, b.Row); ok {
			b.Col, b.Row = c, r
		}
	}
}
