// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package machine

import (
	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/utils"
)

// Constants relating to memory sizes
const (
	PageSize     uint32 = 128
	NumPhysPages uint32 = 1024
	MemorySize   uint32 = PageSize * NumPhysPages
)

// Constants relating to register numbers
const (
	StackReg     int = 29
	RetAddrReg   int = 31
	NumGPRegs    int = 32
	HiReg        int = 32
	LoReg        int = 33
	PCReg        int = 34
	NextPCReg    int = 35
	PrevPCReg    int = 36
	LoadReg      int = 37
	LoadValueReg int = 38
	BadVAddrReg  int = 39
	NumTotalRegs int = 40
)

// Machine is the simulated host workstation hardware, as
// seen by user programs -- the CPU registers, main memory, etc.
// User programs shouldn't be able to tell that they are running on our
// simulator or on the real hardware, except
// * we don't support floating point instructions
// * the system call interface to Nachos is not the same as UNIX
//	 (10 system calls in Nachos vs. 200 in UNIX!)
// If we were to implement more of the UNIX system calls, we ought to be
// able to run Nachos on top of Nachos!
type Machine struct {
	MainMemory [MemorySize]byte
	Registers  [NumTotalRegs]uint32

	KernelPageTable []utils.TranslationEntry

	singleStep   bool
	runUntilTime int
}

// Verify that interface is implemented
var _ interfaces.IMachine = &Machine{}

// Implementation is in translate.go, machine-impl.go and mipssim.go
//   translate.go has functions related to memory access
//   mipssim.go has functionality for emulating the mips architecture
//     and dummy machine
//   machine-impl.go has the rest of the functionality
