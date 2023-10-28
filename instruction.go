package main

import (
	"math/rand"
)

type Instruction interface {
	Execute(core *Chip8Core)
}

type GenericInstruction struct {
	opcode uint16
}

type ClearScreen struct {
	GenericInstruction
}

func (instruction *ClearScreen) Execute(core *Chip8Core) {
	for y := 0; y < len(core.Screen); y++ {
		for x := 0; x < len(core.Screen[y]); x++ {
			core.Screen[y][x] = false
		}
	}
}

type ReturnFromSubroutine struct {
	GenericInstruction
}

func (instruction *ReturnFromSubroutine) Execute(core *Chip8Core) {
	if core.GetSP() == 0 {
		return
	}
	core.SetPC(core.PopStack())
}

type JumpToAddress struct {
	GenericInstruction
}

func (instruction *JumpToAddress) Execute(core *Chip8Core) {
	core.SetPC(instruction.opcode & 0x0FFF)
}

type CallSubroutine struct {
	GenericInstruction
}

func (instruction *CallSubroutine) Execute(core *Chip8Core) {
	core.PushStack(core.PC)
	core.SetPC(instruction.opcode & 0x0FFF)
}

type SkipIfVxEqual struct {
	GenericInstruction
}

func (instruction *SkipIfVxEqual) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	constantValue := uint8(instruction.opcode & 0x00FF)
	if core.GetRegister(registerIndex) == constantValue {
		core.IncrementPC(2)
	}
}

type SkipIfVxNotEqual struct {
	GenericInstruction
}

func (instruction *SkipIfVxNotEqual) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	constantValue := uint8(instruction.opcode & 0x00FF)
	if core.GetRegister(registerIndex) != constantValue {
		core.IncrementPC(2)
	}
}

type SkipIfVxVyEqual struct {
	GenericInstruction
}

func (instruction *SkipIfVxVyEqual) Execute(core *Chip8Core) {
	registerIndexX := uint8((instruction.opcode & 0x0F00) >> 8)
	registerIndexY := uint8((instruction.opcode & 0x00F0) >> 4)

	if core.GetRegister(registerIndexX) == core.GetRegister(registerIndexY) {
		core.IncrementPC(2)
	}
}

type SetVx struct {
	GenericInstruction
}

func (instruction *SetVx) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	value := uint8(instruction.opcode & 0x00FF)
	core.SetRegister(registerIndex, value)
}

type AddToVx struct {
	GenericInstruction
}

func (instruction *AddToVx) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	valueToAdd := uint8(instruction.opcode & 0x00FF)
	core.SetRegister(registerIndex, core.GetRegister(registerIndex)+valueToAdd)
}

type SetVxVy struct {
	GenericInstruction
}

func (instruction *SetVxVy) Execute(core *Chip8Core) {
	registerIndexX := uint8((instruction.opcode & 0x0F00) >> 8)
	registerIndexY := uint8((instruction.opcode & 0x00F0) >> 4)
	core.SetRegister(registerIndexX, core.GetRegister(registerIndexY))
}

type SetVxOrVy struct {
	GenericInstruction
}

func (instruction *SetVxOrVy) Execute(core *Chip8Core) {
	registerIndexX := uint8((instruction.opcode & 0x0F00) >> 8)
	registerIndexY := uint8((instruction.opcode & 0x00F0) >> 4)
	core.SetRegister(registerIndexX, core.GetRegister(registerIndexX)|core.GetRegister(registerIndexY))
}

type SetVxAndVy struct {
	GenericInstruction
}

func (instruction *SetVxAndVy) Execute(core *Chip8Core) {
	registerIndexX := uint8((instruction.opcode & 0x0F00) >> 8)
	registerIndexY := uint8((instruction.opcode & 0x00F0) >> 4)
	core.SetRegister(registerIndexX, core.GetRegister(registerIndexX)&core.GetRegister(registerIndexY))
}

type SetVxXorVy struct {
	GenericInstruction
}

func (instruction *SetVxXorVy) Execute(core *Chip8Core) {
	registerIndexX := uint8((instruction.opcode & 0x0F00) >> 8)
	registerIndexY := uint8((instruction.opcode & 0x00F0) >> 4)
	core.SetRegister(registerIndexX, core.GetRegister(registerIndexX)^core.GetRegister(registerIndexY))
}

