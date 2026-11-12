//go:build cgo && !386 && !windows

// Package commands cmd/release/commands/root_cgo.go c4-coin-tool
package commands

import (
	skyhw "github.com/skycoin/skycoin/cmd/hardware-wallet/commands"
)

func init() {
	RootCmd.AddCommand(skyhw.RootCmd)
	skyhw.RootCmd.Use = "skyhw"
	skyhw.RootCmd.Short = "skycoin hardware wallet utilities"
}
