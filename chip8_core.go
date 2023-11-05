package main

// Chip8Core represents the chip8Core of the Chip-8 machine.
// It contains all the necessary components to emulate a Chip-8 system, including Memory,
// registers, counters, the call stack, and timers.
type Chip8Core struct {
	Memory     [4096]byte   // Memory represents the total Memory of the machine, comprising 4096 bytes.
	V          [16]byte     // V contains the 16 general-purpose registers, named from V0 to VF.
	I          uint16       // I is the index register used for store and load operations.
	PC         uint16       // PC is the program counter that points to the current location in Memory from where to read the next instruction.
	Stack      [16]uint16   // Stack is the call stack that holds return addresses when subroutines are called.
	SP         uint16       // SP is the stack pointer that points to the top of the call stack.
	Keys       [16]bool     // Keys represents the state of the 16-key hex keyboard, where each index corresponds to a specific key.
	Screen     [32][64]bool // Screen represents the current state of the display, with a resolution of 64x32 pixels.
	DelayTimer byte         // DelayTimer is the delay timer that is decremented at a frequency of 60Hz when it's non-zero.
	SoundTimer byte         // SoundTimer is the sound timer that is decremented at a frequency of 60Hz when it's non-zero.
}

func NewChip8Core() *Chip8Core {
	chip8Core := &Chip8Core{}
	sprites := []byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
	chip8Core.PC = 0x200
	chip8Core.SP = 0
	copy(chip8Core.Memory[0x000:], sprites)
	return chip8Core
}

func (chip8Core *Chip8Core) LoadROM(data []byte) {
	for i, d := range data {
		chip8Core.Memory[0x200+i] = d
	}
}

func (chip8Core *Chip8Core) Start() {}

func (chip8Core *Chip8Core) Stop() {}

func (chip8Core *Chip8Core) FetchOpcode() uint16 {
	return uint16(chip8Core.Memory[chip8Core.PC])<<8 | uint16(chip8Core.Memory[chip8Core.PC+1])
}

func (chip8Core *Chip8Core) UpdateTimers() {
	if chip8Core.DelayTimer > 0 {
		chip8Core.DelayTimer--
	}
	if chip8Core.SoundTimer > 0 {
		chip8Core.SoundTimer--
	}
}

func (chip8Core *Chip8Core) SetKey(index uint8, value bool) {
	chip8Core.Keys[index] = value
}

func (chip8Core *Chip8Core) GetKey(index uint8) bool {
	return chip8Core.Keys[index]
}

func (chip8Core *Chip8Core) SetPixel(x uint8, y uint8, value bool) {
	chip8Core.Screen[y][x] = value
}

func (chip8Core *Chip8Core) GetPixel(x uint8, y uint8) bool {
	return chip8Core.Screen[y][x]
}

func (chip8Core *Chip8Core) SetRegister(index uint8, value byte) {
	chip8Core.V[index] = value
}

func (chip8Core *Chip8Core) GetRegister(index uint8) byte {
	return chip8Core.V[index]
}

func (chip8Core *Chip8Core) SetI(value uint16) {
	chip8Core.I = value
}

func (chip8Core *Chip8Core) GetI() uint16 {
	return chip8Core.I
}

func (chip8Core *Chip8Core) IncrementPC(value uint16) {
	chip8Core.PC += value
}

func (chip8Core *Chip8Core) SetPC(value uint16) {
	chip8Core.PC = value
}

func (chip8Core *Chip8Core) GetPC() uint16 {
	return chip8Core.PC
}

func (chip8Core *Chip8Core) GetSP() uint16 {
	return chip8Core.SP
}

func (chip8Core *Chip8Core) PushStack(value uint16) {
	chip8Core.Stack[chip8Core.SP] = value
	chip8Core.SP++
}

func (chip8Core *Chip8Core) PopStack() uint16 {
	chip8Core.SP--
	return chip8Core.Stack[chip8Core.SP]
}

func (chip8Core *Chip8Core) ClearScreen() {
	for i := 0; i < 32; i++ {
		for j := 0; j < 64; j++ {
			chip8Core.Screen[i][j] = false
		}
	}
}
