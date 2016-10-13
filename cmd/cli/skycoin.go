package main

import (
	"fmt"
	"os"

	skycli "github.com/skycoin/skycoin/src/api/cli"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.UsageText = "the skycoin command line interface"
	app.Version = "0.1"
	app.Commands = skycli.Commands
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
