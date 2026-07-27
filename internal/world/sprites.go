package world

import (
	"bytes"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"

	"growtree/assets"
)

// dirtDim darkens the dirt terrain so it reads as a dim background, not a bright
// tiled floor. deadRootDim darkens a root cut off from the Core.
const (
	dirtDim     = 0.4
	deadRootDim = 0.4
)

var (
	sprDirt     *ebiten.Image
	sprCore     *ebiten.Image
	sprRoot     *ebiten.Image
	sprMushroom *ebiten.Image
	sprRot      *ebiten.Image
	sprWater    *ebiten.Image
	sprSpore    *ebiten.Image

	spritesReady bool
	terrainBG    *ebiten.Image // cached dirt background, rebuilt when the tint changes
	terrainTint  [3]float64
)

func ensureSprites() {
	if spritesReady {
		return
	}
	sprDirt = loadSprite("dirt.png")
	sprCore = loadSprite("core.png")
	sprRoot = loadSprite("common koren.png")
	sprMushroom = loadSprite("koren grib.png")
	sprRot = loadSprite("gnil.png")
	sprWater = loadSprite("water.png")
	sprSpore = loadSprite("spora.png")
	spritesReady = true
}

func loadSprite(name string) *ebiten.Image {
	data, err := assets.FS.ReadFile(name)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}

// tileOp scales a sprite to one cell and positions it at (col,row).
func tileOp(img *ebiten.Image, col, row int) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	b := img.Bounds()
	op.GeoM.Scale(float64(CellSize)/float64(b.Dx()), float64(CellSize)/float64(b.Dy()))
	op.GeoM.Translate(float64(col*CellSize), float64(row*CellSize))
	return op
}

func drawTile(screen, img *ebiten.Image, col, row int) {
	screen.DrawImage(img, tileOp(img, col, row))
}

// drawTileAlpha draws a tile faded to alpha a (used for aging rot and depleting
// water, so the dirt shows through as they fade).
func drawTileAlpha(screen, img *ebiten.Image, col, row int, a float32) {
	op := tileOp(img, col, row)
	op.ColorScale.ScaleAlpha(a)
	screen.DrawImage(img, op)
}

// drawTileDim draws a tile with its RGB scaled by f, darkening it (used for a
// dead root, which is the normal root texture dimmed).
func drawTileDim(screen, img *ebiten.Image, col, row int, f float32) {
	op := tileOp(img, col, row)
	op.ColorScale.Scale(f, f, f, 1)
	screen.DrawImage(img, op)
}

var menuBG *ebiten.Image

// DrawDirtBackground tiles the darkened dirt across dst (used by menus). Cached
// per destination size.
func DrawDirtBackground(dst *ebiten.Image) {
	ensureSprites()
	w, h := dst.Bounds().Dx(), dst.Bounds().Dy()
	if menuBG == nil || menuBG.Bounds().Dx() != w || menuBG.Bounds().Dy() != h {
		menuBG = ebiten.NewImage(w, h)
		for y := 0; y < h; y += CellSize {
			for x := 0; x < w; x += CellSize {
				op := tileOp(sprDirt, x/CellSize, y/CellSize)
				op.ColorScale.Scale(dirtDim, dirtDim, dirtDim, 1)
				menuBG.DrawImage(sprDirt, op)
			}
		}
	}
	dst.DrawImage(menuBG, nil)
}

// drawTerrain blits the darkened, season-tinted dirt background covering the
// field. The cache is rebuilt only when the tint changes.
func drawTerrain(screen *ebiten.Image, tr, tg, tb float64) {
	if terrainBG == nil || terrainTint != [3]float64{tr, tg, tb} {
		if terrainBG == nil {
			terrainBG = ebiten.NewImage(Cols*CellSize, Rows*CellSize)
		}
		terrainBG.Clear()
		for row := 0; row < Rows; row++ {
			for col := 0; col < Cols; col++ {
				op := tileOp(sprDirt, col, row)
				op.ColorScale.Scale(float32(dirtDim*tr), float32(dirtDim*tg), float32(dirtDim*tb), 1)
				terrainBG.DrawImage(sprDirt, op)
			}
		}
		terrainTint = [3]float64{tr, tg, tb}
	}
	screen.DrawImage(terrainBG, nil)
}
