package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/udhos/starroute/music"
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
	debug bool

	scenes       []scene
	sceneCurrent int

	defaultScreenWidth  int
	defaultScreenHeight int

	// See comment in game.Layout method.
	screenWidth  int
	screenHeight int

	// used only to debug window size
	windowWidth  int
	windowHeight int

	mouseX, mouseY int
}

func newGame(defaultScreenWidth, defaultScreenHeight int) *game {

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

	audioContext := audio.NewContext(music.SampleRate)

	ts := newTiles(bytes.NewReader(images.Tiles_png), tileSize, sampleLayers, tileLayerXCount)

	scene1 := newScene(ts, sceneTrack1, audioContext)
	scene1.addSprite(50, 50, 0, ebitenImage)
	scene1.addSprite(100, 100, oneEighth, ebitenImage)

	scene2 := newScene(ts, sceneTrack2, audioContext)
	scene2.addSprite(150, 150, 0, ebitenImage)
	scene2.addSprite(200, 200, oneQuarter, ebitenImage)

	const tileEdgeCount = 200 // 200x200=40000
	layers := [][]int{generateLayer(tileEdgeCount)}
	ts3 := newTiles(bytes.NewReader(images.Tiles_png), tileSize, layers, tileEdgeCount)

	scene3 := newScene(ts3, sceneTrack3, audioContext)
	scene3.addSprite(100, 50, 0, ebitenImage)

	g := &game{
		debug: true,

		defaultScreenWidth:  defaultScreenWidth,
		defaultScreenHeight: defaultScreenHeight,

		// See comment in game.Layout method.
		screenWidth:  defaultScreenWidth,
		screenHeight: defaultScreenHeight,

		scenes:       []scene{scene1, scene2, scene3},
		sceneCurrent: 0,
	}

	g.scenes[g.sceneCurrent].musicStart()

	// See comment in game.Layout method.
	log.Printf("Game screen size: %dx%d", g.screenWidth, g.screenHeight)

	return g
}

const camPanStep = 5

// Update is called every tick. Tick is a time unit for logical updating.
// The default value is 1/60 [s], then Update is called 60 times per second by
// default (i.e. an Ebitengine game works in 60 ticks-per-second).
func (g *game) Update() error {

	//
	// handle burst of keys
	//
	keys := inpututil.AppendPressedKeys(nil)
	for _, p := range keys {
		//p := keys[len(keys)-1]

		switch p {
		case ebiten.KeyUp:
			g.scenes[g.sceneCurrent].cam.y = max(g.scenes[g.sceneCurrent].cam.y-camPanStep, 0)
			continue
		case ebiten.KeyDown:
			g.scenes[g.sceneCurrent].cam.y += camPanStep
			continue
		case ebiten.KeyLeft:
			g.scenes[g.sceneCurrent].cam.x = max(g.scenes[g.sceneCurrent].cam.x-camPanStep, 0)
			continue
		case ebiten.KeyRight:
			g.scenes[g.sceneCurrent].cam.x += camPanStep
			continue
		case ebiten.KeyEscape:
			log.Printf("ESC pressed, exiting")
			os.Exit(0)
		}

		zero := p == ebiten.Key0
		plus := p == ebiten.KeyEqual
		minus := p == ebiten.KeyMinus
		if zero {
			g.screenWidth = g.defaultScreenWidth
			g.screenHeight = g.defaultScreenHeight
		}
		if plus {
			if g.screenWidth < 10000 {
				g.screenWidth += g.screenWidth / 10
			}
			if g.screenHeight < 10000 {
				g.screenHeight += g.screenHeight / 10
			}
		}
		if minus {
			if g.screenWidth > 10 {
				g.screenWidth -= g.screenWidth / 10
			}
			if g.screenHeight > 10 {
				g.screenHeight -= g.screenHeight / 10
			}
		}
		if zero || plus || minus {
			log.Printf("Game screen size: %dx%d", g.screenWidth, g.screenHeight)
		}
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyO) {
		g.debug = !g.debug
		log.Printf("Debug: %t", g.debug)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyP) {
		g.pause = !g.pause
		log.Printf("Pause: %t", g.pause)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyBackspace) {
		g.switchScene()
	}
	/*
		if inpututil.IsKeyJustReleased(ebiten.KeyShiftRight) {
			g.screenTrackWindow = !g.screenTrackWindow

			g.screenWidth = defaultScreenWidth
			g.screenHeight = defaultScreenHeight

			log.Printf("Screen track window: %t", g.screenTrackWindow)
		}
	*/

	g.mouseX, g.mouseY = ebiten.CursorPosition()

	if g.pause {
		return nil
	}

	g.scenes[g.sceneCurrent].update()

	return nil
}

func (g *game) switchScene() {
	g.scenes[g.sceneCurrent].musicStop()

	g.sceneCurrent = (g.sceneCurrent + 1) % len(g.scenes)
	log.Printf("Switching to scene %d of %d",
		g.sceneCurrent+1, len(g.scenes))

	g.scenes[g.sceneCurrent].musicStart()
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

	sc := g.scenes[g.sceneCurrent]

	drawnTiles := sc.draw(screen, g.debug)

	if g.debug {
		tileDimX, tileDimY := sc.tiles.dimensions()
		cam := sc.cam
		camLastX := cam.x + g.screenWidth - 1
		camLastY := cam.y + g.screenHeight - 1
		ebitenutil.DebugPrint(screen,
			fmt.Sprintf("TPS:%0.1f FPS:%0.1f tilemap:%dx%d cam:%dx%d-%dx%d mouse:%dx%d screen.Bounds:%dx%d win:%dx%d drawnTiles:%d",
				ebiten.ActualTPS(), ebiten.ActualFPS(),
				tileDimX, tileDimY,
				cam.x, cam.y,
				camLastX, camLastY,
				g.mouseX, g.mouseY,
				screen.Bounds().Dx(),
				screen.Bounds().Dy(),
				g.windowWidth,
				g.windowHeight,
				drawnTiles))

		colorBlue := color.RGBA{0, 0, 0xff, 0xff}
		drawDebugRect(screen, 1, 1, float64(g.screenWidth), float64(g.screenHeight), colorBlue)
	}
}

func drawDebugRect(screen *ebiten.Image, x1, y1, x2, y2 float64, borderColor color.RGBA) {

	const width = 1

	var path vector.Path
	path.MoveTo(float32(x1), float32(y1))
	path.LineTo(float32(x1), float32(y2))
	path.LineTo(float32(x2), float32(y2))
	path.LineTo(float32(x2), float32(y1))
	path.LineTo(float32(x1), float32(y1))

	strokeOp := &vector.StrokeOptions{}
	strokeOp.Width = float32(width)

	drawOp := &vector.DrawPathOptions{}

	drawOp.ColorScale.ScaleWithColor(borderColor)

	drawOp.AntiAlias = false
	vector.StrokePath(screen, &path, strokeOp, drawOp)
}

// Layout takes the outside size (e.g., the window size) and returns the
// (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just
// return a fixed size.
func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// used only to debug window size
	g.windowWidth = outsideWidth
	g.windowHeight = outsideHeight

	// This code ignores the arguments and returns the fixed values.
	// This means that the game screen size is always same,
	// whatever the window's size is.
	return g.screenWidth, g.screenHeight
}
