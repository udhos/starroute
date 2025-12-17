package main

import (
	"bytes"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

func (g *game) getEbitenUI() *ebitenui.UI {
	// load label text font
	face, _ := loadFont(18)

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(5)))),
	)

	// Because this container has a backgroundImage set we track that the ui is hovered over.
	headerContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{R: 200, G: 200, B: 200, A: 100})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  true,
				StretchVertical:    false,
			}),
			// Uncomment this to not track that you are hovering over this header
			// widget.WidgetOpts.TrackHover(false),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	g.headerLbl = widget.NewText(
		widget.TextOpts.Text("Game Demo!", &face, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
			// Uncomment to force tracking hover of this element
			// widget.WidgetOpts.TrackHover(true),
		),
	)
	headerContainer.AddChild(g.headerLbl)
	rootContainer.AddChild(headerContainer)

	hProgressbar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  true,
				StretchVertical:    false,
			}),
			// Set the minimum size for the progress bar.
			// This is necessary if you wish to have the progress bar be larger than
			// the provided track image. In this exampe since we are using NineSliceColor
			// which is 1px x 1px we must set a minimum size.
			widget.WidgetOpts.MinSize(200, 20),
			// Set this parameter to indicate we want do not want to track that this ui element is being hovered over.
			// widget.WidgetOpts.TrackHover(false),
		),
		widget.ProgressBarOpts.Images(
			// Set the track images (Idle, Disabled).
			&widget.ProgressBarImage{
				Idle: eimage.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			// Set the progress images (Idle, Disabled).
			&widget.ProgressBarImage{
				Idle: eimage.NewNineSliceColor(color.NRGBA{255, 255, 100, 255}),
			},
		),
		// Set the min, max, and current values.
		widget.ProgressBarOpts.Values(0, 10, 7),
		// Set how much of the track is displayed when the bar is overlayed.
		widget.ProgressBarOpts.TrackPadding(&widget.Insets{
			Top:    2,
			Bottom: 2,
		}),
	)

	rootContainer.AddChild(hProgressbar)

	// Create a label to show the percentage on top of the progress bar
	label2 := widget.NewText(
		widget.TextOpts.Text("70%", &face, color.Black),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
	)
	rootContainer.AddChild(label2)

	return &ebitenui.UI{
		Container: rootContainer,
		//Call a render method after the rootContainer is drawn but before any ebitenui.Windows are drawn
		//PostRenderHook: g.Render,
	}
}

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}
