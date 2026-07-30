// Package game holds the top-level game state and the Ebiten game loop.
package game

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"growtree/internal/audio"
	"growtree/internal/enemies"
	"growtree/internal/fx"
	"growtree/internal/plants"
	"growtree/internal/resources"
	"growtree/internal/save"
	"growtree/internal/season"
	"growtree/internal/ui"
	"growtree/internal/waves"
	"growtree/internal/world"
)

const (
	fieldPx = world.Cols * world.CellSize // square play field, in logical pixels

	// Initial windowed size (fullscreen overrides this via Layout).
	ScreenWidth  = 1280
	ScreenHeight = 720

	secondsPerTick = 1.0 / 60.0
)

var panelColor = color.RGBA{R: 0x07, G: 0x09, B: 0x0d, A: 0xff}

type state int

const (
	stateMenu state = iota
	statePlaying
	statePaused
	stateOver
)

type Game struct {
	grid        *world.Grid
	res         *resources.Resources
	bugs        *enemies.Manager
	waves       *waves.Manager
	plants      *plants.Manager
	fx          *fx.Manager
	field       *ebiten.Image
	bloom       *ebiten.Image
	lightmap    *ebiten.Image
	lightSprite *ebiten.Image

	screenW, screenH int
	state            state
	high             int
	night            float64 // 0 = day, 1 = full night
}

func New() *Game {
	audio.Init()
	audio.StartMusic()
	g := &Game{screenW: ScreenWidth, screenH: ScreenHeight, high: save.Load()}
	g.startRun()
	g.state = stateMenu
	return g
}

// startRun (re)initializes a fresh playing session.
func (g *Game) startRun() {
	g.grid = world.NewGrid()
	g.res = resources.New()
	g.bugs = enemies.NewManager()
	g.waves = waves.NewManager()
	g.plants = plants.NewManager()
	g.fx = fx.New(world.CellSize)
	g.night = 0
	g.state = statePlaying
}

const (
	lightSpriteSize   = 128
	ambientNight      = 0.06 // brightness of unlit areas at full night
	rootLightRadius   = 34.0
	coreLightRadius   = 48.0
	cursorLightRadius = 95.0
)

// multiplyBlend composites the lightmap over the field: result = light * scene.
var multiplyBlend = ebiten.Blend{
	BlendFactorSourceRGB:        ebiten.BlendFactorDestinationColor,
	BlendFactorSourceAlpha:      ebiten.BlendFactorDestinationAlpha,
	BlendFactorDestinationRGB:   ebiten.BlendFactorZero,
	BlendFactorDestinationAlpha: ebiten.BlendFactorZero,
	BlendOperationRGB:           ebiten.BlendOperationAdd,
	BlendOperationAlpha:         ebiten.BlendOperationAdd,
}

func makeLightSprite(size int) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	c := float64(size) / 2
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			d := math.Hypot(float64(x)+0.5-c, float64(y)+0.5-c) / c
			a := 1 - d
			if a < 0 {
				a = 0
			}
			a *= a // softer falloff
			v := uint8(255 * a)
			img.SetRGBA(x, y, color.RGBA{v, v, v, v})
		}
	}
	return ebiten.NewImageFromImage(img)
}

// applyLight darkens the field at night, carving out light around the cursor,
// the Cores, and the living root network.
func (g *Game) applyLight() {
	if g.night < 0.02 {
		return
	}
	if g.lightSprite == nil {
		g.lightSprite = makeLightSprite(lightSpriteSize)
	}
	if g.lightmap == nil {
		g.lightmap = ebiten.NewImage(fieldPx, fieldPx)
	}
	amb := uint8(255 * (1 - g.night*(1-ambientNight)))
	g.lightmap.Fill(color.RGBA{amb, amb, amb, 255})

	add := func(fx, fy, radius float64) {
		op := &ebiten.DrawImageOptions{Blend: ebiten.BlendLighter, Filter: ebiten.FilterLinear}
		s := radius * 2 / float64(lightSpriteSize)
		op.GeoM.Scale(s, s)
		op.GeoM.Translate(fx-radius, fy-radius)
		g.lightmap.DrawImage(g.lightSprite, op)
	}

	half := float64(world.CellSize) / 2
	g.grid.VisitLit(func(col, row int, strong bool) {
		r := rootLightRadius
		if strong {
			r = coreLightRadius
		}
		add(float64(col*world.CellSize)+half, float64(row*world.CellSize)+half, r)
	})

	cx, cy := ebiten.CursorPosition()
	scale, ox, oy := g.fieldMetrics()
	add((float64(cx)-ox)/scale, (float64(cy)-oy)/scale, cursorLightRadius)

	g.field.DrawImage(g.lightmap, &ebiten.DrawImageOptions{Blend: multiplyBlend})
}

const (
	bloomDiv      = 4    // field is downscaled by this for the blur pass
	bloomStrength = 0.35 // how strongly the glow is added back
)

