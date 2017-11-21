// This file is the actual implementation of mips simulator
// It is not absolutely necesary to understand this fully for
// working on nachos

package machine

import (
	"fmt"
	"runtime"

	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/global"
	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/utils"
)

// OpCodes
const (
	_           = iota
	OP_ADD byte = iota
	OP_ADDI
	OP_ADDIU
	OP_ADDU
	OP_AND
	OP_ANDI
	OP_BEQ
	OP_BGEZ
	OP_BGEZAL
	OP_BGTZ
	OP_BLEZ
	OP_BLTZ
	OP_BLTZAL
	OP_BNE
	_
	OP_DIV
	OP_DIVU
	OP_J
	OP_JAL
	OP_JALR
	OP_JR
	OP_LB
	OP_LBU
	OP_LH
	OP_LHU
	OP_LUI
	OP_LW
	OP_LWL
	OP_LWR
	_
	OP_MFHI
	OP_MFLO
	_
	OP_MTHI
	OP_MTLO
	OP_MULT
	OP_MULTU
	OP_NOR
	OP_OR
	OP_ORI
	OP_RFE
	OP_SB
	OP_SH
	OP_SLL
	OP_SLLV
	OP_SLT
	OP_SLTI
	OP_SLTIU
	OP_SLTU
	OP_SRA
	OP_SRAV
	OP_SRL
	OP_SRLV
	OP_SUB
	OP_SUBU
	OP_SW
	OP_SWL
	OP_SWR
	OP_XOR
	OP_XORI
	OP_SYSCALL
	OP_UNIMP
	OP_RES
	MaxOpcode byte = 63
	SPECIAL   byte = 100
	BCOND     byte = 101
)

// SIGN_BIT has 1 in position of signed bit
const SIGN_BIT uint32 = 0x80000000

const r31 byte = 31

var opTable = [64]OpInfo{
	{SPECIAL, RFMT}, {BCOND, IFMT}, {OP_J, JFMT}, {OP_JAL, JFMT},
	{OP_BEQ, IFMT}, {OP_BNE, IFMT}, {OP_BLEZ, IFMT}, {OP_BGTZ, IFMT},
	{OP_ADDI, IFMT}, {OP_ADDIU, IFMT}, {OP_SLTI, IFMT}, {OP_SLTIU, IFMT},
	{OP_ANDI, IFMT}, {OP_ORI, IFMT}, {OP_XORI, IFMT}, {OP_LUI, IFMT},
	{OP_UNIMP, IFMT}, {OP_UNIMP, IFMT}, {OP_UNIMP, IFMT}, {OP_UNIMP, IFMT},
	{OP_RES, IFMT}, {OP_RES, IFMT}, {OP_RES, IFMT}, {OP_RES, IFMT},
	{OP_RES, IFMT}, {OP_RES, IFMT}, {OP_RES, IFMT}, {OP_RES, IFMT},
	{OP_RES, IFMT}, {OP_RES, IFMT}, {OP_RES, IFMT}, {OP_RES, IFMT},
	{OP_LB, IFMT}, {OP_LH, IFMT}, {OP_LWL, IFMT}, {OP_LW, IFMT},
	{OP_LBU, IFMT}, {OP_LHU, IFMT}, {OP_LWR, IFMT}, {OP_RES, IFMT},
	{OP_SB, IFMT}, {OP_SH, IFMT}, {OP_SWL, IFMT}, {OP_SW, IFMT},
	{OP_RES, IFMT}, {OP_RES, IFMT}, {OP_SWR, IFMT}, {OP_RES, IFMT},
	{OP_UNIMP, IFMT}, {OP_UNIMP, IFMT}, {OP_UNIMP, IFMT}, {OP_UNIMP, IFMT},
	{OP_RES, IFMT}, {OP_RES, IFMT}, {OP_RES, IFMT}, {OP_RES, IFMT},
	{OP_UNIMP, IFMT}, {OP_UNIMP, IFMT}, {OP_UNIMP, IFMT}, {OP_UNIMP, IFMT},
	{OP_RES, IFMT}, {OP_RES, IFMT}, {OP_RES, IFMT}, {OP_RES, IFMT},
}

