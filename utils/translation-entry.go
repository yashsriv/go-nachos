// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package utils

// TranslationEntry defines an entry in a translation table -- either
// in a page table or a TLB.  Each entry defines a mapping from one
// virtual page to one physical page.
// In addition, there are some extra bits for access control (valid and
// read-only) and some bits for usage information (use and dirty).
// NOTE: Used across multiple packages and hence defeined here
type TranslationEntry struct {
	VirtualPage  uint32 // The page number in virtual memory.
	PhysicalPage uint32 // The page number in real memory (relative to the start of "mainMemory")
	Valid        bool   // If this bit is set, the translation is ignored.
	ReadOnly     bool   // If this bit is set, the user program is not allowed to modify the contents of the page.
	Use          bool   // This bit is set by the hardware every time the page is referenced or modified.
	Dirty        bool   // This bit is set by the hardware every time the page is modified.
}
