// Package commands provides CLI commands for interacting with skycoin.
package commands

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/cli"
	"github.com/skycoin/skycoin/src/util/logging"

	"github.com/skycoin/skywire/pkg/skywire-utilities/pkg/calvin"
	"github.com/spf13/cobra"

	// register the supported wallets
	_ "github.com/skycoin/skycoin/src/wallet/bip44wallet"
	_ "github.com/skycoin/skycoin/src/wallet/collection"
	_ "github.com/skycoin/skycoin/src/wallet/deterministic"
	_ "github.com/skycoin/skycoin/src/wallet/xpubwallet"
)

func init() {
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
	RootCmd = skyCLI
	RootCmd.Use = "cli"
	RootCmd.Short = description
	RootCmd.Long = calvin.AsciiFont("skycoin-cli") + "\n" + description

}

var description = "skycoin command line interface"

// RootCmd represents the base command for the application
var RootCmd = &cobra.Command{
	Use:   "cli",
	Short: description,
	Long:  calvin.AsciiFont("skycoin-cli") + "\n" + description,
}
