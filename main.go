package main

func main() {

	// Initialize components
	chip8Core := NewChip8Core()
	opcodeDecoder := NewOpcodeDecoder()
	clock := NewFixedClock()
	//display := NewDisplay()
	//input := NewInput()

	chip8Core.LoadROM([]byte("roms/PONG"))

	// Start CPU and peripherals
	chip8Core.Start()
	clock.Start()
	//display.Start()
	//input.Start()

	// Main execution loop
	for {
		select {
		case <-clock.Tick():
			// Fetch: Retrieve the opcode at the current address from the program counter (PC)
			opcode := chip8Core.FetchOpcode()

			// Decode: Decode the opcode into an instruction
			instruction := opcodeDecoder.Decode(opcode)

			// Execute: Execute the instruction
			instruction.Execute(chip8Core)

			// Update: Update the state of the CPU and peripherals
			chip8Core.UpdateTimers()
			//display.Draw()
			//input.Poll()

		default:
			// Wait for the next tick

		}
	}
	chip8Core.Stop()
	clock.Stop()
	//display.Stop()
	//input.Stop()
}
