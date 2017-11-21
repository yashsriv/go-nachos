// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package machine

import (
	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/utils"
)

// Timer defines a hardware timer
type Timer struct {
	randomize bool
	handler   utils.VoidFunction
	arg       interface{}
}

var _ interfaces.ITimer = &Timer{}

// Implemented in machine/timer-impl.go
