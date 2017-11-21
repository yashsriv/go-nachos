package userprog

// #include "syscall.h"
import "C"
import (
	"fmt"

	"github.com/yashsriv/go-nachos/console"
	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/global"
	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/machine"
	"github.com/yashsriv/go-nachos/threads/synch"
	"github.com/yashsriv/go-nachos/utils"
)

var readAvail interfaces.ISemaphore
var writeDone interfaces.ISemaphore

var initializedConsoleSemaphores = false

func advanceCounters() {
	// Advance program counters.
	global.Machine.WriteRegister(machine.PrevPCReg, global.Machine.ReadRegister(machine.PCReg))
	global.Machine.WriteRegister(machine.PCReg, global.Machine.ReadRegister(machine.NextPCReg))
	global.Machine.WriteRegister(machine.NextPCReg, global.Machine.ReadRegister(machine.NextPCReg)+4)
}

func convertIntToHex(v uint32, console interfaces.IConsole) {
	if v == 0 {
		return
	}
	convertIntToHex(v/16, console)
	x := v % 16
	if x < 10 {
		writeDone.P()
		console.PutChar(byte('0' + x))
	} else {
		writeDone.P()
		console.PutChar(byte('a' + x - 10))
	}

}

func writeDoneFunc(interface{}) {
	writeDone.V()
}

func readAvailFunc(interface{}) {
	readAvail.V()
}

// Init the exception handler
func Init() {
	global.ExceptionHandler = func(which enums.ExceptionType) {
		typeSyscall := global.Machine.ReadRegister(2)
		// int memval, vaddr, printval, tempval, exp;
		// unsigned printvalus;        // Used for printing in hex
		if !initializedConsoleSemaphores {
			readAvail = &synch.Semaphore{}
			readAvail.Init("read avail", 0)
			writeDone = &synch.Semaphore{}
			writeDone.Init("write done", 1)
			initializedConsoleSemaphores = true
		}
		var console interfaces.IConsole = &console.Console{}
		console.Init("", "", readAvailFunc, writeDoneFunc, 0)

		if which == enums.SyscallException {
			switch typeSyscall {
			case C.SysCall_Halt:
				utils.Debug('a', "Shutdown, initiated by user program.\n")
				global.Interrupt.Halt()
			case C.SysCall_PrintInt:
				printval := int32(global.Machine.ReadRegister(4))
				if printval == 0 {
					writeDone.P()
					console.PutChar('0')
				} else {
					if printval < 0 {
						writeDone.P()
						console.PutChar('-')
						printval = -printval
					}
					tempval := printval
					exp := int32(1)
					for tempval != 0 {
						tempval = tempval / 10
						exp = exp * 10
					}
					exp = exp / 10
					for exp > 0 {
						writeDone.P()
						console.PutChar(byte('0' + (printval / exp)))
						printval = printval % exp
						exp = exp / 10
					}
				}
				advanceCounters()
			case C.SysCall_PrintChar:
				writeDone.P()
				console.PutChar(byte(global.Machine.ReadRegister(4))) // echo it!
				advanceCounters()
			case C.SysCall_PrintString:
				vaddr := global.Machine.ReadRegister(4)
				memval, _ := global.Machine.ReadMem(vaddr, 1)
				for memval != 0 {
					writeDone.P()
					console.PutChar(byte(memval))
					vaddr++
					memval, _ = global.Machine.ReadMem(vaddr, 1)
				}
				advanceCounters()
			case C.SysCall_PrintIntHex:
				printval := global.Machine.ReadRegister(4)
				writeDone.P()
				console.PutChar('0')
				writeDone.P()
				console.PutChar('x')
				if printval == 0 {
					writeDone.P()
					console.PutChar('0')
				} else {
					convertIntToHex(printval, console)
				}
			default:
				fmt.Printf("Unexpected user mode exception %v %v\n", which, typeSyscall)
				utils.Assert(false)
			}
		} else {
			fmt.Printf("Unexpected user mode exception %v %v\n", which, typeSyscall)
			utils.Assert(false)
		}
	}
}
