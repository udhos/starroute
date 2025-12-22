package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/udhos/starroute/music"
)

const (
	pi2        = 2 * math.Pi
	maxAngle   = float64(200) // custom number of angles in the circle
	oneQuarter = maxAngle / 4
	oneEighth  = maxAngle / 8
)

// game implements ebiten.Game interface.
type game struct {
	pause bool
	debug bool

	scenes           []*scene
	sceneStart       int
	sceneCurrent     int
	sceneResume      int
	sceneUpdateInput func(sc *scene)

	defaultScreenWidth  int
	defaultScreenHeight int

	// See comment in game.Layout method.
	screenWidth  int
	screenHeight int

	// used only to debug window size
	windowWidth  int
	windowHeight int

	mouseX, mouseY int

	//ui *ebitenui.UI
	//headerLbl     *widget.Text
	//coordinateLbl   *widget.Text
	mplusFaceSource *text.GoTextFaceSource
	//uiCoord         string

	debugui debugui.DebugUI
}

func newGame(defaultScreenWidth, defaultScreenHeight int) *game {

	//
	// Load an image from the embedded image data.
	//

	var ebitenImage *ebiten.Image
	var rotationScene1Sprite2 float64

	if false {
		const scaleAlpha = 0.8
		rotationScene1Sprite2 = oneEighth
		ebitenImage = createImage(bytes.NewReader(images.Ebiten_png), scaleAlpha)
	} else {
		const scaleAlpha = 1
		rotationScene1Sprite2 = -oneQuarter
		ebitenImage = createImage(bytes.NewReader(mustLoadAsset("newGame", "body_01.png")), scaleAlpha)
	}

	mplusFaceSource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	const (
		// FIXME: these should come from tilemap data
		tileSize             = 16
		tileLayerScreenWidth = 240
		tileLayerXCount      = 15
	)

	g := &game{
		debug: true,

		defaultScreenWidth:  defaultScreenWidth,
		defaultScreenHeight: defaultScreenHeight,

		// See comment in game.Layout method.
		screenWidth:  defaultScreenWidth,
		screenHeight: defaultScreenHeight,

		sceneStart:   0,
		sceneCurrent: 0,

		mplusFaceSource: mplusFaceSource,
		//uiCoord:         "? ?",
	}

	// This adds the root container to the UI, so that it will be rendered.
	//g.ui = g.getEbitenUI()

	audioContext := audio.NewContext(music.SampleRate)

	ts := newTiles(bytes.NewReader(images.Tiles_png), tileSize, sampleLayers, tileLayerXCount)

	const (
		cyclicCamera     = false
		centralizeCamera = false
		showCoord        = true
	)

	// scene0: start screen
	var scene0 *scene
	{
		const tileEdgeCount = 10 // 10x10=100
		layers := [][]int{generateLayer(tileEdgeCount)}
		ts := newTiles(bytes.NewReader(images.Tiles_png), tileSize, layers, tileEdgeCount)
		scene0 = newScene(g, ts, sceneTrack1, audioContext, cyclicCamera,
			centralizeCamera, false,
			sceneOptions{banner: "press: [p]lay or [q]uit"})
	}

	// scene1: hard-coded tilemap
	scene1 := newScene(g, ts, sceneTrack1, audioContext, cyclicCamera,
		centralizeCamera, showCoord, sceneOptions{})
	scene1.addSprite(50, 50, 0, ebitenImage)
	scene1.addSprite(100, 100, rotationScene1Sprite2, ebitenImage)

	// scene2: hard-coded tilemap
	scene2 := newScene(g, ts, sceneTrack2, audioContext, cyclicCamera,
		centralizeCamera, showCoord, sceneOptions{})
	scene2.addSprite(150, 150, 0, ebitenImage)
	scene2.addSprite(200, 200, oneQuarter, ebitenImage)

	// scene3: random tilemap

	var scene3 *scene
	{
		const tileEdgeCount = 100 // 100x100=10000
		layers := [][]int{generateLayer(tileEdgeCount)}
		ts3 := newTiles(bytes.NewReader(images.Tiles_png), tileSize, layers, tileEdgeCount)

		scene3 = newScene(g, ts3, sceneTrack3, audioContext, cyclicCamera,
			centralizeCamera, showCoord, sceneOptions{})

		// add a sprite close to top-left corner
		scene3.addSprite(50, 50, -oneQuarter, ebitenImage)

		// add a sprite at center of tilemap
		x := scene3.tiles.tilePixelWidth() / 2
		y := scene3.tiles.tilePixelHeight() / 2
		scene3.addSprite(float64(x), float64(y), -oneQuarter, ebitenImage)
	}

	// scene4: first scene
	var scene4 *scene
	{
		const tileEdgeCount = 100 // 100x100=10000
		layers := [][]int{generateLayer(tileEdgeCount)}
		ts := newTiles(bytes.NewReader(images.Tiles_png), tileSize, layers, tileEdgeCount)

		scene4 = newScene(g, ts, sceneTrack1, audioContext, true, true,
			showCoord, sceneOptions{})

		// add a sprite at center of tilemap
		x := scene3.tiles.tilePixelWidth() / 2
		y := scene3.tiles.tilePixelHeight() / 2
		scene4.addSprite(float64(x), float64(y), -oneQuarter, ebitenImage)
	}

	g.scenes = []*scene{scene0, scene1, scene2, scene3, scene4}

	g.sceneCurrent = 4 // first scene after start screen

	g.switchScene(g.sceneStart)

	// See comment in game.Layout method.
	log.Printf("Game screen size: %dx%d", g.screenWidth, g.screenHeight)

	return g
}

