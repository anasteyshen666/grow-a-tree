package enemies

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	maxBugs = 200

	attackDamagePerSec = 12.0 // damage the live network deals to an adjacent bug
	attackEnergyPerSec = 8.0  // energy the player burns per second of fighting
)

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
		lvl := levelForWave(wave)
		m.bugs = append(m.bugs, &Bug{Col: c, Row: r, level: lvl, hp: levels[lvl].hp})
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
	b.hp -= attackDamagePerSec * dt * (1 - levels[b.level].armor)
	return b.hp <= 0
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
