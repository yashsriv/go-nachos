package machine

import (
	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/global"
	"github.com/yashsriv/go-nachos/utils"
)

// Init helps initialize our timer
func (timer *Timer) Init(handler utils.VoidFunction, callArg interface{}, doRandom bool) {
	timer.randomize = doRandom
	timer.handler = handler
	timer.arg = callArg

	global.Interrupt.Schedule(PendingInterrupt{
		func(v interface{}) {
			v.(*Timer).TimerExpired()
		},
		timer,
		global.Stats.TotalTicks + timer.timeOfNextInterrupt(),
		enums.TimerInt,
	})
}

// TimerExpired is used to simulate the interrupt generated by the hardware
// timer device.  Schedule the next interrupt, and invoke the
// interrupt handler.
func (timer *Timer) TimerExpired() {
	// schedule the next timer device interrupt
	global.Interrupt.Schedule(PendingInterrupt{
		func(v interface{}) {
			v.(*Timer).TimerExpired()
		},
		timer,
		global.Stats.TotalTicks + timer.timeOfNextInterrupt(),
		enums.TimerInt,
	})

	// invoke the Nachos interrupt handler for this device
	(timer.handler)(timer.arg)
}

// Return when the hardware timer device will next cause an interrupt.
// If randomize is turned on, make it a (pseudo-)random delay.
func (timer *Timer) timeOfNextInterrupt() int {
	if timer.randomize {
		return 1 + (utils.Random() % (utils.TimerTicks * 2))
	}
	return utils.TimerTicks
}