var specialTable = []byte{
	OP_SLL, OP_RES, OP_SRL, OP_SRA, OP_SLLV, OP_RES, OP_SRLV, OP_SRAV,
	OP_JR, OP_JALR, OP_RES, OP_RES, OP_SYSCALL, OP_UNIMP, OP_RES, OP_RES,
	OP_MFHI, OP_MTHI, OP_MFLO, OP_MTLO, OP_RES, OP_RES, OP_RES, OP_RES,
	OP_MULT, OP_MULTU, OP_DIV, OP_DIVU, OP_RES, OP_RES, OP_RES, OP_RES,
	OP_ADD, OP_ADDU, OP_SUB, OP_SUBU, OP_AND, OP_OR, OP_XOR, OP_NOR,
	OP_RES, OP_RES, OP_SLT, OP_SLTU, OP_RES, OP_RES, OP_RES, OP_RES,
	OP_RES, OP_RES, OP_RES, OP_RES, OP_RES, OP_RES, OP_RES, OP_RES,
	OP_RES, OP_RES, OP_RES, OP_RES, OP_RES, OP_RES, OP_RES, OP_RES,
}

// For Debugging M.Registers

// Enum for Register Type
type regType int

// regType enum values
const (
	NONE regType = iota
	RS
	RT
	RD
	EXTRA
)

// Struct to store op strings
type opString struct {
	val  string // Printed version of instruction
	args [3]regType
}

