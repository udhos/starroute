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

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Star Route")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := newGame()

	g.pause = pause

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
