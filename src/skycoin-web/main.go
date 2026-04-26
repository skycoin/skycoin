// cmd/skycoin-web/skycoin-web.go
/*
skycoin-web thin client
*/
package main

import (
	"github.com/skycoin/skycoin/src/util/flags"

	"github.com/skycoin/skycoin/cmd/skycoin-web/commands"
)

func init() {
	flags.InitFlags(commands.RootCmd, false)
}

func main() {
	commands.Execute()
}
