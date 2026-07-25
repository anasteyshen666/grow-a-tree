// Package game holds the top-level game state and the Ebiten game loop.
package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"growtree/internal/resources"
	"growtree/internal/ui"
	"growtree/internal/world"
)

const (
	ScreenWidth  = world.Cols * world.CellSize
	ScreenHeight = world.Rows * world.CellSize

	secondsPerTick = 1.0 / 60.0
)

var backgroundColor = color.RGBA{R: 0x0b, G: 0x0d, B: 0x12, A: 0xff}

type Game struct {
	grid *world.Grid
	res  *resources.Resources
}

func New() *Game {
	return &Game{
		grid: world.NewGrid(),
		res:  resources.New(),
	}
}

func (g *Game) Update() error {
	g.res.Update(secondsPerTick)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		col, row := world.CellAt(ebiten.CursorPosition())
		if g.grid.CanGrow(col, row) && g.res.TrySpendEnergy(resources.RootEnergyCost) {
			g.grid.Grow(col, row)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	g.grid.Draw(screen)
	col, row := world.CellAt(ebiten.CursorPosition())
	g.grid.DrawHover(screen, col, row)
	ui.DrawResources(screen, g.res)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
