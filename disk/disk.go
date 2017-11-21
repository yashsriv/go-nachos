package disk

import (
	"os"

	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/utils"
)

// Disk is a data structure to emulate a physical disk.  A physical disk
// can accept (one at a time) requests to read/write a disk sector;
// when the request is satisfied, the CPU gets an interrupt, and
// the next request can be sent to the disk.
//
// Disk contents are preserved across machine crashes, but if
// a file system operation (eg, create a file) is in progress when the
// system shuts down, the file system may be corrupted.
type Disk struct {
	file       *os.File           // UNIX file number for simulated disk
	handler    utils.VoidFunction // Interrupt handler, to be invoked when any disk request finishes
	handlerArg interface{}        // Argument to interrupt handler
	active     bool               // Is a disk operation in progress?
	lastSector int                // The previous disk request
	bufferInit int                // When the track buffer started
}

var _ interfaces.IDisk = (*Disk)(nil)
