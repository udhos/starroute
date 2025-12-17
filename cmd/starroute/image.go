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

	return transformImageScaleAlpha(origEbitenImage, scaleAlpha)
}

func transformImageScaleAlpha(origEbitenImage *ebiten.Image, scaleAlpha float32) *ebiten.Image {

	s := origEbitenImage.Bounds().Size()
	ebitenImage := ebiten.NewImage(s.X, s.Y)

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(scaleAlpha)
	ebitenImage.DrawImage(origEbitenImage, op)

	return ebitenImage
}

/*
// createImageFromFileSVG is broken, cannot draw with transparency.
func createImageFromFileSVG(r io.Reader, scaleAlpha float32) *ebiten.Image {
	c, err := canvas.ParseSVG(r)
	if err != nil {
		log.Fatalf("createImageFromFileSVG error: %v", err)
	}

	w, h := c.Size()
	w *= .2
	h *= .2

	// resize svg image to .2
	scale := canvas.Matrix.Scale(canvas.Identity, .2, .2)

	c.Transform(scale)
	c.Clip(canvas.Rect{X0: 0, Y0: 0, X1: w, Y1: h})

	canvasImg := rasterizer.Draw(c, canvas.DPMM(96.0/25.4), canvas.LinearColorSpace{})

	img := ebiten.NewImageFromImage(canvasImg)

	return img

	//return transformImageScaleAlpha(img, scaleAlpha)
}
*/
