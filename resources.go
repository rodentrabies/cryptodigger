package main

import (
	"fmt"
	"image"
	"image/color"
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

	CamSpeedFactor float64 = 150
	CamSpeed       float64 = CamSpeedFactor * FScale

	FSize float64 = 16
	ASize float64 = FSize * FScale

	DiggerSprite     = "assets/digger.png"
	IconsSprite      = "assets/icons.png"
	BackgroundSprite = "assets/background.png"

	CoinsFont = "assets/tag_font.ttf"
	TextFont  = "assets/somepx/Smart/Smart.ttf"
)

type Alphabet text.Atlas

func LoadAlphabet(filename string, size float64, runeset ...[]rune) *Alphabet {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	ttf, err := truetype.Parse(bytes)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(ttf, &truetype.Options{
		Size: size,
		DPI:  144,
	})
	alphabet := Alphabet(*text.NewAtlas(face, runeset...))
	return &alphabet
}

func (alphabet Alphabet) Draw(t pixel.Target, str string, color color.Color, m pixel.Matrix) {
	atlas := text.Atlas(alphabet)
	textDrawer := text.New(pixel.V(0, 0), &atlas)
	textDrawer.Color = color
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
