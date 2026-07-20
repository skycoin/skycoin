package commands

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/skydex/market"
)

var (
	// uiAddr is the localhost address the trading UI is served on.
	uiAddr string
	// defaultMarket is the market address pre-filled in the UI's connect screen.
	defaultMarket string
)

func init() {
	RootCmd.Flags().StringVar(&uiAddr, "addr", "127.0.0.1:8051", "address to serve the trading UI on")
	RootCmd.Flags().StringVar(&defaultMarket, "market", "", "default market address pre-filled in the UI (standalone: host:port; not auto-connected)")
}

// RootCmd runs the standalone (TCP) SkyDEX trading client.
var RootCmd = &cobra.Command{
	Use:   "skydex-client",
	Short: "SkyDEX client — decentralized exchange trading UI (standalone)",
	Long: `SkyDEX client serves the trading UI and dials a market on demand.

This is the standalone build: it reaches the market over a plain TCP
connection. To run it over skywire — dialing the market by its public key
over dmsg — use the skywire 'skydex-client' app instead.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		return Run(ctx, Config{UIAddr: uiAddr, DefaultPK: defaultMarket}, tcpDialer{}, logrus.New())
	},
}

// Execute runs the root command.
func Execute() {
	if err := RootCmd.ExecuteContext(context.Background()); err != nil {
		log.Fatal(err)
	}
}

// tcpDialer is the standalone MarketDialer: ref is a TCP "host:port" address.
type tcpDialer struct{}

func (tcpDialer) Dial(ref string) (*market.Conn, error) {
	c, err := net.Dial("tcp", ref)
	if err != nil {
		return nil, err
	}
	return market.NewConn(c), nil
}
