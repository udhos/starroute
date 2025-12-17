package main

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
	"github.com/udhos/starroute/music"
)

const (
	sceneTrack1 = iota
	sceneTrack2
	sceneTrack3
)

type scene struct {
	sprites      []*sprite
	tiles        *tiles
	musicPlayer  *music.Player
	musicTrack   int
	audioContext *audio.Context
	cam          *camera
	g            *game
}

func newScene(g *game, ts *tiles, musicTrack int, audioContext *audio.Context) *scene {
	sc := &scene{
		g:            g,
		tiles:        ts,
		musicTrack:   musicTrack,
		audioContext: audioContext,
	}
	sc.cam = newCamera(sc)
	return sc
}

func (sc *scene) musicStart() {
	sc.musicStop()
	//m, err := music.NewPlayer(audioContext, music.TypeOgg, bytes.NewReader(raudio.Ragtime_ogg))
	//m, err := music.NewPlayer(audioContext, music.TypeMP3, bytes.NewReader(raudio.Ragtime_mp3))

	var m *music.Player
	var err error

	switch sc.musicTrack {
	case sceneTrack1:
		const input = "champions-victory-winner-background-music-388566.mp3"
		data := mustLoadAsset("musicStart", input)
		m, err = music.NewPlayer(sc.audioContext, music.TypeMP3, bytes.NewReader(data))
	case sceneTrack2:
		m, err = music.NewPlayer(sc.audioContext, music.TypeMP3, bytes.NewReader(raudio.Ragtime_mp3))
	case sceneTrack3:
		// no music
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

func (sc *scene) draw(screen *ebiten.Image, debug bool) int {
	countTiles := sc.tiles.draw(screen, sc.cam)

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
		op.GeoM.Reset()
		sc.sprites[i].draw(op, screen, camX, camY, debug)
	}

	return countTiles
}
