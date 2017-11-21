package synch

import (
	"container/list"

	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/global"
	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/utils"
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

var _ interfaces.ISemaphore = (*Semaphore)(nil)

// Init is used to initialize a semaphore
func (s *Semaphore) Init(debugName string, initialValue int) {
	s.name = debugName
	s.value = initialValue
	s.queue = list.New()
}

// P waits until semaphore value > 0, then decrements.  Checking the
// value and decrementing must be done atomically, so we
// need to disable interrupts before checking the value.
//
// Note that NachOSThread::PutThreadToSleep assumes that interrupts are disabled
// when it is called.
func (s *Semaphore) P() {
	oldLevel := global.Interrupt.SetLevel(enums.IntOff) // disable interrupts
	utils.Debug('s', "In P(). value = %d\n", s.value)
	for s.value == 0 { // semaphore not available, go to sleep
		s.queue.PushBack(global.CurrentThread)
		global.CurrentThread.PutThreadToSleep()
	}
	s.value-- // semaphore available,
	// consume its value
	utils.Debug('s', "After P(). value = %d\n", s.value)

	global.Interrupt.SetLevel(oldLevel) // re-enable interrupts
}

// V increments a semaphore value, waking up a waiter if necessary.
// As with P(), this operation must be atomic, so we need to disable
// interrupts.  ProcessScheduler::MoveThreadToReadyQueue() assumes that threads
// are disabled when it is called.
func (s *Semaphore) V() {
	oldLevel := global.Interrupt.SetLevel(enums.IntOff)

	utils.Debug('s', "In V(). value = %d\n", s.value)
	if s.queue.Front() != nil {
		thread := s.queue.Remove(s.queue.Front()).(interfaces.IThread)
		global.Scheduler.MoveThreadToReadyQueue(thread)
	}
	s.value++
	utils.Debug('s', "After V(). value = %d\n", s.value)
	global.Interrupt.SetLevel(oldLevel)
}

// Name is a getter for the name field
func (s *Semaphore) Name() string {
	return s.name
}
