// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package machine

// OpFormat is an Enum for Op formats
type OpFormat int

// Set of possible OpFormats
const (
	_             = iota
	IFMT OpFormat = iota
	JFMT
	RFMT
)

// Instruction represented in both
//   undecoded binary form
//   decoded to identify
//	    operation to do
//	    registers to act on
//	    any immediate operand value
type Instruction struct {
	Value      uint32 // binary representation of the instruction 32 bits
	OpCode     byte   // Type of instruction.
	Rs, Rt, Rd byte   // Three registers from instruction.
	Extra      uint32 // Immediate or target or shamt field or offset.
}

// OpInfo stores information for an operation
type OpInfo struct {
	OpCode byte     /* Translated op code. */
	Format OpFormat /* Format type (IFMT or JFMT or RFMT) */
}

// Decode decodes an instruction provided it has some value
func (i *Instruction) Decode() {
	var opPtr *OpInfo

	value := i.Value
	i.Rs = byte((value >> 21) & 0x1f) // 5 bits saved in byte type
	i.Rt = byte((value >> 16) & 0x1f) // 5 bits saved in byte type
	i.Rd = byte((value >> 11) & 0x1f) // 5 bits saved in byte type

	opCode := byte((value >> 26) & 0x3f) // 6 bits saved in byte type
	opPtr = &opTable[opCode]
	i.OpCode = opPtr.OpCode
	if opPtr.Format == IFMT {
		i.Extra = uint32(value & 0xffff)
		// Two's complement if signed bit is not 0:
		if i.Extra&0x8000 != 0 {
			i.Extra |= 0xffff0000
		}
	} else if opPtr.Format == RFMT {
		// shift (shamt) stored in 5 bits (0x1f)
		i.Extra = uint32((value >> 6) & (0x1f))
	} else {
		// pseudo-address stored in 26 bits
		i.Extra = uint32(value & 0x3ffffff)
	}
	// If opCode is 0x00, we need to check funct
	if i.OpCode == SPECIAL {
		funct := value & 0x3f // 6 bits
		i.OpCode = specialTable[funct]
	} else if i.OpCode == BCOND {
		var k = int(value & 0x1f0000)

		if k == 0 {
			i.OpCode = OP_BLTZ
		} else if k == 0x10000 {
			i.OpCode = OP_BGEZ
		} else if k == 0x100000 {
			i.OpCode = OP_BLTZAL
		} else if k == 0x110000 {
			i.OpCode = OP_BGEZAL
		} else {
			i.OpCode = OP_UNIMP
		}
	}
}
