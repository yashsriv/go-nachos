// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package machine

import (
	"encoding/binary"

	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/global"
	"github.com/yashsriv/go-nachos/utils"
)

// ReadMem reads "size" (1, 2, or 4) bytes of virtual memory at "addr" into
// the location pointed to by "value".
//
// Returns FALSE if the translation step from virtual to physical memory
// failed.
//
// "addr" -- the virtual address to read from
// "size" -- the number of bytes to read (1, 2, or 4)
// "value" -- the place to write the result
func (m *Machine) ReadMem(addr uint32, size int) (value uint32, success bool) {
	utils.Debug('a', "Reading VA 0x%d, size %d\n", addr, size)

	physicalAddress, exception := m.Translate(addr, size, false)
	if exception != enums.NoException {
		global.Machine.RaiseException(exception, addr)
		success = false
		return
	}
	switch size {
	case 1:
		data := m.MainMemory[physicalAddress]
		value = uint32(data)
	case 2:
		data := m.MainMemory[physicalAddress : physicalAddress+2]
		value = uint32(binary.LittleEndian.Uint16(data))
	case 4:
		data := m.MainMemory[physicalAddress : physicalAddress+4]
		value = binary.LittleEndian.Uint32(data)
	default:
		utils.Assert(false, "Unsupported size for reading from memory")
	}

	utils.Debug('a', "\tvalue read = %8.8x\n", value)
	success = true

	return
}

// WriteMem writes "size" (1, 2, or 4) bytes of the contents of "value" into
// virtual memory at location "addr".
//
//    Returns FALSE if the translation step from virtual to physical memory
//    failed.
//
// "addr" -- the virtual address to write to
// "size" -- the number of bytes to be written (1, 2, or 4)
// "value" -- the data to be written
func (m *Machine) WriteMem(addr uint32, size int, value uint32) bool {

	utils.Debug('a', "Writing VA 0x%x, size %d, value 0x%x\n", addr, size, value)

	physicalAddress, exception := m.Translate(addr, size, true)
	if exception != enums.NoException {
		global.Machine.RaiseException(exception, addr)
		return false
	}
	switch size {
	case 1:
		m.MainMemory[physicalAddress] = byte(value & 0xff)

	case 2:
		binary.LittleEndian.PutUint16(m.MainMemory[physicalAddress:physicalAddress+2], uint16(value&0xffff))

	case 4:
		binary.LittleEndian.PutUint32(m.MainMemory[physicalAddress:physicalAddress+4], uint32(value))

	default:
		utils.Assert(false, "Unsupported size for writing to memory")
	}

	return true
}

// Translate translates a virtual address into a physical address, using
//	either a page table or a TLB.  Check for alignment and all sorts
//	of other errors, and if everything is ok, set the use/dirty bits in
//	the translation table entry, and store the translated physical
//	address in "physAddr".  If there was an error, returns the type
//	of the exception.
//
//	"virtAddr" -- the virtual address to translate
//	"physAddr" -- the place to store the physical address
//	"size" -- the amount of memory being read or written
// 	"writing" -- if TRUE, check the "read-only" bit in the TLB
func (m *Machine) Translate(virtAddr uint32, size int, writing bool) (physAddr uint32, exception enums.ExceptionType) {

	if writing {
		utils.Debug('a', "\tTranslate 0x%x, %s: ", virtAddr, "write")
	} else {
		utils.Debug('a', "\tTranslate 0x%x, %s: ", virtAddr, "read")
	}

	// check for alignment errors
	if ((size == 4) && (virtAddr&0x3) != 0) || ((size == 2) && (virtAddr&0x1) != 0) {
		utils.Debug('a', "alignment problem at %d, size %d!\n", virtAddr, size)
		exception = enums.AddressErrorException
		return
	}

	// TODO: When I add a TLB if I add one
	// we must have either a TLB or a page table, but not both!
	// ASSERT(tlb == NULL || KernelPageTable == NULL);
	utils.Assert(m.KernelPageTable != nil, "KernelPageTable should not be nil")

	// calculate the virtual page number, and offset within the page,
	// from the virtual address
	vpn := virtAddr / PageSize
	offset := virtAddr % PageSize

	if vpn >= uint32(len(m.KernelPageTable)) {
		utils.Debug('a', "virtual page # %d too large for page table size %d!\n",
			virtAddr, len(m.KernelPageTable))
		exception = enums.AddressErrorException
		return
	} else if !m.KernelPageTable[vpn].Valid {
		utils.Debug('a', "virtual page # %d missing!\n",
			vpn)
		exception = enums.PageFaultException
		return
	}
	entry := &m.KernelPageTable[vpn]

	if entry.ReadOnly && writing { // trying to write to a read-only page
		utils.Debug('a', "%d mapped read-only at %d in page table!\n", virtAddr, entry.PhysicalPage)
		exception = enums.ReadOnlyException
		return
	}
	pageFrame := entry.PhysicalPage

	// if the pageFrame is too big, there is something really wrong!
	// An invalid translation was loaded into the page table or TLB.
	if pageFrame >= NumPhysPages {
		utils.Debug('a', "*** frame %d > %d!\n", pageFrame, NumPhysPages)
		exception = enums.BusErrorException
		return
	}
	entry.Use = true // set the use, dirty bits
	if writing {
		entry.Dirty = true
	}
	physAddr = pageFrame*PageSize + offset
	utils.Assert(physAddr >= 0 && (physAddr+uint32(size) <= MemorySize), "Memory address should be within memory limits")
	utils.Debug('a', "phys addr = 0x%x\n", physAddr)
	exception = enums.NoException
	return
}
