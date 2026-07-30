package enemies

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const maxBugs = 200

type Manager struct {
	bugs []*Bug
}

func NewManager() *Manager { return &Manager{} }

func (m *Manager) Count() int { return len(m.bugs) }

// Spawn adds one bug of a wave-appropriate level at the wave's entry side.
func (m *Manager) Spawn(f Field, wave, side int) {
	if len(m.bugs) >= maxBugs {
		return
	}
	if c, r, ok := f.EdgeSpawn(side); ok {
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

// Draw renders each bug's texture, centered on its cell and rotated so it faces
// the direction it is moving (the sprites are drawn facing up).
func (m *Manager) Draw(screen *ebiten.Image, cell int) {
	ensureSprites()
	for _, b := range m.bugs {
		img := bugSprites[b.level]
		w := float64(img.Bounds().Dx())
		op := &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}
		op.GeoM.Translate(-w/2, -w/2)
		op.GeoM.Rotate(facingAngle(b.fcol, b.frow))
		op.GeoM.Translate(w/2, w/2)
		pad := (float64(cell) - w) / 2
		op.GeoM.Translate(float64(b.Col*cell)+pad, float64(b.Row*cell)+pad)
		screen.DrawImage(img, op)
	}
}

// facingAngle maps a step direction to a clockwise rotation for an up-facing
// sprite.
func facingAngle(dc, dr int) float64 {
	switch {
	case dc == 1:
		return math.Pi / 2
	case dc == -1:
		return -math.Pi / 2
	case dr == 1:
		return math.Pi
	default:
		return 0
	}
}
