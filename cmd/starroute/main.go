// Package main implements the game.
package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	fmt.Println("hint for fullscreen: starroute -resize=full -both 1920x1080")

	var pause bool
	var resize string
	var screen string
	var window string
	var both string
	flag.BoolVar(&pause, "pause", false, "pause game update")
	flag.StringVar(&resize, "resize", "on", "window resize mode: on|off|fullscreen")
	flag.StringVar(&screen, "screen", "800x600", "game logical screen size (should be <= window size)")
	flag.StringVar(&window, "window", "800x600", "outsize window size (should be multiple of screen size)")
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

	isFull := strings.HasPrefix("fullscreen", resize)

	var resizeMode ebiten.WindowResizingModeType
	switch {
	case resize == "on":
		resizeMode = ebiten.WindowResizingModeEnabled
	case resize == "off":
		resizeMode = ebiten.WindowResizingModeDisabled
	case isFull:
		resizeMode = ebiten.WindowResizingModeOnlyFullscreenEnabled
	default:
		log.Fatalf("invalid window resize mode: %s", resize)
	}

	log.Printf("Window size: %dx%d", windowWidth, windowHeight)

	ebiten.SetWindowTitle("Star Route")
	ebiten.SetWindowDecorated(!isFull)
	ebiten.SetFullscreen(isFull)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowResizingMode(resizeMode)

	g := newGame(screenWidth, screenHeight)

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
	valX, errX := strconv.Atoi(x)
	if errX != nil {
		log.Fatalf("bad dimension X: %s: %s: %v", label, dim, errX)
	}
	valY, errY := strconv.Atoi(y)
	if errY != nil {
		log.Fatalf("bad dimension Y: %s: %s: %v", label, dim, errY)
	}
	log.Printf("%s: %dx%d", label, valX, valY)
	return valX, valY
}
