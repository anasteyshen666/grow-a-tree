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
	// NearLiveNetwork reports whether a Core or connected Root sits next to the
	// cell — such bugs are within reach of the network's counterattack.
	NearLiveNetwork(col, row int) bool
}

// EnergyBank is the player's energy, which fighting bugs drains.
type EnergyBank interface {
	TrySpendEnergy(v float64) bool
}

const poisonTime = 0.7 // seconds a bug lingers on rot before dying

type Bug struct {
	Col, Row int
	level    int
	hp       float64
	timer    float64
}

// update advances one bug and reports whether it died (lured onto rot).
func (b *Bug) update(dt float64, f Field) (dead bool) {
	st := &levels[b.level]
	b.timer += dt

	// On bait: the bug stops, gets poisoned, and dies — consuming the rot.
	if f.IsRot(b.Col, b.Row) {
		if b.timer >= poisonTime {
			f.EatRot(b.Col, b.Row)
			return true
		}
		return false
	}

	// Rot lures bugs off the network, so chasing bait takes priority. The BFS
	// only runs on a move tick to stay off the hot path.
	if b.timer >= st.moveInterval {
		if c, r, ok := f.NextRotStep(b.Col, b.Row); ok {
			b.timer = 0
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
	if b.timer >= st.moveInterval {
		b.timer = 0
		if c, r, ok := f.NextStep(b.Col, b.Row); ok {
			b.Col, b.Row = c, r
		}
	}
	return false
}
