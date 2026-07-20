// Package commands is the skydex-market app: the SkyDEX exchange backend.
//
// The exchange engine (SQLite store, matching/escrow, background jobs, the
// client<->market protocol server, and the operator UI) lives here in skycoin
// and carries no skywire dependency. Run is transport-agnostic: it is handed a
// net.Listener and a way to resolve each connection's authenticated identity,
// so the same engine serves over a plain TCP listener (the standalone reference
// command below) or over a skywire appnet/dmsg listener (the thin skywire app
// wrapper, which injects the client's authenticated public key as the identity).
package commands

import (
	"context"
	"net"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/skydex/chain"
	"github.com/skycoin/skycoin/src/skydex/db"
	"github.com/skycoin/skycoin/src/skydex/jobs"
	"github.com/skycoin/skycoin/src/skydex/server"
)

// Host is the runtime environment the market app runs in. Over skywire it is
// backed by the visor app client; the standalone command backs it with a local
// stdout/stderr implementation. Keeping the coupling to these three methods is
// what lets the whole engine live in skycoin with no skywire import.
type Host interface {
	// Log is the logger used by the server engine, background jobs and the
	// operator UI.
	Log() logrus.FieldLogger
	// PublishOTP surfaces a freshly minted operator one-time login code. Over
	// skywire it pushes the code to the visor app list (behind hypervisor auth);
	// standalone it prints it to the log so the operator can read it.
	PublishOTP(otp string)
	// PubKey is the market's own identity, shown in the operator UI: the visor
	// public key over skywire, a configured value (or "") standalone.
	PubKey() string
}

// Config parameterizes a market run.
type Config struct {
	// DBPath is the SQLite database path. Empty falls back to
	// <WorkDir>/skydex-market.db.
	DBPath string
	// WorkDir is the base directory used to derive the default DB path.
	WorkDir string
	// UIAddr is the localhost address the operator UI is served on.
	UIAddr string
	// OnReady, if set, is called once the market is serving. The skywire wrapper
	// uses it to report "Running" to the visor; standalone leaves it nil.
	OnReady func()
}

// Run starts the market engine and blocks until ctx is canceled.
//
// lis is the transport the client<->market protocol is served on and identify
// resolves each accepted connection's authenticated identity (the dmsg remote
// public key over skywire; a TLS/login identity or the remote address over a
// clearnet listener). host supplies logging, the operator OTP sink, and the
// market's own public key.
func Run(ctx context.Context, cfg Config, lis net.Listener, identify server.IdentifyFunc, host Host) error {
	log := host.Log()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	database, err := db.New(cfg.DBPath, cfg.WorkDir)
	if err != nil {
		return err
	}
	defer database.Close() //nolint:errcheck
	log.Infof("skydex-market: database initialized at %s", database.Path())

	if err := database.Migrate(); err != nil {
		return err
	}
	log.Info("skydex-market: database migrations completed")

	if err := database.InitDefaultConfig(); err != nil {
		return err
	}
	log.Info("skydex-market: market configuration initialized")

	// Start the client<->market protocol server on the supplied transport.
	srv := server.New(database, log, 0)
	go func() {
		if err := srv.Accept(ctx, lis, identify); err != nil {
			log.Errorf("skydex-market: server stopped: %v", err)
			cancel()
		}
	}()

	// Serve the operator UI (configuration + monitoring) on the localhost address.
	go serveUI(ctx, host, database, cfg.UIAddr)

	// Start the background jobs (escrow/listing checks, expiry, cleanup, bans).
	// The chain backend reads its per-sell-coin escrow config and per-currency
	// explorer config from the DB, so operator changes take effect without a
	// restart. Coins with no fullnode URL yet simply never confirm.
	chainBackend := chain.New(database)
	if coins, _ := database.AvailableSellCoins(); len(coins) > 0 { //nolint:errcheck
		log.Infof("skydex-market: chain backend ready — sell coins: %s (external payments via configured explorers)", strings.Join(coins, ", "))
	} else {
		log.Info("skydex-market: chain backend ready — no sell coins enabled yet (add one in the operator UI)")
	}
	go jobs.NewRunner(database, chainBackend, log).Run(ctx)

	log.Infof("skydex-market: ready and waiting for client connections; operator UI at http://%s", cfg.UIAddr)
	if cfg.OnReady != nil {
		cfg.OnReady()
	}

	<-ctx.Done()
	log.Info("skydex-market: shutting down")

	return nil
}
