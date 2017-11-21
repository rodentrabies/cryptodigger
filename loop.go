package main

import (
	"math"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func Loop() {
	cfg := pixelgl.WindowConfig{
		Title: "Cryptodigger",
		// Bounds: pixel.R(0, 0, 896, 512),
		Bounds: pixel.R(0, 0, 768, 384),
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

	camStart, camPos, camSpeed := winCenter, winCenter, 500.0
	minView := pixel.V(0, win.Bounds().H()).Sub(camPos).ScaledXY(pixel.V(1, -1))
	maxView := pixel.V(win.Bounds().W(), 0).Sub(camPos).ScaledXY(pixel.V(1, -1))

	last := time.Now()

	diggerFrame := pixel.V(0, 2*FSize)
	step := time.Tick(time.Millisecond * 100)
	scale := pixel.V(4, 4)
	for !win.Closed() {
		win.Clear(colornames.Skyblue)

		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Moved(winCenter.Sub(camPos))
		deltaCam := camPos.Sub(camStart)
		win.SetMatrix(cam)

		plax.Draw(win, pixel.IM.Scaled(pixel.ZV, 4).Moved(camPos))

		minBlocks, maxBocks := minView.Add(deltaCam), maxView.Add(deltaCam)
		blocks := VisibleBlocks(world, minBlocks, maxBocks)
		for v, block := range blocks {
			if block == nil {
				continue
			}
			v := v.ScaledXY(pixel.V(1, -1)).Add(camStart)
			blockMat := pixel.IM.Scaled(pixel.ZV, FScale).Moved(v)
			minX, minY := float64(block.Type)*FSize, 3*FSize
			maxX, maxY := float64(block.Type+1)*FSize, 4*FSize
			r := pixel.R(minX, minY, maxX, maxY)
			blockSprite.Set(blockPic, r)
			blockSprite.Draw(win, blockMat)
		}

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
				camPos.X += camSpeed * dt
			default:
			}
		} else if win.Pressed(pixelgl.KeyA) {
			diggerFrame.Y = 2 * FSize
			scale = pixel.V(-FScale, FScale)
			select {
			case <-step:
				diggerFrame.X = float64(int(diggerFrame.X+FSize) % int(4*FSize))
				camPos.X -= camSpeed * dt
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
	minBlockX, minBlockY := int(math.Ceil(min.X/ASize)), int(math.Ceil(min.Y/ASize))
	maxBlockX := minBlockX + int(math.Ceil((max.X-min.X)/ASize)) + 1
	maxBlockY := minBlockY + int(math.Ceil((max.Y-min.Y)/ASize)) + 1
	view := world.GridView(Cell{X: minBlockX, Y: minBlockY}, Cell{X: maxBlockX, Y: maxBlockY})
	for i := 0; i < len(view); i++ {
		for j := 0; j < len(view[i]); j++ {
			x, y := float64(j*int(ASize))+min.X, float64(i*int(ASize))+min.Y
			res[pixel.V(x, y)] = view[i][j]
		}
	}
	return res
}
