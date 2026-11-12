// Package main cmd/skydex-market/skydex-market.go c4-app-dex
//
/*
skydex-market — the SkyDEX decentralized-exchange backend (standalone server).
*/
package main

import (
	"github.com/skycoin/skycoin/cmd/skydex-market/commands"
)

func main() {
	commands.Execute()
}
