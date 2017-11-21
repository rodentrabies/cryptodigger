package main

import (
	"image"
	"math"
	"os"

	"github.com/faiface/pixel"
)

var (
	DiggerSprite     = "assets/digger.png"
	BlockSprite      = "assets/blocks.png"
	BackgroundSprite = "assets/background.png"

	FSize  float64 = 16
	FScale float64 = 4
	ASize  float64 = FSize * FScale
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func Float64Mod(a, b float64) float64 {
	return a - b*math.Floor(a/b)
}
