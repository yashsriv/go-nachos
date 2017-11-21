// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/yashsriv/go-nachos/enums"
	"github.com/yashsriv/go-nachos/global"
	"github.com/yashsriv/go-nachos/machine"
	"github.com/yashsriv/go-nachos/threads"
	"github.com/yashsriv/go-nachos/userprog"
	"github.com/yashsriv/go-nachos/utils"
)

func initialize() {
	var randomYield = false
	var debugArgs utils.StringFlag
	var seed utils.Int64Flag
	flag.Var(&debugArgs, "d", "set debug flags")
	flag.Var(&seed, "rs", "seed random number generator")
	var singleStep = flag.Bool("s", false, "debug the user program step by step")

	flag.Parse()

	if debugArgs.IsSet {
		if debugArgs.Value == "" {
			utils.InitDebug("+")
		} else {
			utils.InitDebug(debugArgs.Value)
		}
	}

	if seed.IsSet {
		utils.RandomInit(seed.Value)
		randomYield = true
	}

	global.Stats = utils.Statistics{}

	global.Interrupt = &machine.Interrupt{}
	global.Interrupt.Init()

	global.Scheduler = &threads.Scheduler{}
	global.Scheduler.Init()

	global.Timer = &machine.Timer{}
	global.Timer.Init(func(arg interface{}) {
		if global.Interrupt.GetStatus() != enums.IdleMode {
			global.Interrupt.YieldOnReturn()
		}
	}, nil, randomYield)

	userprog.Init()

	global.Machine = &machine.Machine{}
	if *singleStep {
		global.Machine.EnableDebugging()
	}

	// We didn't explicitly allocate the current thread we are running in.
	// But if it ever tries to yield, we better have a thread object to save
	// its state.
	global.CurrentThread = &threads.Thread{}
	global.CurrentThread.Init("main")
	global.CurrentThread.SetStatus(enums.RUNNING)

	global.Interrupt.Enable()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			// sig is a ^C, handle it
			utils.Cleanup()
		}
	}()
}

func main() {
	var program utils.StringFlag
	flag.Var(&program, "x", "runs a user program")
	initialize()
	if program.IsSet {
		userprog.LaunchUserProcess(program.Value)
	}
	global.CurrentThread.FinishThread()
	utils.Assert(false, "This point of code can never be reached")
}
