package userprog

// NOFFMAGIC is a magic number denoting Nachos object code file
const NOFFMAGIC uint32 = 0xbadfad

// Segment of the file
type Segment struct {
	VirtualAddr uint32 /* location of segment in virt addr space */
	InFileAddr  uint32 /* location of segment in this file */
	Size        uint32 /* size of segment */
}

// NoffHeader contains info about the file
type NoffHeader struct {
	NoffMagic  uint32  /* should be NOFFMAGIC */
	Code       Segment /* executable code segment */
	InitData   Segment /* initialized data segment */
	UninitData Segment /* uninitialized data segment --
	 * should be zero'ed before use
	 */
}

