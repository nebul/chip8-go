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
	core.ClearScreen()
	core.IncrementPC(2)
}

type ReturnFromSubroutine struct {
	GenericInstruction
}

func (instruction *ReturnFromSubroutine) Execute(core *Chip8Core) {
	if core.GetSP() == 0 {
		return
	}
	core.SetPC(core.PopStack())
	core.IncrementPC(2)
}

type JumpToAddress struct {
	GenericInstruction
}

func (instruction *JumpToAddress) Execute(core *Chip8Core) {
	address := instruction.opcode & 0x0FFF
	core.SetPC(address)
}

type CallSubroutine struct {
	GenericInstruction
}

func (instruction *CallSubroutine) Execute(core *Chip8Core) {
	address := instruction.opcode & 0x0FFF
	core.PushStack(core.PC)
	core.SetPC(address)
}

type SkipIfVxEqual struct {
	GenericInstruction
}

func (instruction *SkipIfVxEqual) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	value := uint8(instruction.opcode & 0x00FF)
	registerValue := core.GetRegister(registerIndex)
	if registerValue == value {
		core.IncrementPC(4)
	} else {
		core.IncrementPC(2)
	}
}

type SkipIfVxNotEqual struct {
	GenericInstruction
}

func (instruction *SkipIfVxNotEqual) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	value := uint8(instruction.opcode & 0x00FF)
	registerValue := core.GetRegister(registerIndex)
	if registerValue != value {
		core.IncrementPC(4)
	} else {
		core.IncrementPC(2)
	}
}

type SkipIfVxVyEqual struct {
	GenericInstruction
}

func (instruction *SkipIfVxVyEqual) Execute(core *Chip8Core) {
	xRegisterIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	yRegisterIndex := uint8((instruction.opcode & 0x00F0) >> 4)
	xRegisterValue := core.GetRegister(xRegisterIndex)
	yRegisterValue := core.GetRegister(yRegisterIndex)
	if xRegisterValue == yRegisterValue {
		core.IncrementPC(4)
	} else {
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
	core.IncrementPC(2)
}

type AddToVx struct {
	GenericInstruction
}

func (instruction *AddToVx) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	value := uint8(instruction.opcode & 0x00FF)
	registerValue := core.GetRegister(registerIndex)
	core.SetRegister(registerIndex, registerValue+value)
	core.IncrementPC(2)
}

type SetVxVy struct {
	GenericInstruction
}

func (instruction *SetVxVy) Execute(core *Chip8Core) {
	xRegisterIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	yRegisterIndex := uint8((instruction.opcode & 0x00F0) >> 4)
	yRegisterValue := core.GetRegister(yRegisterIndex)
	core.SetRegister(xRegisterIndex, yRegisterValue)
	core.IncrementPC(2)
}

type SetVxOrVy struct {
	GenericInstruction
}

func (instruction *SetVxOrVy) Execute(core *Chip8Core) {
	xRegisterIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	yRegisterIndex := uint8((instruction.opcode & 0x00F0) >> 4)
	xRegisterValue := core.GetRegister(xRegisterIndex)
	yRegisterValue := core.GetRegister(yRegisterIndex)
	core.SetRegister(xRegisterIndex, xRegisterValue|yRegisterValue)
	core.IncrementPC(2)
}

type SetVxAndVy struct {
	GenericInstruction
}

func (instruction *SetVxAndVy) Execute(core *Chip8Core) {
	xRegisterIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	yRegisterIndex := uint8((instruction.opcode & 0x00F0) >> 4)
	xRegisterValue := core.GetRegister(xRegisterIndex)
	yRegisterValue := core.GetRegister(yRegisterIndex)
	core.SetRegister(xRegisterIndex, xRegisterValue&yRegisterValue)
	core.IncrementPC(2)
}

type SetVxXorVy struct {
	GenericInstruction
}

func (instruction *SetVxXorVy) Execute(core *Chip8Core) {
	xRegisterIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	yRegisterIndex := uint8((instruction.opcode & 0x00F0) >> 4)
	xRegisterValue := core.GetRegister(xRegisterIndex)
	yRegisterValue := core.GetRegister(yRegisterIndex)
	core.SetRegister(xRegisterIndex, xRegisterValue^yRegisterValue)
	core.IncrementPC(2)
}

type AddVyToVx struct {
	GenericInstruction
}

