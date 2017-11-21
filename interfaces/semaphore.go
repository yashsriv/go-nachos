package interfaces

// ISemaphore defines the interface for a semaphore
type ISemaphore interface {
	Init(string, int)
	Name() string

	P()
	V()
}

// Concrete implementation in threads/synchro/semaphore.go
