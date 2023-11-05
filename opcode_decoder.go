package main

type OpcodeDecoder struct {
}

func NewOpcodeDecoder() *OpcodeDecoder {
	return &OpcodeDecoder{}
}

func (opcodeDecoder *OpcodeDecoder) Decode(opcode uint16) Instruction {
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x00FF {
		case 0x00E0:
			return &ClearScreen{GenericInstruction{opcode}}
		case 0x00EE:
			return &ReturnFromSubroutine{GenericInstruction{opcode}}
		default:
			return &UnknownInstruction{GenericInstruction{opcode}}
		}
	case 0x1000:
		return &JumpToAddress{GenericInstruction{opcode}}
	case 0x2000:
		return &CallSubroutine{GenericInstruction{opcode}}
	case 0x3000:
		return &SkipIfVxEqual{GenericInstruction{opcode}}
	case 0x4000:
		return &SkipIfVxNotEqual{GenericInstruction{opcode}}
	case 0x5000:
		return &SkipIfVxVyEqual{GenericInstruction{opcode}}
	case 0x6000:
		return &SetVx{GenericInstruction{opcode}}
	case 0x7000:
		return &AddToVx{GenericInstruction{opcode}}
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000:
			return &SetVxVy{GenericInstruction{opcode}}
		case 0x0001:
			return &SetVxOrVy{GenericInstruction{opcode}}
		case 0x0002:
			return &SetVxAndVy{GenericInstruction{opcode}}
		case 0x0003:
			return &SetVxXorVy{GenericInstruction{opcode}}
		case 0x0004:
			return &AddVyToVx{GenericInstruction{opcode}}
		case 0x0005:
			return &SubtractVyFromVx{GenericInstruction{opcode}}
		case 0x0006:
			return &ShiftVxRight{GenericInstruction{opcode}}
		case 0x0007:
			return &SetVxVyMinusVx{GenericInstruction{opcode}}
		case 0x000E:
			return &ShiftVxLeft{GenericInstruction{opcode}}
		default:
			return &UnknownInstruction{GenericInstruction{opcode}}
		}
	case 0x9000:
		return &SkipIfVxVyNotEqual{GenericInstruction{opcode}}
	case 0xA000:
		return &SetI{GenericInstruction{opcode}}
	case 0xB000:
		return &JumpToAddressPlusV0{GenericInstruction{opcode}}
	case 0xC000:
		return &SetVxRandom{GenericInstruction{opcode}}
	case 0xD000:
		return &DrawSprite{GenericInstruction{opcode}}
	case 0xE000:
		switch opcode & 0x00FF {
		case 0x009E:
			return &SkipIfKeyPressed{GenericInstruction{opcode}}
		case 0x00A1:
			return &SkipIfKeyNotPressed{GenericInstruction{opcode}}
		default:
			return &UnknownInstruction{GenericInstruction{opcode}}
		}
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x0007:
			return &SetVxDelayTimer{GenericInstruction{opcode}}
		case 0x000A:
			return &WaitForKeyPress{GenericInstruction{opcode}}
		case 0x0015:
			return &SetDelayTimer{GenericInstruction{opcode}}
		case 0x0018:
			return &SetSoundTimer{GenericInstruction{opcode}}
		case 0x001E:
			return &SetIPlusVx{GenericInstruction{opcode}}
		case 0x0029:
			return &SetISprite{GenericInstruction{opcode}}
		case 0x0033:
			return &StoreBCD{GenericInstruction{opcode}}
		case 0x0055:
			return &Storeregisters{GenericInstruction{opcode}}
		case 0x0065:
			return &Fillregisters{GenericInstruction{opcode}}
		default:
			return &UnknownInstruction{GenericInstruction{opcode}}
		}
	default:
		return &UnknownInstruction{GenericInstruction{opcode}}
	}
}
