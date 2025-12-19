package main

import (
	"image"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type tiles struct {
	tilesImage      *ebiten.Image
	tileSize        int
	layers          [][]int
	tileLayerXCount int
}

func (ts tiles) tilePixelDimensions() (int, int) {
	return ts.tilePixelWidth(), ts.tilePixelHeight()
}

func (ts tiles) tilePixelWidth() int {
	return ts.tileSize * ts.tileLayerXCount
}

func (ts tiles) tilePixelHeight() int {
	return ts.tileSize * len(ts.layers[0]) / ts.tileLayerXCount
}

func newTiles(r io.Reader, tileSize int, layers [][]int, tileLayerXCount int) *tiles {

	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(r)
	if err != nil {
		log.Fatal(err)
	}

	ts := &tiles{
		tilesImage:      ebiten.NewImageFromImage(img),
		tileSize:        tileSize,
		layers:          layers,
		tileLayerXCount: tileLayerXCount,
	}

	log.Printf("Tile size: %d", tileSize)

	log.Printf("Tiles image size: %dx%d", ts.tilesImage.Bounds().Dx(), ts.tilesImage.Bounds().Dy())

	log.Printf("Tile layer X count: %d", tileLayerXCount)

	dimX, dimY := ts.tilePixelDimensions()

	log.Printf("Tile layer size: %dx%d", dimX, dimY)

	tilesImageXCount := ts.tilesImage.Bounds().Dx() / tileSize
	tilesImageYCount := ts.tilesImage.Bounds().Dy() / tileSize

	log.Printf("Tiles image has %dx%d tiles", tilesImageXCount, tilesImageYCount)

	return ts
}

func (ts *tiles) draw(screen *ebiten.Image, cam *camera) int {
	tileSize := ts.tileSize

	// Draw each tile with each DrawImage call.
	// As the source images of all DrawImage calls are always same,
	// this rendering is done very efficiently.
	// For more detail, see https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage

	// number of tiles per row defined in the tile layer
	xCount := ts.tileLayerXCount

	// number of tiles in the tiles image
	tileImageXCount := ts.tilesImage.Bounds().Dx() / tileSize

	/*
		for _, l := range ts.layers {
			for i, t := range l {
				op := &ebiten.DrawImageOptions{}
				screenX := (i % xCount) * tileSize
				screenY := (i / xCount) * tileSize
				op.GeoM.Translate(float64(screenX-cam.x), float64(screenY-cam.y))

				sx := (t % tileImageXCount) * tileSize
				sy := (t / tileImageXCount) * tileSize
				subImage := ts.tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image)
				screen.DrawImage(subImage, op)
			}
		}
	*/

	screenWidth := screen.Bounds().Dx()
	screenHeight := screen.Bounds().Dy()

	var sum int

	if cam.cyclic {
		// cyclic

	} else {
		// non-cyclic

		for _, l := range ts.layers {
			offset, xAmount, yAmount := findTilemapWindow(len(l), ts.tileLayerXCount, ts.tileSize,
				cam.x, cam.y, screenWidth, screenHeight)

			i := offset
			for range yAmount {
				for range xAmount {
					t := l[i]

					op := &ebiten.DrawImageOptions{}
					screenX := (i % xCount) * tileSize
					screenY := (i / xCount) * tileSize
					op.GeoM.Translate(float64(screenX-cam.x), float64(screenY-cam.y))

					sx := (t % tileImageXCount) * tileSize
					sy := (t / tileImageXCount) * tileSize
					subImage := ts.tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image)
					screen.DrawImage(subImage, op)

					sum++
					i++
				}
				i += xCount - xAmount
			}
		}
	} // non-cyclic

	return sum
}

var sampleLayers = [][]int{
	{
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 244, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 243, 244, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 219, 243, 243, 243, 219, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 244, 243, 243, 243,
		243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
	},

	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 26, 27, 28, 29, 30, 31, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 51, 52, 53, 54, 55, 56, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 76, 77, 78, 79, 80, 81, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 101, 102, 103, 104, 105, 106, 0, 0, 0, 0,

		0, 0, 0, 0, 0, 126, 127, 128, 129, 130, 131, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 303, 303, 245, 242, 303, 303, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,

		0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
	},
}

// findTilemapWindow finds within a layer, composed of layerTiles tiles with
// layerTileWidth tiles per row and width tilePixelWidth pixels, the subset
// of tiles that must be drawn to completely fill the windown at
// winX, winY of size winWidth x winHeight. Only tiles
// that intersect the window with a least one pixel should be drawn.
// Result is:
// tileOffset is the offset of the first tile to be drawn.
// tileXAmount is the amount of tiles to be drawn per row.
// tileYAmount is the amount of rows of tiles to be drawn.
// The caller must then draw tileXAmount tiles starting from tileOffset,
// then skip tileXAmount-layerTileWidth tiles, and so forth,
// up to tileYAmount rows.
func findTilemapWindow(layerTiles, layerTileWidth, tilePixelWidth,
	winX, winY, winWidth, winHeight int) (tileOffset, tileXAmount,
	tileYAmount int) {
	// Calculate the starting column and row
	startCol := winX / tilePixelWidth
	startRow := winY / tilePixelWidth

	// Calculate the ending column and row (inclusive)
	// We subtract 1 from the sum to get the last pixel, then divide
	endCol := (winX + winWidth - 1) / tilePixelWidth
	endRow := (winY + winHeight - 1) / tilePixelWidth

	// Calculate the number of tiles in each dimension
	tileXAmount = endCol - startCol + 1
	tileYAmount = endRow - startRow + 1

	// Calculate the layer height based on total tiles and width
	layerTileHeight := layerTiles / layerTileWidth
	if layerTiles%layerTileWidth != 0 {
		layerTileHeight++
	}

	// Clamp to layer boundaries
	if tileXAmount > layerTileWidth-startCol {
		tileXAmount = layerTileWidth - startCol
	}
	if tileYAmount > layerTileHeight-startRow {
		tileYAmount = layerTileHeight - startRow
	}

	// Calculate the linear offset of the first tile
	tileOffset = startRow*layerTileWidth + startCol

	return
}
