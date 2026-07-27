package plants

import (
	"bytes"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"

	"growtree/assets"
)

var (
	petSprites   [kindCount]*ebiten.Image
	spritesReady bool
)

func ensureSprites() {
	if spritesReady {
		return
	}
	petSprites[Battery] = loadSprite("pet1.png")
	petSprites[Moss] = loadSprite("pet2.png")
	petSprites[Thorn] = loadSprite("pet3.png")
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
