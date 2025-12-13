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
	var resize string
	flag.BoolVar(&pause, "pause", false, "pause game update")
	flag.StringVar(&resize, "resize", "on", "window resize mode: on|off|full")
	flag.Parse()

	var resizeMode ebiten.WindowResizingModeType
	switch resize {
	case "on":
		resizeMode = ebiten.WindowResizingModeEnabled
	case "off":
		resizeMode = ebiten.WindowResizingModeDisabled
	case "full":
		resizeMode = ebiten.WindowResizingModeOnlyFullscreenEnabled
	default:
		log.Fatalf("invalid window resize mode: %s", resize)
	}

	log.Printf("Window size: %dx%d", windowWidth, windowHeight)

	isFull := resize == "full"
	ebiten.SetFullscreen(isFull)
	ebiten.SetWindowResizingMode(resizeMode)
	ebiten.SetWindowDecorated(!isFull)
	if !isFull {
		ebiten.SetWindowTitle("Star Route")
		ebiten.SetWindowSize(windowWidth, windowHeight)
	}

	g := newGame(isFull)

	g.pause = pause

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
