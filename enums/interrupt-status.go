// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package enums

// IntStatus is an enumeration for interrupt status
type IntStatus bool

// IntStatus enums
const (
	IntOff IntStatus = false
	IntOn  IntStatus = true
)

func (s IntStatus) String() string {
	switch s {
	case IntOff:
		return "off"
	case IntOn:
		return "on"
	}
	return "unknown status"
}
