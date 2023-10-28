package main

// OpcodeDecoder is responsible for decoding opcodes into instructions
type OpcodeDecoder struct {
}

func NewOpcodeDecoder() *OpcodeDecoder {
	return &OpcodeDecoder{}
}

// Decode takes an opcode and returns a corresponding Instruction
func (opcodeDecoder *OpcodeDecoder) Decode(opcode uint16) Instruction {
	switch opcode & 0xF000; { // Mask the high nibble to identify the instruction category
	case 0x0000: // Multi-purpose category, more decoding required
		switch opcode & 0x00FF; { // Mask the low byte for further decoding
		case 0x00E0: // Opcode for clearing the screen
			return &ClearScreen{GenericInstruction{opcode}}
		case 0x00EE: // Opcode for returning from subroutine
			return &ReturnFromSubroutine{GenericInstruction{opcode}}
		default: // Unrecognized 0x0??? opcode
			return &UnknownInstruction{opcode}
		}
	case 0x1000: // Opcode for jumping to address NNN
		return &JumpToAddress{GenericInstruction{opcode}}
	case 0x2000: // Opcode for calling subroutine at NNN
		return &CallSubroutine{GenericInstruction{opcode}}
	case 0x3000: // Opcode for skipping next instruction if Vx == NN
		return &SkipIfVxEqual{GenericInstruction{opcode}}
	case 0x4000: // Opcode for skipping next instruction if Vx != NN
		return &SkipIfVxNotEqual{GenericInstruction{opcode}}
	case 0x5000: // Opcode for skipping next instruction if Vx == Vy
		return &SkipIfVxVyEqual{GenericInstruction{opcode}}
	case 0x6000: // Opcode for setting Vx = NN
		return &SetVx{GenericInstruction{opcode}}
	case 0x7000: // Opcode for adding NN to Vx
		return &AddToVx{GenericInstruction{opcode}}
	case 0x8000: // Multi-purpose category, more decoding required
		switch opcode & 0x000F; { // Mask the low nibble for further decoding
		case 0x0000: // Opcode for setting Vx = Vy
			return &SetVxVy{GenericInstruction{opcode}}
		case 0x0001: // Opcode for setting Vx = Vx OR Vy
			return &SetVxOrVy{GenericInstruction{opcode}}
		case 0x0002: // Opcode for setting Vx = Vx AND Vy
			return &SetVxAndVy{GenericInstruction{opcode}}
		case 0x0003: // Opcode for setting Vx = Vx XOR Vy
			return &SetVxXorVy{GenericInstruction{opcode}}
		case 0x0004: // Opcode for adding Vy to Vx
			return &AddVyToVx{GenericInstruction{opcode}}
		case 0x0005: // Opcode for subtracting Vy from Vx
			return &SubtractVyFromVx{GenericInstruction{opcode}}
		case 0x0006: // Opcode for shifting Vx right by 1
			return &ShiftVxRight{GenericInstruction{opcode}}
		case 0x0007: // Opcode for setting Vx = Vy - Vx
			return &SetVxVyMinusVx{GenericInstruction{opcode}}
		case 0x000E: // Opcode for shifting Vx left by 1
			return &ShiftVxLeft{GenericInstruction{opcode}}
		default: // Unrecognized 0x8??? opcode
			return &UnknownInstruction{opcode}
		}
	case 0x9000: // Opcode for skipping next instruction if Vx != Vy
		return &SkipIfVxVyNotEqual{GenericInstruction{opcode}}
	case 0xA000: // Opcode for setting I = NNN
		return &SetI{GenericInstruction{opcode}}
	case 0xB000: // Opcode for jumping to address NNN + V0
		return &JumpToAddressPlusV0{GenericInstruction{opcode}}
	case 0xC000: // Opcode for setting Vx = random byte AND NN
		return &SetVxRandom{GenericInstruction{opcode}}
	case 0xD000: // Opcode for drawing a sprite at (Vx, Vy) with width 8 and height N
		return &DrawSprite{GenericInstruction{opcode}}
	case 0xE000: // Multi-purpose category, more decoding required
		switch opcode & 0x00FF; { // Mask the low byte for further decoding
		case 0x009E: // Opcode for skipping next instruction if key with value Vx is pressed
			return &SkipIfKeyPressed{GenericInstruction{opcode}}
		case 0x00A1: // Opcode for skipping next instruction if key with value Vx is not pressed
			return &SkipIfKeyNotPressed{GenericInstruction{opcode}}
		default: // Unrecognized 0xE??? opcode
			return &UnknownInstruction{opcode}
		}
	case 0xF000: // Multi-purpose category, more decoding required
		switch opcode & 0x00FF; { // Mask the low byte for further decoding
		case 0x0007: // Opcode for setting Vx = delay timer value
			return &SetVxDelayTimer{GenericInstruction{opcode}}
		case 0x000A: // Opcode for waiting for a key press and storing the value in Vx
			return &WaitForKeyPress{GenericInstruction{opcode}}
		case 0x0015: // Opcode for setting delay timer = Vx
			return &SetDelayTimer{GenericInstruction{opcode}}
		case 0x0018: // Opcode for setting sound timer = Vx
			return &SetSoundTimer{GenericInstruction{opcode}}
		case 0x001E: // Opcode for setting I = I + Vx
			return &SetIPlusVx{GenericInstruction{opcode}}
		case 0x0029: // Opcode for setting I = location of sprite for digit Vx
			return &SetISprite{GenericInstruction{opcode}}
		case 0x0033: // Opcode for storing BCD representation of Vx in memory locations I, I+1, and I+2
			return &StoreBCD{GenericInstruction{opcode}}
		case 0x0055: // Opcode for storing V0 to Vx in memory starting at location I
			return &StoreRegisters{GenericInstruction{opcode}}
		case 0x0065: // Opcode for filling V0 to Vx with values from memory starting at location I
			return &FillRegisters{GenericInstruction{opcode}}
		default: // Unrecognized 0xF??? opcode
			return &UnknownInstruction{opcode}
		}
	default: // Unrecognized opcode
		return &UnknownInstruction{opcode} // Default case for unrecognized instructions
	}
}
