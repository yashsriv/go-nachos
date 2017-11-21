package userprog

import (
	"github.com/yashsriv/go-nachos/global"
	"github.com/yashsriv/go-nachos/utils"
)

// LaunchUserProcess runs a user program.  Open the executable, load it into
//	memory, and jump to it.
func LaunchUserProcess(filename string) {
	var space = &ProcessAddressSpace{}
	space.Init(filename)

	global.CurrentThread.SetSpace(space)

	space.InitUserModeCPURegisters() // set the initial register values
	space.RestoreContextOnSwitch()   // load page table register

	global.Machine.Run() // jump to the user progam
	utils.Assert(false)  // machine->Run never returns;
	// the address space exits
	// by doing the syscall "exit"
}

var forkFunction = func(arg interface{}) {
	utils.Debug('t', "Now in thread %q\n", global.CurrentThread.Name())

	if global.ThreadToBeDestroyed != nil {
		global.ThreadToBeDestroyed = nil
	}

	if global.CurrentThread.Space() != nil { // if there is an address space
		global.CurrentThread.RestoreUserState() // to restore, do it.
		global.CurrentThread.Space().RestoreContextOnSwitch()
	}

	global.Machine.Run()
}
