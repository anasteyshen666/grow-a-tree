// Command game is the entry point for "Grow a Tree".
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"growtree/internal/game"
)

func main() {
	g := game.New()

	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Grow a Tree")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
