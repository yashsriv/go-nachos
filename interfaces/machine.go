package interfaces

import (
	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/utils"
)

// IMachine defines the interface for our machine
type IMachine interface {
	EnableDebugging()
	Run()
	ReadRegister(num int) uint32
	WriteRegister(num int, value uint32)

	OneInstruction(instr IInstruction)
	DelayedLoad(nextReg byte, nextVal uint32)

	ReadMem(addr uint32, size int) (uint32, bool)
	WriteMem(addr uint32, size int, value uint32) bool

	Translate(virtAddr uint32, size int, writing bool) (uint32, enums.ExceptionType)
	RaiseException(which enums.ExceptionType, badVAddr uint32)

	PageTable() []utils.TranslationEntry
	SetPageTable([]utils.TranslationEntry)
	GetMainMemory() []byte

	Debugger()
	DumpState()
}

// Concrete implementation in machine/machine.go