var opStrings = []opString{
	{"Shouldn't happen", [3]regType{NONE, NONE, NONE}},
	{"ADD r%d,r%d,r%d", [3]regType{RD, RS, RT}},
	{"ADDI r%d,r%d,%d", [3]regType{RT, RS, EXTRA}},
	{"ADDIU r%d,r%d,%d", [3]regType{RT, RS, EXTRA}},
	{"ADDU r%d,r%d,r%d", [3]regType{RD, RS, RT}},
	{"AND r%d,r%d,r%d", [3]regType{RD, RS, RT}},
	{"ANDI r%d,r%d,%d", [3]regType{RT, RS, EXTRA}},
	{"BEQ r%d,r%d,%d", [3]regType{RS, RT, EXTRA}},
	{"BGEZ r%d,%d", [3]regType{RS, EXTRA, NONE}},
	{"BGEZAL r%d,%d", [3]regType{RS, EXTRA, NONE}},
	{"BGTZ r%d,%d", [3]regType{RS, EXTRA, NONE}},
	{"BLEZ r%d,%d", [3]regType{RS, EXTRA, NONE}},
	{"BLTZ r%d,%d", [3]regType{RS, EXTRA, NONE}},
	{"BLTZAL r%d,%d", [3]regType{RS, EXTRA, NONE}},
	{"BNE r%d,r%d,%d", [3]regType{RS, RT, EXTRA}},
	{"Shouldn't happen", [3]regType{NONE, NONE, NONE}},
	{"DIV r%d,r%d", [3]regType{RS, RT, NONE}},
	{"DIVU r%d,r%d", [3]regType{RS, RT, NONE}},
	{"J %d", [3]regType{EXTRA, NONE, NONE}},
	{"JAL %d", [3]regType{EXTRA, NONE, NONE}},
	{"JALR r%d,r%d", [3]regType{RD, RS, NONE}},
	{"JR r%d,r%d", [3]regType{RD, RS, NONE}},
	{"LB r%d,%d(r%d)", [3]regType{RT, EXTRA, RS}},
	{"LBU r%d,%d(r%d)", [3]regType{RT, EXTRA, RS}},
	{"LH r%d,%d(r%d)", [3]regType{RT, EXTRA, RS}},
	{"LHU r%d,%d(r%d)", [3]regType{RT, EXTRA, RS}},
	{"LUI r%d,%d", [3]regType{RT, EXTRA, NONE}},
	{"LW r%d,%d(r%d)", [3]regType{RT, EXTRA, RS}},
	{"LWL r%d,%d(r%d)", [3]regType{RT, EXTRA, RS}},
	{"LWR r%d,%d(r%d)", [3]regType{RT, EXTRA, RS}},
	{"Shouldn't happen", [3]regType{NONE, NONE, NONE}},
	{"MFHI r%d", [3]regType{RD, NONE, NONE}},
	{"MFLO r%d", [3]regType{RD, NONE, NONE}},
	{"Shouldn't happen", [3]regType{NONE, NONE, NONE}},
	{"MTHI r%d", [3]regType{RS, NONE, NONE}},
	{"MTLO r%d", [3]regType{RS, NONE, NONE}},
	{"MULT r%d,r%d", [3]regType{RS, RT, NONE}},
	{"MULTU r%d,r%d", [3]regType{RS, RT, NONE}},
	{"NOR r%d,r%d,r%d", [3]regType{RD, RS, RT}},
	{"OR r%d,r%d,r%d", [3]regType{RD, RS, RT}},
	{"ORI r%d,r%d,%d", [3]regType{RT, RS, EXTRA}},
	{"RFE", [3]regType{NONE, NONE, NONE}},
	{"SB r%d,%d(r%d)", [3]regType{RT, EXTRA, RS}},
	{"SH r%d,%d(r%d)", [3]regType{RT, EXTRA, RS}},
	{"SLL r%d,r%d,%d", [3]regType{RD, RT, EXTRA}},
	{"SLLV r%d,r%d,r%d", [3]regType{RD, RT, RS}},
	{"SLT r%d,r%d,r%d", [3]regType{RD, RS, RT}},
	{"SLTI r%d,r%d,%d", [3]regType{RT, RS, EXTRA}},
	{"SLTIU r%d,r%d,%d", [3]regType{RT, RS, EXTRA}},
	{"SLTU r%d,r%d,r%d", [3]regType{RD, RS, RT}},
	{"SRA r%d,r%d,%d", [3]regType{RD, RT, EXTRA}},
	{"SRAV r%d,r%d,r%d", [3]regType{RD, RT, RS}},
	{"SRL r%d,r%d,%d", [3]regType{RD, RT, EXTRA}},
	{"SRLV r%d,r%d,r%d", [3]regType{RD, RT, RS}},
	{"SUB r%d,r%d,r%d", [3]regType{RD, RS, RT}},
	{"SUBU r%d,r%d,r%d", [3]regType{RD, RS, RT}},
	{"SW r%d,%d(r%d)", [3]regType{RT, EXTRA, RS}},
	{"SWL r%d,%d(r%d)", [3]regType{RT, EXTRA, RS}},
	{"SWR r%d,%d(r%d)", [3]regType{RT, EXTRA, RS}},
	{"XOR r%d,r%d,r%d", [3]regType{RD, RS, RT}},
	{"XORI r%d,r%d,%d", [3]regType{RT, RS, EXTRA}},
	{"SYSCALL", [3]regType{NONE, NONE, NONE}},
	{"Unimplemented", [3]regType{NONE, NONE, NONE}},
	{"Reserved", [3]regType{NONE, NONE, NONE}},
}

func typeToReg(reg regType, instr *Instruction) int {
	switch reg {
	case RS:
		return int(instr.Rs)
	case RT:
		return int(instr.Rt)
	case RD:
		return int(instr.Rd)
	case EXTRA:
		return int(instr.Extra)
	default:
		return -1
	}
}

func mult(a int32, b int32, signedArith bool, hiPtr *uint32, loPtr *uint32) {
	if (a == 0) || (b == 0) {
		*hiPtr = 0
		*loPtr = 0
		return
	}

	// Compute the sign of the result, then make everything positive
	// so unsigned computation can be done in the main loop.
	var negative = false
	if signedArith {
		if a < 0 {
			negative = !negative
			a = -a
		}
		if b < 0 {
			negative = !negative
			b = -b
		}
	}

	// Compute the result in unsigned arithmetic (check a's bits one at
	// a time, and add in a shifted value of b).
	var bLo = uint32(b)
	var bHi uint32
	var lo uint32
	var hi uint32
	for i := 0; i < 32; i++ {
		if a&1 != 0 {
			lo += bLo
			if lo < bLo { // Carry out of the low bits?
				hi++
			}
			hi += bHi
			if (uint32(a) & 0xfffffffe) == 0 {
				break
			}
		}
		bHi <<= 1
		if bLo&SIGN_BIT != 0 {
			bHi |= 1
		}

		bLo <<= 1
		a >>= 1
	}

	// If the result is supposed to be negative, compute the two's
	// complement of the double-word result.
	if negative {
		hi = ^hi
		lo = ^lo
		lo++
		if lo == 0 {
			hi++
		}
	}

	*hiPtr = hi
	*loPtr = lo
}

