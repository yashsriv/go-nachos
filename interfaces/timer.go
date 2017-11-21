package interfaces

import "github.com/yashsriv/go-nachos/utils"

// ITimer defines an interface for the timer
type ITimer interface {
	Init(utils.VoidFunction, interface{}, bool)
	TimerExpired()
}

// Concrete implementation in machine/timer.go