func (g *game) switchScene(newScene int) {

	log.Printf("switch scene to: %d", newScene)

	g.getCurrentScene().musicStop()

	if newScene == g.sceneStart {
		// save current scene to resume
		if g.sceneCurrent != g.sceneStart {
			// but not if we are on start
			g.sceneResume = g.sceneCurrent
		}
		g.sceneUpdateInput = sceneUpdateInputStart // resume or exit controls
	} else {
		g.sceneUpdateInput = sceneUpdateInputDefault // all controls
	}
	g.sceneCurrent = newScene

	g.getCurrentScene().musicStart()
}

func sceneUpdateInputStart(sc *scene) {
	g := sc.g
	if inpututil.IsKeyJustReleased(ebiten.KeyP) {
		g.switchScene(g.sceneResume)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyQ) {
		log.Printf("Q pressed, quitting")
		os.Exit(0)
	}
}

func sceneUpdateInputDefault(sc *scene) {

	g := sc.g

	//
	// handle burst of keys
	//
	keys := inpututil.AppendPressedKeys(nil)
	for _, p := range keys {
		//p := keys[len(keys)-1]

		switch p {
		case ebiten.KeyUp:
			g.getCurrentScene().cam.stepUp()
			continue
		case ebiten.KeyDown:
			g.getCurrentScene().cam.stepDown()
			continue
		case ebiten.KeyLeft:
			g.getCurrentScene().cam.stepRight()
			continue
		case ebiten.KeyRight:
			g.getCurrentScene().cam.stepLeft()
			continue
		case ebiten.KeyEscape:
			log.Printf("ESC pressed, switching to start screen")
			g.switchScene(g.sceneStart)
			continue
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
		next := (g.sceneCurrent + 1) % len(g.scenes)
		if next == g.sceneStart {
			next = (next + 1) % len(g.scenes)
		}
		g.switchScene(next)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyPeriod) {
		// toggle camera cyclic
		cam := g.getCurrentScene().cam
		cam.cyclic = !cam.cyclic
		log.Printf("Camera cyclic: %t", cam.cyclic)
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

	//
	// ui
	//
	//g.ui.Update()
	// Update the Label text to indicate if the ui is currently being hovered over or not
	//g.headerLbl.Label = fmt.Sprintf("Game Demo!\nUI is hovered: %t", input.UIHovered)
	//g.coordinateLbl.Label = g.getCurrentScene().getWorldCoordinates() // placeholder for coordinate display
	// Log out if we have clicked on the gamefield and NOT the ui
	/*
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !input.UIHovered {
			log.Println("Mouse clicked on gamefield")
		}
	*/

}

// Update is called every tick. Tick is a time unit for logical updating.
// The default value is 1/60 [s], then Update is called 60 times per second by
// default (i.e. an Ebitengine game works in 60 ticks-per-second).
func (g *game) Update() (err error) {

	g.sceneUpdateInput(g.getCurrentScene())

	//g.uiCoord = g.getCurrentScene().getWorldCoordinates()

	if _, e := g.debugui.Update(func(ctx *debugui.Context) error {
		x, y := 350, 30
		dx := x + 320
		dy := y + 150
		ctx.Window("Debugui Window", image.Rect(x, y, dx, dy), func(_ debugui.ContainerLayout) {
			// Place all your widgets inside a ctx.Window's callback.
			ctx.Text("test")

			// Use Loop if you ever need to make a loop to make widgets.
			const loopCount = 4
			ctx.Loop(loopCount, func(index int) {
				// Specify a presssing-button event handler by On.
				ctx.Button(fmt.Sprintf("Button %d", index)).On(func() {
					fmt.Printf("Button %d is pressed\n", index)
				})
			})
		})
		return nil
	}); e != nil {
		err = e
	}

	if g.pause {
		return
	}

	g.getCurrentScene().update()

	return
}

func (g *game) getCurrentScene() *scene {
	return g.scenes[g.sceneCurrent]
}

/*
func (g *game) switchScene() {
	g.getCurrentScene().musicStop()

	g.sceneCurrent = (g.sceneCurrent + 1) % len(g.scenes)
	log.Printf("Switching to scene %d of %d",
		g.sceneCurrent+1, len(g.scenes))

	g.getCurrentScene().musicStart()
}
*/

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

	sc := g.getCurrentScene()

	drawnTiles := sc.draw(screen, g.debug)

	//g.ui.Draw(screen)

	//g.drawSimpleUI(screen)

	/*
		if true {
			red := color.RGBA{0xff, 0, 0, 0xff}
			drawDebugAxis(screen, float32(sc.cam.x), float32(sc.cam.y),
				float32(sc.tiles.tilePixelWidth()), float32(sc.tiles.tilePixelHeight()), red)
		}
	*/

	if g.debug {
		tileDimX, tileDimY := sc.tiles.tilePixelDimensions()
		cam := sc.cam
		camLastX := cam.x + g.screenWidth - 1
		camLastY := cam.y + g.screenHeight - 1
		ebitenutil.DebugPrint(screen,
			fmt.Sprintf("TPS:%0.1f FPS:%0.1f tilemap:%dx%d cam:%dx%d-%dx%d camMax:%dx%d mouse:%dx%d win:%dx%d drawnTiles:%d",
				ebiten.ActualTPS(), ebiten.ActualFPS(),
				tileDimX, tileDimY,
				cam.x, cam.y,
				camLastX, camLastY,
				cam.maxX(), cam.maxY(),
				g.mouseX, g.mouseY,
				g.windowWidth, g.windowHeight,
				drawnTiles))

		//colorBlue := color.RGBA{0, 0, 0xff, 0xff}
		//drawDebugRect(screen, 1, 1, float32(g.screenWidth), float32(g.screenHeight), colorBlue)

		g.debugui.Draw(screen)
	}
}

/*
func drawDebugAxis(screen *ebiten.Image, camX, camY, width, height float32, axisColor color.RGBA) {

	const lineWidth = 1

	strokeOp := &vector.StrokeOptions{}
	strokeOp.Width = float32(lineWidth)
	drawOp := &vector.DrawPathOptions{}
	drawOp.ColorScale.ScaleWithColor(axisColor)
	drawOp.AntiAlias = false

	// vertical line
	x := 1 - camX
	var pathV vector.Path
	pathV.MoveTo(x, 1-camY)
	pathV.LineTo(x, height-camY)
	vector.StrokePath(screen, &pathV, strokeOp, drawOp)

	// horizontal line
	y := 1 - camY
	var pathH vector.Path
	pathH.MoveTo(1-camX, y)
	pathH.LineTo(width-camX, y)
	vector.StrokePath(screen, &pathH, strokeOp, drawOp)
}
*/

func drawDebugRect(screen *ebiten.Image, x1, y1, x2, y2 float32, borderColor color.RGBA) {

	const width = 1

	var path vector.Path
	path.MoveTo(x1, y1)
	path.LineTo(x1, y2)
	path.LineTo(x2, y2)
	path.LineTo(x2, y1)
	path.LineTo(x1, y1)
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
