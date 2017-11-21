// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package threads

import (
	"container/list"

	"github.com/yashsriv/go-nachos/interfaces"
)

// Scheduler defines the scheduler/dispatcher abstraction --
// the data structures and operations needed to keep track of which
// thread is running, and which threads are ready but not running.
type Scheduler struct {
	listOfReadyThreads *list.List
}

var _ interfaces.IScheduler = (*Scheduler)(nil)
