package enums

// ThreadStatus is an enumeration for a thread's status
type ThreadStatus int

// Enum values
const (
	JUST_CREATED ThreadStatus = iota
	RUNNING
	READY
	BLOCKED
)
