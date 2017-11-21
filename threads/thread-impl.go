// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package threads

import (
	"fmt"

	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/global"
	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/machine"
	"github.com/yashsriv/go-nachos/utils"
)

const stackFencepost int = 0xdeadbeef

// NO_PARENT is the ppid of a process having no parent
const NO_PARENT = -66

// FinishThread is called by ThreadRoot when a thread is done executing the
//	forked procedure.
//
// 	NOTE: we don't immediately de-allocate the thread data structure
//	or the execution stack, because we're still running in the thread
//	and we're still on the stack!  Instead, we set "threadToBeDestroyed",
//	so that ProcessScheduler::ScheduleThread() will call the destructor, once we're
//	running in the context of a different thread.
//
// 	NOTE: we disable interrupts, so that we don't get a time slice
//	between setting threadToBeDestroyed, and going to sleep.
func (t *Thread) FinishThread() {
	global.Interrupt.SetLevel(enums.IntOff)
	utils.Assert(t == global.CurrentThread.(*Thread), "Only currently running thread can be finished")

	utils.Debug('t', "Finishing thread %q\n", t.Name())

	global.ThreadToBeDestroyed = global.CurrentThread.(*Thread)
	t.PutThreadToSleep() // invokes SWITCH
}

// Print prints name of thread to stdout
func (t *Thread) Print() {
	fmt.Printf("%s, ", t.Name())
}

// PutThreadToSleep relinquishes the CPU, because the current thread is blocked
//	waiting on a synchronization variable (Semaphore, Lock, or Condition).
//	Eventually, some thread will wake this thread up, and put it
//	back on the ready queue, so that it can be re-scheduled.
//
//	NOTE: if there are no threads on the ready queue, that means
//	we have no thread to run.  "Interrupt::Idle" is called
//	to signify that we should idle the CPU until the next I/O interrupt
//	occurs (the only thing that could cause a thread to become
//	ready to run).
//
//	NOTE: we assume interrupts are already disabled, because it
//	is called from the synchronization routines which must
//	disable interrupts for atomicity.   We need interrupts off
//	so that there can't be a time slice between pulling the first thread
//	off the ready list, and switching to it.
func (t *Thread) PutThreadToSleep() {
	var nextThread interfaces.IThread

	utils.Assert(t == global.CurrentThread.(*Thread), "Only current thread can be put to sleep")
	utils.Assert(global.Interrupt.GetLevel() == enums.IntOff, "Interrupts should be off when putting thread to sleep")

	utils.Debug('t', "Sleeping thread %q\n", t.Name())

	t.status = enums.BLOCKED
	for nextThread = global.Scheduler.SelectNextReadyThread(); nextThread == nil; {
		global.Interrupt.Idle() // no one to run, wait for an interrupt
		nextThread = global.Scheduler.SelectNextReadyThread()
	}
	// Found new thread

	global.Scheduler.ScheduleThread(nextThread) // returns when we've been signalled

}

// SetStatus sets the thread's status
func (t *Thread) SetStatus(st enums.ThreadStatus) {
	t.status = st
}

// RestoreUserState restores the CPU state of a user program on a context switch.
//
//	Note that a user program thread has *two* sets of CPU registers --
//	one for its state while executing user code, one for its state
//	while executing kernel code.  This routine restores the former.
func (t *Thread) RestoreUserState() {
	for i := 0; i < machine.NumTotalRegs; i++ {
		global.Machine.WriteRegister(i, t.userRegisters[i])
	}
	t.stateRestored = true
}

// SaveUserState saves the CPU state of a user program on a context switch.
//
//	Note that a user program thread has *two* sets of CPU registers --
//	one for its state while executing user code, one for its state
//	while executing kernel code.  This routine saves the former.
func (t *Thread) SaveUserState() {
	if t.stateRestored {
		for i := 0; i < machine.NumTotalRegs; i++ {
			t.userRegisters[i] = global.Machine.ReadRegister(i)
		}
		t.stateRestored = false
	}
}

// ThreadFork invokes function, allowing caller and callee to execute
//	concurrently.
//
//	NOTE: although our definition allows only a single integer argument
//	to be passed to the procedure, it is possible to pass multiple
//	arguments by making them fields of a structure, and passing a pointer
//	to the structure as "arg".
//
// 	Implemented as the following steps:
//		1. Allocate a stack
//		2. Initialize the stack so that a call to SWITCH will
//		cause it to run the procedure
//		3. Put the thread on the ready queue
//
//	"func" is the procedure to run concurrently.
//	"arg" is a single argument to be passed to the procedure.
//----------------------------------------------------------------------
func (t *Thread) ThreadFork(function utils.VoidFunction, arg interface{}) {
	utils.Debug('t', "Forking thread %q with func = 0x%v, arg = %v\n",
		t.Name(), function, arg)

	t.createThreadStack(function, arg)

	oldLevel := global.Interrupt.SetLevel(enums.IntOff)
	global.Scheduler.MoveThreadToReadyQueue(t)
	// MoveThreadToReadyQueue assumes that interrupts
	// are disabled!
	global.Interrupt.SetLevel(oldLevel)
}

// YieldCPU relinquishes the CPU if any other thread is ready to run.
//	If so, put the thread on the end of the ready list, so that
//	it will eventually be re-scheduled.
//
//	NOTE: returns immediately if no other thread on the ready queue.
//	Otherwise returns when the thread eventually works its way
//	to the front of the ready list and gets re-scheduled.
//
//	NOTE: we disable interrupts, so that looking at the thread
//	on the front of the ready list, and switching to it, can be done
//	atomically.  On return, we re-set the interrupt level to its
//	original state, in case we are called with interrupts disabled.
//
// 	Similar to PutThreadToSleep(), but a little different.
func (t *Thread) YieldCPU() {
	oldLevel := global.Interrupt.SetLevel(enums.IntOff)

	utils.Assert(t == global.CurrentThread.(*Thread), "Only current thread can yield CPU")

	utils.Debug('t', "Yielding thread %q\n", t.Name())

	nextThread := global.Scheduler.SelectNextReadyThread()
	if nextThread != nil {
		global.Scheduler.MoveThreadToReadyQueue(t)
		global.Scheduler.ScheduleThread(nextThread)
	}
	global.Interrupt.SetLevel(oldLevel)
}

// Init initializes our thread
func (t *Thread) Init(name string) {
	t.name = name
	t.stateRestored = true
	t.pid = 0
	t.ppid = NO_PARENT
}

// CreateThreadStack allocates and initializes an execution stack.  The stack is
//	initialized with an initial stack frame for ThreadRoot, which:
//		enables interrupts
//		calls (*func)(arg)
//		calls NachOSThread::FinishThread
//
//	"func" is the procedure to be forked
//	"arg" is the parameter to be passed to the procedure
func (t *Thread) createThreadStack(function utils.VoidFunction, arg interface{}) {
	go func() {
		global.Interrupt.Enable()
		function(arg)
		global.CurrentThread.FinishThread()
	}()
}

// Name getter
func (t *Thread) Name() string {
	return t.name
}

// PID getter
func (t *Thread) PID() int {
	return t.pid
}

// PPID getter
func (t *Thread) PPID() int {
	return t.ppid
}

// Space getter
func (t *Thread) Space() interfaces.IProcessAddressSpace {
	return t.space
}

// SetSpace setter
func (t *Thread) SetSpace(space interfaces.IProcessAddressSpace) {
	t.space = space
}
