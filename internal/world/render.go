package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	colorGridLine = color.RGBA{0x1b, 0x21, 0x2c, 0xff}
	colorGrowOK   = color.RGBA{0x5a, 0xff, 0xc4, 0x55}
	colorCut      = color.RGBA{0xff, 0x6e, 0x4c, 0x66}
)

// gridLines caches the static grid overlay so we blit one image per frame
// instead of stroking ~70 lines every time.
var gridLines *ebiten.Image

func (g *Grid) Draw(screen *ebiten.Image) {
	ensureSprites()
	tr, tg, tb := g.season.DirtTint()
	drawTerrain(screen, tr, tg, tb)
	drawLines(screen)

	for row := 0; row < Rows; row++ {
		for col := 0; col < Cols; col++ {
			switch g.cells[row][col] {
			case Core:
				drawTile(screen, sprCore, col, row)
			case Root:
				switch {
				case !g.connected[row][col]:
					drawTileDim(screen, sprRoot, col, row, deadRootDim) // dead root = dimmed
				case g.mushroom[row][col]:
					drawTile(screen, sprMushroom, col, row)
				default:
					drawTile(screen, sprRoot, col, row)
				}
			case Spore:
				drawTile(screen, sprSpore, col, row)
			}
		}
	}
	g.drawRots(screen)
	g.drawSources(screen)
}

// drawRots draws rot tiles, fading them out as they age toward crumbling away.
func (g *Grid) drawRots(screen *ebiten.Image) {
	for _, rt := range g.rots {
		a := float32(0.35 + 0.65*(rt.ttl/rotLifetime))
		drawTileAlpha(screen, sprRot, rt.col, rt.row, a)
	}
}

// drawSources draws water tiles, fading them as they deplete.
func (g *Grid) drawSources(screen *ebiten.Image) {
	for _, s := range g.sources {
		max := s.max
		if max <= 0 {
			max = SourceMaxAmount
		}
		a := float32(0.45 + 0.55*(s.amount/max))
		drawTileAlpha(screen, sprWater, s.col, s.row, a)
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
