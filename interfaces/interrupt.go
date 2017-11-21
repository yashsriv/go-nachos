// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package interfaces

import "github.com/yashsriv/go-nachos/enums"

// IInterrupt defines the interface for the interrupt struct
type IInterrupt interface {
	Init()
	SetLevel(level enums.IntStatus) enums.IntStatus // Disable or enable interrupts and return previous setting.

	Enable()                   // Enable interrupts.
	GetLevel() enums.IntStatus // Return whether interrupts are enabled or disabled
	Idle()                     // The ready queue is empty, roll simulated time forward until the next interrupt

	Halt() // quit and print out stats

	YieldOnReturn() // cause a context switch on return from an interrupt handler

	GetStatus() enums.MachineStatus // idle, kernel, user
	SetStatus(st enums.MachineStatus)

	DumpState() // Print interrupt state

	// NOTE: the following are internal to the hardware simulation code.
	// DO NOT call these directly. They are called by the
	// hardware device simulators.

	Schedule(pending IPendingInterrupt)

	OneTick() // Advance simulated time

}

// Concrete implementation in machine/interrupt.go
