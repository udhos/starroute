// Package main implements the game.
package main

import (
	"bytes"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	pi2       = 2 * math.Pi
	maxAngle  = float64(100) // custom number of angles in the circle
	oneEighth = maxAngle / 8
)

type sprite struct {
	x, y          float64
	width, height int
	angle         float64
	angleNative   float64 // undo this intrinsic rotate of image to point image to zero angle (right)
	image         *ebiten.Image
}

func (s *sprite) update() {
	// move sprint etc
	s.angle = math.Mod(s.angle+1, maxAngle)
	//s.angle = oneEighth
}

// game implements ebiten.Game interface.
type game struct {
	op      ebiten.DrawImageOptions
	sprites []*sprite
	//ebitenImage *ebiten.Image
	tiles *tiles

	// See comment in game.Layout method.
	screenWidth  int
	screenHeight int
}

func (g *game) addSprite(x, y, angleNative float64, spriteImage *ebiten.Image) {
	w, h := spriteImage.Bounds().Dx(), spriteImage.Bounds().Dy()
	spr := sprite{
		x:           x,
		y:           y,
		width:       w,
		height:      h,
		angleNative: angleNative,
		image:       spriteImage,
	}
	g.sprites = append(g.sprites, &spr)
}

func newGame() *game {

	//
	// Load an image from the embedded image data.
	//

	const scaleAlpha = 0.8

	ebitenImage := createImage(bytes.NewReader(images.Ebiten_png), scaleAlpha)

	const (
		// FIXME: these should come from tilemap data
		tileSize             = 16
		tileLayerScreenWidth = 240
	)

	g := &game{
		tiles: newTiles(bytes.NewReader(images.Tiles_png), tileSize, layers, tileLayerScreenWidth),

		// See comment in game.Layout method.
		screenWidth:  320,
		screenHeight: 240,
	}

	// See comment in game.Layout method.
	log.Printf("Game screen size: %dx%d", g.screenWidth, g.screenHeight)

	//
	// Add sprites.
	//

	g.addSprite(50, 50, 0, ebitenImage)
	g.addSprite(100, 100, oneEighth, ebitenImage)

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

	g.tiles.draw(screen)

	// Draw each sprite.
	// DrawImage can be called many many times, but in the implementation,
	// the actual draw call to GPU is very few since these calls satisfy
	// some conditions e.g. all the rendering sources and targets are same.
	// For more detail, see:
	// https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage

	for i := 0; i < len(g.sprites); i++ {
		s := g.sprites[i]
		g.op.GeoM.Reset()

		w, h := s.image.Bounds().Dx(), s.image.Bounds().Dy()

		centerX := float64(w) / 2
		centerY := float64(h) / 2

		// move the rotation center to the origin
		g.op.GeoM.Translate(-centerX, -centerY)

		// rotate around the origin

		// undo intrinsic image rotation
		angleNativeRad := pi2 * float64(s.angleNative) / maxAngle
		g.op.GeoM.Rotate(-angleNativeRad)

		// apply actual rotation
		angleRad := pi2 * float64(s.angle) / maxAngle
		g.op.GeoM.Rotate(angleRad)

		// undo the translation used to move the rotation center
		g.op.GeoM.Translate(centerX, centerY)

		// apply the actual object's position
		g.op.GeoM.Translate(float64(s.x), float64(s.y))

		screen.DrawImage(s.image, &g.op)

		//
		// Red show how much the sprite was rotated back (counter clockwise) to
		// make its front point to the right (zero angle).
		//
		colorRed := color.RGBA{0xff, 0, 0, 0xff}
		drawDebugArrow(screen, float64(s.x+centerX), float64(s.y+centerY),
			angleRad-angleNativeRad, 20, 3, colorRed)

		//
		// Yellow show the sprint front direction and should point to
		// right (zero angle).
		//
		colorYellow := color.RGBA{0xff, 0xff, 0, 0xff}
		drawDebugArrow(screen, float64(s.x+centerX), float64(s.y+centerY),
			angleRad, 30, 1, colorYellow)
	}

}

func drawDebugArrow(screen *ebiten.Image, x, y, angle, lenght, width float64, arrowColor color.RGBA) {
	arrowX := x + lenght*math.Cos(angle)
	arrowY := y + lenght*math.Sin(angle)

	var path vector.Path
	path.MoveTo(float32(x), float32(y))
	path.LineTo(float32(arrowX), float32(arrowY))

	strokeOp := &vector.StrokeOptions{}
	strokeOp.Width = float32(width)

	drawOp := &vector.DrawPathOptions{}

	drawOp.ColorScale.ScaleWithColor(arrowColor)

	drawOp.AntiAlias = false
	vector.StrokePath(screen, &path, strokeOp, drawOp)
}

// Layout accepts an outside size, which is a window size on desktop, and
// returns the game's logical screen size. This code ignores the arguments
// and returns the fixed values. This means that the game screen size is
// always same, whatever the window's size is. Layout will be more meaningful
// e.g., when the window is resizable.
func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenWidth, g.screenHeight
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
