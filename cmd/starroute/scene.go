package main

import (
	"bytes"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/udhos/starroute/music"
)

const (
	sceneTrack1 = 0
	sceneTrack2 = 1
)

type camera struct {
	x, y int
}

type scene struct {
	sprites      []*sprite
	tiles        *tiles
	musicPlayer  *music.Player
	musicTrack   int
	audioContext *audio.Context
	cam          *camera
}

func newScene(ts *tiles, musicTrack int, audioContext *audio.Context) scene {
	sc := scene{
		tiles:        ts,
		musicTrack:   musicTrack,
		audioContext: audioContext,
		cam:          &camera{},
	}
	return sc
}

func (sc *scene) musicStart() {
	sc.musicStop()
	//m, err := music.NewPlayer(audioContext, music.TypeOgg, bytes.NewReader(raudio.Ragtime_ogg))
	//m, err := music.NewPlayer(audioContext, music.TypeMP3, bytes.NewReader(raudio.Ragtime_mp3))

	var m *music.Player
	var err error

	if sc.musicTrack == sceneTrack1 {
		const input = "assets/champions-victory-winner-background-music-388566.mp3"
		data, errRead := os.ReadFile(input)
		if errRead != nil {
			log.Fatalf("scene.musicStart: error: open file: %s: %v", input, errRead)
		}
		m, err = music.NewPlayer(sc.audioContext, music.TypeMP3, bytes.NewReader(data))
	} else {
		m, err = music.NewPlayer(sc.audioContext, music.TypeMP3, bytes.NewReader(raudio.Ragtime_mp3))
	}

	if err != nil {
		log.Fatalf("scene.musicStart error: %v", err)
	}
	sc.musicPlayer = m
}

func (sc *scene) musicStop() {
	if sc.musicPlayer == nil {
		return
	}
	sc.musicPlayer.Close()
	sc.musicPlayer = nil
}

func (sc *scene) addSprite(x, y, angleNative float64, spriteImage *ebiten.Image) {
	w, h := spriteImage.Bounds().Dx(), spriteImage.Bounds().Dy()
	spr := sprite{
		x:           x,
		y:           y,
		width:       w,
		height:      h,
		angleNative: angleNative,
		image:       spriteImage,
	}
	sc.sprites = append(sc.sprites, &spr)
}

func (sc *scene) update() {
	// Update all sprites.
	for _, spr := range sc.sprites {
		spr.update()
	}

	if sc.musicPlayer != nil {
		if err := sc.musicPlayer.Update(); err != nil {
			log.Printf("scene.update: music player error: %v", err)
		}
	}
}

func (sc *scene) draw(screen *ebiten.Image, debug bool) {
	sc.tiles.draw(screen, sc.cam)

	// Draw each sprite.
	// DrawImage can be called many many times, but in the implementation,
	// the actual draw call to GPU is very few since these calls satisfy
	// some conditions e.g. all the rendering sources and targets are same.
	// For more detail, see:
	// https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage

	camX := float64(sc.cam.x)
	camY := float64(sc.cam.y)

	var op ebiten.DrawImageOptions

	for i := 0; i < len(sc.sprites); i++ {
		s := sc.sprites[i]
		op.GeoM.Reset()

		w, h := s.image.Bounds().Dx(), s.image.Bounds().Dy()

		centerX := float64(w) / 2
		centerY := float64(h) / 2

		// move the rotation center to the origin
		op.GeoM.Translate(-centerX, -centerY)

		// rotate around the origin

		// undo intrinsic image rotation
		angleNativeRad := pi2 * float64(s.angleNative) / maxAngle
		op.GeoM.Rotate(-angleNativeRad)

		// apply actual rotation
		angleRad := pi2 * float64(s.angle) / maxAngle
		op.GeoM.Rotate(angleRad)

		// undo the translation used to move the rotation center
		op.GeoM.Translate(centerX, centerY)

		// apply the actual object's position
		op.GeoM.Translate(s.x-camX, s.y-camY)

		screen.DrawImage(s.image, &op)

		if debug {
			x := s.x + centerX - camX
			y := s.y + centerY - camY
			//
			// Red show how much the sprite was rotated back (counter clockwise) to
			// make its front point to the right (zero angle).
			//
			colorRed := color.RGBA{0xff, 0, 0, 0xff}
			drawDebugArrow(screen, x, y,
				angleRad-angleNativeRad, 20, 3, colorRed)

			//
			// Yellow show the sprint front direction and should point to
			// right (zero angle).
			//
			colorYellow := color.RGBA{0xff, 0xff, 0, 0xff}
			drawDebugArrow(screen, x, y,
				angleRad, 30, 1, colorYellow)
		}
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
