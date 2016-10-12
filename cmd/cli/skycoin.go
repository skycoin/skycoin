package main

import (
	"os"

	skycli "github.com/skycoin/skycoin/src/api/cli"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Usage = "the skycoin command line interface"
	app.Version = "0.1"
	app.Commands = skycli.Commands
	app.Run(os.Args)
}
