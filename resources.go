package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"math"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
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

	CoinsFont = "assets/somepx/Runners/Runners.ttf"
	TextFont  = "assets/somepx/Runners/Runners.ttf"
)

type Alphabet text.Atlas

func LoadAlphabet(filename string, runeset ...[]rune) *Alphabet {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	font, err := truetype.Parse(bytes)
	if err != nil {
		panic(err)
	}
	alphabet := Alphabet(*text.NewAtlas(truetype.NewFace(font, nil), runeset...))
	return &alphabet
}

func (alphabet Alphabet) Draw(t pixel.Target, str string, m pixel.Matrix) {
	atlas := text.Atlas(alphabet)
	textDrawer := text.New(pixel.V(0, 0), &atlas)
	fmt.Fprintf(textDrawer, "%s", str)
	textDrawer.Draw(t, m)
}

func LoadPicture(path string) (pixel.Picture, error) {
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

func LoadSprite(file string) (pixel.Picture, *pixel.Sprite) {
	pic, err := LoadPicture(file)
	if err != nil {
		panic(err)
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())
	return pic, sprite
}

func Float64Mod(a, b float64) float64 {
	return a - b*math.Floor(a/b)
}
