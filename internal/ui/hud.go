package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"growtree/internal/resources"
)

var (
	colorLabel  = color.RGBA{0xc8, 0xcf, 0xd8, 0xff}
	colorDanger = color.RGBA{0xff, 0x6e, 0x4c, 0xff}
	colorEnergy = color.RGBA{0x5a, 0xff, 0xc4, 0xff}
	colorWater  = color.RGBA{0x4c, 0xa8, 0xff, 0xff}
	colorSeeds  = color.RGBA{0xff, 0xd7, 0x4c, 0xff}

	barBorder = color.RGBA{0xd0, 0xd7, 0xe0, 0xff}
	barTrack  = color.RGBA{0x14, 0x18, 0x20, 0xff}
	barSeg    = color.RGBA{0x00, 0x00, 0x00, 0x66}
)

const (
	bodySize = 16.0
	barW     = 232
	barH     = 16
	segStep  = 12
)

var frameTick int

func blinkOn() bool { return (frameTick/18)%2 == 0 }

// drawBar renders a chunky, segmented pixel-style progress bar.
func drawBar(dst *ebiten.Image, x, y int, frac float64, fill, border color.Color) {
	if frac < 0 {
		frac = 0
	} else if frac > 1 {
		frac = 1
	}
	vector.DrawFilledRect(dst, float32(x-2), float32(y-2), float32(barW+4), float32(barH+4), border, false)
	vector.DrawFilledRect(dst, float32(x), float32(y), float32(barW), float32(barH), barTrack, false)
	if fw := int(float64(barW) * frac); fw > 0 {
		vector.DrawFilledRect(dst, float32(x), float32(y), float32(fw), float32(barH), fill, false)
	}
	for sx := x + segStep; sx < x+barW; sx += segStep {
		vector.StrokeLine(dst, float32(sx), float32(y), float32(sx), float32(y+barH), 1, barSeg, false)
	}
}

// DrawResources draws the three resource bars in the left panel at (x,y).
func DrawResources(dst *ebiten.Image, r *resources.Resources, x, y int) {
	frameTick++
	drawResource(dst, x, y+0, "ENERGY", r.Energy, colorEnergy, true)
	drawResource(dst, x, y+56, "WATER", r.Water, colorWater, true)
	drawResource(dst, x, y+112, "SEEDS", r.Seeds, colorSeeds, false)
}

// drawResource draws one labeled bar. When warnEmpty is set and the pool is
// empty, the bar border blinks red.
func drawResource(dst *ebiten.Image, x, y int, label string, p resources.Pool, fill color.Color, warnEmpty bool) {
	drawText(dst, label, x, y, bodySize, colorLabel)
	val := fmt.Sprintf("%d/%d", int(p.Cur+0.5), int(p.Max+0.5))
	drawText(dst, val, x+barW-textWidth(val, bodySize), y, bodySize, colorLabel)
	frac := 0.0
	if p.Max > 0 {
		frac = p.Cur / p.Max
	}
	var border color.Color = barBorder
	if warnEmpty && p.Cur < 1 && blinkOn() {
		border = colorDanger
	}
	drawBar(dst, x, y+22, frac, fill, border)
}

// DrawStatus shows the number of Cores and their total HP.
func DrawStatus(dst *ebiten.Image, cores, hp int, x, y int) {
	drawText(dst, fmt.Sprintf("CORES  %d", cores), x, y, bodySize, colorLabel)
	drawText(dst, fmt.Sprintf("HP     %d", hp), x, y+24, bodySize, colorLabel)
}

// DrawWaveInfo shows the current wave, its status, and the season, right-aligned
// to rightEdge so it always sits inside the right panel.
func DrawWaveInfo(dst *ebiten.Image, rightEdge, y, wave int, inPrep bool, prep float64, seasonName string) {
	line := func(s string, yy int, size float64, clr color.Color) {
		drawText(dst, s, rightEdge-textWidth(s, size), yy, size, clr)
	}
	shown := wave
	if inPrep {
		shown = wave + 1
	}
	line("WAVE", y, bodySize, colorLabel)
	line(fmt.Sprintf("%d", shown), y+24, 40, colorSeeds)
	if inPrep {
		line(fmt.Sprintf("NEXT IN %ds", int(prep+0.5)), y+76, bodySize, colorLabel)
	} else {
		line("ATTACKING", y+76, bodySize, colorDanger)
	}
	line(seasonName, y+108, bodySize, colorLabel)
}

// DrawControls lists the input hints along the bottom of the left panel.
func DrawControls(dst *ebiten.Image, x, y int) {
	lines := []string{
		"LMB  GROW ROOT",
		"RMB  CUT / ROT",
		"1 2 3  PLANT (80)",
		"M  NEW CORE (100)",
	}
	for i, s := range lines {
		drawText(dst, s, x, y+i*22, bodySize, colorLabel)
	}
}

var colorGameOverVeil = color.RGBA{0x00, 0x00, 0x00, 0xc0}

// DrawGameOver dims the screen and shows the result with a restart hint.
func DrawGameOver(dst *ebiten.Image, w, h, wave, best int) {
	vector.DrawFilledRect(dst, 0, 0, float32(w), float32(h), colorGameOverVeil, false)
	center(dst, "GAME OVER", h/2-70, 48, colorDanger, w)
	center(dst, fmt.Sprintf("YOU SURVIVED TO WAVE %d", wave), h/2-10, bodySize, colorLabel, w)
	center(dst, fmt.Sprintf("BEST WAVE  %d", best), h/2+16, bodySize, colorSeeds, w)
	center(dst, "PRESS R TO RESTART", h/2+46, bodySize, colorLabel, w)
}

func center(dst *ebiten.Image, s string, y int, size float64, clr color.Color, w int) {
	drawText(dst, s, (w-textWidth(s, size))/2, y, size, clr)
}
