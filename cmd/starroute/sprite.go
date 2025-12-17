package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type sprite struct {
	x, y          float64
	width, height int
	angle         float64
	angleNative   float64 // undo this intrinsic rotate of image to point image to zero angle (right)
	image         *ebiten.Image
}

func (s *sprite) update() {
	s.angle = math.Mod(s.angle+1, maxAngle)
}

func (s *sprite) draw(op ebiten.DrawImageOptions, screen *ebiten.Image, camX, camY float64, debug bool) {

	w, h := s.image.Bounds().Dx(), s.image.Bounds().Dy()

	centerX := float64(w) / 2
	centerY := float64(h) / 2

	// move the rotation center to the origin
	op.GeoM.Translate(-centerX, -centerY)

	// scale
	//op.GeoM.Scale(s.rescale, s.rescale)

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
