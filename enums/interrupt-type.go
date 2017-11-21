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
