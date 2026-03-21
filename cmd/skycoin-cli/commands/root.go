// Package commands provides CLI commands for interacting with skycoin.
package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/cli"
	"github.com/skycoin/skycoin/src/fiber"
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

	// Determine coin name from FIBER_TOML or default
	coinName := "skycoin"
	rpcDefault := "http://127.0.0.1:6420"
	if fiberTomlPath := os.Getenv("FIBER_TOML"); fiberTomlPath != "" {
		if absPath, err := filepath.Abs(fiberTomlPath); err == nil {
			fiberTomlPath = absPath
		}
		if fiberCfg, err := fiber.NewConfig(filepath.Base(fiberTomlPath), filepath.Dir(fiberTomlPath)); err == nil {
			if fiberCfg.Node.DisplayName != "" {
				coinName = fiberCfg.Node.DisplayName
			}
			if fiberCfg.Node.WebInterfacePort != 0 {
				rpcDefault = fmt.Sprintf("http://127.0.0.1:%d", fiberCfg.Node.WebInterfacePort)
			}
		}
	}
	coinNameLower := strings.ToLower(coinName)

	description := fmt.Sprintf("%s command line interface", coinNameLower)

	// Configure the RootCmd with the CLI subcommands
	RootCmd.Use = "cli"
	RootCmd.Short = description
	RootCmd.Long = calvin.AsciiFont(coinNameLower+"-cli") + "\n" + description + "\n" + fmt.Sprintf(`
ENVIRONMENT VARIABLES:
  RPC_ADDR: Address of RPC node. Must be in scheme://host format. Default "%s"
  RPC_USER: Username for RPC API, if enabled in the RPC.
  RPC_PASS: Password for RPC API, if enabled in the RPC.
  COIN: Name of the coin. Default "%s"
  DATA_DIR: Directory where everything is stored. Default "$HOME/.$COIN/"
`, rpcDefault, coinNameLower)

	// Add all CLI subcommands to RootCmd
	for _, cmd := range cliCmd.Commands() {
		RootCmd.AddCommand(cmd)
	}

	// Use flags.InitFlags for consistent help formatting
	flags.InitFlags(RootCmd, true)
}

// RootCmd represents the base command for the application
var RootCmd = &cobra.Command{
	Use:   "cli",
	Short: "command line interface",
}
