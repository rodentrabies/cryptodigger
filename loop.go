package main

import (
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func Loop() {
	cfg := pixelgl.WindowConfig{
		Title:  "Cryptodigger",
		Bounds: pixel.R(0, 0, 768, 432),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	backgroundPic, err := loadPicture(BackgroundSprite)
	if err != nil {
		panic(err)
	}
	plax := pixel.NewSprite(backgroundPic, backgroundPic.Bounds())

	diggerPic, err := loadPicture(DiggerSprite)
	if err != nil {
		panic(err)
	}
	diggerSprite := pixel.NewSprite(diggerPic, diggerPic.Bounds())

	diggerX, diggerY := 0, 2*FSize
	step := time.Tick(time.Millisecond * 100)
	scale := pixel.V(4, 4)
	for !win.Closed() {
		win.Clear(colornames.Skyblue)
		plax.Draw(win, pixel.IM.Scaled(pixel.ZV, 4).Moved(win.Bounds().Center()))
		rect := pixel.R(
			float64(diggerX), float64(diggerY),
			float64(diggerX+FSize), float64(diggerY+FSize),
		)
		diggerSprite.Set(diggerPic, rect)

		goldenRatio := pixel.V(win.Bounds().Center().X, win.Bounds().H()/2.61)
		mat := pixel.IM.ScaledXY(pixel.ZV, scale).Moved(goldenRatio)

		diggerSprite.Draw(win, mat)

		if win.Pressed(pixelgl.KeyD) {
			diggerY = 2 * FSize
			scale = pixel.V(FScale, FScale)
			select {
			case <-step:
				diggerX = (diggerX + FSize) % (4 * FSize)
			default:
			}
		} else if win.Pressed(pixelgl.KeyA) {
			diggerY = 2 * FSize
			scale = pixel.V(-FScale, FScale)
			select {
			case <-step:
				diggerX = (diggerX + FSize) % (4 * FSize)
			default:
			}
		} else if win.Pressed(pixelgl.MouseButtonLeft) {
			diggerY = 1 * FSize
			select {
			case <-step:
				diggerX = (diggerX + FSize) % (2 * FSize)
			default:
			}
		} else if win.Pressed(pixelgl.MouseButtonRight) {
			diggerY = 1 * FSize
			select {
			case <-step:
				diggerX = (diggerX+FSize)%(2*FSize) + 2*FSize
			default:
			}
		} else if win.JustReleased(pixelgl.KeyA) ||
			win.JustReleased(pixelgl.KeyD) ||
			win.JustReleased(pixelgl.MouseButtonLeft) ||
			win.JustReleased(pixelgl.MouseButtonRight) {
			diggerX, diggerY = 0, 2*FSize
		}

		win.Update()
	}
}
