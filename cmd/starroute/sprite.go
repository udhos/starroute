package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
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
