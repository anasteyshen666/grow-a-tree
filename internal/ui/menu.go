package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	colorTitle = color.RGBA{0x6c, 0xe0, 0x4c, 0xff}
	btnBorder  = color.RGBA{0xd0, 0xd7, 0xe0, 0xff}
	btnBg      = color.RGBA{0x18, 0x24, 0x1a, 0xff}
	btnBgHover = color.RGBA{0x2a, 0x3e, 0x2c, 0xff}
	btnText    = color.RGBA{0xe8, 0xee, 0xf5, 0xff}
	veilColor  = color.RGBA{0x00, 0x00, 0x00, 0xb4}
)

const (
	TitleSize  = 48.0
	buttonSize = 22.0
)

// Button is a clickable pixel-styled rectangle.
type Button struct {
	X, Y, W, H int
	Label      string
}

func (b Button) Contains(mx, my int) bool {
	return mx >= b.X && mx < b.X+b.W && my >= b.Y && my < b.Y+b.H
}

func (b Button) Draw(dst *ebiten.Image, hovered bool) {
	bg := btnBg
	if hovered {
		bg = btnBgHover
	}
	vector.DrawFilledRect(dst, float32(b.X-3), float32(b.Y-3), float32(b.W+6), float32(b.H+6), btnBorder, false)
	vector.DrawFilledRect(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), bg, false)
	tw := textWidth(b.Label, buttonSize)
	drawText(dst, b.Label, b.X+(b.W-tw)/2, b.Y+(b.H-int(buttonSize))/2, buttonSize, btnText)
}

// DrawTitle draws the green game title centered horizontally at y.
func DrawTitle(dst *ebiten.Image, screenW, y int) {
	s := "GROW A TREE"
	drawText(dst, s, (screenW-textWidth(s, TitleSize))/2, y, TitleSize, colorTitle)
}

// DrawCenteredLabel draws a label centered horizontally at y.
func DrawCenteredLabel(dst *ebiten.Image, screenW, y int, s string, size float64, clr color.Color) {
	drawText(dst, s, (screenW-textWidth(s, size))/2, y, size, clr)
}

// DrawVeil dims the whole screen for the pause overlay.
func DrawVeil(dst *ebiten.Image, w, h int) {
	vector.DrawFilledRect(dst, 0, 0, float32(w), float32(h), veilColor, false)
}

// DrawVersion prints the build version in the bottom-left corner.
func DrawVersion(dst *ebiten.Image, w, h int) {
	drawText(dst, "beta v1.0", 14, h-28, bodySize, colorLabel)
}

// LabelColor exposes the neutral label color for callers.
var LabelColor color.Color = colorLabel
