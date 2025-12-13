// Package main implements the game.
package main

import (
	"flag"
	"log"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
// screen and window dimensions should be proportional,
// then the screen would fit nicely within the window
// without a padding border.

//defaultScreenWidth  = 400
//defaultScreenHeight = 300

// windowWidth  = 2 * defaultScreenWidth
// windowHeight = 2 * defaultScreenHeight
)

func main() {

	var pause bool
	var resize string
	var screen string
	var window string
	var both string
	flag.BoolVar(&pause, "pause", false, "pause game update")
	flag.StringVar(&resize, "resize", "on", "window resize mode: on|off|full")
	flag.StringVar(&screen, "screen", "400x300", "screen size")
	flag.StringVar(&window, "window", "800x600", "window size")
	flag.StringVar(&both, "both", "", "screen and window size")
	flag.Parse()

	var screenWidth, screenHeight, windowWidth, windowHeight int

	if both == "" {
		screenWidth, screenHeight = parseDim("screen", screen)
		windowWidth, windowHeight = parseDim("window", window)
	} else {
		screenWidth, screenHeight = parseDim("screen", both)
		windowWidth, windowHeight = parseDim("window", both)
	}

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
	ebiten.SetWindowTitle("Star Route")
	ebiten.SetWindowDecorated(!isFull)
	ebiten.SetFullscreen(isFull)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowResizingMode(resizeMode)

	g := newGame(isFull, screenWidth, screenHeight)

	g.pause = pause

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func parseDim(label, dim string) (int, int) {
	s := strings.Split(dim, "x")
	if len(s) != 2 {
		log.Fatalf("bad dimensions: %s: %s", label, dim)
	}
	x := s[0]
	y := s[1]
	valX, _ := strconv.Atoi(x)
	valY, _ := strconv.Atoi(y)
	log.Printf("%s: %dx%d", label, valX, valY)
	return valX, valY
}
