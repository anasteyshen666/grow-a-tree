// Package waves drives the endless assault (GDD §7): a prep timer between waves,
// then releasing a wave of bugs in a steady trickle, tracking the wave number
// and scaling the wave size with it. The next prep does not begin until every
// bug from the current wave has been cleared.
package waves

import "math/rand"

const (
	prepTime     = 30.0 // seconds of calm before a wave
	spawnCadence = 0.6  // seconds between bug releases within a wave
	baseWaveSize = 3
)

type phase int

const (
	prep phase = iota
	spawning
	clearing
)

type Manager struct {
	wave      int
	phase     phase
	timer     float64
	remaining int
	side      int // edge the current/next swarm comes from (0..3)
}

func NewManager() *Manager {
	return &Manager{phase: prep, timer: prepTime, side: rand.Intn(4)}
}

// waveSize is how many bugs a wave releases. The count steps up every 2 waves
// (in step with the bug-level unlocks), and each 10-wave decade grows faster
// than the last, so later decades have more bugs per 2-wave step.
func waveSize(n int) int {
	if n < 1 {
		n = 1
	}
	decade := (n - 1) / 10
	step := ((n - 1) % 10) / 2
	return baseWaveSize + (decade+1)*(step+1)*2
}

// Update advances the wave clock and returns how many bugs to spawn this tick.
// aliveBugs is the current number of live bugs; the prep for the next wave only
// starts once it reaches zero.
func (m *Manager) Update(dt float64, aliveBugs int) int {
	switch m.phase {
	case prep:
		m.timer -= dt
		if m.timer <= 0 {
			m.wave++
			m.remaining = waveSize(m.wave)
			m.phase = spawning
			m.timer = 0
		}
		return 0

	case spawning:
		m.timer -= dt
		spawned := 0
		for m.timer <= 0 && m.remaining > 0 {
			spawned++
			m.remaining--
			m.timer += spawnCadence
		}
		if m.remaining == 0 {
			m.phase = clearing
		}
		return spawned

	case clearing:
		if aliveBugs == 0 {
			m.phase = prep
			m.timer = prepTime
			m.side = rand.Intn(4) // pick the next swarm's side during prep
		}
		return 0
	}
	return 0
}

func (m *Manager) Number() int { return m.wave }

// Side is the edge (0..3) the current/next swarm enters from.
func (m *Manager) Side() int { return m.side }

func (m *Manager) InPrep() bool { return m.phase == prep }

// PrepRemaining is the seconds left before the next wave, or 0 while spawning
// or waiting for the field to be cleared.
func (m *Manager) PrepRemaining() float64 {
	if m.phase == prep {
		return m.timer
	}
	return 0
}
