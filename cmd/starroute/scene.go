package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type scene struct {
	sprites []*sprite
	tiles   *tiles
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
}

func (sc *scene) draw(screen *ebiten.Image) {
	sc.tiles.draw(screen)

	// Draw each sprite.
	// DrawImage can be called many many times, but in the implementation,
	// the actual draw call to GPU is very few since these calls satisfy
	// some conditions e.g. all the rendering sources and targets are same.
	// For more detail, see:
	// https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage

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
		op.GeoM.Translate(float64(s.x), float64(s.y))

		screen.DrawImage(s.image, &op)

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