type AddVyToVx struct {
	GenericInstruction
}

func (instruction *AddVyToVx) Execute(core *Chip8Core) {
	registerIndexX := uint8((instruction.opcode & 0x0F00) >> 8)
	registerIndexY := uint8((instruction.opcode & 0x00F0) >> 4)
	sum := uint16(core.GetRegister(registerIndexX)) + uint16(core.GetRegister(registerIndexY))
	if sum > 0xFF {
		core.SetRegister(0xF, 1)
	} else {
		core.SetRegister(0xF, 0)
	}
	core.SetRegister(registerIndexX, uint8(sum&0xFF))
}

type SubtractVyFromVx struct {
	GenericInstruction
}

func (instruction *SubtractVyFromVx) Execute(core *Chip8Core) {
	registerIndexX := uint8((instruction.opcode & 0x0F00) >> 8)
	registerIndexY := uint8((instruction.opcode & 0x00F0) >> 4)
	if core.GetRegister(registerIndexX) >= core.GetRegister(registerIndexY) {
		core.SetRegister(0xF, 1)
	} else {
		core.SetRegister(0xF, 0)
	}
	core.SetRegister(registerIndexX, core.GetRegister(registerIndexX)-core.GetRegister(registerIndexY))
}

type ShiftVxRight struct {
	GenericInstruction
}

func (instruction *ShiftVxRight) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	core.SetRegister(0xF, core.GetRegister(registerIndex)&0x01)
	core.SetRegister(registerIndex, core.GetRegister(registerIndex)>>1)
}

type SetVxVyMinusVx struct {
	GenericInstruction
}

func (instruction *SetVxVyMinusVx) Execute(core *Chip8Core) {
	registerIndexX := uint8((instruction.opcode & 0x0F00) >> 8)
	registerIndexY := uint8((instruction.opcode & 0x00F0) >> 4)
	result := int(core.GetRegister(registerIndexY)) - int(core.GetRegister(registerIndexX))
	if result >= 0 {
		core.SetRegister(0xF, 1)
	} else {
		core.SetRegister(0xF, 0)
		result += 256
	}
	core.SetRegister(registerIndexX, uint8(result&0xFF))
}

type ShiftVxLeft struct {
	GenericInstruction
}

func (instruction *ShiftVxLeft) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	originalValue := core.GetRegister(registerIndex)
	core.SetRegister(registerIndex, originalValue<<1)
	core.SetRegister(0xF, (originalValue&0x80)>>7)
}

type SkipIfVxVyNotEqual struct {
	GenericInstruction
}

func (instruction *SkipIfVxVyNotEqual) Execute(core *Chip8Core) {
	registerIndexX := uint8((instruction.opcode & 0x0F00) >> 8)
	registerIndexY := uint8((instruction.opcode & 0x00F0) >> 4)
	valueVx := core.GetRegister(registerIndexX)
	valueVy := core.GetRegister(registerIndexY)
	if valueVx != valueVy {
		core.IncrementPC(4)
	} else {
		core.IncrementPC(2)
	}
}

type SetI struct {
	GenericInstruction
}

func (instruction *SetI) Execute(core *Chip8Core) {
	value := instruction.opcode & 0x0FFF
	core.SetI(value)
	core.IncrementPC(2)
}

type JumpToAddressPlusV0 struct {
	GenericInstruction
}

func (instruction *JumpToAddressPlusV0) Execute(core *Chip8Core) {
	address := instruction.opcode & 0x0FFF
	jumpAddress := address + uint16(core.GetRegister(0))
	core.PC = jumpAddress
}

type SetVxRandom struct {
	GenericInstruction
}

func (instruction *SetVxRandom) Execute(core *Chip8Core) {
	randomByte := uint8(rand.Intn(256))
	constant := uint8(instruction.opcode & 0x00FF)
	result := randomByte & constant
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	core.SetRegister(registerIndex, result)
}

type DrawSprite struct {
	GenericInstruction
}

