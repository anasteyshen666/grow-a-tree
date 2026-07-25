package enemies

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	maxBugs = 200

	bugHP              = 30.0 // hits a bug can take before dying
	attackDamagePerSec = 12.0 // damage the live network deals to an adjacent bug
	attackEnergyPerSec = 8.0  // energy the player burns per second of fighting
)

var colorBug = color.RGBA{0xff, 0x4c, 0x5e, 0xff}

type Manager struct {
	bugs []*Bug
}

func NewManager() *Manager { return &Manager{} }

func (m *Manager) Count() int { return len(m.bugs) }

// Spawn adds one bug at a random map edge (driven by the wave manager).
func (m *Manager) Spawn(f Field) {
	if len(m.bugs) >= maxBugs {
		return
	}
	if c, r, ok := f.RandomEdgeSpawn(); ok {
		m.bugs = append(m.bugs, &Bug{Col: c, Row: r, hp: bugHP})
	}
}

func (m *Manager) Update(dt float64, f Field, bank EnergyBank) {
	kept := m.bugs[:0]
	for _, b := range m.bugs {
		if b.update(dt, f) {
			continue // died on rot bait
		}
		if attack(dt, b, f, bank) {
			continue // killed by the network
		}
		kept = append(kept, b)
	}
	m.bugs = kept
}

// attack lets the live network chip an adjacent bug, spending energy to do so.
// Returns whether the bug died.
func attack(dt float64, b *Bug, f Field, bank EnergyBank) bool {
	if !f.NearLiveNetwork(b.Col, b.Row) {
		return false
	}
	if !bank.TrySpendEnergy(attackEnergyPerSec * dt) {
		return false // no energy to fight back
	}
	b.hp -= attackDamagePerSec * dt
	return b.hp <= 0
}

// Draw renders each bug as a square a little smaller than a cell.
func (m *Manager) Draw(screen *ebiten.Image, cell int) {
	size := float32(cell) * 0.7
	pad := (float32(cell) - size) / 2
	for _, b := range m.bugs {
		x := float32(b.Col*cell) + pad
		y := float32(b.Row*cell) + pad
		vector.DrawFilledRect(screen, x, y, size, size, colorBug, false)
	}
}
