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
	colorRootDead = color.RGBA{0x3a, 0x44, 0x3d, 0xff}
	colorRot      = color.RGBA{0x6b, 0x4a, 0x2f, 0xff}
	colorSource   = color.RGBA{0x3c, 0x9a, 0xff, 0xff}
	colorGrowOK   = color.RGBA{0x5a, 0xff, 0xc4, 0x55}
	colorCut      = color.RGBA{0xff, 0x6e, 0x4c, 0x66}
)

// gridLines caches the static grid overlay so we blit one image per frame
// instead of stroking ~70 lines every time.
var gridLines *ebiten.Image

func (g *Grid) Draw(screen *ebiten.Image) {
	drawLines(screen)
	for row := 0; row < Rows; row++ {
		for col := 0; col < Cols; col++ {
			switch g.cells[row][col] {
			case Core:
				fillCell(screen, col, row, colorCore)
			case Root:
				if g.connected[row][col] {
					fillCell(screen, col, row, colorRoot)
				} else {
					fillCell(screen, col, row, colorRootDead)
				}
			case Rot:
				fillCell(screen, col, row, colorRot)
			}
		}
	}
	g.drawSources(screen)
}

func (g *Grid) drawSources(screen *ebiten.Image) {
	for _, s := range g.sources {
		f := 0.4 + 0.6*(s.amount/SourceMaxAmount)
		fillCell(screen, s.col, s.row, scale(colorSource, f))
	}
}

// DrawHover previews the cell under the cursor: green where a root can grow,
// red over a cuttable root.
func (g *Grid) DrawHover(screen *ebiten.Image, col, row int) {
	if !g.InBounds(col, row) {
		return
	}
	switch {
	case g.CanGrow(col, row):
		fillCell(screen, col, row, colorGrowOK)
	case g.cells[row][col] == Root:
		fillCell(screen, col, row, colorCut)
	}
}

func fillCell(screen *ebiten.Image, col, row int, c color.Color) {
	vector.DrawFilledRect(screen, float32(col*CellSize), float32(row*CellSize), CellSize, CellSize, c, false)
}

func drawLines(screen *ebiten.Image) {
	if gridLines == nil {
		gridLines = ebiten.NewImage(Cols*CellSize, Rows*CellSize)
		w, h := float32(Cols*CellSize), float32(Rows*CellSize)
		for col := 0; col <= Cols; col++ {
			x := float32(col * CellSize)
			vector.StrokeLine(gridLines, x, 0, x, h, 1, colorGridLine, false)
		}
		for row := 0; row <= Rows; row++ {
			y := float32(row * CellSize)
			vector.StrokeLine(gridLines, 0, y, w, y, 1, colorGridLine, false)
		}
	}
	screen.DrawImage(gridLines, nil)
}

func scale(c color.RGBA, f float64) color.RGBA {
	return color.RGBA{uint8(float64(c.R) * f), uint8(float64(c.G) * f), uint8(float64(c.B) * f), c.A}
}
