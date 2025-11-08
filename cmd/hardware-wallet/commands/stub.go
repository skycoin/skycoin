//go:build !linux

// Package commands provides a stub for non-Linux platforms where hardware wallet is not supported.
package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// RootCmd is a stub command for non-Linux platforms
var RootCmd = &cobra.Command{
	Use:   "hardware-wallet",
	Short: "Hardware wallet utilities (not supported on this platform)",
	Long:  "Hardware wallet functionality is only supported on Linux",
	Run: func(_ *cobra.Command, _ []string) {
		log.Fatal("Hardware wallet is not supported on this platform")
	},
}

// Execute executes root CLI command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
