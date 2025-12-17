// Package commands implements the skycoin wallet commands.
package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/skycoin/skywire/pkg/skywire-utilities/pkg/buildinfo"
	"github.com/spf13/cobra"

	explorer "github.com/skycoin/skycoin/cmd/explorer/commands"
	newcoin "github.com/skycoin/skycoin/cmd/newcoin/commands"
	cli "github.com/skycoin/skycoin/cmd/skycoin-cli/commands"
	web "github.com/skycoin/skycoin/cmd/skycoin-web/commands"
	skycoin "github.com/skycoin/skycoin/cmd/skycoin/commands"
)

var (
	bv bool
	di bool
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
	if fmt.Sprintf("%v", buildinfo.DebugBuildInfo()) != "" {
		RootCmd.Flags().BoolVarP(&di, "info", "d", false, "print runtime/debug.BuildInfo")
	}
	if fmt.Sprintf("%v", buildinfo.DBIVersion()) != "" {
		RootCmd.Flags().BoolVarP(&bv, "bv", "b", false, "print runtime/debug.BuildInfo.Main.Version")
	}
}

// RootCmd contains every daemon, cli, & newcoin
var RootCmd = &cobra.Command{
	Use: func() string {
		return strings.Split(filepath.Base(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v", os.Args), "[", ""), "]", "")), " ")[0]
	}(),
	Long: func() (ret string) {
		ret = `
    ┌─┐┬┌─┬ ┬┌─┐┌─┐┬┌┐┌
    └─┐├┴┐└┬┘│  │ │││││
    └─┘┴ ┴ ┴ └─┘└─┘┴┘└┘`
		if buildinfo.DBIVersion() != "" {
			ret += fmt.Sprintf("\n%v", buildinfo.DBIVersion())
		} else {
			ret += fmt.Sprintf("\nskycoin version %v", buildinfo.Version())
		}
		if buildinfo.Go() != "unknown" && buildinfo.Go() != "" {
			ret += "\nbuilt with " + buildinfo.Go()
		}
		return ret
	}(),
	SilenceErrors:         true,
	SilenceUsage:          true,
	DisableSuggestions:    true,
	DisableFlagsInUseLine: true,
	Version:               buildinfo.Version(),
	Run: func(cmd *cobra.Command, _ []string) {
		if di {
			fmt.Printf("%v\n", buildinfo.DebugBuildInfo())
			return
		}
		if bv {
			fmt.Printf("%v\n", buildinfo.DBIVersion())
			return
		}
		if err := cmd.Help(); err != nil {
			log.Printf("Failed to print help: %v", err)
		}
	},
}

// Execute executes root CLI command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal("Failed to execute command: ", err)
	}
}
