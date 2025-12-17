package main

import (
	"log"
	"os"
)

func loadAsset(filename string) ([]byte, error) {
	input := "assets/" + filename
	return os.ReadFile(input)
}

func mustLoadAsset(caller, filename string) []byte {
	data, err := loadAsset(filename)
	if err != nil {
		log.Fatalf("%s: mustLoadAsset: error: open file: %s: %v", caller, filename, err)
	}
	return data
}
