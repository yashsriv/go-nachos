// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package interfaces

import (
	"fmt"

	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/utils"
)

// IThread defines the interface for Thread
type IThread interface {
	Init(string)
	ThreadFork(utils.VoidFunction, interface{})
	YieldCPU()
	PutThreadToSleep()
	FinishThread()
	SetStatus(enums.ThreadStatus)
	fmt.Stringer // Can be used to print thread for debugging

	SaveUserState()
	RestoreUserState()
	Space() IProcessAddressSpace
	SetSpace(IProcessAddressSpace)
	PID() int
	PPID() int
}

// Concrete implementation in threads/thread.go
