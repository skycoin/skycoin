//go:build !386 && !(windows && arm64)
// +build !386
// +build !windows !arm64

// Package commands implements the skycoin hardware wallet daemon and CLI.
//
// Note: Hardware wallet support is disabled on 386 architecture due to
// limitations in the github.com/google/gousb library which lacks proper
// 32-bit support. The library's CGO bindings to libusb fail to compile
// on 386 with undefined type errors.
package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/hardware-wallet-daemon/api"
	"github.com/skycoin/skycoin/src/hardware-wallet-daemon/daemon"
	cli "github.com/skycoin/skycoin/src/hardware-wallet/cli"

	"github.com/skycoin/skywire/pkg/skywire-utilities/pkg/buildinfo"

	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	bv bool
	di bool

	appConfig = daemon.NewAppConfig(9510, "$HOME/.skycoin")
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Hardware wallet daemon for Skycoin",
	Long:  `A daemon that provides an HTTP API for interfacing with Skycoin hardware wallets`,
	Run: func(_ *cobra.Command, _ []string) {
		logger := logging.MustGetLogger("hw-daemon")

		d := daemon.NewDaemon(daemon.Config{
			App: appConfig,
			Build: api.BuildInfo{
				Version: buildinfo.Version(),
				Commit:  buildinfo.Commit(),
				Branch:  "", // buildinfo package doesn't track branch
			},
		}, logger)

		if err := d.ParseConfig(); err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		if err := d.Run(); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(
		daemonCmd,
		cli.RootCmd,
	)
	cli.RootCmd.Use = "cli"
	daemonCmd.Flags().IntVarP(&appConfig.WebInterfacePort, "web-interface-port", "p", appConfig.WebInterfacePort, "port to serve web interface on")
	daemonCmd.Flags().StringVarP(&appConfig.WebInterfaceAddr, "web-interface-addr", "a", appConfig.WebInterfaceAddr, "addr to serve web interface on")
	daemonCmd.Flags().BoolVar(&appConfig.EnableCSRF, "enable-csrf", appConfig.EnableCSRF, "enable CSRF check")
	daemonCmd.Flags().BoolVar(&appConfig.DisableHeaderCheck, "disable-header-check", appConfig.DisableHeaderCheck, "disables the host, origin and referer header checks")
	daemonCmd.Flags().StringVar(&appConfig.HostWhitelist, "host-whitelist", appConfig.HostWhitelist, "Hostnames to whitelist in the Host header check. Only applies when the web interface is bound to localhost.")

	daemonCmd.Flags().BoolVar(&appConfig.ColorLog, "color-log", appConfig.ColorLog, "Add terminal colors to log output")
	daemonCmd.Flags().StringVarP(&appConfig.LogLevel, "log-level", "l", appConfig.LogLevel, "Choices are: debug, info, warn, error, fatal, panic")
	daemonCmd.Flags().BoolVar(&appConfig.LogToFile, "logtofile", appConfig.LogToFile, "log to file")

	daemonCmd.Flags().BoolVar(&appConfig.ProfileCPU, "profile-cpu", appConfig.ProfileCPU, "enable cpu profiling")
	daemonCmd.Flags().StringVar(&appConfig.ProfileCPUFile, "profile-cpu-file", appConfig.ProfileCPUFile, "where to write the cpu profile file")
	daemonCmd.Flags().BoolVar(&appConfig.HTTPProf, "http-prof", appConfig.HTTPProf, "run the HTTP profiling interface")
	daemonCmd.Flags().StringVar(&appConfig.HTTPProfHost, "http-prof-host", appConfig.HTTPProfHost, "hostname to bind the HTTP profiling interface to")

	daemonCmd.Flags().StringVarP(&appConfig.DataDirectory, "data-dir", "d", appConfig.DataDirectory, "directory to store app data (defaults to ~/.skycoin)")

	daemonCmd.Flags().StringVarP(&appConfig.DaemonMode, "daemon-mode", "m", appConfig.DaemonMode, "Choices are: USB or EMULATOR")

	if fmt.Sprintf("%v", buildinfo.DebugBuildInfo()) != "" {
		RootCmd.Flags().BoolVarP(&di, "info", "d", false, "print runtime/debug.BuildInfo")
	}
	if fmt.Sprintf("%v", buildinfo.DBIVersion()) != "" {
		RootCmd.Flags().BoolVarP(&bv, "bv", "b", false, "print runtime/debug.BuildInfo.Main.Version")
	}
	daemonCmd.Version = buildinfo.Version()
}

// RootCmd contains hardware wallet daemon and CLI commands
var RootCmd = &cobra.Command{
	Use: func() string {
		return strings.Split(filepath.Base(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v", os.Args), "[", ""), "]", "")), " ")[0]
	}(),
	Long: func() (ret string) {
		ret = `
    ┌─┐┬┌─┬ ┬┌─┐┌─┐┬┌┐┌
    └─┐├┴┐└┬┘│  │ │││││
    └─┘┴ ┴ ┴ └─┘└─┘┴┘└┘
    skycoin hardware wallet utilities`
		if buildinfo.DBIVersion() != "" {
			ret += fmt.Sprintf("\n%v", buildinfo.DBIVersion())
		} else {
			ret += fmt.Sprintf("\nversion %v", buildinfo.Version())
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
