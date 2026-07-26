// Package fx renders short-lived visual effects — spawn pops, hit flashes, and
// pixel bursts — in the play field's coordinate space.
package fx

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type particle struct {
	x, y, vx, vy float64
	life, max    float64
	size         float64
	c            color.RGBA
}

type overlay struct {
	col, row  int
	life, max float64
	c         color.RGBA
	ring      bool // true = expanding outline (pop), false = filling flash
}

type Manager struct {
	cell      int
	particles []*particle
	overlays  []*overlay
}

func New(cellSize int) *Manager { return &Manager{cell: cellSize} }

// Burst throws a spray of pixels out of a cell (root/core destroyed).
func (m *Manager) Burst(col, row int, c color.RGBA, count int) {
	cx := float64(col*m.cell) + float64(m.cell)/2
	cy := float64(row*m.cell) + float64(m.cell)/2
	for i := 0; i < count; i++ {
		ang := rand.Float64() * 2 * math.Pi
		spd := 40 + rand.Float64()*90
		life := 0.35 + rand.Float64()*0.35
		m.particles = append(m.particles, &particle{
			x: cx, y: cy,
			vx:   spd * math.Cos(ang),
			vy:   spd * math.Sin(ang),
			life: life, max: life,
			size: 2 + rand.Float64()*2,
			c:    c,
		})
	}
}

// Pop plays an expanding-outline "appear" effect on a cell.
func (m *Manager) Pop(col, row int, c color.RGBA) {
	m.overlays = append(m.overlays, &overlay{col: col, row: row, life: 0.3, max: 0.3, c: c, ring: true})
}

// Flash briefly tints a cell (a core taking a hit).
func (m *Manager) Flash(col, row int, c color.RGBA) {
	m.overlays = append(m.overlays, &overlay{col: col, row: row, life: 0.18, max: 0.18, c: c, ring: false})
}

func (m *Manager) Update(dt float64) {
	kp := m.particles[:0]
	for _, p := range m.particles {
		p.life -= dt
		if p.life <= 0 {
			continue
		}
		p.x += p.vx * dt
		p.y += p.vy * dt
		p.vy += 140 * dt // slight gravity
		kp = append(kp, p)
	}
	m.particles = kp

	ko := m.overlays[:0]
	for _, o := range m.overlays {
		if o.life -= dt; o.life > 0 {
			ko = append(ko, o)
		}
	}
	m.overlays = ko
}

func (m *Manager) Draw(dst *ebiten.Image) {
	for _, o := range m.overlays {
		t := o.life / o.max // 1 -> 0
		x := float32(o.col * m.cell)
		y := float32(o.row * m.cell)
		s := float32(m.cell)
		if o.ring {
			grow := float32(1-t) * float32(m.cell) * 0.5
			c := withAlpha(o.c, t)
			vector.StrokeRect(dst, x-grow, y-grow, s+2*grow, s+2*grow, 2, c, false)
		} else {
			vector.DrawFilledRect(dst, x, y, s, s, withAlpha(o.c, t*0.8), false)
		}
	}
	for _, p := range m.particles {
		a := p.life / p.max
		vector.DrawFilledRect(dst, float32(p.x), float32(p.y), float32(p.size), float32(p.size), withAlpha(p.c, a), false)
	}
}

func withAlpha(c color.RGBA, f float64) color.RGBA {
	if f < 0 {
		f = 0
	} else if f > 1 {
		f = 1
	}
	return color.RGBA{c.R, c.G, c.B, uint8(float64(c.A) * f)}
}
