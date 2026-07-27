// Package plants holds the companion plants ("pets") the player seeds near a
// Core (GDD §5). They do not conduct water; each radiates an aura:
//
//   - Battery (yellow): +20% energy regen.
//   - Moss (blue):      -10% water cost of energy regen.
//   - Thorn (white):    keeps nearby water sources from freezing in winter
//     (takes effect in Stage 14).
//
// Auras are applied globally as economy multipliers, since the economy is a
// single pool; each plant's radius is drawn for feedback.
package plants

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Kind int

const (
	Battery Kind = iota
	Moss
	Thorn
	kindCount
)

const (
	SeedCost   = 80.0
	auraRadius = 3

	batteryEnergyBonus = 0.20
	mossWaterSaving    = 0.10
	waterMultFloor     = 0.5
)

type Plant struct {
	Col, Row int
	Kind     Kind
}

type Manager struct {
	plants []Plant
}

func NewManager() *Manager { return &Manager{} }

func (m *Manager) Add(col, row int, k Kind) {
	m.plants = append(m.plants, Plant{Col: col, Row: row, Kind: k})
}

func (m *Manager) Count() int { return len(m.plants) }

// NearThorn reports whether a Winter Thornbush sits within its aura radius of
// the cell — it shields nearby water sources from the winter drain penalty.
func (m *Manager) NearThorn(col, row int) bool {
	for _, p := range m.plants {
		if p.Kind == Thorn && abs(col-p.Col) <= auraRadius && abs(row-p.Row) <= auraRadius {
			return true
		}
	}
	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Modifiers aggregates the plant auras into economy multipliers: energy regen
// speed and water cost of that regen.
func (m *Manager) Modifiers() (energyMult, waterMult float64) {
	energyMult, waterMult = 1, 1
	for _, p := range m.plants {
		switch p.Kind {
		case Battery:
			energyMult += batteryEnergyBonus
		case Moss:
			waterMult -= mossWaterSaving
		}
	}
	if waterMult < waterMultFloor {
		waterMult = waterMultFloor
	}
	return energyMult, waterMult
}

var auraColor = [kindCount]color.RGBA{
	Battery: {0xff, 0xd7, 0x4c, 0xaa},
	Moss:    {0x4c, 0xa8, 0xff, 0xaa},
	Thorn:   {0xe8, 0xee, 0xf5, 0xaa},
}

func (m *Manager) Draw(screen *ebiten.Image, cell int) {
	ensureSprites()
	half := float32(cell) / 2
	for _, p := range m.plants {
		cx := float32(p.Col*cell) + half
		cy := float32(p.Row*cell) + half
		vector.StrokeCircle(screen, cx, cy, float32(auraRadius*cell), 1.5, auraColor[p.Kind], true)

		img := petSprites[p.Kind]
		op := &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}
		b := img.Bounds()
		op.GeoM.Scale(float64(cell)/float64(b.Dx()), float64(cell)/float64(b.Dy()))
		op.GeoM.Translate(float64(p.Col*cell), float64(p.Row*cell))
		screen.DrawImage(img, op)
	}
}
