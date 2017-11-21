package interfaces

import "github.com/yashsriv/go-nachos/utils"

// IConsole defines the interface for a console device
type IConsole interface {
	Init(string, string, utils.VoidFunction, utils.VoidFunction, interface{})
	PutChar(byte)
	GetChar() (byte, error)

	WriteDone()
	CheckCharAvail()
	Close()
}

// Concrete implementation in console/console.go
