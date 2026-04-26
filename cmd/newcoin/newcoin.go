// package main cmd/newcoin/newcoin.go
/*
newcoin generates a new coin cmd from a toml configuration file
*/

// Package main implements the newcoin template generator.
package main

import (
	"log"

	"github.com/skycoin/skycoin/src/util/flags"

	"github.com/skycoin/skycoin/cmd/newcoin/commands"
)

func init() {
	flags.InitFlags(commands.RootCmd, true)
}

func main() {
	err := commands.RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
