// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package global

import (
	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/utils"
)

// MaxProcesses that can be spawned
const MaxProcesses = 1024

// ControlChannel are channels for synchronization of threads
var ControlChannel [MaxProcesses]chan int

// Following are *all* the global instances of interfaces required by
// NachOS

// CurrentThread is a pointer to the current thread
var CurrentThread interfaces.IThread

// ThreadToBeDestroyed is a pointer to a thread to be destroyed
var ThreadToBeDestroyed interfaces.IThread

// Scheduler is an instance of scheduler
var Scheduler interfaces.IScheduler

// Interrupt is an instance of Interrupt
var Interrupt interfaces.IInterrupt

// Stats is an instance of Statistics
var Stats utils.Statistics

// Machine is an instance of Machine
var Machine interfaces.IMachine

// Timer is an instance of Timer
var Timer interfaces.ITimer

// ExceptionHandler is invoked to handle all exceptions
// Defined in package userprog
var ExceptionHandler func(which enums.ExceptionType)

// Console is an instance of console
var Console interfaces.IConsole
