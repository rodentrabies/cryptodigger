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

	// goldenRatio := win.Bounds().H() / 2.61
	cameraPos := winCenter
	cameraSpeed := 500.0
	last := time.Now()

	diggerPos := pixel.ZV

	diggerFrame := pixel.V(0, 2*FSize)
	step := time.Tick(time.Millisecond * 100)
	scale := pixel.V(4, 4)
	for !win.Closed() {
		diggerPos = cameraPos.Sub(winCenter)
		fmt.Printf("cam: %v\ndigger: %v\n", cameraPos, diggerPos)

		win.Clear(colornames.Skyblue)

		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Moved(winCenter.Sub(cameraPos))
		win.SetMatrix(cam)

		blocksMin := pixel.V(diggerPos.X-winCenter.X, diggerPos.Y-winCenter.Y)
		blocksMax := pixel.V(diggerPos.X+winCenter.X, diggerPos.Y+winCenter.Y)
		blocks := VisibleBlocks(world, blocksMin, blocksMax)

		plax.Draw(win, pixel.IM.Scaled(pixel.ZV, 4).Moved(cameraPos))
		rect := pixel.Rect{diggerFrame, diggerFrame.Add(pixel.V(FSize, FSize))}
		diggerSprite.Set(diggerPic, rect)

		mat := pixel.IM.ScaledXY(pixel.ZV, scale).Moved(pixel.V(cameraPos.X, cameraPos.Y))

		diggerSprite.Draw(win, mat)
		for v, block := range blocks {
			if block == nil {
				continue
			}
			v := pixel.V(v.X+winCenter.X, v.Y-winCenter.Y)
			blockMat := pixel.IM.Scaled(pixel.ZV, FScale).Moved(v)
			minX, minY := float64(block.Type)*FSize, 3*FSize
			maxX, maxY := float64(block.Type+1)*FSize, 4*FSize
			r := pixel.R(minX, minY, maxX, maxY)
			blockSprite.Set(blockPic, r)
			blockSprite.Draw(win, blockMat)
		}

		if win.Pressed(pixelgl.KeyD) {
			diggerFrame.Y = 2 * FSize
			scale = pixel.V(FScale, FScale)
			select {
			case <-step:
				diggerFrame.X = float64(int(diggerFrame.X+FSize) % int(4*FSize))
				cameraPos.X += cameraSpeed * dt
			default:
			}
		} else if win.Pressed(pixelgl.KeyA) {
			diggerFrame.Y = 2 * FSize
			scale = pixel.V(-FScale, FScale)
			select {
			case <-step:
				diggerFrame.X = float64(int(diggerFrame.X+FSize) % int(4*FSize))
				cameraPos.X -= cameraSpeed * dt
			default:
			}
		} else if win.Pressed(pixelgl.MouseButtonLeft) {
			diggerFrame.Y = 1 * FSize
			select {
			case <-step:
				diggerFrame.X = float64(int(diggerFrame.X+FSize) % int(2*FSize))
			default:
			}
		} else if win.Pressed(pixelgl.MouseButtonRight) {
			diggerFrame.Y = 1 * FSize
			select {
			case <-step:
				diggerFrame.X = float64(int(diggerFrame.X+FSize)%int(2*FSize)) + 2*FSize
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

func VisibleBlocks(world World, min pixel.Vec, max pixel.Vec) map[pixel.Vec]*Block {
	res := make(map[pixel.Vec]*Block)
	fmt.Printf("min: %v, max: %v\n", min, max)
	view := world.GridView(
		Cell{X: int(min.X)/ASize - 1, Y: int(min.Y)/ASize - 1},
		Cell{X: int(max.X)/ASize + 1, Y: int(max.Y)/ASize + 1},
	)

	fmt.Printf("\nview:\n")
	for i := 0; i < len(view); i++ {
		for j := 0; j < len(view[i]); j++ {
			x, y := float64(j*ASize)+min.X, float64(i*ASize)+min.Y

			if view[i][j] == nil {
				fmt.Printf(" ____  ____ ")
			} else {
				fmt.Printf("(%4d, %4d)", int(x), int(y))
			}

			res[pixel.V(x, y)] = view[i][j]
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")

	return res
}
