// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package interfaces

// IProcessAddressSpace defines the interface for an address space
type IProcessAddressSpace interface {
	Init(string)
	InitUserModeCPURegisters()
	RestoreContextOnSwitch()
	SaveContextOnSwitch()
}

// Concrete implementation in userprog/addrspace.go
