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

	// register the supported wallets
	_ "github.com/SkycoinProject/skycoin/src/wallet/bip44wallet"
	_ "github.com/SkycoinProject/skycoin/src/wallet/collection"
	_ "github.com/SkycoinProject/skycoin/src/wallet/deterministic"
	_ "github.com/SkycoinProject/skycoin/src/wallet/xpubwallet"
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
