// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package interfaces

// IScheduler defines the interface for a scheduler
type IScheduler interface {
	Init()
	MoveThreadToReadyQueue(IThread)
	SelectNextReadyThread() IThread
	ScheduleThread(IThread)
	Print()
}

// Concrete implementation in threads/scheduler.go
