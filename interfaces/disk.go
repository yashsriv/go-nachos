// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package interfaces

import "github.com/yashsriv/go-nachos/utils"

// IDisk defines the interface for a disk
type IDisk interface {
	Init(string, utils.VoidFunction, interface{})
	Close()

	ReadRequest(int, []byte)
	WriteRequest(int, []byte)

	HandleInterrupt()

	ComputeLatency(int, bool) int
}
