// cmd/skycoin-web/skycoin-web.go
/*
skycoin-web thin client
*/
package main

import (
	"github.com/skycoin/skycoin/src/util/flags"

	"github.com/skycoin/skycoin/cmd/skycoin-web/commands"
	// Serves the skycoin-lite cipher wasm at /assets/scripts/. The command
	// itself no longer embeds it, so a binary that wants it says so.
	_ "github.com/skycoin/skycoin/cmd/skycoin-web/wasmassets"
)

func init() {
	flags.InitFlags(commands.RootCmd, false)
}

func main() {
	commands.Execute()
}
