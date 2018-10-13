// +build testrunmain

/*
CODE GENERATED AUTOMATICALLY WITH FIBER COIN CREATOR
AVOID EDITING THIS MANUALLY
*/

// This file allows us to run the entire program with test coverage enabled, useful for integration tests
package main

import "testing"

func init() {
	// Disable flag.Parse because go's test runner calls it instead
	parseFlags = false
}

func TestRunMain(t *testing.T) {
	main()
}
