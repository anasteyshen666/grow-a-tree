// Package game holds the top-level game state and the Ebiten game loop.
package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"growtree/internal/enemies"
	"growtree/internal/plants"
	"growtree/internal/resources"
	"growtree/internal/ui"
	"growtree/internal/waves"
	"growtree/internal/world"
)

const (
	fieldPx = world.Cols * world.CellSize // square play field (720)

	// The logical canvas is 16:9; the square field sits centered with a HUD
	// panel on each side. Fullscreen scales this whole canvas.
	ScreenWidth  = 1280
	ScreenHeight = 720

	fieldOX = (ScreenWidth - fieldPx) / 2
	fieldOY = (ScreenHeight - fieldPx) / 2

	secondsPerTick = 1.0 / 60.0
)

var panelColor = color.RGBA{R: 0x07, G: 0x09, B: 0x0d, A: 0xff}

type Game struct {
	grid   *world.Grid
	res    *resources.Resources
	bugs   *enemies.Manager
	waves  *waves.Manager
	plants *plants.Manager
	field  *ebiten.Image
	over   bool
}

func New() *Game {
	return &Game{
		grid:   world.NewGrid(),
		res:    resources.New(),
		bugs:   enemies.NewManager(),
		waves:  waves.NewManager(),
		plants: plants.NewManager(),
	}
}

var plantKeys = [...]struct {
	key  ebiten.Key
	kind plants.Kind
}{
	{ebiten.Key1, plants.Battery},
	{ebiten.Key2, plants.Moss},
	{ebiten.Key3, plants.Thorn},
}

// fieldCell maps the cursor to a grid cell, or ok=false when it is outside the
// play field (over a side panel).
func (g *Game) fieldCell() (col, row int, ok bool) {
	cx, cy := ebiten.CursorPosition()
	fx, fy := cx-fieldOX, cy-fieldOY
	if fx < 0 || fy < 0 || fx >= fieldPx || fy >= fieldPx {
		return 0, 0, false
	}
	col, row = world.CellAt(fx, fy)
	return col, row, true
}

func (g *Game) Update() error {
	if g.over {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			*g = *New()
		}
		return nil
	}

	g.grid.Update(secondsPerTick, g.waves.Number())
	g.res.SetCoreLinks(g.grid.CoreMerges())
	g.res.AddWater(g.grid.MineWater(secondsPerTick))
	energyMult, waterMult := g.plants.Modifiers()
	g.res.Update(secondsPerTick, energyMult, waterMult)

	for n := g.waves.Update(secondsPerTick, g.bugs.Count()); n > 0; n-- {
		g.bugs.Spawn(g.grid, g.waves.Number())
	}
	g.bugs.Update(secondsPerTick, g.grid)

	if g.grid.CoreCount() == 0 {
		g.over = true
		return nil
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if col, row, ok := g.fieldCell(); ok && g.grid.CanGrow(col, row) && g.res.TrySpendEnergy(resources.RootEnergyCost) {
			g.grid.Grow(col, row)
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if col, row, ok := g.fieldCell(); ok && g.grid.Cut(col, row) {
			g.res.AddEnergy(resources.RootRefund)
		}
	}
	for _, pk := range plantKeys {
		if inpututil.IsKeyJustPressed(pk.key) {
			if col, row, ok := g.fieldCell(); ok && g.grid.CanPlant(col, row) && g.res.TrySpendSeeds(plants.SeedCost) {
				g.grid.SetPlant(col, row)
				g.plants.Add(col, row, pk.kind)
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) && g.res.Seeds.Cur >= world.MatureSeedCost {
		if g.grid.Mature() {
			g.res.TrySpendSeeds(world.MatureSeedCost)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.field == nil {
		g.field = ebiten.NewImage(fieldPx, fieldPx)
	}
	screen.Fill(panelColor)

	g.field.Clear()
	g.grid.Draw(g.field)
	g.plants.Draw(g.field, world.CellSize)
	if col, row, ok := g.fieldCell(); ok {
		g.grid.DrawHover(g.field, col, row)
	}
	g.bugs.Draw(g.field, world.CellSize)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(fieldOX, fieldOY)
	screen.DrawImage(g.field, op)

	rightX := fieldOX + fieldPx + 16
	ui.DrawResources(screen, g.res, 16, 40)
	ui.DrawStatus(screen, g.grid.CoreCount(), g.grid.CoreHP(), 16, 220)
	ui.DrawControls(screen, 16, ScreenHeight-140)
	ui.DrawWaveInfo(screen, rightX, 40, g.waves.Number(), g.waves.InPrep(), g.waves.PrepRemaining())

	if g.over {
		ui.DrawGameOver(screen, ScreenWidth, ScreenHeight, g.waves.Number())
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
