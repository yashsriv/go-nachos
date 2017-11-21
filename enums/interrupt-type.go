// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package enums

// IntType is an enumeration for type of interrupt
type IntType int

// IntType enums
const (
	TimerInt IntType = iota
	DiskInt
	ConsoleWriteInt
	ConsoleReadInt
	NetworkSendInt
	NetworkRecvInt
)

func (i IntType) String() string {
	switch i {
	case TimerInt:
		return "timer"
	case DiskInt:
		return "disk"
	case ConsoleWriteInt:
		return "console write"
	case ConsoleReadInt:
		return "console read"
	case NetworkSendInt:
		return "network send"
	case NetworkRecvInt:
		return "network recv"
	}
	return "unknown interrupt"
}
