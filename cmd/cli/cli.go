/*
cli is a command line client for interacting with a skycoin node and offline wallet management
*/
package main

import (
	"fmt"
	"os"

	"github.com/skycoin/skycoin/src/cli"
)

func main() {
	cfg, err := cli.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cli, err := cli.NewCLI(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
