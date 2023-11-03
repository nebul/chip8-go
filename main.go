package main

import (
	"image"
	"image/color"
	"image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image"
	"image/color"
	"image/draw"
)

func main() {

	chip8Core := NewChip8Core()
	opcodeDecoder := NewOpcodeDecoder()
	clock := NewFixedClock()

	application := app.New()
	window := application.NewWindow("Chip-8 Emulator")
	window.Resize(fyne.NewSize(640, 320))

	displayCanvas := canvas.NewRaster(func(w, h int) image.Image {
		pixel := image.NewNRGBA(image.Rect(0, 0, 64*10, 32*10))
		for y := 0; y < 32; y++ {
			for x := 0; x < 64; x++ {
				column := color.White
				if chip8Core.GetPixel(uint8(x), uint8(y)) {
					column = color.Black
				}
				draw.Draw(pixel, image.Rect(x*10, y*10, (x+1)*10, (y+1)*10), &image.Uniform{column}, image.Point{}, draw.Src)
			}
		}
		return pixel
	})

	window.SetContent(container.NewWithoutLayout(displayCanvas))

	window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		var keyMap = map[fyne.KeyName]byte{
			fyne.KeyX: 0x0,
			fyne.Key1: 0x1,
			fyne.Key2: 0x2,
			fyne.Key3: 0x3,
			fyne.KeyQ: 0x4,
			fyne.KeyW: 0x5,
			fyne.KeyE: 0x6,
			fyne.KeyA: 0x7,
			fyne.KeyS: 0x8,
			fyne.KeyD: 0x9,
			fyne.KeyZ: 0xA,
			fyne.KeyC: 0xB,
			fyne.Key4: 0xC,
			fyne.KeyR: 0xD,
			fyne.KeyF: 0xE,
			fyne.KeyV: 0xF,
		}
		if chipKey, exists := keyMap[key.Name]; exists {
			chip8Core.SetKey(chipKey, true)
		}
	})

	updateDisplay := func() {
		window.Canvas().Refresh(displayCanvas)
	}

	chip8Core.LoadROM([]byte("roms/PONG"))

	chip8Core.Start()
	clock.Start()

	go func() {
		for {
			select {
			case <-clock.Tick():
				opcode := chip8Core.FetchOpcode()
				instruction := opcodeDecoder.Decode(opcode)
				instruction.Execute(chip8Core)
				chip8Core.UpdateTimers()
				updateDisplay()

			default:
			}
		}
	}()

	window.ShowAndRun()

	chip8Core.Stop()
	clock.Stop()
}
