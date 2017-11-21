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

// ExceptionNames are human readable names for exceptions for debugging purposes
var ExceptionNames = map[ExceptionType]string{
	NoException:           "no exception",
	SyscallException:      "syscall",
	PageFaultException:    "page fault",
	ReadOnlyException:     "page read only",
	BusErrorException:     "bus error",
	AddressErrorException: "address error",
	OverflowException:     "overflow",
	IllegalInstrException: "illegal instruction",
}
