//go:build 386
// +build 386

// Package commands provides a stub for 386 architecture.
//
// Hardware wallet support is disabled on 386 architecture due to
// limitations in the github.com/google/gousb library which lacks proper
// 32-bit support.
package commands

import (
	"github.com/spf13/cobra"
)

// RootCmd is a stub for 386 architecture - hardware wallet is not supported
var RootCmd = &cobra.Command{
	Use:    "skyhw",
	Short:  "hardware wallet utilities (not available on 386)",
	Long:   "Hardware wallet support is not available on 386 architecture",
	Hidden: true,
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Println("Hardware wallet support is not available on 386 architecture")
	},
}

// Execute is a stub for 386 architecture
func Execute() {
	// No-op on 386
}
