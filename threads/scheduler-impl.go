// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package threads

import (
	"container/list"
	"fmt"

	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/global"
	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/utils"
)

// Init initialises the data structures of this scheduler
func (s *Scheduler) Init() {
	s.listOfReadyThreads = list.New()
}

// MoveThreadToReadyQueue marks a thread as ready, but not running.
//	Put it on the ready list, for later scheduling onto the CPU.
//
//	"thread" is the thread to be put on the ready list.
func (s *Scheduler) MoveThreadToReadyQueue(thread interfaces.IThread) {
	utils.Debug('t', "Putting thread %q on ready list.\n", thread.Name())

	thread.SetStatus(enums.READY)
	s.listOfReadyThreads.PushBack(thread)
}

// Print prints the ready list
func (s *Scheduler) Print() {
	fmt.Println("Ready list contents")
	for e := s.listOfReadyThreads.Front(); e != nil; e = e.Next() {
		e.Value.(*Thread).Print()
	}
}

// ScheduleThread dispatches the CPU to nextThread. Save the state of the old thread,
//	and load the state of the new thread, by calling the machine
//	dependent context switch routine, SWITCH.
//
//      Note: we assume the state of the previously running thread has
//	already been changed from running to blocked or ready (depending).
// Side effect:
//	The global variable currentThread becomes nextThread.
//
//	"nextThread" is the thread to be put into the CPU.
//----------------------------------------------------------------------
func (s *Scheduler) ScheduleThread(nextThread interfaces.IThread) {

	oldThread := global.CurrentThread

	if global.CurrentThread.Space() != nil { // if this thread is a user program,
		global.CurrentThread.SaveUserState() // save the user's CPU registers
		global.CurrentThread.Space().SaveContextOnSwitch()
	}

	global.CurrentThread = nextThread             // switch to the next thread
	global.CurrentThread.SetStatus(enums.RUNNING) // nextThread is now running

	utils.Debug('t', "Switching from thread %q to thread %q\n",
		oldThread.Name(), nextThread.Name())

	// This is a machine-dependent assembly language routine defined
	// in switch.s.  You may have to think
	// a bit to figure out what happens after this, both from the point
	// of view of the thread and from the perspective of the "outside world".

	_switch(oldThread, nextThread)

	utils.Debug('t', "Now in thread %q\n", global.CurrentThread.Name())

	// If the old thread gave up the processor because it was finishing,
	// we need to delete its carcass.  Note we cannot delete the thread
	// before now (for example, in NachOSThread::FinishThread()), because up to this
	// point, we were still running on the old thread's stack!
	if global.ThreadToBeDestroyed != nil {
		global.ThreadToBeDestroyed = nil
	}

	if global.CurrentThread.Space != nil { // if there is an address space
		global.CurrentThread.RestoreUserState() // to restore, do it.
		global.CurrentThread.Space().RestoreContextOnSwitch()
	}

}

// SelectNextReadyThread returns the next available thread in the ready queue
func (s *Scheduler) SelectNextReadyThread() interfaces.IThread {
	if s.listOfReadyThreads.Front() == nil {
		utils.Debug('t', "No threads in ready queue\n")
		return nil
	}
	return s.listOfReadyThreads.Remove(s.listOfReadyThreads.Front()).(interfaces.IThread)
}

func _switch(oldThread, nextThread interfaces.IThread) {
	// If I don't this, then the old thread loses context
	// and never switches back again
	if oldThread != nextThread {
		global.ControlChannel[oldThread.PID()] <- 0
		global.ControlChannel[nextThread.PID()] <- 1
	}
}
