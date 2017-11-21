// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package enums

// ThreadStatus is an enumeration for a thread's status
type ThreadStatus int

// Enum values
const (
	JustCreated ThreadStatus = iota
	Running
	Ready
	Blocked
)
