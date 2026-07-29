// Command cover renders the itch.io cover image: the game's dirt-tile field with
// the title, using the same embedded dirt texture and Dogica font.
package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"growtree/assets"
)

const (
	width    = 630
	height   = 500
	tileSize = 42
	dirtDim  = 0.42
)

var (
	green   = color.RGBA{0x6c, 0xe0, 0x4c, 0xff}
	outline = color.RGBA{0x00, 0x00, 0x00, 0xff}
)

func main() {
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	tileDirt(canvas)
	drawTitle(canvas, "GROW A TREE", 40)

	f, err := os.Create("dist/cover.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, canvas); err != nil {
		log.Fatal(err)
	}
	log.Println("wrote dist/cover.png")
}

func tileDirt(canvas *image.RGBA) {
	src := decode("dirt.png")

	// darken to a tile
	b := src.Bounds()
	dark := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, a := src.At(x, y).RGBA()
			dark.SetRGBA(x, y, color.RGBA{
				uint8(float64(r>>8) * dirtDim),
				uint8(float64(g>>8) * dirtDim),
				uint8(float64(bl>>8) * dirtDim),
				uint8(a >> 8),
			})
		}
	}

	// scale one tile up (nearest, keep it pixelated)
	tile := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	xdraw.NearestNeighbor.Scale(tile, tile.Bounds(), dark, dark.Bounds(), xdraw.Over, nil)

	for y := 0; y < height; y += tileSize {
		for x := 0; x < width; x += tileSize {
			draw.Draw(canvas, image.Rect(x, y, x+tileSize, y+tileSize), tile, image.Point{}, draw.Src)
		}
	}
}

func drawTitle(canvas *image.RGBA, s string, size float64) {
	data, err := assets.FS.ReadFile("dogicabold.ttf")
	if err != nil {
		log.Fatal(err)
	}
	tt, err := opentype.Parse(data)
	if err != nil {
		log.Fatal(err)
	}
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		log.Fatal(err)
	}

	d := &font.Drawer{Dst: canvas, Face: face}
	w := d.MeasureString(s).Round()
	m := face.Metrics()
	x := (width - w) / 2
	baseline := (height+(m.Ascent-m.Descent).Round())/2 - 20

	offsets := [][2]int{{-2, 0}, {2, 0}, {0, -2}, {0, 2}, {-2, -2}, {2, -2}, {-2, 2}, {2, 2}}
	d.Src = image.NewUniform(outline)
	for _, o := range offsets {
		d.Dot = fixed.P(x+o[0], baseline+o[1])
		d.DrawString(s)
	}
	d.Src = image.NewUniform(green)
	d.Dot = fixed.P(x, baseline)
	d.DrawString(s)
}

func decode(name string) image.Image {
	b, err := assets.FS.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	return img
}
