// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package interfaces

// ISemaphore defines the interface for a semaphore
type ISemaphore interface {
	Init(string, int)
	Name() string

	P()
	V()
}

// Concrete implementation in threads/synchro/semaphore.go
