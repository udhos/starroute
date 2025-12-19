package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func (g *game) drawSimpleUI(screen *ebiten.Image) {
	// coordinates
	op := &text.DrawOptions{}
	op.GeoM.Translate(13, 13)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, g.uiCoord, &text.GoTextFace{
		Source: g.mplusFaceSource,
		Size:   16,
	}, op)
}
