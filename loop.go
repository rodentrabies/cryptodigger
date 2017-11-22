package main

import (
	"fmt"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

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

	blockPic, err := loadPicture(BlockSprite)
	if err != nil {
		panic(err)
	}
	blockSprite := pixel.NewSprite(blockPic, blockPic.Bounds())

	world := NewWorld()

	camStart, camPos := winCenter, winCenter
	minView := pixel.V(0, win.Bounds().H()).Sub(camPos).ScaledXY(pixel.V(1, -1))
	maxView := pixel.V(win.Bounds().W(), 0).Sub(camPos).ScaledXY(pixel.V(1, -1))

	last := time.Now()

	diggerFrame := pixel.V(0, 2*FSize)
	step := time.Tick(time.Millisecond * 100)
	scale := pixel.V(FScale, FScale)
	for !win.Closed() {
		win.Clear(colornames.Skyblue)

		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Moved(winCenter.Sub(camPos))
		deltaCam := camPos.Sub(camStart)
		// fmt.Printf("delta: %f\n", deltaCam)
		win.SetMatrix(cam)

		plax.Draw(win, pixel.IM.Scaled(pixel.ZV, FScale).Moved(camPos))

		minBlocks := minView.Add(deltaCam.ScaledXY(pixel.V(1, -1)))
		maxBlocks := maxView.Add(deltaCam.ScaledXY(pixel.V(1, -1)))
		blockDelta := pixel.V(Float64Mod(deltaCam.X, ASize), Float64Mod(deltaCam.Y, ASize))
		fmt.Printf("block delta: %f\n", blockDelta)
		blocks := world.VisibleBlocks(minBlocks, maxBlocks)
		for v, block := range blocks {
			if block == nil {
				continue
			}
			v := v.ScaledXY(pixel.V(1, -1)).Add(camStart.Sub(blockDelta))
			blockMat := pixel.IM.Scaled(pixel.ZV, FScale).Moved(v)
			minX, minY := float64(block.Type)*FSize, 3*FSize
			maxX, maxY := float64(block.Type+1)*FSize, 4*FSize
			r := pixel.R(minX, minY, maxX, maxY)
			blockSprite.Set(blockPic, r)
			blockSprite.Draw(win, blockMat)
		}

		diggerCell := CellFromVec(camPos.Sub(camStart).ScaledXY(pixel.V(1, -1)))
		if !world.ContainsBlock(diggerCell.Down()) {
			camPos.Y -= ASize
		}

		fmt.Printf("dcell: %v\n", diggerCell)

		rect := pixel.Rect{diggerFrame, diggerFrame.Add(pixel.V(FSize, FSize))}
		diggerSprite.Set(diggerPic, rect)
		mat := pixel.IM.ScaledXY(pixel.ZV, scale).Moved(pixel.V(camPos.X, camPos.Y))
		diggerSprite.Draw(win, mat)

		if win.Pressed(pixelgl.KeyD) {
			diggerFrame.Y = 2 * FSize
			scale = pixel.V(FScale, FScale)
			select {
			case <-step:
				diggerFrame.X = float64(int(diggerFrame.X+FSize) % int(4*FSize))
				if !world.ContainsBlock(diggerCell.Right()) {
					camPos.X += CamSpeed * dt
				}

			default:
			}
		} else if win.Pressed(pixelgl.KeyA) {
			diggerFrame.Y = 2 * FSize
			scale = pixel.V(-FScale, FScale)
			select {
			case <-step:
				diggerFrame.X = float64(int(diggerFrame.X+FSize) % int(4*FSize))
				if !world.ContainsBlock(diggerCell.Left()) {
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
					world.HammerBlock(diggerCell.Left())
				} else {
					world.HammerBlock(diggerCell.Right())
				}
			default:
			}
		} else if win.Pressed(pixelgl.MouseButtonRight) {
			diggerFrame.Y = 1 * FSize
			select {
			case <-step:
				diggerFrame.X = float64(int(diggerFrame.X+FSize)%int(2*FSize)) + 2*FSize
				world.HammerBlock(diggerCell.Down())
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
