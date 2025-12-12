// Package commands provides CLI commands for interacting with skycoin.
package commands

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/cli"
	"github.com/skycoin/skycoin/src/util/logging"

	"github.com/skycoin/skywire/pkg/skywire-utilities/pkg/calvin"
	"github.com/skycoin/skywire/pkg/skywire-utilities/pkg/flags"
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
	cliCmd, err := cli.NewCLI(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Configure the RootCmd with the CLI subcommands
	RootCmd.Use = "cli"
	RootCmd.Short = description
	RootCmd.Long = calvin.AsciiFont("skycoin-cli") + "\n" + description + "\n" + `
ENVIRONMENT VARIABLES:
  RPC_ADDR: Address of RPC node. Must be in scheme://host format. Default "http://127.0.0.1:6420"
  RPC_USER: Username for RPC API, if enabled in the RPC.
  RPC_PASS: Password for RPC API, if enabled in the RPC.
  COIN: Name of the coin. Default "skycoin"
  DATA_DIR: Directory where everything is stored. Default "$HOME/.$COIN/"
`

	// Add all CLI subcommands to RootCmd
	for _, cmd := range cliCmd.Commands() {
		RootCmd.AddCommand(cmd)
	}

	// Use flags.InitFlags for consistent help formatting
	flags.InitFlags(RootCmd, true)
}

var description = "skycoin command line interface"

// RootCmd represents the base command for the application
var RootCmd = &cobra.Command{
	Use:   "cli",
	Short: description,
	Long: calvin.AsciiFont("skycoin") + " cli\n" + description + "\n" + `
ENVIRONMENT VARIABLES:
  RPC_ADDR: Address of RPC node. Must be in scheme://host format. Default "http://127.0.0.1:6420"
  RPC_USER: Username for RPC API, if enabled in the RPC.
  RPC_PASS: Password for RPC API, if enabled in the RPC.
  COIN: Name of the coin. Default "skycoin"
  DATA_DIR: Directory where everything is stored. Default "$HOME/.$COIN/"
`,
}
