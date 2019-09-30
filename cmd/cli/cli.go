/*
cli is a command line client for interacting with a skycoin node and offline wallet management
*/
package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/SkycoinProject/skycoin/src/cli"
	"github.com/SkycoinProject/skycoin/src/util/logging"
)

func main() {
	logging.SetLevel(logrus.WarnLevel)

	cfg, err := cli.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	skyCLI, err := cli.NewCLI(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := skyCLI.Execute(); err != nil {
		os.Exit(1)
	}
}