func (instruction *AddVyToVx) Execute(core *Chip8Core) {
	xRegisterIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	yRegisterIndex := uint8((instruction.opcode & 0x00F0) >> 4)
	xRegisterValue := core.GetRegister(xRegisterIndex)
	yRegisterValue := core.GetRegister(yRegisterIndex)
	result := xRegisterValue + yRegisterValue
	if result > 0xFF {
		core.SetRegister(0xF, 1)
	} else {
		core.SetRegister(0xF, 0)
	}
	core.SetRegister(xRegisterIndex, uint8(result&0xFF))
	core.IncrementPC(2)
}

type SubtractVyFromVx struct {
	GenericInstruction
}

func (instruction *SubtractVyFromVx) Execute(core *Chip8Core) {
	xRegisterIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	yRegisterIndex := uint8((instruction.opcode & 0x00F0) >> 4)
	xRegisterValue := core.GetRegister(xRegisterIndex)
	yRegisterValue := core.GetRegister(yRegisterIndex)
	if xRegisterValue >= yRegisterValue {
		core.SetRegister(0xF, 1)
	} else {
		core.SetRegister(0xF, 0)
	}
	core.SetRegister(xRegisterIndex, xRegisterValue-yRegisterValue)
	core.IncrementPC(2)
}

type ShiftVxRight struct {
	GenericInstruction
}

func (instruction *ShiftVxRight) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	core.SetRegister(0xF, core.GetRegister(registerIndex)&0x01)
	core.SetRegister(registerIndex, core.GetRegister(registerIndex)>>1)
	core.IncrementPC(2)
}

type SetVxVyMinusVx struct {
	GenericInstruction
}

func (instruction *SetVxVyMinusVx) Execute(core *Chip8Core) {
	xRegisterIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	yRegisterIndex := uint8((instruction.opcode & 0x00F0) >> 4)
	xRegisterValue := core.GetRegister(xRegisterIndex)
	yRegisterValue := core.GetRegister(yRegisterIndex)

	result := yRegisterValue - xRegisterValue
	if result >= 0 {
		core.SetRegister(0xF, 1)
	} else {
		core.SetRegister(0xF, 0)
	}
	core.SetRegister(xRegisterIndex, uint8(result&0xFF))
	core.IncrementPC(2)
}

type ShiftVxLeft struct {
	GenericInstruction
}

func (instruction *ShiftVxLeft) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	originalValue := core.GetRegister(registerIndex)
	core.SetRegister(0xF, originalValue>>7)
	core.SetRegister(registerIndex, originalValue<<1)
	core.IncrementPC(2)
}

type SkipIfVxVyNotEqual struct {
	GenericInstruction
}

func (instruction *SkipIfVxVyNotEqual) Execute(core *Chip8Core) {
	xRegisterIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	yRegisterIndex := uint8((instruction.opcode & 0x00F0) >> 4)
	xRegisterValue := core.GetRegister(xRegisterIndex)
	yRegisterValue := core.GetRegister(yRegisterIndex)
	if xRegisterValue != yRegisterValue {
		core.IncrementPC(4)
	} else {
		core.IncrementPC(2)
	}
}

type SetI struct {
	GenericInstruction
}

func (instruction *SetI) Execute(core *Chip8Core) {
	iRegisterValue := instruction.opcode & 0x0FFF
	core.SetI(iRegisterValue)
	core.IncrementPC(2)
}

type JumpToAddressPlusV0 struct {
	GenericInstruction
}

func (instruction *JumpToAddressPlusV0) Execute(core *Chip8Core) {
	address := instruction.opcode & 0x0FFF
	jumpAddress := address + uint16(core.GetRegister(0))
	core.SetPC(jumpAddress)
}

type SetVxRandom struct {
	GenericInstruction
}

func (instruction *SetVxRandom) Execute(core *Chip8Core) {
	randomByte := uint8(rand.Intn(256))
	constant := uint8(instruction.opcode & 0x00FF)
	randomValue := randomByte & constant
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	core.SetRegister(registerIndex, randomValue)
	core.IncrementPC(2)
}

type DrawSprite struct {
	GenericInstruction
}

