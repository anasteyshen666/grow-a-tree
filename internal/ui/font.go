package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	xfont "golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"growtree/assets"
)

var (
	fontSource *opentype.Font
	faceCache  = map[float64]text.Face{}
)

func loadFontSource() {
	if fontSource != nil {
		return
	}
	data, err := assets.FS.ReadFile("dogica.ttf")
	if err != nil {
		panic(err)
	}
	f, err := opentype.Parse(data)
	if err != nil {
		panic(err)
	}
	fontSource = f
}

// face returns a cached Dogica face at the given pixel size (best at multiples
// of the font's native 8px, e.g. 16 or 24).
func face(size float64) text.Face {
	if f, ok := faceCache[size]; ok {
		return f
	}
	loadFontSource()
	ff, err := opentype.NewFace(fontSource, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: xfont.HintingFull})
	if err != nil {
		panic(err)
	}
	f := text.NewGoXFace(ff)
	faceCache[size] = f
	return f
}

// drawText draws s with its top-left at (x,y) in the given size and color.
func drawText(dst *ebiten.Image, s string, x, y int, size float64, clr color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(dst, s, face(size), op)
}

// textWidth measures s at the given size.
func textWidth(s string, size float64) int {
	w, _ := text.Measure(s, face(size), 0)
	return int(w)
}
