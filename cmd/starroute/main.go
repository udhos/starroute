// Package main implements the game.
package main

import (
	"bytes"
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
)

type sprite struct {
	x, y          float64
	width, height int
	angle         float64
}

const maxAngle = 100 // custom number of angles in the circle

func (s *sprite) update() {
	// move sprint etc
	s.angle = math.Mod(s.angle+1, maxAngle)
}

// game implements ebiten.Game interface.
type game struct {
	op          ebiten.DrawImageOptions
	sprites     []*sprite
	ebitenImage *ebiten.Image
}

func newGame() *game {

	//
	// Load an image from the embedded image data.
	//

	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Ebiten_png))
	if err != nil {
		log.Fatalf("newGame: %v", err)
	}
	origEbitenImage := ebiten.NewImageFromImage(img)

	s := origEbitenImage.Bounds().Size()
	ebitenImage := ebiten.NewImage(s.X, s.Y)

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(0.5)
	ebitenImage.DrawImage(origEbitenImage, op)

	g := &game{ebitenImage: ebitenImage}

	//
	// Add one sprite to the game.
	//

	w, h := ebitenImage.Bounds().Dx(), ebitenImage.Bounds().Dy()
	x, y := 50, 50

	spr := sprite{
		x:      float64(x),
		y:      float64(y),
		width:  w,
		height: h,
	}
	g.sprites = append(g.sprites, &spr)

	return g
}

// Update is called every tick. Tick is a time unit for logical updating.
// The default value is 1/60 [s], then Update is called 60 times per second by
// default (i.e. an Ebitengine game works in 60 ticks-per-second).
func (g *game) Update() error {

	// Update all sprites.
	for _, spr := range g.sprites {
		spr.update()
	}

	return nil
}

// Draw is called every frame. Frame is a time unit for rendering and this
// depends on the display's refresh rate. If the monitor's refresh rate
// is 60 [Hz], Draw is called 60 times per second.
//
// Draw takes an argument screen, which is a pointer to an ebiten.Image.
// In Ebitengine, all images like images created from image files, offscreen
// images (temporary render target), and the screen are represented as
// ebiten.Image objects. screen argument is the final destination of
// rendering. The window shows the final state of screen every frame.
func (g *game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")

	// Draw each sprite.
	// DrawImage can be called many many times, but in the implementation,
	// the actual draw call to GPU is very few since these calls satisfy
	// some conditions e.g. all the rendering sources and targets are same.
	// For more detail, see:
	// https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
	w, h := g.ebitenImage.Bounds().Dx(), g.ebitenImage.Bounds().Dy()
	for i := 0; i < len(g.sprites); i++ {
		s := g.sprites[i]
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		g.op.GeoM.Rotate(2 * math.Pi * float64(s.angle) / maxAngle)
		g.op.GeoM.Translate(float64(w)/2, float64(h)/2)
		g.op.GeoM.Translate(float64(s.x), float64(s.y))
		screen.DrawImage(g.ebitenImage, &g.op)
	}

}

// Layout accepts an outside size, which is a window size on desktop, and
// returns the game's logical screen size. This code ignores the arguments
// and returns the fixed values. This means that the game screen size is
// always same, whatever the window's size is. Layout will be more meaningful
// e.g., when the window is resizable.
func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := newGame()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