// DelayedLoad simulates the effects of a delayed load.
//
// 	NOTE -- RaiseException/CheckInterrupts must also call DelayedLoad,
//	since any delayed load must get applied before we trap to the kernel.
func (m *Machine) DelayedLoad(nextReg byte, nextValue uint32) {
	m.Registers[m.Registers[LoadReg]] = m.Registers[LoadValueReg]
	m.Registers[LoadReg] = uint32(nextReg)
	m.Registers[LoadValueReg] = nextValue
	m.Registers[0] = 0 // and always make sure R0 stays zero.
}

func indexToAddr(x uint32) uint32 {
	return x << 2
}

// OneInstruction executes one instruction from a user-level program
//
// 	If there is any kind of exception or interrupt, we invoke the
//	exception handler, and when it returns, we return to Run(), which
//	will re-invoke us in a loop.  This allows us to
//	re-start the instruction execution from the beginning, in
//	case any of our state has changed.  On a syscall,
// 	the OS software must increment the PC so execution begins
// 	at the instruction immediately after the syscall.
//
//	This routine is re-entrant, in that it can be called multiple
//	times concurrently -- one for each thread executing user code.
//	We get re-entrancy by never caching any data -- we always re-start the
//	simulation from scratch each time we are called (or after trapping
//	back to the Nachos kernel on an exception or interrupt), and we always
//	store all data back to the machine m.Registers and memory before
//	leaving.  This allows the Nachos kernel to control our behavior
//	by controlling the contents of memory, the translation table,
//	and the register set.
func (m *Machine) OneInstruction(instruction interfaces.IInstruction) {
	var raw uint32
	var nextLoadReg byte
	var nextLoadValue uint32 // record delayed load operation, to apply in the future

	// Fetch instruction
	raw, success := m.ReadMem(m.Registers[PCReg], 4)
	if !success {
		return // exception occurred
	}

	instr := instruction.(*Instruction)

	instr.Value = raw
	instr.Decode()

	if utils.DebugIsEnabled('m') {
		x := opStrings[instr.OpCode]
		utils.Assert(instr.OpCode <= MaxOpcode)
		fmt.Printf("At PC = 0x%x: ", m.Registers[PCReg])
		fmt.Printf(x.val, typeToReg(x.args[0], instr),
			typeToReg(x.args[1], instr), typeToReg(x.args[2], instr))
		fmt.Printf("\n")
	}

	// Compute next pc, but don't install in case there's an error or branch.
	pcAfter := m.Registers[NextPCReg] + 4

	// Execute the instruction (cf. Kane's book)
	switch instr.OpCode {

	case OP_ADD:
		sum := int32(m.Registers[instr.Rs]) + int32(m.Registers[instr.Rt])
		if (m.Registers[instr.Rs]^m.Registers[instr.Rt])&SIGN_BIT == 0 &&
			(m.Registers[instr.Rs]^uint32(sum)&SIGN_BIT != 0) {
			m.RaiseException(enums.OverflowException, 0)
			return
		}
		m.Registers[instr.Rd] = uint32(sum)

	case OP_ADDI:
		sum := int32(m.Registers[instr.Rs]) + int32(instr.Extra)
		if (m.Registers[instr.Rs]^instr.Extra&SIGN_BIT) == 0 &&
			(instr.Extra^uint32(sum)&SIGN_BIT) != 0 {
			m.RaiseException(enums.OverflowException, 0)
			return
		}
		m.Registers[instr.Rt] = uint32(sum)

	case OP_ADDIU:
		m.Registers[instr.Rt] = m.Registers[instr.Rs] + instr.Extra

	case OP_ADDU:
		m.Registers[instr.Rd] = m.Registers[instr.Rs] + m.Registers[instr.Rt]

	case OP_AND:
		m.Registers[instr.Rd] = m.Registers[instr.Rs] & m.Registers[instr.Rt]

	case OP_ANDI:
		m.Registers[instr.Rt] = m.Registers[instr.Rs] & instr.Extra & 0xffff

	case OP_BEQ:
		if m.Registers[instr.Rs] == m.Registers[instr.Rt] {
			pcAfter = m.Registers[NextPCReg] + indexToAddr(instr.Extra)
		}

	case OP_BGEZAL:
		m.Registers[r31] = m.Registers[NextPCReg] + 4
		if (m.Registers[instr.Rs] & SIGN_BIT) == 0 {
			pcAfter = m.Registers[NextPCReg] + indexToAddr(instr.Extra)
		}
	case OP_BGEZ:
		if (m.Registers[instr.Rs] & SIGN_BIT) == 0 {
			pcAfter = m.Registers[NextPCReg] + indexToAddr(instr.Extra)
		}

	case OP_BGTZ:
		if int32(m.Registers[instr.Rs]) > 0 {
			pcAfter = m.Registers[NextPCReg] + indexToAddr(instr.Extra)
		}
	case OP_BLEZ:
		if int32(m.Registers[instr.Rs]) <= 0 {
			pcAfter = m.Registers[NextPCReg] + indexToAddr(instr.Extra)
		}
	case OP_BLTZAL:
		m.Registers[r31] = m.Registers[NextPCReg] + 4
		if m.Registers[instr.Rs]&SIGN_BIT != 0 {
			pcAfter = m.Registers[NextPCReg] + indexToAddr(instr.Extra)
		}
	case OP_BLTZ:
		if m.Registers[instr.Rs]&SIGN_BIT != 0 {
			pcAfter = m.Registers[NextPCReg] + indexToAddr(instr.Extra)
		}

	case OP_BNE:
		if m.Registers[instr.Rs] != m.Registers[instr.Rt] {
			pcAfter = m.Registers[NextPCReg] + indexToAddr(instr.Extra)
		}
	case OP_DIV:
		if m.Registers[instr.Rt] == 0 {
			m.Registers[LoReg] = 0
			m.Registers[HiReg] = 0
		} else {
			m.Registers[LoReg] = uint32(int32(m.Registers[instr.Rs]) / int32(m.Registers[instr.Rt]))
			m.Registers[HiReg] = uint32(int32(m.Registers[instr.Rs]) % int32(m.Registers[instr.Rt]))
		}

	case OP_DIVU:
		rs := m.Registers[instr.Rs]
		rt := m.Registers[instr.Rt]
		if rt == 0 {
			m.Registers[LoReg] = 0
			m.Registers[HiReg] = 0
		} else {
			tmp := rs / rt
			m.Registers[LoReg] = tmp
			tmp = rs % rt
			m.Registers[HiReg] = tmp
		}
	case OP_JAL:
		m.Registers[r31] = m.Registers[NextPCReg] + 4
		pcAfter = (pcAfter & 0xf0000000) | indexToAddr(instr.Extra)
	case OP_J:
		pcAfter = (pcAfter & 0xf0000000) | indexToAddr(instr.Extra)

	case OP_JALR:
		m.Registers[instr.Rd] = m.Registers[NextPCReg] + 4
		pcAfter = m.Registers[instr.Rs]
	case OP_JR:
		pcAfter = m.Registers[instr.Rs]

	case OP_LB:
		tmp := int32(m.Registers[instr.Rs]) + int32(instr.Extra)
		value, status := m.ReadMem(uint32(tmp), 1)
		if !status {
			return
		}

		if (value&0x80) != 0 && (instr.OpCode == OP_LB) {
			value |= 0xffffff00
		} else {
			value &= 0xff
		}
		nextLoadReg = instr.Rt
		nextLoadValue = value
	case OP_LBU:
		tmpi := m.Registers[instr.Rs] + instr.Extra
		value, status := m.ReadMem(tmpi, 1)
		if !status {
			return
		}

		if (value&0x80) != 0 && (instr.OpCode == OP_LB) {
			value |= 0xffffff00
		} else {
			value &= 0xff
		}
		nextLoadReg = instr.Rt
		nextLoadValue = value
	case OP_LH:
		tmpi := int32(m.Registers[instr.Rs]) + int32(instr.Extra)
		if tmpi&0x1 != 0 {
			m.RaiseException(enums.AddressErrorException, uint32(tmpi))
			return
		}
		value, status := m.ReadMem(uint32(tmpi), 2)
		if !status {
			return
		}

		if (value&0x8000) != 0 && (instr.OpCode == OP_LH) {
			value |= 0xffff0000
		} else {
			value &= 0xffff
		}
		nextLoadReg = instr.Rt
		nextLoadValue = value
	case OP_LHU:
		tmpi := m.Registers[instr.Rs] + instr.Extra
		if tmpi&0x1 != 0 {
			m.RaiseException(enums.AddressErrorException, tmpi)
			return
		}
		value, status := m.ReadMem(tmpi, 2)
		if !status {
			return
		}

		if (value&0x8000) != 0 && (instr.OpCode == OP_LH) {
			value |= 0xffff0000
		} else {
			value &= 0xffff
		}
		nextLoadReg = instr.Rt
		nextLoadValue = value
	case OP_LUI:
		utils.Debug('m', "Executing: LUI r%d,%d\n", instr.Rt, instr.Extra)
		m.Registers[instr.Rt] = instr.Extra << 16
	case OP_LW:
		tmpi := m.Registers[instr.Rs] + instr.Extra
		if tmpi&0x3 != 0 {
			m.RaiseException(enums.AddressErrorException, tmpi)
			return
		}
		value, status := m.ReadMem(tmpi, 4)
		if !status {
			return
		}
		nextLoadReg = instr.Rt
		nextLoadValue = value
	case OP_LWL:
		tmpi := m.Registers[instr.Rs] + instr.Extra

		// ReadMem assumes all 4 byte requests are aligned on an even
		// word boundary.  Also, the little endian/big endian swap code would
		// fail (I think) if the other cases are ever exercised.
		utils.Assert((tmpi & 0x3) == 0)

		value, status := m.ReadMem(tmpi, 4)
		if !status {
			return
		}
		if byte(m.Registers[LoadReg]) == instr.Rt {
			nextLoadValue = m.Registers[LoadValueReg]
		} else {
			nextLoadValue = m.Registers[instr.Rt]
		}
		switch tmpi & 0x3 {
		case 0:
			nextLoadValue = value
		case 1:
			nextLoadValue = (nextLoadValue & 0xff) | (value << 8)
		case 2:
			nextLoadValue = (nextLoadValue & 0xffff) | (value << 16)
		case 3:
			nextLoadValue = (nextLoadValue & 0xffffff) | (value << 24)
		}
		nextLoadReg = instr.Rt
	case OP_LWR:
		tmpi := m.Registers[instr.Rs] + instr.Extra

		// ReadMem assumes all 4 byte requests are aligned on an even
		// word boundary.  Also, the little endian/big endian swap code would
		// fail (I think) if the other cases are ever exercised.
		utils.Assert((tmpi & 0x3) == 0)

		value, status := m.ReadMem(tmpi, 4)
		if !status {
			return
		}
		if byte(m.Registers[LoadReg]) == instr.Rt {
			nextLoadValue = m.Registers[LoadValueReg]
		} else {
			nextLoadValue = m.Registers[instr.Rt]
		}
		switch tmpi & 0x3 {
		case 0:
			nextLoadValue = (nextLoadValue & 0xffffff00) |
				((value >> 24) & 0xff)
		case 1:
			nextLoadValue = (nextLoadValue & 0xffff0000) |
				((value >> 16) & 0xffff)
		case 2:
			nextLoadValue = (nextLoadValue & 0xff000000) | ((value >> 8) & 0xffffff)
		case 3:
			nextLoadValue = value
		}
		nextLoadReg = instr.Rt
	case OP_MFHI:
		m.Registers[instr.Rd] = m.Registers[HiReg]
		break

	case OP_MFLO:
		m.Registers[instr.Rd] = m.Registers[LoReg]
		break

	case OP_MTHI:
		m.Registers[HiReg] = m.Registers[instr.Rs]
		break

	case OP_MTLO:
		m.Registers[LoReg] = m.Registers[instr.Rs]
		break

	case OP_MULT:
		mult(int32(m.Registers[instr.Rs]), int32(m.Registers[instr.Rt]), true,
			&m.Registers[HiReg], &m.Registers[LoReg])

	case OP_MULTU:
		mult(int32(m.Registers[instr.Rs]), int32(m.Registers[instr.Rt]), false,
			&m.Registers[HiReg], &m.Registers[LoReg])

	case OP_NOR:
		m.Registers[instr.Rd] = ^(m.Registers[instr.Rs] | m.Registers[instr.Rt])

	case OP_OR:
		// TODO: Look into this interesting bug
		m.Registers[instr.Rd] = m.Registers[instr.Rs] | m.Registers[instr.Rt]
		break

	case OP_ORI:
		m.Registers[instr.Rt] = m.Registers[instr.Rs] | (instr.Extra & 0xffff)
		break

	case OP_SB:
		if !m.WriteMem(m.Registers[instr.Rs]+instr.Extra, 1, m.Registers[instr.Rt]) {
			return
		}
	case OP_SH:
		if !m.WriteMem(m.Registers[instr.Rs]+instr.Extra, 2, m.Registers[instr.Rt]) {
			return
		}

	case OP_SLL:
		m.Registers[instr.Rd] = m.Registers[instr.Rt] << instr.Extra
	case OP_SLLV:
		m.Registers[instr.Rd] = m.Registers[instr.Rt] << m.Registers[instr.Rs] & 0x1f
	case OP_SLT:
		if int32(m.Registers[instr.Rs]) < int32(m.Registers[instr.Rt]) {
			m.Registers[instr.Rd] = 1
		} else {
			m.Registers[instr.Rd] = 0
		}

	case OP_SLTI:
		if int32(m.Registers[instr.Rs]) < int32(instr.Extra) {
			m.Registers[instr.Rt] = 1
		} else {
			m.Registers[instr.Rt] = 0
		}

	case OP_SLTIU:
		rs := m.Registers[instr.Rs]
		imm := instr.Extra
		if rs < imm {
			m.Registers[instr.Rt] = 1
		} else {
			m.Registers[instr.Rt] = 0
		}
	case OP_SLTU:
		rs := m.Registers[instr.Rs]
		rt := m.Registers[instr.Rt]
		if rs < rt {
			m.Registers[instr.Rd] = 1
		} else {
			m.Registers[instr.Rd] = 0
		}
	case OP_SRA:
		m.Registers[instr.Rd] = m.Registers[instr.Rt] >> uint(instr.Extra)
	case OP_SRAV:
		m.Registers[instr.Rd] = m.Registers[instr.Rt] >> uint(m.Registers[instr.Rs]&0x1f)
	case OP_SRL:
		tmpi := m.Registers[instr.Rt]
		tmpi >>= instr.Extra
		m.Registers[instr.Rd] = tmpi
	case OP_SRLV:
		tmpi := m.Registers[instr.Rt]
		tmpi >>= m.Registers[instr.Rs] & 0x1f
		m.Registers[instr.Rd] = tmpi
	case OP_SUB:
		diff := int32(m.Registers[instr.Rs]) - int32(m.Registers[instr.Rt])
		if ((m.Registers[instr.Rs]^m.Registers[instr.Rt])&SIGN_BIT) != 0 &&
			((m.Registers[instr.Rs]^uint32(diff))&SIGN_BIT) != 0 {
			m.RaiseException(enums.OverflowException, 0)
			return
		}
		m.Registers[instr.Rd] = uint32(diff)
	case OP_SUBU:
		m.Registers[instr.Rd] = m.Registers[instr.Rs] - m.Registers[instr.Rt]
	case OP_SW:
		if !m.WriteMem(m.Registers[instr.Rs]+instr.Extra, 4, m.Registers[instr.Rt]) {
			return
		}
	case OP_SWL:
		tmpi := m.Registers[instr.Rs] + instr.Extra

		// The little endian/big endian swap code would
		// fail (I think) if the other cases are ever exercised.

		utils.Assert((tmpi & 0x3) == 0)

		value, status := m.ReadMem((tmpi & ^uint32(0x3)), 4)
		if !status {
			return
		}
		switch tmpi & 0x3 {
		case 0:
			value = m.Registers[instr.Rt]
		case 1:
			value = (value & 0xff000000) | ((m.Registers[instr.Rt] >> 8) & 0xffffff)
		case 2:
			value = (value & 0xffff0000) | ((m.Registers[instr.Rt] >> 16) & 0xffff)
		case 3:
			value = (value & 0xffffff00) | ((m.Registers[instr.Rt] >> 24) & 0xff)
		}
		if !m.WriteMem((tmpi & ^uint32(0x3)), 4, value) {
			return
		}

	case OP_SWR:
		tmpi := m.Registers[instr.Rs] + instr.Extra
		// The little endian/big endian swap code would
		// fail (I think) if the other cases are ever exercised.

		utils.Assert((tmpi & 0x3) == 0)
		value, status := m.ReadMem((tmpi & ^uint32(0x3)), 4)
		if !status {
			return
		}
		switch tmpi & 0x3 {
		case 0:
			value = (value & 0xffffff) | (m.Registers[instr.Rt] << 24)
		case 1:
			value = (value & 0xffff) | (m.Registers[instr.Rt] << 16)
		case 2:
			value = (value & 0xff) | (m.Registers[instr.Rt] << 8)
		case 3:
			value = m.Registers[instr.Rt]
		}
		if !m.WriteMem((tmpi & ^uint32(0x3)), 4, value) {
			return
		}
		break

	case OP_SYSCALL:
		m.RaiseException(enums.SyscallException, 0)
		return

	case OP_XOR:
		m.Registers[instr.Rd] = m.Registers[instr.Rs] ^ m.Registers[instr.Rt]

	case OP_XORI:
		m.Registers[instr.Rt] = m.Registers[instr.Rs] ^ (instr.Extra & 0xffff)

	case OP_RES:
		m.RaiseException(enums.IllegalInstrException, 0)
		return

	case OP_UNIMP:
		m.RaiseException(enums.IllegalInstrException, 0)
		return

	default:
		utils.Assert(false)
	}

	// Now we have successfully executed the instruction.

	// Do any delayed load operation
	m.DelayedLoad(nextLoadReg, nextLoadValue)

	// Advance program counters.
	m.Registers[PrevPCReg] = m.Registers[PCReg] // for debugging, in case we
	// are jumping into lala-land
	m.Registers[PCReg] = m.Registers[NextPCReg]
	m.Registers[NextPCReg] = pcAfter
}

