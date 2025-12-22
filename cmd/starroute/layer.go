package main

import "math/rand/v2"

func generateLayer(tileEdgeCount int) []int {
	size := tileEdgeCount * tileEdgeCount
	layer := make([]int, size)
	for i := range size {
		if randOneIn(20) {
			layer[i] = 218
			continue
		}
		layer[i] = 243
	}
	return layer
}

func randOneIn(n int) bool {
	return rand.IntN(n) == n-1
}

func generateLayerSingleTile(tileEdgeCount, index int) []int {
	size := tileEdgeCount * tileEdgeCount
	layer := make([]int, size)
	for i := range size {
		layer[i] = index
	}
	return layer
}
