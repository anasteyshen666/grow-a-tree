package enemies

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"

	"growtree/assets"
)

var (
	bugSprites   [maxLevel + 1]*ebiten.Image
	spritesReady bool
)

func ensureSprites() {
	if spritesReady {
		return
	}
	for lvl := 1; lvl <= maxLevel; lvl++ {
		bugSprites[lvl] = loadSprite(fmt.Sprintf("fly%d.png", lvl))
	}
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
