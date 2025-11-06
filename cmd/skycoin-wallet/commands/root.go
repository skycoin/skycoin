// Package commands implements the skycoin wallet commands.
package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	explorer "github.com/skycoin/skycoin/cmd/explorer/commands"
	newcoin "github.com/skycoin/skycoin/cmd/newcoin/commands"
	cli "github.com/skycoin/skycoin/cmd/skycoin-cli/commands"
	web "github.com/skycoin/skycoin/cmd/skycoin-web/commands"
	skycoin "github.com/skycoin/skycoin/cmd/skycoin/commands"
)

func init() {

	RootCmd.AddCommand(
		skycoin.RootCmd,
		web.RootCmd,
		cli.RootCmd,
		newcoin.RootCmd,
		explorer.RootCmd,
	)
	skycoin.RootCmd.Use = "daemon"
	web.RootCmd.Use = "web"
	web.RootCmd.Short = "skycoin thin client web wallet"
	explorer.RootCmd.Use = "explorer"
}

// RootCmd contains every daemon, cli, & newcoin
var RootCmd = &cobra.Command{
	Use: func() string {
		return strings.Split(filepath.Base(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v", os.Args), "[", ""), "]", "")), " ")[0]
	}(),
	Long: `
	┌─┐┬┌─┬ ┬┌─┐┌─┐┬┌┐┌
	└─┐├┴┐└┬┘│  │ │││││
	└─┘┴ ┴ ┴ └─┘└─┘┴┘└┘`,
	SilenceUsage:          true,
	DisableSuggestions:    true,
	DisableFlagsInUseLine: true,
}

// Execute executes root CLI command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal("Failed to execute command: ", err)
	}
}
