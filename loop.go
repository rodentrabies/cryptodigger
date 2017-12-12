package main

import (
	"strconv"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

var SM = pixel.IM.Scaled(pixel.ZV, FScale)

func Loop() {
	cfg := pixelgl.WindowConfig{
		Title:  "Cryptodigger",
		Bounds: pixel.R(0, 0, WinWidth, WinHeight),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	winCenter := win.Bounds().Center()
	winTopLeft := pixel.V(50, WinHeight-50)

	_, backgroundSprite := LoadSprite(BackgroundSprite)
	_, buttonsSprite := LoadSprite(ButtonsSprite)
	popupPic, popupSprite := LoadSprite(PopupSprite)
	diggerPic, diggerSprite := LoadSprite(DiggerSprite)
	iconsPic, iconsSprite := LoadSprite(IconsSprite)
	coinAlphabet := LoadAlphabet(CoinsFont, 8, 144, text.ASCII)
	eventAlphabet := LoadAlphabet(CoinsFont, 4, 144, text.ASCII)

	world, digger := NewWorld(), NewDigger(100)

	camStart, camPos := winCenter, winCenter
	minView := pixel.V(0, win.Bounds().H()).Sub(camPos).ScaledXY(pixel.V(1, -1))
	maxView := pixel.V(win.Bounds().W(), 0).Sub(camPos).ScaledXY(pixel.V(1, -1))

	last := time.Now()

	diggerFrame := pixel.V(0, 2*FSize)
	step := time.Tick(time.Millisecond * 100)
	scale := pixel.V(1, 1)
	for !win.Closed() {
		win.Clear(colornames.Skyblue)

		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Moved(winCenter.Sub(camPos))
		deltaCam := camPos.Sub(camStart)
		invDeltaCam := deltaCam.ScaledXY(pixel.V(1, -1))
		// fmt.Printf("delta: %f\n", deltaCam)
		win.SetMatrix(cam)

		backgroundSprite.Draw(win, SM.Moved(camPos))

		minBlocks := minView.Add(invDeltaCam)
		maxBlocks := maxView.Add(invDeltaCam)
		blockDelta := pixel.V(Float64Mod(deltaCam.X, ASize), Float64Mod(deltaCam.Y, ASize))
		blocks := world.VisibleBlocks(minBlocks, maxBlocks)
		for v, block := range blocks {
			if block == nil {
				continue
			}
			v := v.ScaledXY(pixel.V(1, -1)).Add(camStart.Sub(blockDelta))
			blockMat := SM.Moved(v.Add(pixel.V(ASize, 0)))
			minX, minY := float64(block.Type)*FSize, 3*FSize
			maxX, maxY := float64(block.Type+1)*FSize, 4*FSize
			r := pixel.R(minX, minY, maxX, maxY)
			iconsSprite.Set(iconsPic, r)
			iconsSprite.Draw(win, blockMat)
		}

		diggerCell := CellFromVec(invDeltaCam)
		if !world.ContainsBlock(diggerCell.Down()) {
			camPos.Y -= ASize
			digger.Depth++
		}

		// Draw digger
		rect := pixel.Rect{diggerFrame, diggerFrame.Add(pixel.V(FSize, FSize))}
		diggerSprite.Set(diggerPic, rect)
		dMat := SM.ScaledXY(pixel.ZV, scale).Moved(pixel.V(camPos.X, camPos.Y))
		diggerSprite.Draw(win, dMat)

		// Draw coin count icon
		topLeft := winTopLeft.Add(deltaCam)
		// Coin icon
		iconsSprite.Set(iconsPic, pixel.R(0, 2*FSize, FSize, 3*FSize))
		iconsSprite.Draw(win, SM.Moved(topLeft))
		// Coin count
		str, col := strconv.Itoa(digger.Coins), colornames.White
		topLeftText := topLeft.Add(pixel.V(0.6*ASize, -0.3*ASize))
		coinAlphabet.Draw(win, str, col, SM.Moved(topLeftText))

		if win.Pressed(pixelgl.KeyD) {
			diggerFrame.Y = 2 * FSize
			scale = pixel.V(1, 1)
			select {
			case <-step:
				diggerFrame.X = float64(int(diggerFrame.X+FSize) % int(4*FSize))
				if !world.ContainsBlock(diggerCell.Right(invDeltaCam)) {
					camPos.X += CamSpeed * dt
				}

			default:
			}
		} else if win.Pressed(pixelgl.KeyA) {
			diggerFrame.Y = 2 * FSize
			scale = pixel.V(-1, 1)
			select {
			case <-step:
				diggerFrame.X = float64(int(diggerFrame.X+FSize) % int(4*FSize))
				if !world.ContainsBlock(diggerCell.Left(invDeltaCam)) {
					camPos.X -= CamSpeed * dt
				}
			default:
			}
		} else if win.Pressed(pixelgl.MouseButtonLeft) {
			diggerFrame.Y = 1 * FSize
			select {
			case <-step:
				diggerFrame.X = float64(int(diggerFrame.X+FSize) % int(2*FSize))
				if scale.X < 0 {
					digger.DigCell(world, diggerCell.Left(invDeltaCam))
				} else {
					digger.DigCell(world, diggerCell.Right(invDeltaCam))
				}
			default:
			}
		} else if win.Pressed(pixelgl.MouseButtonRight) {
			diggerFrame.Y = 1 * FSize
			select {
			case <-step:
				diggerFrame.X = float64(int(diggerFrame.X+FSize)%int(2*FSize)) + 2*FSize
				digger.DigCell(world, diggerCell.Down())
			default:
			}
		} else if win.JustReleased(pixelgl.KeyA) ||
			win.JustReleased(pixelgl.KeyD) ||
			win.JustReleased(pixelgl.MouseButtonLeft) ||
			win.JustReleased(pixelgl.MouseButtonRight) {
			diggerFrame = pixel.V(0, 2.0*FSize)
		}

		win.Update()
	}
}
