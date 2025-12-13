package main

import "math/rand/v2"

func generateLayer(tileEdgeCount int) []int {
	size := tileEdgeCount * tileEdgeCount
	layer := make([]int, size)
	for i := range size {
		if rand.IntN(20) == 19 {
			layer[i] = 218
			continue
		}
		layer[i] = 243
	}
	return layer
}