func (instruction *DrawSprite) Execute(core *Chip8Core) {
	registerIndexX := uint8((instruction.opcode & 0x0F00) >> 8)
	registerIndexY := uint8((instruction.opcode & 0x00F0) >> 4)
	height := instruction.opcode & 0x000F
	core.SetRegister(0xF, 0)
	for row := uint16(0); row < height; row++ {
		spriteData := core.Memory[core.GetI()+row]
		for column := uint16(0); column < 8; column++ {
			positionX := (core.GetRegister(registerIndexX) + uint8(column)) % 64
			positionY := (core.GetRegister(registerIndexY) + uint8(row)) % 32
			pixel := spriteData & (0x80 >> column)
			if pixel != 0 {
				core.Screen[positionY][positionX] = true
				if core.Screen[positionY][positionX] == false {
					core.SetRegister(0xF, 1)
				}
			}
		}
	}
}

type SkipIfKeyPressed struct {
	GenericInstruction
}

func (instruction *SkipIfKeyPressed) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	key := core.GetRegister(registerIndex)
	if core.Keys[key] == true {
		core.IncrementPC(4)
	} else {
		core.IncrementPC(2)
	}
}

type SkipIfKeyNotPressed struct {
	GenericInstruction
}

func (instruction *SkipIfKeyNotPressed) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	key := core.GetRegister(registerIndex)
	if core.Keys[key] == false {
		core.IncrementPC(4)
	} else {
		core.IncrementPC(2)
	}
}

type SetVxDelayTimer struct {
	GenericInstruction
}

func (instruction *SetVxDelayTimer) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	core.DelayTimer = core.GetRegister(registerIndex)
	core.IncrementPC(2)
}

type WaitForKeyPress struct {
	GenericInstruction
}

func (instruction *WaitForKeyPress) Execute(core *Chip8Core) {
	keyPressed := false
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	for index, isPressed := range core.Keys {
		if isPressed {
			core.SetRegister(registerIndex, byte(index))
			keyPressed = true
			core.Keys[index] = false
			break
		}
	}
	if !keyPressed {
		return
	}
	core.IncrementPC(2)
}

type SetDelayTimer struct {
	GenericInstruction
}

func (instruction *SetDelayTimer) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	core.DelayTimer = core.GetRegister(registerIndex)
	core.IncrementPC(2)
}

type SetSoundTimer struct {
	GenericInstruction
}

func (instruction *SetSoundTimer) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	core.SoundTimer = core.GetRegister(registerIndex)
	core.IncrementPC(2)
}

type SetIPlusVx struct {
	GenericInstruction
}

func (instruction *SetIPlusVx) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	core.SetI(uint16(core.GetRegister(registerIndex)))
	core.IncrementPC(2)
}

type SetISprite struct {
	GenericInstruction
}

func (instruction *SetISprite) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	core.SetI(uint16(core.GetRegister(registerIndex) * 5))
	core.IncrementPC(2)
}

type StoreBCD struct {
	GenericInstruction
}

func (instruction *StoreBCD) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	vx := core.GetRegister(registerIndex)
	hundreds := vx / 100
	tens := (vx % 100) / 10
	ones := vx % 10
	core.Memory[core.GetI()] = hundreds
	core.Memory[core.GetI()+1] = tens
	core.Memory[core.GetI()+2] = ones
	core.IncrementPC(2)
}

type Storeregisters struct {
	GenericInstruction
}

func (instruction *Storeregisters) Execute(core *Chip8Core) {
	x := (instruction.opcode & 0x0F00) >> 8
	for j := uint16(0); j <= x; j++ {
		core.Memory[core.GetI()+j] = core.GetRegister(j)
	}
	// core.SetI( core.GetI() + x + 1)
	core.IncrementPC(2)
}

type Fillregisters struct {
	GenericInstruction
}

func (instruction *Fillregisters) Execute(core *Chip8Core) {
	x := (instruction.opcode & 0x0F00) >> 8
	for j := uint16(0); j <= x; j++ {
		core.SetRegister(j, core.Memory[core.GetI()+j])
	}
	// core.SetI( core.GetI() + x + 1)
	core.IncrementPC(2)
}

// UnknownInstruction handles any unrecognized opcodes
type UnknownInstruction struct {
	GenericInstruction
}

func (instruction *UnknownInstruction) Execute(core *Chip8Core) {
	core.IncrementPC(2)
}
