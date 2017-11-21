package interfaces

import (
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
	Print()

	SaveUserState()
	RestoreUserState()
	Name() string
	Space() IProcessAddressSpace
	SetSpace(IProcessAddressSpace)
	PID() int
	PPID() int
}

// Concrete implementation in threads/thread.go
