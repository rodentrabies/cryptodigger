package main

import (
	"image"
	"os"

	"github.com/faiface/pixel"
)

var (
	DiggerSprite     = "assets/digger.png"
	BlockSprite      = "assets/blocks.png"
	BackgroundSprite = "assets/background.png"

	FSize  int     = 16
	FScale float64 = 4
	ASize  int     = FSize * int(FScale)
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
