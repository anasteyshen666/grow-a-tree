package enemies

// Field is the slice of the world the bug AI needs: where to enter, how to step
// toward the network or toward rot bait, what to gnaw, and how to damage it.
// world.Grid implements it.
type Field interface {
	RandomEdgeSpawn() (col, row int, ok bool)
	NextStep(col, row int) (col2, row2 int, ok bool)
	NextRotStep(col, row int) (col2, row2 int, ok bool)
	AdjacentTarget(col, row int) (col2, row2 int, ok bool)
	IsRot(col, row int) bool
	EatRot(col, row int)
	Damage(col, row, dmg int)
	// SlowsBug reports whether a mushroom root is next to the cell, halving the
	// bug's crawl speed there.
	SlowsBug(col, row int) bool
}

const (
	poisonTime     = 0.7 // seconds a bug lingers on a rot tile before eating it
	poisonSlowTime = 3.0 // seconds a bug crawls at half speed after eating rot
)

type Bug struct {
	Col, Row   int
	fcol, frow int // facing direction (0,0 = up), set to the last step
	level      int
	eaten      int     // rot tiles consumed so far
	poison     float64 // seconds of poison slowdown remaining
	timer      float64
}

// face records the direction of the step from the bug's current cell to (c,r).
func (b *Bug) face(c, r int) { b.fcol, b.frow = c-b.Col, r-b.Row }

// update advances one bug and reports whether it died. Rot bait is the only
// thing that kills a bug: a level-N bug must eat N rot tiles.
func (b *Bug) update(dt float64, f Field) (dead bool) {
	st := &levels[b.level]
	b.timer += dt
	if b.poison > 0 {
		b.poison -= dt
	}

	// On bait: the bug stops, then eats the rot. Enough rot poisons it to death;
	// otherwise it is left poisoned (slowed) as it seeks the next tile.
	if f.IsRot(b.Col, b.Row) {
		if b.timer >= poisonTime {
			f.EatRot(b.Col, b.Row)
			b.eaten++
			if b.eaten >= b.level {
				return true
			}
			b.poison = poisonSlowTime
			b.timer = 0
		}
		return false
	}

	// Crawling next to a mushroom root, or while poisoned, is twice as slow.
	move := st.moveInterval
	if f.SlowsBug(b.Col, b.Row) {
		move *= 2
	}
	if b.poison > 0 {
		move *= 2
	}

	// Rot lures bugs off the network, so chasing bait takes priority. The BFS
	// only runs on a move tick to stay off the hot path.
	if b.timer >= move {
		if c, r, ok := f.NextRotStep(b.Col, b.Row); ok {
			b.timer = 0
			b.face(c, r)
			b.Col, b.Row = c, r
			return false
		}
	}

	// Gnaw an adjacent root/core.
	if tc, tr, ok := f.AdjacentTarget(b.Col, b.Row); ok {
		if b.timer >= st.gnawTime {
			b.timer = 0
			f.Damage(tc, tr, st.biteDamage)
		}
		return false
	}

	// Otherwise crawl toward the network.
	if b.timer >= move {
		b.timer = 0
		if c, r, ok := f.NextStep(b.Col, b.Row); ok {
			b.face(c, r)
			b.Col, b.Row = c, r
		}
	}
	return false
}
