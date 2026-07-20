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
)

var (
	// dbPath is the path to the SQLite database file.
	dbPath string
	// uiAddr is the localhost address the operator UI is served on.
	uiAddr string
	// listenAddr is the TCP address the client<->market protocol is served on
	// in standalone mode.
	listenAddr string
	// marketPK is the market identity shown in the operator UI. Standalone it is
	// operator-supplied (or empty); over skywire the wrapper injects the visor PK.
	marketPK string
)

func init() {
	RootCmd.Flags().StringVar(&dbPath, "db", "", "path to SQLite database file (default: ./skydex-market.db)")
	RootCmd.Flags().StringVar(&uiAddr, "addr", "127.0.0.1:8050", "address to serve the operator UI on")
	RootCmd.Flags().StringVar(&listenAddr, "listen", "127.0.0.1:8060", "TCP address to serve the client<->market protocol on")
	RootCmd.Flags().StringVar(&marketPK, "pubkey", "", "market identity shown in the operator UI (optional; set to your visor PK when fronted by skywire)")
}

// RootCmd runs the standalone (TCP) SkyDEX market.
var RootCmd = &cobra.Command{
	Use:   "skydex-market",
	Short: "SkyDEX market — decentralized exchange backend (standalone reference server)",
	Long: `SkyDEX market runs the exchange backend: the SQLite store, the
matching/escrow engine, background jobs, the operator UI, and the
client<->market protocol served over a plain TCP listener.

This is the standalone reference server. To run it over skywire — an
authenticated, no-DNS, public-key-addressed transport that injects each
client's authenticated public key as its identity — use the skywire
'skydex-market' app instead.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		lis, err := net.Listen("tcp", listenAddr)
		if err != nil {
			return err
		}

		host := stdHost{log: logrus.New(), pubKey: marketPK}

		// Standalone identity is the client's TCP remote address. Over skywire the
		// app wrapper injects the authenticated dmsg public key instead; a future
		// standalone auth (TLS client cert / login) would replace this.
		identify := func(c net.Conn) (string, error) { return c.RemoteAddr().String(), nil }

		return Run(ctx, Config{DBPath: dbPath, UIAddr: uiAddr}, lis, identify, host)
	},
}

// Execute runs the root command.
func Execute() {
	if err := RootCmd.ExecuteContext(context.Background()); err != nil {
		log.Fatal(err)
	}
}

// stdHost is the standalone Host: it logs to stdout, prints the operator OTP to
// the log, and reports the operator-supplied market identity.
type stdHost struct {
	log    logrus.FieldLogger
	pubKey string
}

func (h stdHost) Log() logrus.FieldLogger { return h.log }
func (h stdHost) PublishOTP(otp string) {
	h.log.Infof("skydex-market: operator login code (type it into the operator UI): %s", otp)
}
func (h stdHost) PubKey() string { return h.pubKey }
