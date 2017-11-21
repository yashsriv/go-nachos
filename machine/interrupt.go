// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package machine

import (
	"container/list"

	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/utils"
)

// PendingInterrupt defines an interrupt that is scheduled
// to occur in the future.
type PendingInterrupt struct {
	Handler utils.VoidFunction
	Param   interface{}
	When    int
	TypeInt enums.IntType
}

// Interrupt defines the data structures for the simulation
// of hardware interrupts.  We record whether interrupts are enabled
// or disabled, and any hardware interrupts that are scheduled to occur
// in the future.
type Interrupt struct {
	pending       *list.List          // the list of interrupts scheduled to occur in the future
	level         enums.IntStatus     // are interrupts enabled or disabled?
	inHandler     bool                // TRUE if we are running an interrupt handler
	yieldOnReturn bool                // TRUE if we are to context switch on return from the interrupt handler
	status        enums.MachineStatus // idle, kernel mode, user mode

}

// internal interface for interrupt
type internalInterrupt interface {
	checkIfDue(advanceClock bool) bool // Check if an interrupt is supposed to occur now
	changeLevel(old, now enums.IntStatus)
}

// Check if interface is implemented by our Interrupt
// Implemented in interrupt-impl.go
var _ interfaces.IInterrupt = (*Interrupt)(nil)
var _ internalInterrupt = (*Interrupt)(nil)

// Various human-readable names for clearer debug messages
var intStatusNames = map[enums.IntStatus]string{
	enums.IntOff: "off",
	enums.IntOn:  "on",
}

var intTypeNames = map[enums.IntType]string{
	enums.TimerInt:        "timer",
	enums.DiskInt:         "disk",
	enums.ConsoleWriteInt: "console write",
	enums.ConsoleReadInt:  "console read",
	enums.NetworkSendInt:  "network send",
	enums.NetworkRecvInt:  "network recv",
}
