// Package main implements the game.
package main

import (
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	var pause bool
	flag.BoolVar(&pause, "pause", false, "pause game update")
	flag.Parse()

	const (
		windowWidth  = 640
		windowHeight = 480
	)

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
