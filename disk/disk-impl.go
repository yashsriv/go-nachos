// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package disk

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/global"
	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/machine"
	"github.com/yashsriv/go-nachos/utils"
)

const magicNumber int32 = 0x456789ab

// Constants for the disk
const (
	SectorSize      int = 128 // number of bytes per disk sector
	SectorsPerTrack int = 32  // number of sectors per disk track
	NumTracks       int = 32  // number of tracks per disk
	NumSectors      int = (SectorsPerTrack * NumTracks)
)

var diskDone = func(arg interface{}) {
	var disk = arg.(interfaces.IDisk)
	disk.HandleInterrupt()
}

// Init initializes a simulated disk.  Open the UNIX file (creating it
//	if it doesn't exist), and check the magic number to make sure it's
// 	ok to treat it as Nachos disk storage.
//
//	"name" -- text name of the file simulating the Nachos disk
//	"callWhenDone" -- interrupt handler to be called when disk read/write
//	   request completes
//	"callArg" -- argument to pass the interrupt handler
func (d *Disk) Init(name string, callWhenDone utils.VoidFunction, callArg interface{}) {

	utils.Debug('d', "Initializing the disk, 0x%v %v\n", callWhenDone, callArg)
	d.handler = callWhenDone
	d.handlerArg = callArg
	d.lastSector = 0
	d.bufferInit = 0

	var err, err1 error
	if d.file, err = os.OpenFile(name, os.O_RDWR, 0); err == nil {
		var magicNum int32

		err1 = binary.Read(d.file, binary.LittleEndian, &magicNum)
		if err1 != nil {
			utils.Panic(err1)
		}
		utils.Assert(magicNum == magicNumber)
		d.active = false
		return
	}

	d.file, err1 = os.Create(name)
	if err1 != nil {
		utils.Panic(err1)
	}
	binary.Write(d.file, binary.LittleEndian, magicNumber)
	// need to write at end of file, so that reads will not return EOF
	d.file.Seek(-8, 2)
	binary.Write(d.file, binary.LittleEndian, 0)
	d.active = false
}

// Close performs cleanup operations
func (d *Disk) Close() {
	if d != nil {
		d.file.Close()
	}
}

// printSector dumps the data in a disk
// read/write request, for debugging.
func printSector(writing bool, sector int, data []byte) {
	var p = make([]int32, len(data)/binary.Size(int32(0)))
	for i := 0; i < len(data); i += 8 {
		p[i/8] = int32(binary.LittleEndian.Uint32(data[i : i+8]))
	}

	if writing {
		fmt.Printf("Writing sector: %d\n", sector)
	} else {
		fmt.Printf("Reading sector: %d\n", sector)
	}
	for i := 0; i < len(p); i++ {
		fmt.Printf("%x ", p[i])
	}
	fmt.Printf("\n")
}

// ReadRequest simulates a request to read a single disk sector
//	   Do the read immediately to the UNIX file
//	   Set up an interrupt handler to be called later,
//	      that will notify the caller when the simulator says
//	      the operation has completed.
//
//	Note that a disk only allows an entire sector to be read,
//	not part of a sector.
//
//	"sectorNumber" -- the disk sector to read
//	"data" -- the buffer to hold the incoming bytes
func (d *Disk) ReadRequest(sectorNumber int, data []byte) {
	var magicSize = int(binary.Size(magicNumber))
	ticks := d.ComputeLatency(sectorNumber, false)

	utils.Assert(!d.active) // only one request at a time
	utils.Assert((sectorNumber >= 0) && (sectorNumber < NumSectors))

	utils.Debug('d', "Reading from sector %d\n", sectorNumber)
	d.file.Seek(int64(SectorSize*sectorNumber+magicSize), 0)
	if n, err := d.file.Read(data); err != nil {
		utils.Panic(err)
	} else if n != SectorSize {
		utils.Panic(errors.New("Number of bytes read should be equal to SectorSize"))
	}
	if utils.DebugIsEnabled('d') {
		printSector(false, sectorNumber, data)
	}

	d.active = true
	d.updateLast(sectorNumber)
	global.Stats.NumDiskReads++
	global.Interrupt.Schedule(machine.PendingInterrupt{
		diskDone,
		d,
		global.Stats.TotalTicks + ticks,
		enums.DiskInt,
	})
}

