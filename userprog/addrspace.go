// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package userprog

import (
	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/utils"
)

// UserStackSize is the size of the user stack
const UserStackSize = 1024

// ProcessAddressSpace is a data structure to keep track of existing user programs
type ProcessAddressSpace struct {
	kernelPageTable []utils.TranslationEntry
	numVirtualPages uint32
}

// Check if ProcessAddressSpace implements IProcessAddressSpace
var _ interfaces.IProcessAddressSpace = &ProcessAddressSpace{}

// Implemented in addrspace-impl.go
