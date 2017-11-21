package console

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/global"
	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/machine"
	"github.com/yashsriv/go-nachos/utils"
)

var consoleReadPoll = func(arg interface{}) {
	var console = arg.(interfaces.IConsole)
	console.CheckCharAvail()
}

var consoleWriteDone = func(arg interface{}) {
	var console = arg.(interfaces.IConsole)
	console.WriteDone()
}

// Init initializes the simulation of a hardware console device.
//
//	"readFile" -- UNIX file simulating the keyboard (NULL -> use stdin)
//	"writeFile" -- UNIX file simulating the display (NULL -> use stdout)
// 	"readAvail" is the interrupt handler called when a character arrives
//		from the keyboard
// 	"writeDone" is the interrupt handler called when a character has
//		been output, so that it is ok to request the next char be
//		output
func (c *Console) Init(readFile string, writeFile string, readAvail utils.VoidFunction, writeDone utils.VoidFunction, callArg interface{}) {
	var err error
	if readFile == "" {
		c.readFile = os.Stdin
	} else {
		c.readFile, err = os.Open(readFile)
		if err != nil {
			utils.Panic(err)
		}
	}
	if writeFile == "" {
		c.writeFile = os.Stdout
	} else {
		c.writeFile, err = os.Create(writeFile)
		if err != nil {
			utils.Panic(err)
		}
	}

	c.writeHandler = writeDone
	c.readHandler = readAvail
	c.handlerArg = callArg
	c.putBusy = false
	c.incoming = -1
	c.nextPoll = make(chan byte)

	go pollFile(c)

	global.Interrupt.Schedule(machine.PendingInterrupt{
		consoleReadPoll,
		c,
		global.Stats.TotalTicks + utils.ConsoleTime,
		enums.ConsoleReadInt,
	})
}

// Close acts as explicit destructor fr this interface
func (c *Console) Close() {
	if c.readFile != nil && c.readFile != os.Stdin {
		c.readFile.Close()
	}
	if c.writeFile != nil && c.writeFile != os.Stdout {
		c.writeFile.Close()
	}
}

// CheckCharAvail is periodically called to check if a character is available for
//	input from the simulated keyboard (eg, has it been typed?).
//
//	Only read it in if there is buffer space for it (if the previous
//	character has been grabbed out of the buffer by the Nachos kernel).
//	Invoke the "read" interrupt handler, once the character has been
//	put into the buffer.
func (c *Console) CheckCharAvail() {

	// schedule the next time to poll for a char
	global.Interrupt.Schedule(machine.PendingInterrupt{
		consoleReadPoll,
		c,
		global.Stats.TotalTicks + utils.ConsoleTime,
		enums.ConsoleReadInt,
	})

	c.nextPoll <- 1

}

// WriteDone is an internal routine called when it is time to invoke the interrupt
//	handler to tell the Nachos kernel that the output character has
//	completed.
func (c *Console) WriteDone() {
	c.putBusy = false
	global.Stats.NumConsoleCharsWritten++
	c.writeHandler(c.handlerArg)
}

// GetChar is used to read a character from the input buffer, if there is any there.
//	Either return the character, or EOF if none buffered.
func (c *Console) GetChar() (byte, error) {
	if c.incoming == -1 {
		return 0, io.EOF
	}
	ch := byte(c.incoming)

	c.incoming = -1
	return ch, nil
}

// PutChar is used to write a character to the simulated display, schedule an interrupt
//	to occur in the future, and return.
func (c *Console) PutChar(ch byte) {
	utils.Assert(c.putBusy == false, "Console should not be busy putting another character")
	var charray = make([]byte, 1)
	charray[0] = ch
	if _, err := c.writeFile.Write(charray); err != nil {
		utils.Panic(err)
	}
	c.putBusy = true
	global.Interrupt.Schedule(machine.PendingInterrupt{
		consoleWriteDone,
		c,
		global.Stats.TotalTicks + utils.ConsoleTime,
		enums.ConsoleWriteInt,
	})
}

func readFile(ch chan byte, f *os.File) {
	var char = make([]byte, 1)

	for {
		// otherwise, read character and tell user about it
		if n, err := f.Read(char); err != nil {
			utils.Panic(err)
		} else if n != 1 {
			utils.Panic(errors.New("Not enough characters to read from file"))
		}
		ch <- char[0]
	}
}

func pollFile(c *Console) {
	fileChan := make(chan byte)
	go readFile(fileChan, c.readFile)
	for {
		<-c.nextPoll
		// do nothing if character is already buffered
		if c.incoming == -1 {
			select {
			case x := <-fileChan:
				c.incoming = int(x)
				global.Stats.NumConsoleCharsRead++
				c.readHandler(c.handlerArg)
			case <-time.After(0):

			}
		}
	}
}
