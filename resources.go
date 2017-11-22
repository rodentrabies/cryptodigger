package main

import (
	"image"
	"math"
	"os"

	"github.com/faiface/pixel"
)

var (
	WinWidth, WinHeight float64 = 768, 384
	// WinWidth, WinHeight float64 = 896, 512
	// WinWidth, WinHeight float64 = 1536, 768

	FScale float64 = 4
	// FScale float64 = 8

	CamSpeedFactor float64 = 125
	CamSpeed       float64 = CamSpeedFactor * FScale

	FSize float64 = 16
	ASize float64 = FSize * FScale

	DiggerSprite     = "assets/digger.png"
	BlockSprite      = "assets/blocks.png"
	BackgroundSprite = "assets/background.png"
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
