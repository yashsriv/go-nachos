// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package utils

import (
	"errors"
	"fmt"
	"strings"
)

var debugSym = ""

// General tools to help in debugging.
// Specific debug symbols can be enabled
// based on requirements so that only selective
// debug statements are executed.
//
// For example, if only the debug symbol "g" was enabled,
// utils.Debug('g', "Hello World")  // Would print message on command line
// utils.Debug('f', "Bye!")  // Would not print message on command line

// InitDebug initializes flags which are enabled
//   symbols is the set of symbols which are enabled
func InitDebug(symbols string) {
	debugSym = symbols
}

// DebugIsEnabled checks if debug is enabled for a particular rune
func DebugIsEnabled(symbol rune) bool {
	if debugSym == "+" {
		return true
	}
	return strings.ContainsRune(debugSym, symbol)
}

// Assert is a function used for assertions. If the assertion fails,
// an error message is printed with a stack trace to enable debugging
func Assert(assertion bool, message string) {
	if !assertion {
		Panic(errors.New(message))
	}
}

// Panic frees memory stuff to avoid leaks and then panics
// Can be called from the kernel itself for panicking.
func Panic(err error) {
	freeStuff()
	panic(err)
}

// Debug prints statement if debugging is enabled for the particular rune
// The fmtString argument is passed as is to fmt.Printf function alongwith
// the remaining arguments, so the the usage remains same as fmt.Printf
func Debug(symbol rune, fmtString string, args ...interface{}) {
	if DebugIsEnabled(symbol) {
		fmt.Printf(fmtString, args...)
	}
}
