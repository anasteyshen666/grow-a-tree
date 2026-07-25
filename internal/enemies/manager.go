package enemies

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const maxBugs = 200

type Manager struct {
	bugs []*Bug
}

func NewManager() *Manager { return &Manager{} }

func (m *Manager) Count() int { return len(m.bugs) }

// Spawn adds one bug of a wave-appropriate level at a random map edge.
func (m *Manager) Spawn(f Field, wave int) {
	if len(m.bugs) >= maxBugs {
		return
	}
	if c, r, ok := f.RandomEdgeSpawn(); ok {
		m.bugs = append(m.bugs, &Bug{Col: c, Row: r, level: levelForWave(wave)})
	}
}

func (m *Manager) Update(dt float64, f Field) {
	kept := m.bugs[:0]
	for _, b := range m.bugs {
		if b.update(dt, f) {
			continue // died on rot bait
		}
		kept = append(kept, b)
	}
	m.bugs = kept
}

// Draw renders each bug as a square sized and colored by its level.
func (m *Manager) Draw(screen *ebiten.Image, cell int) {
	for _, b := range m.bugs {
		st := &levels[b.level]
		size := float32(cell) * st.size
		pad := (float32(cell) - size) / 2
		x := float32(b.Col*cell) + pad
		y := float32(b.Row*cell) + pad
		vector.DrawFilledRect(screen, x, y, size, size, st.col, false)
	}
}
