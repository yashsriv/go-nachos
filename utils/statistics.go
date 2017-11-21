// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package utils

import "fmt"

// Statistics defines the statistics that are to be kept
// about Nachos behavior -- how much time (ticks) elapsed, how
// many user instructions executed, etc.
type Statistics struct {
	TotalTicks             int // Total time running Nachos
	IdleTicks              int // Time spent idle (no threads to run)
	SystemTicks            int // Time spent executing system code
	UserTicks              int // Time spent executing user code
	NumDiskReads           int // number of disk read requests
	NumDiskWrites          int // number of disk write requests
	NumConsoleCharsRead    int // number of characters read from the keyboard
	NumConsoleCharsWritten int // number of characters written to the display
	NumPageFaults          int // number of virtual memory page faults
	NumPacketsSent         int // number of packets sent over the network
	NumPacketsRecvd        int // number of packets received over the network
}

// Print performance metrics, when we've finished everything
// at system shutdown.
func (stats *Statistics) Print() {
	fmt.Printf("Ticks: total %d, idle %d, system %d, user %d\n", stats.TotalTicks,
		stats.IdleTicks, stats.SystemTicks, stats.UserTicks)
	fmt.Printf("Disk I/O: reads %d, writes %d\n", stats.NumDiskReads, stats.NumDiskWrites)
	fmt.Printf("Console I/O: reads %d, writes %d\n", stats.NumConsoleCharsRead,
		stats.NumConsoleCharsWritten)
	fmt.Printf("Paging: faults %d\n", stats.NumPageFaults)
	fmt.Printf("Network I/O: packets received %d, sent %d\n", stats.NumPacketsRecvd,
		stats.NumPacketsSent)
}

// Constants used to reflect the relative time an operation would
// take in a real system.  A "tick" is a just a unit of time -- if you
// like, a microsecond.
//
// Since Nachos kernel code is directly executed, and the time spent
// in the kernel measured by the number of calls to enable interrupts,
// these time constants are none too exact.
const (
	UserTick     int = 1   // advance for each user-level instruction
	SystemTick   int = 10  // advance each time interrupts are enabled
	RotationTime int = 500 // time disk takes to rotate one sector
	SeekTime     int = 500 // time disk takes to seek past one track
	ConsoleTime  int = 100 // time to read or write one character
	NetworkTime  int = 100 // time to send or receive one packet
	TimerTicks   int = 100 // (average) time between timer interrupts
)