// applyBloom adds a soft neon glow: downscale the field (a cheap blur), then
// add it back on top with additive blending so bright elements bloom.
func (g *Game) applyBloom() {
	if g.bloom == nil {
		g.bloom = ebiten.NewImage(fieldPx/bloomDiv, fieldPx/bloomDiv)
	}
	g.bloom.Clear()
	down := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	down.GeoM.Scale(1.0/bloomDiv, 1.0/bloomDiv)
	g.bloom.DrawImage(g.field, down)

	up := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear, Blend: ebiten.BlendLighter}
	up.GeoM.Scale(bloomDiv, bloomDiv)
	up.ColorScale.ScaleAlpha(bloomStrength)
	g.field.DrawImage(g.bloom, up)
}

// fieldMetrics places the square field: it fills the screen height and is
// centered horizontally, leaving a HUD panel on each side.
func (g *Game) fieldMetrics() (scale, ox, oy float64) {
	scale = float64(g.screenH) / fieldPx
	draw := float64(g.screenH)
	ox = (float64(g.screenW) - draw) / 2
	return scale, ox, 0
}

var (
	fxRoot     = color.RGBA{0x2f, 0xa8, 0x6b, 0xff}
	fxWater    = color.RGBA{0x3c, 0x9a, 0xff, 0xff}
	fxMushroom = color.RGBA{0x35, 0xe0, 0xc4, 0xff}
	fxCore     = color.RGBA{0x5a, 0xff, 0xc4, 0xff}
	fxHit      = color.RGBA{0xff, 0xff, 0xff, 0xff}
	fxDestroy  = color.RGBA{0xff, 0x6e, 0x4c, 0xff}
	fxCoreDead = color.RGBA{0xff, 0x4c, 0x5e, 0xff}
	fxRotTone  = color.RGBA{0x6b, 0x4a, 0x2f, 0xff}
)

// applyFx turns the grid's reported events into visual effects.
func (g *Game) applyFx() {
	for _, e := range g.grid.DrainFx() {
		switch e.Kind {
		case world.FxPlaceRoot:
			g.fx.Pop(e.Col, e.Row, fxRoot)
			audio.PlayGrow()
		case world.FxPlaceWater:
			g.fx.Pop(e.Col, e.Row, fxWater)
		case world.FxPlaceMushroom:
			g.fx.Pop(e.Col, e.Row, fxMushroom)
		case world.FxPlaceCore:
			g.fx.Pop(e.Col, e.Row, fxCore)
			g.fx.Pop(e.Col+1, e.Row+1, fxCore)
		case world.FxCutRot:
			g.fx.Pop(e.Col, e.Row, fxRotTone)
			audio.PlayRot()
		case world.FxEatRot:
			audio.PlayEat()
		case world.FxHitCore:
			g.fx.Flash(e.Col, e.Row, fxHit)
		case world.FxDestroyRoot:
			g.fx.Burst(e.Col, e.Row, fxDestroy, 12)
		case world.FxDestroyCore:
			g.fx.Burst(e.Col, e.Row, fxCoreDead, 24)
			g.fx.Burst(e.Col+1, e.Row+1, fxCoreDead, 24)
		}
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
	scale, ox, oy := g.fieldMetrics()
	fx := (float64(cx) - ox) / scale
	fy := (float64(cy) - oy) / scale
	if fx < 0 || fy < 0 || fx >= fieldPx || fy >= fieldPx {
		return 0, 0, false
	}
	return int(fx) / world.CellSize, int(fy) / world.CellSize, true
}

// twoButtons lays out two stacked, centered buttons with the given labels.
func (g *Game) twoButtons(a, b string) (ui.Button, ui.Button) {
	w, h := 260, 56
	cx, cy := g.screenW/2, g.screenH/2
	top := ui.Button{X: cx - w/2, Y: cy - 8, W: w, H: h, Label: a}
	bot := ui.Button{X: cx - w/2, Y: cy + 64, W: w, H: h, Label: b}
	return top, bot
}

func clicked() bool { return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) }

func (g *Game) Update() error {
	switch g.state {
	case stateMenu:
		start, exit := g.twoButtons("START", "EXIT")
		if clicked() {
			mx, my := ebiten.CursorPosition()
			switch {
			case start.Contains(mx, my):
				audio.PlayClick()
				g.startRun()
			case exit.Contains(mx, my):
				audio.PlayClick()
				return ebiten.Termination
			}
		}
		return nil

	case statePaused:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = statePlaying
			return nil
		}
		cont, menu := g.twoButtons("CONTINUE", "MENU")
		if clicked() {
			mx, my := ebiten.CursorPosition()
			switch {
			case cont.Contains(mx, my):
				audio.PlayClick()
				g.state = statePlaying
			case menu.Contains(mx, my):
				audio.PlayClick()
				g.state = stateMenu
			}
		}
		return nil

	case stateOver:
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.startRun()
		}
		return nil
	}

	// statePlaying
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.state = statePaused
		return nil
	}
	g.updatePlaying()
	return nil
}

