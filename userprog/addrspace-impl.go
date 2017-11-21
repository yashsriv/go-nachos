package userprog

import (
	"encoding/binary"
	"math"
	"os"

	"github.com/yashsriv/go-nachos/global"
	"github.com/yashsriv/go-nachos/machine"
	"github.com/yashsriv/go-nachos/utils"
)

var mainMemoryOffset uint32

func bzero(memory []byte) {
	for i := 0; i < len(memory); i++ {
		memory[i] = 0
	}
}

// Init should be called on a process address space before anything else
// acts as a constructor
func (addrspace *ProcessAddressSpace) Init(filename string) {
	var executable *os.File

	// Open the File
	executable, err := os.Open(filename)
	if err != nil {
		utils.Panic(err)
	}
	defer executable.Close()

	var noffH = NoffHeader{}

	binary.Read(executable, binary.LittleEndian, &noffH)

	utils.Assert(noffH.NoffMagic == NOFFMAGIC, "Executable magic value to confirm if its the right format")

	var offset = mainMemoryOffset
	var size = noffH.Code.Size + noffH.InitData.Size + noffH.UninitData.Size + UserStackSize
	var numVirtualPages = uint32(math.Ceil(float64(size) / float64(machine.PageSize)))
	size = numVirtualPages * machine.PageSize
	utils.Assert(numVirtualPages <= (machine.NumPhysPages-(offset/machine.PageSize)), "There should be enough number of free physical pages")
	utils.Debug('a', "Initializing address space, num pages %d, size %d\n",
		numVirtualPages, size)
	addrspace.kernelPageTable = make([]utils.TranslationEntry, numVirtualPages)
	addrspace.numVirtualPages = numVirtualPages
	for i := uint32(0); i < numVirtualPages; i++ {
		addrspace.kernelPageTable[i] = utils.TranslationEntry{
			VirtualPage:  i,
			PhysicalPage: (offset / machine.PageSize) + i,
			Valid:        true,
			ReadOnly:     false,
			Use:          false,
			Dirty:        false,
		}
	}

	mainMemory := global.Machine.GetMainMemory()
	// Zero out memory
	bzero(mainMemory[offset : offset+size])

	if noffH.Code.Size > 0 {
		utils.Debug('a', "Initializing code segment, at 0x%x, size %d\n", noffH.Code.VirtualAddr, noffH.Code.Size)
		start := noffH.Code.VirtualAddr + offset
		executable.ReadAt(mainMemory[start:start+noffH.Code.Size], int64(noffH.Code.InFileAddr))
	}

	if noffH.InitData.Size > 0 {
		utils.Debug('a', "Initializing data segment, at 0x%x, size %d\n", noffH.InitData.VirtualAddr, noffH.InitData.Size)
		start := noffH.InitData.VirtualAddr + offset
		executable.ReadAt(mainMemory[start:start+noffH.InitData.Size], int64(noffH.InitData.InFileAddr))
	}

	mainMemoryOffset += size

}

// InitUserModeCPURegisters initializes registers
func (addrspace *ProcessAddressSpace) InitUserModeCPURegisters() {
	for i := 0; i < machine.NumTotalRegs; i++ {
		global.Machine.WriteRegister(i, 0)
	}
	// Initial program counter -- must be location of "Start"
	global.Machine.WriteRegister(machine.PCReg, 0)

	// Need to also tell MIPS where next instruction is, because
	// of branch delay possibility
	global.Machine.WriteRegister(machine.NextPCReg, 4)

	// Set the stack register to the end of the address space, where we
	// allocated the stack; but subtract off a bit, to make sure we don't
	// accidentally reference off the end!
	global.Machine.WriteRegister(machine.StackReg, addrspace.numVirtualPages*machine.PageSize-16)
	utils.Debug('a', "Initializing stack register to %d\n", addrspace.numVirtualPages*machine.PageSize-16)

}

// SaveContextOnSwitch saves machine space specific to this addrspace that needs saving
func (addrspace *ProcessAddressSpace) SaveContextOnSwitch() {
}

// RestoreContextOnSwitch restores the machine state so that
//	this address space can run.
//
//      For now, tell the machine where to find the page table.
func (addrspace *ProcessAddressSpace) RestoreContextOnSwitch() {
	global.Machine.SetPageTable(addrspace.kernelPageTable)
}
