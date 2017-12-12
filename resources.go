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
	PopupSprite      = "assets/popup.png"
	ButtonsSprite    = "assets/buttons.png"

	CoinsFont = "assets/tag_font.ttf"
	TextFont  = "assets/text_font.ttf"

	LineLength = 30
)

//------------------------------------------------------------------------------
// Alphabet resource
type Alphabet text.Atlas

func LoadAlphabet(filename string, size, dpi float64, runeset ...[]rune) *Alphabet {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	ttf, err := truetype.Parse(bytes)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(ttf, &truetype.Options{Size: size, DPI: dpi})
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

//------------------------------------------------------------------------------
// Button resource
type Button struct {
	sprite *pixel.Sprite
	matrix pixel.Matrix

	upFrame, downFrame, currentFrame pixel.Rect
}

func NewButton(sprite *pixel.Sprite, uFrame, dFrame pixel.Rect, matrix pixel.Matrix) *Button {
	return &Button{
		sprite:       sprite,
		upFrame:      uFrame,
		downFrame:    dFrame,
		currentFrame: uFrame,
		matrix:       matrix,
	}
}

func (b Button) Bounds() pixel.Rect {
	r := b.sprite.Frame()
	return pixel.Rect{
		Min: b.matrix.Project(r.Center().Sub(r.Max)),
		Max: b.matrix.Project(r.Center()),
	}
}

func (b *Button) Push() {
	b.currentFrame = b.downFrame
}

func (b *Button) Release() {
	b.currentFrame = b.upFrame
}

func (b Button) Draw(t pixel.Target) {
	b.sprite.Set(b.sprite.Picture(), b.currentFrame)
	b.sprite.Draw(t, b.matrix)
}

func (b Button) Within(v pixel.Vec) bool {
	r := b.Bounds()
	return v.X > r.Min.X && v.X < r.Max.X && v.Y > r.Min.Y && v.Y < r.Max.Y
}

func (b *Button) Register(registry map[*Button]func(), action func()) {
	registry[b] = action
}

func (b *Button) Unregister(registry map[*Button]func()) {
	delete(registry, b)
}

//------------------------------------------------------------------------------
// Event box
type EventBox struct {
	message  string
	position pixel.Vec
	button   *Button
}

//------------------------------------------------------------------------------
// Utils
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
