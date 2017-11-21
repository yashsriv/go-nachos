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
