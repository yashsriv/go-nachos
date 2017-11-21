// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package machine

import (
	"errors"
	"fmt"

	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/global"
	"github.com/yashsriv/go-nachos/utils"
)

// EnableDebugging allows enabling debugging for userprogs
func (m *Machine) EnableDebugging() {
	m.singleStep = true
}

// RaiseException transfers control to the Nachos kernel from user mode, because
//	the user program either invoked a system call, or some exception
//	occured (such as the address translation failed).
//
//	"which" -- the cause of the kernel trap
//	"badVaddr" -- the virtual address causing the trap, if appropriate
func (m *Machine) RaiseException(which enums.ExceptionType, badVAddr uint32) {
	utils.Debug('m', "Exception: %q\n", which)
	m.Registers[BadVAddrReg] = badVAddr
	m.DelayedLoad(0, 0) // finish anything in progress
	global.Interrupt.SetStatus(enums.SystemMode)
	global.ExceptionHandler(which) // interrupts are enabled at this point
	global.Interrupt.SetStatus(enums.UserMode)
}

// Debugger is a primitive debugger for user programs.  Note that we can't use
//	gdb to debug user programs, since gdb doesn't run on top of Nachos.
//	It could, but you'd have to implement *a lot* more system calls
//	to get it to work!
//
//	So just allow single-stepping, and printing the contents of memory.
func (m *Machine) Debugger() {
	var num int
	var buf string

	global.Interrupt.DumpState()
	m.DumpState()
	fmt.Printf("%d> ", global.Stats.TotalTicks)
	if _, err := fmt.Scanf("%s", &buf); err != nil {
		utils.Panic(errors.New("Error while trying to scan from stdio"))
	}
	if n, err := fmt.Sscanf(buf, "%d", &num); err != nil {
		utils.Panic(errors.New("Error while trying to scan from string"))
	} else if n == 1 {
		m.runUntilTime = num
	} else {
		m.runUntilTime = 0
	}
	switch buf {
	case "\n":
		break
	case "c":
		m.singleStep = false
	case "?":
		fmt.Printf("Machine commands:\n")
		fmt.Printf("    <return>  execute one instruction\n")
		fmt.Printf("    <number>  run until the given timer tick\n")
		fmt.Printf("    c         run until completion\n")
		fmt.Printf("    ?         print help message\n")
		break
	}
	// }
}

// DumpState prints the user program's CPU state.  We might print the contents
//	of memory, but that seemed like overkill.
func (m *Machine) DumpState() {
	fmt.Println("Machine registers:")
	for i := 0; i < NumGPRegs; i++ {
		endchar := ""
		if i%4 == 3 {
			endchar = "\n"
		}
		switch i {
		case StackReg:
			fmt.Printf("\tSP(%d):\t0x%x%s", i, m.Registers[i], endchar)
		case RetAddrReg:
			fmt.Printf("\tRA(%d):\t0x%x%s", i, m.Registers[i], endchar)
		default:
			fmt.Printf("\t%d:\t0x%x%s", i, m.Registers[i], endchar)
		}
	}
	fmt.Printf("\tHi:\t0x%x", m.Registers[HiReg])
	fmt.Printf("\tLo:\t0x%x\n", m.Registers[LoReg])
	fmt.Printf("\tPC:\t0x%x", m.Registers[PCReg])
	fmt.Printf("\tNextPC:\t0x%x", m.Registers[NextPCReg])
	fmt.Printf("\tPrevPC:\t0x%x\n", m.Registers[PrevPCReg])
	fmt.Printf("\tLoad:\t0x%x", m.Registers[LoadReg])
	fmt.Printf("\tLoadV:\t0x%x\n", m.Registers[LoadValueReg])
	fmt.Printf("\n")
}

// ReadRegister fetches the contents of a user program register.
func (m *Machine) ReadRegister(num int) uint32 {
	utils.Assert(num >= 0 && num < NumTotalRegs, "Register number should be within range")
	return m.Registers[num]
}

// WriteRegister writes the contents of a user program register.
func (m *Machine) WriteRegister(num int, value uint32) {
	utils.Assert(num >= 0 && num < NumTotalRegs, "Register number should be within range")
	m.Registers[num] = value
}

// PageTable allows getting the page table
func (m *Machine) PageTable() []utils.TranslationEntry {
	return m.KernelPageTable
}

// SetPageTable allows setting the page table
func (m *Machine) SetPageTable(table []utils.TranslationEntry) {
	m.KernelPageTable = table
}

// GetMainMemory returns the main memory
func (m *Machine) GetMainMemory() []byte {
	return m.MainMemory[:]
}
