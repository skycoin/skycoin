// Package apputil src/util/apputil/apputil.go c0-com-util
//
/*
Package apputil provides utility methods for cmd applications
*/
package apputil

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

// CatchInterrupt catches CTRL-C and closes the quit channel if it occurs.
// If CTRL-C is called again, every goroutine stack is dumped and the process exits,
// so that shutdown hangs can be diagnosed.
func CatchInterrupt(quit chan<- struct{}) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	signal.Stop(sigchan)
	close(quit)

	// If ctrl-c is called again, dump the goroutines so that the program state can be examined.
	// Ctrl-c would be called again if program shutdown was stuck.
	go CatchInterruptPanic()
}

// CatchInterruptPanic catches os.Interrupt, prints every goroutine stack and exits.
func CatchInterruptPanic() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	signal.Stop(sigchan)
	PrintProgramStatus()
	os.Exit(2)
}

// CatchDebug catches SIGUSR1 and prints internal program state
func CatchDebug() {
	sigchan := make(chan os.Signal, 1)
	//signal.Notify(sigchan, syscall.SIGUSR1)
	signal.Notify(sigchan, syscall.Signal(0xa)) // SIGUSR1 = Signal(0xa)
	for range sigchan {
		PrintProgramStatus()
	}
}

// PrintProgramStatus prints all goroutine data to stdout
func PrintProgramStatus() {
	p := pprof.Lookup("goroutine")
	if err := p.WriteTo(os.Stdout, 2); err != nil {
		fmt.Println("ERROR:", err)
		return
	}
}