func (instruction *DrawSprite) Execute(core *Chip8Core) {
	xRegisterIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	yRegisterIndex := uint8((instruction.opcode & 0x00F0) >> 4)
	xRegisterValue := core.GetRegister(xRegisterIndex)
	yRegisterValue := core.GetRegister(yRegisterIndex)
	iRegisterValue := core.GetI()
	height := instruction.opcode & 0x000F
	core.SetRegister(0xF, 0)
	for row := uint16(0); row < height; row++ {
		spriteData := core.Memory[iRegisterValue+row]
		for column := uint16(0); column < 8; column++ {
			pixel := spriteData & (0x80 >> column)
			if pixel != 0 {
				positionX := (xRegisterValue + uint8(column)) % 64
				positionY := (yRegisterValue + uint8(row)) % 32
				if !core.GetPixel(positionX, positionY) {
					core.SetRegister(0xF, 1)
				}
				pixel := (core.GetPixel(positionX, positionY) || true) && !(core.GetPixel(positionX, positionY) && true)
				core.SetPixel(positionX, positionY, pixel)
			}
		}
	}
	core.IncrementPC(2)
}

type SkipIfKeyPressed struct {
	GenericInstruction
}

func (instruction *SkipIfKeyPressed) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	key := core.GetRegister(registerIndex)
	if core.Keys[key] {
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
	if !core.Keys[key] {
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
	registerValue := core.GetRegister(registerIndex)
	core.DelayTimer = registerValue
	core.IncrementPC(2)
}

type WaitForKeyPress struct {
	GenericInstruction
}

func (instruction *WaitForKeyPress) Execute(core *Chip8Core) {
	keyPressed := false
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	for keyIndex, isPressed := range core.Keys {
		if isPressed {
			core.SetRegister(registerIndex, byte(keyIndex))
			keyPressed = true
			core.Keys[keyIndex] = false
			break
		}
	}
	if !keyPressed {
		return
	} else {
		core.IncrementPC(2)
	}
}

type SetDelayTimer struct {
	GenericInstruction
}

func (instruction *SetDelayTimer) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	registerValue := core.GetRegister(registerIndex)
	core.DelayTimer = registerValue
	core.IncrementPC(2)
}

type SetSoundTimer struct {
	GenericInstruction
}

func (instruction *SetSoundTimer) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	registerValue := core.GetRegister(registerIndex)
	core.SoundTimer = registerValue
	core.IncrementPC(2)
}

type SetIPlusVx struct {
	GenericInstruction
}

func (instruction *SetIPlusVx) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	registerValue := core.GetRegister(registerIndex)
	core.SetI(uint16(registerValue))
	core.IncrementPC(2)
}

type SetISprite struct {
	GenericInstruction
}

func (instruction *SetISprite) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	registerValue := core.GetRegister(registerIndex)
	core.SetI(uint16(registerValue * 5))
	core.IncrementPC(2)
}

type StoreBCD struct {
	GenericInstruction
}

func (instruction *StoreBCD) Execute(core *Chip8Core) {
	registerIndex := uint8((instruction.opcode & 0x0F00) >> 8)
	registerValue := core.GetRegister(registerIndex)
	iRegisterValue := core.GetI()
	core.Memory[iRegisterValue] = registerValue / 10
	core.Memory[iRegisterValue+1] = (registerValue / 10) % 10
	core.Memory[iRegisterValue+2] = (registerValue / 100) % 10
	core.IncrementPC(2)
}

type Storeregisters struct {
	GenericInstruction
}

func (instruction *Storeregisters) Execute(core *Chip8Core) {
	registersNumber := (instruction.opcode & 0x0F00) >> 8
	iRegisterValue := core.GetI()

	for registerIndex := uint16(0); registerIndex <= registersNumber; registerIndex++ {
		core.Memory[iRegisterValue+registerIndex] = core.GetRegister(uint8(registerIndex))
	}
	core.SetI(iRegisterValue + registersNumber + 1)
	core.IncrementPC(2)
}

type Fillregisters struct {
	GenericInstruction
}

func (instruction *Fillregisters) Execute(core *Chip8Core) {
	registersNumber := (instruction.opcode & 0x0F00) >> 8
	iRegisterValue := core.GetI()

	for registerIndex := uint16(0); registerIndex <= registersNumber; registerIndex++ {
		core.SetRegister(uint8(registerIndex), core.Memory[iRegisterValue+registerIndex])
	}
	core.SetI(iRegisterValue + registersNumber + 1)
	core.IncrementPC(2)
}

type UnknownInstruction struct {
	GenericInstruction
}

func (instruction *UnknownInstruction) Execute(core *Chip8Core) {
	core.IncrementPC(2)
}
