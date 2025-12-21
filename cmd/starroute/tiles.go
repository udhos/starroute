package main

import (
	"image"
	"image/color"
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

// quad represents one of the four quadrants needed to draw with a cyclic camera.
type quad struct {
	draw bool

	// quadrant offset relative to camera
	camOffsetX, camOffsetY int

	// quadrant world position to get tiles from tilemap
	worldX, worldY int

	width, height int
}

// getQuadrants returns the four quadrants needed to draw with a cyclic camera.
//
// world size is given by tilemap size in pixels: tilePixelWidth(), tilePixelHeight()
// every tile is tileSize x tileSize pixels
// world width in tiles is tileLayerXCount
// screenWidth,screenHeight is screen dimensions in pixels
// cam gives the viewport coordinates (region of the world drawn on the screen) in pixels
// a tile value in a layer gives the index of the tile graphic in the tiles image,
// encoded as tileX + tileY*tileImageXCount
//
// the cyclic camera is drawn in 4 quadrants to cover all cases
// quadrant 1: always drawn
// quadrant 2: when part of the view is beyond the right edge of the tilemap
// quadrant 3: when part of the view is beyond the bottom edge of the tilemap
// quadrant 4: when part of the view is beyond both the right and bottom edges of the tilemap
func (ts *tiles) getQuadrants(cam *camera, screenWidth, screenHeight int) [4]quad {

	tilemapWidth := ts.tilePixelWidth()
	tilemapHeight := ts.tilePixelHeight()

	widthQuads1and3 := min(screenWidth, tilemapWidth-cam.x)
	heightQuads1and2 := min(screenHeight, tilemapHeight-cam.y)
	widthQuads2and4 := screenWidth - (tilemapWidth - cam.x)
	heightQuads3and4 := screenHeight - (tilemapHeight - cam.y)

	drawQuadrant2 := cam.x+screenWidth > tilemapWidth
	drawQuadrant3 := cam.y+screenHeight > tilemapHeight
	drawQuadrant4 := drawQuadrant2 && drawQuadrant3

	return [4]quad{
		// quadrant 1: top-left
		{
			draw:       true,
			camOffsetX: 0, camOffsetY: 0,
			worldX: cam.x, worldY: cam.y,
			width: widthQuads1and3, height: heightQuads1and2,
		},

		// quadrant 2: top-right
		{
			draw:       drawQuadrant2,
			camOffsetX: widthQuads1and3, camOffsetY: 0,
			worldX: 0, worldY: cam.y,
			width: widthQuads2and4, height: heightQuads1and2,
		},

		// quadrant 3: bottom-left
		{
			draw:       drawQuadrant3,
			camOffsetX: 0, camOffsetY: heightQuads1and2,
			worldX: cam.x, worldY: 0,
			width: widthQuads1and3, height: heightQuads3and4,
		},

		// quadrant 4: bottom-right
		{
			draw:       drawQuadrant4,
			camOffsetX: widthQuads1and3, camOffsetY: heightQuads1and2,
			worldX: 0, worldY: 0,
			width: widthQuads2and4, height: heightQuads3and4,
		},
	}
}

func (ts *tiles) draw(screen *ebiten.Image, cam *camera, debug bool, quads *[4]quad) int {

	// Draw each tile with each DrawImage call.
	// As the source images of all DrawImage calls are always same,
	// this rendering is done very efficiently.
	// For more detail, see https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage

	screenWidth := screen.Bounds().Dx()
	screenHeight := screen.Bounds().Dy()

	var sum int

	if cam.cyclic {
		// cyclic

		for _, q := range quads {
			if !q.draw {
				continue
			}
			sum += ts.drawQuadrant(screen,
				q.worldX, q.worldY,
				q.width, q.height,
				q.camOffsetX, q.camOffsetY)
		}

		if debug {
			yellow := color.RGBA{0xff, 0xff, 0x00, 0xff}
			red := color.RGBA{0xff, 0x00, 0x00, 0xff}
			green := color.RGBA{0x00, 0xff, 0x00, 0xff}
			blue := color.RGBA{0x00, 0x00, 0xff, 0xff}
			colors := []color.RGBA{
				yellow,
				red,
				green,
				blue,
			}
			for i, q := range quads {
				if q.draw {
					drawDebugRect(screen,
						float32(1+q.camOffsetX), float32(1+q.camOffsetY),
						float32(q.camOffsetX+q.width), float32(q.camOffsetY+q.height),
						colors[i])
				}
			}
		}

	} else {
		// non-cyclic

		const camOffsetX = 0
		const camOffsetY = 0

		sum = ts.drawQuadrant(screen,
			cam.x, cam.y,
			screenWidth, screenHeight,
			camOffsetX, camOffsetY)

		if debug {
			yellow := color.RGBA{0xff, 0xff, 0x00, 0xff}
			drawDebugRect(screen,
				float32(1+camOffsetX), float32(1+camOffsetY),
				float32(camOffsetX+screenWidth), float32(camOffsetY+screenHeight),
				yellow)
		}

	} // non-cyclic

	return sum
}

func (ts *tiles) drawQuadrant(screen *ebiten.Image,
	worldX, worldY,
	width, height,
	camOffsetX, camOffsetY int) int {

	var sum int

	tileSize := ts.tileSize
	xCount := ts.tileLayerXCount
	tileImageXCount := ts.tilesImage.Bounds().Dx() / tileSize

	for _, l := range ts.layers {
		offset, xAmount, yAmount := findTilemapWindow(len(l), ts.tileLayerXCount, ts.tileSize,
			worldX, worldY, width, height)

		i := offset
		for range yAmount {
			for range xAmount {
				t := l[i]

				op := &ebiten.DrawImageOptions{}
				// screenX,screenY is the position on the screen where the tile must be drawn
				// i % xCount gives the column of the tile in the tile layer
				// i / xCount gives the row of the tile in the tile layer
				// We translate by -worldX and -worldY to account for the camera
				// position, and add the quadrant camera offset so this quadrant
				// is drawn at the correct place on the screen when wrapping.
				screenX := (i % xCount) * tileSize
				screenY := (i / xCount) * tileSize
				op.GeoM.Translate(float64(screenX-worldX+camOffsetX), float64(screenY-worldY+camOffsetY))

				// sx,sy is the position within the tiles image where the tile graphic is located
				// t % tileImageXCount gives the column of the tile in the tiles image
				// t / tileImageXCount gives the row of the tile in the tiles image
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

	/*
		if debug {
			drawDebugRect(screen, float32(1+beginX-camX), float32(1+beginY-camY),
				float32(beginX+width-camX), float32(beginY+height-camY), debugColor)
		}
	*/

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
