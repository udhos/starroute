package main

import (
	"fmt"
	"testing"
)

type tyleTest struct {
	name string

	layerTiles     int
	layerTileWidth int
	tilePixelWidth int
	winX           int
	winY           int
	winWidth       int
	winHeight      int

	expectTileOffset  int
	expectTileXAmount int
	expectTileYAmount int
}

var tyleTestTable = []tyleTest{
	{
		name:              "window equal single tile",
		layerTiles:        1,
		layerTileWidth:    1,
		tilePixelWidth:    8,
		winX:              0,
		winY:              0,
		winWidth:          8,
		winHeight:         8,
		expectTileOffset:  0,
		expectTileXAmount: 1,
		expectTileYAmount: 1,
	},
	{
		name:              "window within single tile",
		layerTiles:        1,
		layerTileWidth:    1,
		tilePixelWidth:    8,
		winX:              1,
		winY:              1,
		winWidth:          6,
		winHeight:         6,
		expectTileOffset:  0,
		expectTileXAmount: 1,
		expectTileYAmount: 1,
	},
	{
		name:              "window bigger than single tile",
		layerTiles:        1,
		layerTileWidth:    1,
		tilePixelWidth:    8,
		winX:              1,
		winY:              1,
		winWidth:          10,
		winHeight:         10,
		expectTileOffset:  0,
		expectTileXAmount: 1,
		expectTileYAmount: 1,
	},
	{
		name:              "window picks only exactly central tile of 9-tiles",
		layerTiles:        9,
		layerTileWidth:    3,
		tilePixelWidth:    8,
		winX:              8,
		winY:              8,
		winWidth:          8,
		winHeight:         8,
		expectTileOffset:  4,
		expectTileXAmount: 1,
		expectTileYAmount: 1,
	},
	{
		name:              "window picks only, from within, central tile of 9-tiles",
		layerTiles:        9,
		layerTileWidth:    3,
		tilePixelWidth:    8,
		winX:              9,
		winY:              9,
		winWidth:          6,
		winHeight:         6,
		expectTileOffset:  4,
		expectTileXAmount: 1,
		expectTileYAmount: 1,
	},
	{
		name:              "window picks all 9-tiles from within",
		layerTiles:        9,
		layerTileWidth:    3,
		tilePixelWidth:    8,
		winX:              7,
		winY:              7,
		winWidth:          10,
		winHeight:         10,
		expectTileOffset:  0,
		expectTileXAmount: 3,
		expectTileYAmount: 3,
	},
	{
		name:              "window picks all 9-tiles exactly",
		layerTiles:        9,
		layerTileWidth:    3,
		tilePixelWidth:    8,
		winX:              0,
		winY:              0,
		winWidth:          24,
		winHeight:         24,
		expectTileOffset:  0,
		expectTileXAmount: 3,
		expectTileYAmount: 3,
	},
}

// go test -count 1 -run '^TestFindTilemapWindow$' ./...
func TestFindTilemapWindow(t *testing.T) {
	for i, data := range tyleTestTable {
		name := fmt.Sprintf("%02d of %02d: %s", i+1, len(tyleTestTable), data.name)
		t.Run(name, func(t *testing.T) {
			offset, xAmount, yAmount := findTilemapWindow(data.layerTiles,
				data.layerTileWidth, data.tilePixelWidth,
				data.winX, data.winY, data.winWidth, data.winHeight)
			if offset != data.expectTileOffset {
				t.Errorf("wrong offset: expected %d got %d", data.expectTileOffset, offset)
			}
			if xAmount != data.expectTileXAmount {
				t.Errorf("wrong xAmount: expected %d got %d", data.expectTileXAmount, xAmount)
			}
			if yAmount != data.expectTileYAmount {
				t.Errorf("wrong yAmount: expected %d got %d", data.expectTileYAmount, yAmount)
			}
		})
	}
}