// WriteRequest simulates a request to write a single disk sector
//	   Do the write immediately to the UNIX file
//	   Set up an interrupt handler to be called later,
//	      that will notify the caller when the simulator says
//	      the operation has completed.
//
//	Note that a disk only allows an entire sector to be written,
//	not part of a sector.
//
//	"sectorNumber" -- the disk sector to write
//	"data" -- the bytes to be written, the buffer to hold the incoming bytes
func (d *Disk) WriteRequest(sectorNumber int, data []byte) {
	var magicSize = binary.Size(magicNumber)
	ticks := d.ComputeLatency(sectorNumber, true)

	utils.Assert(!d.active)
	utils.Assert((sectorNumber >= 0) && (sectorNumber < NumSectors))

	utils.Debug('d', "Writing to sector %d\n", sectorNumber)
	d.file.Seek(int64(SectorSize*sectorNumber+magicSize), 0)
	if n, err := d.file.Write(data); err != nil {
		utils.Panic(err)
	} else if n != SectorSize {
		utils.Panic(errors.New("Number of bytes written should be equal to SectorSize"))
	}
	if utils.DebugIsEnabled('d') {
		printSector(true, sectorNumber, data)
	}

	d.active = true
	d.updateLast(sectorNumber)
	global.Stats.NumDiskWrites++
	global.Interrupt.Schedule(machine.PendingInterrupt{
		diskDone,
		d,
		global.Stats.TotalTicks + ticks,
		enums.DiskInt,
	})
}

// HandleInterrupt is called when it is time to invoke
// the disk interrupt handler, to tell the Nachos
// kernel that the disk request is done.
func (d *Disk) HandleInterrupt() {
	d.active = false
	d.handler(d.handlerArg)
}

// timeToSeek returns how long it will take to position the disk head over the correct
//	track on the disk.  Since when we finish seeking, we are likely
//	to be in the middle of a sector that is rotating past the head,
//	we also return how long until the head is at the next sector boundary.
//
//   	Disk seeks at one track per SeekTime ticks (cf. stats.h)
//   	and rotates at one sector per RotationTime ticks
func (d *Disk) timeToSeek(newSector int) (int, int) {
	newTrack := newSector / SectorsPerTrack
	oldTrack := d.lastSector / SectorsPerTrack
	var seek int
	if newTrack > oldTrack {
		seek = (newTrack - oldTrack) * utils.SeekTime
	} else {
		seek = (oldTrack - newTrack) * utils.SeekTime
	}
	// how long will seek take?
	over := (global.Stats.TotalTicks + seek) % utils.RotationTime
	// will we be in the middle of a sector when
	// we finish the seek?

	if over > 0 { // if so, need to round up to next full sector
		return seek, utils.RotationTime - over
	}
	return seek, 0
}

// updateLast keep track of the most recently requested sector. So we can know
// what is in the track buffer.
func (d *Disk) updateLast(newSector int) {
	seek, rotate := d.timeToSeek(newSector)
	if seek != 0 {
		d.bufferInit = global.Stats.TotalTicks + seek + rotate
	}
	d.lastSector = newSector
	utils.Debug('d', "Updating last sector = %d, %d\n", d.lastSector, d.bufferInit)
}

// ComputeLatency returns how long will it take to read/write a disk sector, from
//	the current position of the disk head.
//
//   	Latency = seek time + rotational latency + transfer time
//   	Disk seeks at one track per SeekTime ticks (cf. stats.h)
//   	and rotates at one sector per RotationTime ticks
//
//   	To find the rotational latency, we first must figure out where the
//   	disk head will be after the seek (if any).  We then figure out
//   	how long it will take to rotate completely past newSector after
//	that point.
//
//   	The disk also has a "track buffer"; the disk continuously reads
//   	the contents of the current disk track into the buffer.  This allows
//   	read requests to the current track to be satisfied more quickly.
//   	The contents of the track buffer are discarded after every seek to
//   	a new track.
func (d *Disk) ComputeLatency(newSector int, writing bool) int {
	seek, rotation := d.timeToSeek(newSector)
	timeAfter := global.Stats.TotalTicks + seek + rotation

	// check if track buffer applies
	if (writing == false) && (seek == 0) &&
		(((timeAfter - d.bufferInit) / utils.RotationTime) >
			moduloDiff(newSector, d.bufferInit/utils.RotationTime)) {
		utils.Debug('d', "Request latency = %d\n", utils.RotationTime)
		return utils.RotationTime // time to transfer sector from the track buffer
	}

	rotation += moduloDiff(newSector, timeAfter/utils.RotationTime) * utils.RotationTime

	utils.Debug('d', "Request latency = %d\n", seek+rotation+utils.RotationTime)
	return (seek + rotation + utils.RotationTime)
}

func moduloDiff(to int, from int) int {
	toOffset := to % SectorsPerTrack
	fromOffset := from % SectorsPerTrack

	return ((toOffset - fromOffset) + SectorsPerTrack) % SectorsPerTrack
}
