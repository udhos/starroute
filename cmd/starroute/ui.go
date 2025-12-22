package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

/*
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
*/

func (sc *scene) drawSimpleUI(screen *ebiten.Image) {
	g := sc.g

	// coordinates
	if sc.showCoord {
		op := &text.DrawOptions{}
		op.GeoM.Translate(13, 13)
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, sc.uiCoord, &text.GoTextFace{
			Source: g.mplusFaceSource,
			Size:   16,
		}, op)
	}

	if sc.opt.banner != "" {
		op := &text.DrawOptions{}
		op.GeoM.Translate(50, 300)
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, sc.opt.banner, &text.GoTextFace{
			Source: g.mplusFaceSource,
			Size:   40,
		}, op)
	}
}
