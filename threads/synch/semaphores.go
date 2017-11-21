// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package synch

import (
	"container/list"

	"github.com/yashsriv/go-nachos/interfaces"
)

// Semaphore is class whose value is a non-negative
// integer.  The semaphore has only two operations P() and V():
//
//	P() -- waits until value > 0, then decrement
//
//	V() -- increment, waking up a thread waiting in P() if necessary
//
// Note that the interface does *not* allow a thread to read the value of
// the semaphore directly -- even if you did read the value, the
// only thing you would know is what the value used to be.  You don't
// know what the value is now, because by the time you get the value
// into a register, a context switch might have occurred,
// and some other thread might have called P or V, so the true value might
// now be different.
type Semaphore struct {
	name  string
	value int
	queue *list.List
}

var _ interfaces.ISemaphore = &Semaphore{}

// Implemented in threads/synch/semaphores-impl.go
