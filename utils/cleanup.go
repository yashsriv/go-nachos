package utils

import (
	"fmt"
	"os"
)

// General tools to help during program cleanup
// Registering cleanup functions allows us to do
// a variety of tasks on cleanup depending upon
// requirements

// CleanupFunc is a function which takes no input, has no output and
// has the side effect of freeing resources
type CleanupFunc func()

var cleanupFuncs = make([]CleanupFunc, 0, 5)

// Cleanup is called at the very end of simulation
func Cleanup() {
	fmt.Printf("\nCleaning up...\n")
	freeStuff()
	os.Exit(0)
}

// RegisterCleanup registers a function to be called during cleanup
func RegisterCleanup(a CleanupFunc) {
	cleanupFuncs = append(cleanupFuncs, a)
}

func freeStuff() {
	for _, v := range cleanupFuncs {
		v()
	}
}
