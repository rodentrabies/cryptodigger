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

	var newEvent Event
	var pendingEvent *EventBox
	buttons := make(map[*Button]func())

	for !win.Closed() {
		win.Clear(colornames.White)

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

		topLeft := winTopLeft.Add(deltaCam)

		// Coin icon
		iconsSprite.Set(iconsPic, pixel.R(0, 2*FSize, FSize, 3*FSize))
		iconsSprite.Draw(win, SM.Moved(topLeft))
		// Coin count
		coinStr, depthStr := strconv.Itoa(digger.Coins), strconv.Itoa(digger.Depth)
		coinTextPos := topLeft.Add(pixel.V(0.6*ASize, -0.3*ASize))
		depthTextPos := coinTextPos.Add(pixel.V(0, -ASize))
		coinAlphabet.Draw(win, coinStr, colornames.White, SM.Moved(coinTextPos))
		coinAlphabet.Draw(win, depthStr, colornames.White, SM.Moved(depthTextPos))

		// When digger has zero ballance, game ends
		if digger.Coins < 0 {
			win.Clear(colornames.Black)
			str, col := "You're broke", colornames.White
			coinAlphabet.Draw(win, str, col, SM.Moved(topLeft))
			win.Update()
			continue
		}

		// If there was some event, wait until player closes it
		if newEvent != nil {
			digger.Coins = newEvent.Consequence(digger.Coins)
			px, py := popupPic.Bounds().Max.XY()
			textPos := camPos.Add(pixel.V((-px*FScale+ASize)/2, py*FScale/2-ASize))

			buttonP := camPos.Add(pixel.V(px*FScale/2-20, py*FScale/2-20))
			buttonUpFrame := pixel.R(0, FSize*3, FSize, FSize*4)
			buttonDownFrame := pixel.R(FSize, FSize*3, 2*FSize, FSize*4)
			button := NewButton(buttonsSprite, buttonUpFrame, buttonDownFrame,
				SM.Moved(buttonP))
			button.Register(buttons, func() {
				button.Unregister(buttons)
				pendingEvent = nil
			})
			newEvent, pendingEvent = nil, &EventBox{
				message:  newEvent.Description(),
				position: textPos,
				button:   button,
			}
		}

		if pendingEvent != nil {
			popupSprite.Draw(win, SM.Moved(camPos))
			str, col := pendingEvent.message, colornames.Black
			eventAlphabet.Draw(win, str, col, SM.Moved(pendingEvent.position))
			pendingEvent.button.Draw(win)
			if win.Pressed(pixelgl.MouseButtonLeft) {
				for button, _ := range buttons {
					if button.Within(win.MousePosition().Add(deltaCam)) {
						button.Push()
					} else {
						button.Release()
					}
				}
			}

			if win.JustReleased(pixelgl.MouseButtonLeft) {
				for button, actionFunc := range buttons {
					if button.Within(win.MousePosition().Add(deltaCam)) {
						actionFunc()
					}
				}
			}
			win.Update()
			continue
		}

		if win.Pressed(pixelgl.MouseButtonLeft) && win.Pressed(pixelgl.KeyA) {
			diggerFrame.Y = 1 * FSize
			select {
			case <-step:
				diggerFrame.X = float64(int(diggerFrame.X+FSize) % int(2*FSize))
				if int(diggerFrame.X) == 0 {
					digger.Coins--
				}
				newEvent = digger.DigCell(world, diggerCell.Left(invDeltaCam))
			default:
			}
		} else if win.Pressed(pixelgl.MouseButtonLeft) && win.Pressed(pixelgl.KeyD) {
			diggerFrame.Y = 1 * FSize
			select {
			case <-step:
				diggerFrame.X = float64(int(diggerFrame.X+FSize) % int(2*FSize))
				if int(diggerFrame.X) == 0 {
					digger.Coins--
				}
				newEvent = digger.DigCell(world, diggerCell.Right(invDeltaCam))
			default:
			}
		} else if win.Pressed(pixelgl.MouseButtonLeft) && win.Pressed(pixelgl.KeyS) {
			diggerFrame.Y = 1 * FSize
			select {
			case <-step:

				diggerFrame.X = float64(int(diggerFrame.X+FSize)%int(2*FSize)) + 2*FSize
				if int(diggerFrame.X) == int(2*FSize) {
					digger.Coins--
				}
				newEvent = digger.DigCell(world, diggerCell.Down())
			default:
			}
		} else if win.Pressed(pixelgl.KeyD) {
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
		} else if win.JustReleased(pixelgl.KeyA) ||
			win.JustReleased(pixelgl.KeyD) ||
			win.JustReleased(pixelgl.MouseButtonLeft) ||
			win.JustReleased(pixelgl.MouseButtonRight) {
			diggerFrame = pixel.V(0, 2.0*FSize)
		}

		win.Update()
	}
}
