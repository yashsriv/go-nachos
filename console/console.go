// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package console

import (
	"os"

	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/utils"
)

// Console defines a hardware console device.
// Input and output to the device is simulated by reading
// and writing to UNIX files ("readFile" and "writeFile").
//
// Since the device is asynchronous, the interrupt handler "readAvail"
// is called when a character has arrived, ready to be read in.
// The interrupt handler "writeDone" is called when an output character
// has been "put", so that the next character can be written.
type Console struct {
	readFile     *os.File
	writeFile    *os.File           // UNIX file emulating the display
	writeHandler utils.VoidFunction // Interrupt handler to call when
	// the PutChar I/O completes
	readHandler utils.VoidFunction // Interrupt handler to call when
	// a character arrives from the keyboard
	handlerArg interface{} // argument to be passed to the
	// interrupt handlers
	putBusy bool // Is a PutChar operation in progress?
	// If so, you can't do another one!
	incoming int // Contains the character to be read,

	nextPoll chan byte
}

// Test if our console implements the necessary interface
var _ interfaces.IConsole = &Console{}
