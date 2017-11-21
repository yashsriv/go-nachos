package interfaces

// IProcessAddressSpace defines the interface for an address space
type IProcessAddressSpace interface {
	Init(string)
	InitUserModeCPURegisters()
	RestoreContextOnSwitch()
	SaveContextOnSwitch()
}

// Concrete implementation in userprog/addrspace.go