// Run simulates the execution of a user-level program on Nachos.
// Called by the kernel when the program starts up; never returns.
// This routine is re-entrant, in that it can be called multiple
// times concurrently -- one for each thread executing user code.
func (m *Machine) Run() {
	instr := &Instruction{} // storage for decoded instruction

	if utils.DebugIsEnabled('m') {
		fmt.Printf("Starting thread %q at time %d\n", global.CurrentThread.Name(), global.Stats.TotalTicks)
	}

	pid := global.CurrentThread.PID()
	global.ControlChannel[pid] = make(chan int, 1)
	ws := global.ControlChannel[pid]
	state := 1 // Running
	global.Interrupt.SetStatus(enums.UserMode)
	for {
		select {
		case state = <-ws:
			switch state {
			case 1:
				utils.Debug('t', "Thread %d: Running", pid)
			case 0:
				utils.Debug('t', "Thread %d: Paused", pid)
			}

		default:
			// We use runtime.Gosched() to prevent a deadlock in this case.
			// It will not be needed of work is performed here which yields
			// to the scheduler.
			runtime.Gosched()

			if state == 0 {
				break
			}

			m.OneInstruction(instr)
			global.Interrupt.OneTick()
			if m.singleStep && (m.runUntilTime <= global.Stats.TotalTicks) {
				m.Debugger()
			}
			// Do actual work here.
		}
	}
}
