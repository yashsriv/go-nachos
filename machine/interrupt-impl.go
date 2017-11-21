package machine

import (
	"container/list"
	"fmt"

	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/global"
	"github.com/yashsriv/go-nachos/interfaces"
	"github.com/yashsriv/go-nachos/utils"
)

// Init is used to set up this interrupt
func (interrupt *Interrupt) Init() {
	interrupt.pending = list.New()
	interrupt.level = enums.IntOff
	interrupt.inHandler = false
	interrupt.yieldOnReturn = false
	interrupt.status = enums.SystemMode
}

// changeLevel
// Change interrupts to be enabled or disabled, without advancing
// the simulated time (normally, enabling interrupts advances the time).
//
// Used internally.
//
// "old" -- the old interrupt status
// "now" -- the new interrupt status
func (interrupt *Interrupt) changeLevel(old enums.IntStatus, now enums.IntStatus) {
	interrupt.level = now
	utils.Debug('i', "\tinterrupts: %q -> %q\n", intStatusNames[old], intStatusNames[now])
}

// SetLevel is used to set the interrupt level for our machine
func (interrupt *Interrupt) SetLevel(level enums.IntStatus) enums.IntStatus {
	old := interrupt.level
	utils.Assert(level == enums.IntOff || !interrupt.inHandler)
	interrupt.changeLevel(old, level)
	if level == enums.IntOn && old == enums.IntOff {
		interrupt.OneTick()
	}
	return old
}

// Enable allows enabling interrupts
func (interrupt *Interrupt) Enable() {
	interrupt.SetLevel(enums.IntOn)
}

// OneTick advances simulated time and checks if there are any pending
// interrupts to be called.
//
// Two things can cause OneTick to be called:
//   interrupts are re-enabled
//   a user instruction is executed
func (interrupt *Interrupt) OneTick() {
	old := interrupt.status
	// advance simulated time
	if interrupt.status == enums.SystemMode {
		global.Stats.TotalTicks += utils.SystemTick
		global.Stats.SystemTicks += utils.SystemTick
	} else { // USER_PROGRAM
		global.Stats.TotalTicks += utils.UserTick
		global.Stats.UserTicks += utils.UserTick
	}
	utils.Debug('i', "\n== Tick %d ==\n", global.Stats.TotalTicks)

	// check any pending interrupts are now ready to fire
	interrupt.changeLevel(enums.IntOn, enums.IntOff) // first, turn off interrupts
	// (interrupt handlers run with interrupts disabled)

	for interrupt.checkIfDue(false) {
	} // check for pending interrupts

	interrupt.changeLevel(enums.IntOff, enums.IntOn) // re-enable interrupts
	if interrupt.yieldOnReturn {                     // if the timer device handler asked
		// for a context switch, ok to do it now
		interrupt.yieldOnReturn = false
		interrupt.status = enums.SystemMode // yield is a kernel routine
		global.CurrentThread.YieldCPU()
		interrupt.status = old
	}
}

// YieldOnReturn is called from within an interrupt handler, to cause a context switch
// (for example, on a time slice) in the interrupted thread,
// when the handler returns.
//
// We can't do the context switch here, because that would switch
// out the interrupt handler, and we want to switch out the
// interrupted thread.
func (interrupt *Interrupt) YieldOnReturn() {
	utils.Assert(interrupt.inHandler == true)
	interrupt.yieldOnReturn = true
}

// Idle Routine is called when there is nothing in the ready queue.
//
// Since something has to be running in order to put a thread
// on the ready queue, the only thing to do is to advance
// simulated time until the next scheduled hardware interrupt.
// If there are no pending interrupts, stop.  There's nothing
// more for us to do.
func (interrupt *Interrupt) Idle() {
	utils.Debug('i', "Machine idling; checking for interrupts.\n")
	interrupt.status = enums.IdleMode
	if interrupt.checkIfDue(true) {
		for interrupt.checkIfDue(false) {
		}
		interrupt.yieldOnReturn = false
		interrupt.status = enums.SystemMode
		return
	}
	// if there are no pending interrupts, and nothing is on the ready
	// queue, it is time to stop.   If the console or the network is
	// operating, there are *always* pending interrupts, so this code
	// is not reached.  Instead, the halt must be invoked by the user program.

	utils.Debug('i', "Machine idle.  No interrupts to do.\n")
	fmt.Println("No threads ready or runnable, and no pending interrupts.")
	fmt.Println("Assuming the program completed.")
	interrupt.Halt()
}

// Halt shuts down nachos cleanly, printing out performance statistics
func (interrupt *Interrupt) Halt() {
	fmt.Printf("Machine Halting\n\n")
	global.Stats.Print()
	utils.Cleanup()
}

