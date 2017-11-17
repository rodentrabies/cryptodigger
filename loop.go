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
		Bounds: pixel.R(0, 0, 768, 432),
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

	cameraPos := winCenter
	cameraSpeed := 500.0
	last := time.Now()

	diggerPos := pixel.ZV

	diggerX, diggerY := 0, 2*FSize
	step := time.Tick(time.Millisecond * 100)
	scale := pixel.V(4, 4)
	for !win.Closed() {
		diggerPos = cameraPos.Sub(winCenter)

		win.Clear(colornames.Skyblue)

		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Moved(winCenter.Sub(cameraPos))
		win.SetMatrix(cam)

		blocks := VisibleBlocks(
			world,
			pixel.V(diggerPos.X-winCenter.X/2, diggerPos.Y-winCenter.Y/2),
			pixel.V(diggerPos.X+winCenter.X/2, diggerPos.Y+winCenter.Y/2),
		)
		fmt.Printf("blocks: %v\n", blocks)

		plax.Draw(win, pixel.IM.Scaled(pixel.ZV, 4).Moved(cameraPos))
		rect := pixel.R(
			float64(diggerX), float64(diggerY),
			float64(diggerX+FSize), float64(diggerY+FSize),
		)
		diggerSprite.Set(diggerPic, rect)

		goldenRatio := pixel.V(cameraPos.X, win.Bounds().H()/2.61)
		mat := pixel.IM.ScaledXY(pixel.ZV, scale).Moved(goldenRatio)

		diggerSprite.Draw(win, mat)
		for v, block := range blocks {
			if block == nil {
				continue
			}
			v := pixel.V(v.X+winCenter.X, v.Y)
			blockMat := pixel.IM.Scaled(pixel.ZV, FScale).Moved(v)
			minX, minY, maxX, maxY := block.Type*FSize, 3*FSize, (block.Type+1)*FSize, 4*FSize
			r := pixel.R(float64(minX), float64(minY), float64(maxX), float64(maxY))
			blockSprite.Set(blockPic, r)
			blockSprite.Draw(win, blockMat)
		}

		if win.Pressed(pixelgl.KeyD) {
			diggerY = 2 * FSize
			scale = pixel.V(FScale, FScale)
			select {
			case <-step:
				diggerX = (diggerX + FSize) % (4 * FSize)
				cameraPos.X += cameraSpeed * dt
			default:
			}
		} else if win.Pressed(pixelgl.KeyA) {
			diggerY = 2 * FSize
			scale = pixel.V(-FScale, FScale)
			select {
			case <-step:
				diggerX = (diggerX + FSize) % (4 * FSize)
				cameraPos.X -= cameraSpeed * dt
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

func VisibleBlocks(world World, min pixel.Vec, max pixel.Vec) map[pixel.Vec]*Block {
	res := make(map[pixel.Vec]*Block)

	view := world.GridView(
		Cell{X: int(min.X)/ASize - 1, Y: int(min.Y)/ASize - 1},
		Cell{X: int(max.X)/ASize + 1, Y: int(max.Y)/ASize + 1},
	)

	fmt.Printf("\nview:\n")
	for i := 0; i < len(view); i++ {
		for j := 0; j < len(view[i]); j++ {
			if view[i][j] == nil {
				fmt.Printf("_")
			} else {
				fmt.Printf("%d", view[i][j].Type)
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")

	for i := 0; i < len(view); i++ {
		for j := 0; j < len(view[i]); j++ {
			res[pixel.V(float64(j*ASize)+min.X, float64(i*ASize)+min.Y)] = view[i][j]
		}
	}

	return res
}
