package main

import (
	"image"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func createImage(r io.Reader, scaleAlpha float32) *ebiten.Image {
	img, _, err := image.Decode(r)
	if err != nil {
		log.Fatalf("createImage error: %v", err)
	}
	origEbitenImage := ebiten.NewImageFromImage(img)

	s := origEbitenImage.Bounds().Size()
	ebitenImage := ebiten.NewImage(s.X, s.Y)

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(scaleAlpha)
	ebitenImage.DrawImage(origEbitenImage, op)

	return ebitenImage
}
