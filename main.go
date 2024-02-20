package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("CHIP-8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 640, 320, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	chip8Core := NewChip8Core()
	opcodeDecoder := NewOpcodeDecoder()
	clock := NewFixedClock()

	// TODO - Improvement Load ROM another way
	data, err := os.ReadFile("roms/PONG")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	chip8Core.LoadROM(data)
	chip8Core.Start()
	defer chip8Core.Stop()

	clock.Start()
	defer clock.Stop()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				keyMap := map[sdl.Keycode]uint8{
					sdl.K_1: 0x1, sdl.K_2: 0x2, sdl.K_3: 0x3, sdl.K_4: 0xC,
					sdl.K_q: 0x4, sdl.K_w: 0x5, sdl.K_e: 0x6, sdl.K_r: 0xD,
					sdl.K_a: 0x7, sdl.K_s: 0x8, sdl.K_d: 0x9, sdl.K_f: 0xE,
					sdl.K_z: 0xA, sdl.K_x: 0x0, sdl.K_c: 0xB, sdl.K_v: 0xF,
				}
				if keyIndex, exists := keyMap[e.Keysym.Sym]; exists {
					chip8Core.SetKey(keyIndex, e.Type == sdl.KEYDOWN)
				}
			}
		}
		<-clock.Tick()
		opcode := chip8Core.FetchOpcode()
		instruction := opcodeDecoder.Decode(opcode)
		instruction.Execute(chip8Core)
		chip8Core.UpdateTimers()

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		pixelSize := int32(10)

		for positionY := 0; positionY < 32; positionY++ {
			for positionX := 0; positionX < 64; positionX++ {
				color := sdl.Color{R: 0, G: 0, B: 0, A: 255}
				if chip8Core.GetPixel(uint8(positionX), uint8(positionY)) {
					color = sdl.Color{R: 255, G: 255, B: 255, A: 255}
				}
				renderer.SetDrawColor(color.R, color.G, color.B, color.A)
				rectangle := sdl.Rect{X: int32(positionX) * pixelSize, Y: int32(positionY) * pixelSize, W: pixelSize, H: pixelSize}
				renderer.FillRect(&rectangle)
			}
		}
		renderer.Present()
	}
}