func (g *Game) updatePlaying() {
	// Day/night alternates per wave: odd waves (and their prep) are day, even
	// waves are night. era is the wave the current phase belongs to.
	era := g.waves.Number()
	if g.waves.InPrep() {
		era++
	}
	target := 0.0
	if era%2 == 0 {
		target = 1.0
	}
	g.night += (target - g.night) * 2.5 * secondsPerTick

	g.grid.Update(secondsPerTick, g.waves.Number())
	g.res.SetCoreLinks(g.grid.CoreMerges())
	g.res.AddWater(g.grid.MineWater(secondsPerTick, g.plants.NearThorn))
	energyMult, waterMult := g.plants.Modifiers()
	upkeep := float64(g.grid.CoreCount()-1) * resources.WaterUpkeepPerCore
	g.res.Update(secondsPerTick, energyMult, waterMult, upkeep)

	for n := g.waves.Update(secondsPerTick, g.bugs.Count()); n > 0; n-- {
		g.bugs.Spawn(g.grid, g.waves.Number(), g.waves.Side())
	}
	g.bugs.Update(secondsPerTick, g.grid)

	g.applyFx()
	g.fx.Update(secondsPerTick)

	if g.grid.CoreCount() == 0 {
		g.state = stateOver
		if w := g.waves.Number(); w > g.high {
			g.high = w
			save.Save(w)
		}
		return
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if col, row, ok := g.fieldCell(); ok && g.grid.CanGrow(col, row) && g.res.TrySpendEnergy(resources.RootEnergyCost) {
			g.grid.Grow(col, row)
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if col, row, ok := g.fieldCell(); ok && g.grid.Cut(col, row) {
			g.res.AddEnergy(resources.RefundFor(g.waves.Number()))
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
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.state == stateMenu {
		g.drawMenu(screen)
		return
	}
	g.drawPlaying(screen)
	switch g.state {
	case statePaused:
		g.drawPause(screen)
	case stateOver:
		ui.DrawGameOver(screen, g.screenW, g.screenH, g.waves.Number(), g.high)
	}
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	world.DrawDirtBackground(screen)
	ui.DrawTitle(screen, g.screenW, g.screenH/2-150)
	start, exit := g.twoButtons("START", "EXIT")
	mx, my := ebiten.CursorPosition()
	start.Draw(screen, start.Contains(mx, my))
	exit.Draw(screen, exit.Contains(mx, my))
	ui.DrawCenteredLabel(screen, g.screenW, g.screenH/2+150, fmt.Sprintf("BEST WAVE  %d", g.high), 20, ui.LabelColor)
	ui.DrawVersion(screen, g.screenW, g.screenH)
}

func (g *Game) drawPause(screen *ebiten.Image) {
	ui.DrawVeil(screen, g.screenW, g.screenH)
	ui.DrawCenteredLabel(screen, g.screenW, g.screenH/2-150, "PAUSED", ui.TitleSize, ui.LabelColor)
	cont, menu := g.twoButtons("CONTINUE", "MENU")
	mx, my := ebiten.CursorPosition()
	cont.Draw(screen, cont.Contains(mx, my))
	menu.Draw(screen, menu.Contains(mx, my))
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	if g.field == nil {
		g.field = ebiten.NewImage(fieldPx, fieldPx)
	}
	screen.Fill(panelColor)

	g.field.Clear()
	g.grid.Draw(g.field)
	g.plants.Draw(g.field, world.CellSize)
	if g.state == statePlaying {
		if col, row, ok := g.fieldCell(); ok {
			g.grid.DrawHover(g.field, col, row)
		}
	}
	g.bugs.Draw(g.field, world.CellSize)
	g.fx.Draw(g.field)
	if g.state == statePlaying && g.waves.InPrep() {
		ui.DrawTelegraph(g.field, g.waves.Side(), fieldPx)
	}
	g.applyBloom()
	g.applyLight()

	scale, ox, oy := g.fieldMetrics()
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(ox, oy)
	screen.DrawImage(g.field, op)

	leftX := 16
	rightEdge := g.screenW - 24
	ui.DrawResources(screen, g.res, leftX, 40)
	ui.DrawStatus(screen, g.grid.CoreCount(), g.grid.CoreHP(), leftX, 220)
	ui.DrawControls(screen, leftX, g.screenH-140)
	shownWave := g.waves.Number()
	if g.waves.InPrep() {
		shownWave++
	}
	ui.DrawWaveInfo(screen, rightEdge, 40, g.waves.Number(), g.waves.InPrep(), g.waves.PrepRemaining(), season.Of(shownWave).Name())
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.screenW, g.screenH = outsideWidth, outsideHeight
	return outsideWidth, outsideHeight
}
