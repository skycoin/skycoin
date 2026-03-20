//go:build cgo && !386 && !(windows && arm64)

package commands

import (
	skyhw "github.com/skycoin/skycoin/cmd/hardware-wallet/commands"
)

func init() {
	RootCmd.AddCommand(skyhw.RootCmd)
	skyhw.RootCmd.Use = "skyhw"
	skyhw.RootCmd.Short = "skycoin hardware wallet utilities"
}
