//go:build testrunmain
// +build testrunmain

// This file allows us to run the entire program with test coverage enabled, useful for integration tests
package main

import (
	"flag"
	"os"
	"strings"
	"testing"
)

func init() {
	// Parse and remove test flags from os.Args before Cobra sees them
	// This prevents Cobra from complaining about unknown --test.* flags
	testFlags := flag.NewFlagSet("test", flag.ContinueOnError)

	// Register all the test flags we might receive
	testFlags.String("test.coverprofile", "", "")
	testFlags.String("test.run", "", "")
	testFlags.String("test.v", "", "")
	testFlags.Bool("test.short", false, "")
	testFlags.String("test.timeout", "", "")
	testFlags.String("test.cpu", "", "")
	testFlags.String("test.parallel", "", "")
	testFlags.String("test.bench", "", "")
	testFlags.String("test.benchmem", "", "")
	testFlags.String("test.benchtime", "", "")
	testFlags.String("test.cpuprofile", "", "")
	testFlags.String("test.memprofile", "", "")
	testFlags.String("test.memprofilerate", "", "")
	testFlags.String("test.blockprofile", "", "")
	testFlags.String("test.blockprofilerate", "", "")
	testFlags.String("test.mutexprofile", "", "")
	testFlags.String("test.mutexprofilefraction", "", "")
	testFlags.String("test.trace", "", "")
	testFlags.String("test.outputdir", "", "")
	testFlags.String("test.gocoverdir", "", "")

	// Parse test flags and remove them from os.Args
	var newArgs []string
	newArgs = append(newArgs, os.Args[0]) // Keep the program name

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if strings.HasPrefix(arg, "--test.") || strings.HasPrefix(arg, "-test.") {
			// Skip this test flag
			// If it's in the form --flag=value, it's already complete
			// If it's in the form --flag value, skip the next arg too
			if !strings.Contains(arg, "=") && i+1 < len(os.Args) && !strings.HasPrefix(os.Args[i+1], "-") {
				i++ // Skip the value
			}
		} else {
			newArgs = append(newArgs, arg)
		}
	}

	os.Args = newArgs
}

func TestRunMain(t *testing.T) {
	main()
}
