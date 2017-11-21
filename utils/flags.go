package utils

import (
	"fmt"
	"strconv"
)

// General tools to help with dealing with command line flags.
// Custom flag types store special IsSet variable to help decide
// which flags have been set and which haven't

// StringFlag represents a custom type for string flags
type StringFlag struct {
	IsSet bool
	Value string
}

// Int64Flag represents a custom type for int64 flags
type Int64Flag struct {
	IsSet bool
	Value int64
}

// Set allows setting the value
func (sf *StringFlag) Set(x string) error {
	sf.Value = x
	sf.IsSet = true
	return nil
}

// String give a string representation of this flag
func (sf *StringFlag) String() string {
	return sf.Value
}

// Set allows setting the value
func (sf *Int64Flag) Set(x string) error {
	var err error
	sf.Value, err = strconv.ParseInt(x, 10, 64)
	if err != nil {
		return err
	}
	sf.IsSet = true
	return nil
}

// String give a string representation of this flag
func (sf *Int64Flag) String() string {
	return fmt.Sprintf("%d", sf.Value)
}
