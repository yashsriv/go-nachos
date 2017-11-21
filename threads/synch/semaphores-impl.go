package synch

import (
	"container/list"

	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/global"
	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/utils"
)

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
