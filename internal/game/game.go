// Package game holds the top-level game state and the Ebiten game loop.
package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"growtree/internal/enemies"
	"growtree/internal/resources"
	"growtree/internal/ui"
	"growtree/internal/waves"
	"growtree/internal/world"
)

const (
	ScreenWidth  = world.Cols * world.CellSize
	ScreenHeight = world.Rows * world.CellSize

	secondsPerTick = 1.0 / 60.0
)

var backgroundColor = color.RGBA{R: 0x0b, G: 0x0d, B: 0x12, A: 0xff}

type Game struct {
	grid  *world.Grid
	res   *resources.Resources
	bugs  *enemies.Manager
	waves *waves.Manager
}

func New() *Game {
	return &Game{
		grid:  world.NewGrid(),
		res:   resources.New(),
		bugs:  enemies.NewManager(),
		waves: waves.NewManager(),
	}
}

func (g *Game) Update() error {
	g.res.AddWater(g.grid.MineWater(secondsPerTick))
	g.res.Update(secondsPerTick)

	for n := g.waves.Update(secondsPerTick, g.bugs.Count()); n > 0; n-- {
		g.bugs.Spawn(g.grid)
	}
	g.bugs.Update(secondsPerTick, g.grid)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		col, row := world.CellAt(ebiten.CursorPosition())
		if g.grid.CanGrow(col, row) && g.res.TrySpendEnergy(resources.RootEnergyCost) {
			g.grid.Grow(col, row)
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		col, row := world.CellAt(ebiten.CursorPosition())
		if g.grid.Cut(col, row) {
			g.res.AddEnergy(resources.RootRefund)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	g.grid.Draw(screen)
	col, row := world.CellAt(ebiten.CursorPosition())
	g.grid.DrawHover(screen, col, row)
	g.bugs.Draw(screen, world.CellSize)
	ui.DrawResources(screen, g.res)
	ui.DrawCoreHP(screen, g.grid.CoreHP())
	ui.DrawWaveInfo(screen, ScreenWidth-160, g.waves.Number(), g.waves.InPrep(), g.waves.PrepRemaining())
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
