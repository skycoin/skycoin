//go:build tinygo

// Package commands implements the skycoin explorer.
//
// The explorer renders its pages with html/template. TinyGo's reflect package
// does not implement reflect.Type.NumOut, which html/template requires when it
// validates template functions during initialization, so importing it panics
// at startup. The explorer is therefore stubbed out in TinyGo builds; this
// RootCmd reports the limitation when invoked.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RootCmd is the skycoin blockchain explorer command. Under TinyGo it is a
// stub that reports that the explorer is unavailable in this build.
var RootCmd = &cobra.Command{
	Use:   "explorer",
	Short: "skycoin blockchain explorer (unavailable in TinyGo builds)",
	RunE: func(_ *cobra.Command, _ []string) error {
		return fmt.Errorf("the blockchain explorer is not available in TinyGo builds")
	},
}
