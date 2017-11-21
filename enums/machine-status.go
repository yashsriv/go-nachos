// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package enums

// MachineStatus is an enumeration for machine status
type MachineStatus int

// MachineStatus enums
const (
	IdleMode MachineStatus = iota
	SystemMode
	UserMode
)
