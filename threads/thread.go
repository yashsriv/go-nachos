// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package threads

import (
	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/machine"
)

// Constants
const (
	MachineStateSize int = 18
	StackSize        int = 4 * 1024
)

// Thread defines a "thread control block" -- which
// represents a single thread of execution.
//
//  Every thread has:
//     an execution stack for activation records ("stackTop" and "stack")
//     space to save CPU registers while not running ("machineState")
//     a "status" (running/ready/blocked)
//
//  Some threads also belong to a user address space; threads
//  that only run in the kernel have a NULL address space.
type Thread struct {
	name   string
	stack  []int
	status enums.ThreadStatus
	pid    int
	ppid   int

	userRegisters [machine.NumTotalRegs]uint32
	stateRestored bool
	space         interfaces.IProcessAddressSpace
}

// Check if Thread implements IThread
var _ interfaces.IThread = (*Thread)(nil)
