package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	colorGridLine = color.RGBA{0x16, 0x1b, 0x24, 0xff}
	colorCore     = color.RGBA{0x5a, 0xff, 0xc4, 0xff}
	colorRoot     = color.RGBA{0x2f, 0xa8, 0x6b, 0xff}
	colorGrowOK   = color.RGBA{0x5a, 0xff, 0xc4, 0x55}
	colorGrowBad  = color.RGBA{0xff, 0x4c, 0x5e, 0x3c}
)

func (g *Grid) Draw(screen *ebiten.Image) {
	g.drawLines(screen)
	for row := 0; row < Rows; row++ {
		for col := 0; col < Cols; col++ {
			switch g.cells[row][col] {
			case Core:
				fillCell(screen, col, row, colorCore)
			case Root:
				fillCell(screen, col, row, colorRoot)
			}
		}
	}
}

// DrawHover previews the cell under the cursor: green if a root can grow there,
// red if not.
func (g *Grid) DrawHover(screen *ebiten.Image, col, row int) {
	if !g.InBounds(col, row) || g.cells[row][col] != Empty {
		return
	}
	if g.touchesNetwork(col, row) {
		fillCell(screen, col, row, colorGrowOK)
	} else {
		fillCell(screen, col, row, colorGrowBad)
	}
}

func fillCell(screen *ebiten.Image, col, row int, c color.Color) {
	vector.DrawFilledRect(screen, float32(col*CellSize), float32(row*CellSize), CellSize, CellSize, c, false)
}

func (g *Grid) drawLines(screen *ebiten.Image) {
	w, h := float32(Cols*CellSize), float32(Rows*CellSize)
	for col := 0; col <= Cols; col++ {
		x := float32(col * CellSize)
		vector.StrokeLine(screen, x, 0, x, h, 1, colorGridLine, false)
	}
	for row := 0; row <= Rows; row++ {
		y := float32(row * CellSize)
		vector.StrokeLine(screen, 0, y, w, y, 1, colorGridLine, false)
	}
}
