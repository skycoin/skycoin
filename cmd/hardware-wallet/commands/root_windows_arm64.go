//go:build windows && arm64
// +build windows,arm64

// Package commands provides a stub for Windows ARM64 architecture.
//
// Hardware wallet support is disabled on Windows ARM64 due to
// limitations in the github.com/google/gousb library which lacks
// proper CGO/libusb support for this platform.
package commands

import (
	"github.com/spf13/cobra"
)

// RootCmd is a stub for Windows ARM64 - hardware wallet is not supported
var RootCmd = &cobra.Command{
	Use:    "skyhw",
	Short:  "hardware wallet utilities (not available on Windows ARM64)",
	Long:   "Hardware wallet support is not available on Windows ARM64",
	Hidden: true,
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Println("Hardware wallet support is not available on Windows ARM64")
	},
}

// Execute is a stub for Windows ARM64
func Execute() {
	// No-op on Windows ARM64
}
