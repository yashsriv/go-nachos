// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package enums

// ExceptionType is an enum to handle exceptions
type ExceptionType int

// All possible exceptions
const (
	NoException ExceptionType = iota
	SyscallException
	PageFaultException
	ReadOnlyException     // Write attempted to page marked "read-only"
	BusErrorException     // Translation resulted in an invalid physical address
	AddressErrorException // Unaligned reference or one that was beyond the end of the address space
	OverflowException     // Integer overflow in add or sub.
	IllegalInstrException // Unimplemented or reserved instr.
	NumExceptionTypes
)

func (t ExceptionType) String() string {
	switch t {
	case NoException:
		return "no exception"
	case SyscallException:
		return "syscall"
	case PageFaultException:
		return "page fault"
	case ReadOnlyException:
		return "page read only"
	case BusErrorException:
		return "bus error"
	case AddressErrorException:
		return "address error"
	case OverflowException:
		return "overflow"
	case IllegalInstrException:
		return "illegal instruction"
	}
	return "unexpected exception"
}
