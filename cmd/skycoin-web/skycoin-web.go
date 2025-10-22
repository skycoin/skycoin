// cmd/skycoin-web/skycoin-web.go
/*
skycoin-web thin client
*/
package main

import (
	"github.com/skycoin/skycoin/cmd/skycoin-web/commands"
	"github.com/skycoin/skywire/pkg/skywire-utilities/pkg/flags"
)

func init() {
	flags.InitFlags(commands.RootCmd, false)
}


func main() {
	commands.Execute()
}
