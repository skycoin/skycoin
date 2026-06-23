//go:build tinygo

// Package commands provides commands for the skycoin web interface.
//
// The full web interface is built on github.com/gin-gonic/gin, which
// unconditionally imports github.com/quic-go/quic-go/http3 for its RunQUIC
// helper. quic-go relies on QUIC TLS APIs that TinyGo's crypto/tls does not
// implement, so the web interface cannot be compiled under TinyGo. This stub
// provides a RootCmd that reports the limitation when invoked.
package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is the root web command. Under TinyGo it is a stub that reports that
// the web interface is unavailable in this build.
var RootCmd = &cobra.Command{
	Use:   "web",
	Short: "thin client web wallet (unavailable in TinyGo builds)",
	RunE: func(_ *cobra.Command, _ []string) error {
		return fmt.Errorf("the web interface is not available in TinyGo builds")
	},
}

// Execute runs the root web command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
