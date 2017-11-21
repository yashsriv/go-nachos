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
