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

// DrawCoreHP shows how many Cores stand and their total health (top-left, under
// the bars).
func DrawCoreHP(screen *ebiten.Image, cores, hp int) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Cores: %d   HP: %d", cores, hp), labelX, topY+3*rowH-2)
}

// DrawPlantHint shows the planting and mature keys along the bottom-left.
func DrawPlantHint(screen *ebiten.Image, screenH int) {
	ebitenutil.DebugPrintAt(screen, "Plant near Core (100 seeds): 1 Battery  2 Moss  3 Thorn      M: mature new Core (80 seeds)", labelX, screenH-18)
}

var colorGameOverVeil = color.RGBA{0x00, 0x00, 0x00, 0xb0}

// DrawGameOver dims the field and shows the result with a restart hint.
func DrawGameOver(screen *ebiten.Image, w, h, wave int) {
	vector.DrawFilledRect(screen, 0, 0, float32(w), float32(h), colorGameOverVeil, false)
	ebitenutil.DebugPrintAt(screen, "GAME OVER", w/2-28, h/2-16)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("You survived to wave %d", wave), w/2-72, h/2)
	ebitenutil.DebugPrintAt(screen, "Press R to restart", w/2-56, h/2+16)
}

func drawBar(screen *ebiten.Image, y int, label string, p resources.Pool, c color.Color) {
	ebitenutil.DebugPrintAt(screen, label, labelX, y-2)
	vector.DrawFilledRect(screen, barX, float32(y), barW, barH, colorBarBg, false)
	if p.Max > 0 {
		vector.DrawFilledRect(screen, barX, float32(y), float32(barW)*float32(p.Cur/p.Max), barH, c, false)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d/%d", int(p.Cur), int(p.Max)), barX+barW+8, y-2)
}
