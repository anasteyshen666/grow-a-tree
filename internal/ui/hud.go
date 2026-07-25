package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"growtree/internal/resources"
)

var (
	colorEnergy = color.RGBA{0x5a, 0xff, 0xc4, 0xff}
	colorWater  = color.RGBA{0x4c, 0xa8, 0xff, 0xff}
	colorSeeds  = color.RGBA{0xff, 0xd7, 0x4c, 0xff}
	colorBarBg  = color.RGBA{0x1a, 0x20, 0x2b, 0xff}
)

const (
	labelX = 12
	barX   = 70
	barW   = 150
	barH   = 12
	rowH   = 20
	topY   = 14
)

func DrawResources(screen *ebiten.Image, r *resources.Resources) {
	drawBar(screen, topY+0*rowH, "Energy", r.Energy, colorEnergy)
	drawBar(screen, topY+1*rowH, "Water", r.Water, colorWater)
	drawBar(screen, topY+2*rowH, "Seeds", r.Seeds, colorSeeds)
}

// DrawWaveInfo shows the current wave and its status (top-right).
func DrawWaveInfo(screen *ebiten.Image, x, wave int, inPrep bool, prep float64) {
	if inPrep {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Wave %d", wave+1), x, topY-2)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("starts in %.0fs", prep), x, topY+rowH-2)
	} else {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Wave %d", wave), x, topY-2)
		ebitenutil.DebugPrintAt(screen, "attacking!", x, topY+rowH-2)
	}
}

// DrawCoreHP shows the Core's health as a plain number (top-left, under the bars).
func DrawCoreHP(screen *ebiten.Image, hp int) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Core HP: %d", hp), labelX, topY+3*rowH-2)
}

func drawBar(screen *ebiten.Image, y int, label string, p resources.Pool, c color.Color) {
	ebitenutil.DebugPrintAt(screen, label, labelX, y-2)
	vector.DrawFilledRect(screen, barX, float32(y), barW, barH, colorBarBg, false)
	if p.Max > 0 {
		vector.DrawFilledRect(screen, barX, float32(y), float32(barW)*float32(p.Cur/p.Max), barH, c, false)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d/%d", int(p.Cur), int(p.Max)), barX+barW+8, y-2)
}