// Schedule arranges for the CPU to be interrupted when simulated time
// reaches "when".
//
// Implementation: just put it on a sorted list.
//
// NOTE: the Nachos kernel should not call this routine directly.
// Instead, it is only called by the hardware device simulators.
//
// "handler" is the procedure to call when the interrupt occurs
// "arg" is the argument to pass to the procedure
// "when" is how far in the future (in simulated time) the
//   interrupt is to occur
// "type" is the hardware device that generated the interrupt
func (interrupt *Interrupt) Schedule(ipending interfaces.IPendingInterrupt) {
	pending := ipending.(PendingInterrupt)
	utils.Debug('i', "Scheduling interrupt handler the %q at time = %d\n",
		intTypeNames[pending.TypeInt], pending.When)
	utils.Assert(pending.When > 0)
	interrupt.sortedInsert(pending)
}

// GetLevel allows getting level
func (interrupt *Interrupt) GetLevel() enums.IntStatus {
	return interrupt.level
}

// SetStatus is used to set the machine status
func (interrupt *Interrupt) SetStatus(status enums.MachineStatus) {
	interrupt.status = status
}

// GetStatus allows getting machine status
func (interrupt *Interrupt) GetStatus() enums.MachineStatus {
	return interrupt.status
}

// DumpState prints the complete interrupt state - the status, and all interrupts
// that are scheduled to occur in the future.
func (interrupt *Interrupt) DumpState() {
	fmt.Printf("Time: %d, interrupts %v\n", global.Stats.TotalTicks, interrupt.level)
	fmt.Println("Pending interrupts:")
	element := interrupt.pending.Front()
	for element != nil {
		x := element.Value.(PendingInterrupt)
		fmt.Printf("Interrupt handler %q, scheduled at %d\n", intTypeNames[x.TypeInt], x.When)
		element = element.Next()
	}
	fmt.Println("\nEnd of pending interrupts")
}

func (interrupt *Interrupt) sortedInsert(pending PendingInterrupt) {
	element := interrupt.pending.Front()
	if element == nil {
		interrupt.pending.PushFront(pending)
		return
	}
	x := element.Value.(PendingInterrupt)
	for x.When < pending.When {
		element = element.Next()
		if element == nil {
			interrupt.pending.PushBack(pending)
			return
		}
		x = element.Value.(PendingInterrupt)
	}
	interrupt.pending.InsertBefore(pending, element)
}

// checkIfDue checks if an interrupt is scheduled to occur, and if so, fires it off.
//
// Returns:
//   TRUE, if we fired off any interrupt handlers
// Params:
//   "advanceClock" -- if TRUE, there is nothing in the ready queue,
//		so we should simply advance the clock to when the next
//		pending interrupt would occur (if any).  If the pending
//		interrupt is just the time-slice daemon, however, then
//		we're done!
func (interrupt *Interrupt) checkIfDue(advanceClock bool) bool {
	old := interrupt.status

	utils.Assert(interrupt.level == enums.IntOff) // interrupts need to be disabled, to invoke an interrupt handler
	if utils.DebugIsEnabled('i') {
		interrupt.DumpState()
	}
	if interrupt.pending.Front() == nil {
		// no pending interrupts
		return false
	}
	var toOccur = interrupt.pending.Remove(interrupt.pending.Front()).(PendingInterrupt)
	when := toOccur.When

	if advanceClock && when > global.Stats.TotalTicks { // advance the clock
		global.Stats.IdleTicks += (when - global.Stats.TotalTicks)
		global.Stats.TotalTicks = when
	} else if when > global.Stats.TotalTicks { // not time yet, put it back
		interrupt.pending.PushFront(toOccur)
		return false
	}

	// Check if there is nothing more to do, and if so, quit
	if (interrupt.status == enums.IdleMode) && (toOccur.TypeInt == enums.TimerInt) && interrupt.pending.Len() == 0 {
		interrupt.pending.PushFront(toOccur)
		return false
	}

	utils.Debug('i', "Invoking interrupt handler for the %q at time %d\n", intTypeNames[toOccur.TypeInt], toOccur.When)
	if global.Machine != nil {
		global.Machine.DelayedLoad(0, 0)
	}
	interrupt.inHandler = true
	interrupt.status = enums.SystemMode // whatever we were doing,
	// we are now going to be
	// running in the kernel
	(toOccur.Handler)(toOccur.Param) // call the interrupt handler
	interrupt.status = old           // restore the machine status
	interrupt.inHandler = false
	return true
}
