// Package main implements the game.
package main

import (
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// screen and window dimensions should be proportional,
	// then the screen would fit nicely within the window
	// without a padding border.

	defaultScreenWidth  = 400
	defaultScreenHeight = 300

	windowWidth  = 2 * defaultScreenWidth
	windowHeight = 2 * defaultScreenHeight
)

func main() {

	var pause bool
	flag.BoolVar(&pause, "pause", false, "pause game update")
	flag.Parse()

	log.Printf("Window size: %dx%d", windowWidth, windowHeight)

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Star Route")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := newGame()

	g.pause = pause

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
