package main

import (
	"bytes"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	pi2        = 2 * math.Pi
	maxAngle   = float64(100) // custom number of angles in the circle
	oneQuarter = maxAngle / 4
	oneEighth  = maxAngle / 8
)

// game implements ebiten.Game interface.
type game struct {
	pause bool

	scenes       []scene
	sceneCurrent int

	/*
		op      ebiten.DrawImageOptions
		sprites []*sprite
		tiles   *tiles
	*/

	// See comment in game.Layout method.
	screenWidth  int
	screenHeight int

	//keys []ebiten.Key
}

/*
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
*/

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
		tileLayerXCount      = 15
	)

	ts := newTiles(bytes.NewReader(images.Tiles_png), tileSize, layers, tileLayerXCount)

	scene1 := scene{tiles: ts}
	scene1.addSprite(50, 50, 0, ebitenImage)
	scene1.addSprite(100, 100, oneEighth, ebitenImage)

	scene2 := scene{tiles: ts}
	scene2.addSprite(150, 150, 0, ebitenImage)
	scene2.addSprite(200, 200, oneQuarter, ebitenImage)

	g := &game{
		//tiles: newTiles(bytes.NewReader(images.Tiles_png), tileSize, layers, tileLayerXCount),

		// See comment in game.Layout method.
		screenWidth:  320,
		screenHeight: 320,

		scenes:       []scene{scene1, scene2},
		sceneCurrent: 0,
	}

	// See comment in game.Layout method.
	log.Printf("Game screen size: %dx%d", g.screenWidth, g.screenHeight)

	//
	// Add sprites.
	//

	//g.addSprite(50, 50, 0, ebitenImage)
	//g.addSprite(100, 100, oneEighth, ebitenImage)

	return g
}

// Update is called every tick. Tick is a time unit for logical updating.
// The default value is 1/60 [s], then Update is called 60 times per second by
// default (i.e. an Ebitengine game works in 60 ticks-per-second).
func (g *game) Update() error {

	/*
		keys := inpututil.AppendPressedKeys(nil)
		if len(keys) > 0 {
			p := keys[len(keys)-1]
			switch p {
			case ebiten.KeyP:
				g.pause = !g.pause
				log.Printf("Pause: %t", g.pause)
			case ebiten.KeyBackspace:
				g.sceneCurrent = (g.sceneCurrent + 1) % len(g.scenes)
				log.Printf("Switching to scene %d of %d",
					g.sceneCurrent+1, len(g.scenes))
			}
		}
	*/

	if inpututil.IsKeyJustReleased(ebiten.KeyP) {
		g.pause = !g.pause
		log.Printf("Pause: %t", g.pause)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyBackspace) {
		g.sceneCurrent = (g.sceneCurrent + 1) % len(g.scenes)
		log.Printf("Switching to scene %d of %d",
			g.sceneCurrent+1, len(g.scenes))
	}

	if g.pause {
		return nil
	}

	g.scenes[g.sceneCurrent].update()

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

	backgroundColor := color.RGBA{R: 128, G: 128, B: 128, A: 255}
	screen.Fill(backgroundColor)

	g.scenes[g.sceneCurrent].draw(screen)
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
