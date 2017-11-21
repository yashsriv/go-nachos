package enums

// MachineStatus is an enumeration for machine status
type MachineStatus int

// MachineStatus enums
const (
	IdleMode MachineStatus = iota
	SystemMode
	UserMode
)
