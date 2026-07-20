// Package commands operator UI for the skydex-market app.
package commands

import (
	"context"
	"embed"
	"net/http"
	"time"

	"github.com/skycoin/skycoin/src/skydex/db"
)

//go:embed static
var uiFS embed.FS

// serveUI runs the market operator UI (configuration + monitoring) and its
// operator API on addr until ctx is canceled. See skydex-design.md §2.
func serveUI(ctx context.Context, host Host, database *db.Database, addr string) {
	mux := http.NewServeMux()

	// The operator API is gated by a one-time code the market publishes to the
	// host (surfaced on the hypervisor app list over skywire). Register the API
	// on its own mux so a single wrapper covers every route — including any added
	// later, which is the point: a new handler is authenticated by default.
	gate, err := newAuthGate(host.PublishOTP)
	if err != nil {
		host.Log().Errorf("skydex-market: failed to initialize operator auth: %v", err)
		return
	}

	apiMux := http.NewServeMux()
	registerOperatorAPI(apiMux, database, host)

	// Exact patterns beat the "/api/" prefix in net/http's mux, so login stays
	// reachable without a token while everything else requires one.
	mux.Handle("/api/login", gate.loginHandler())
	mux.Handle("/api/", gate.requireAuth(apiMux))

	// Serve the embedded single-page operator UI. It is a single index.html, so
	// any non-/api path returns it.
	index, err := uiFS.ReadFile("static/index.html")
	if err != nil {
		host.Log().Errorf("skydex-market: failed to load operator UI: %v", err)
		return
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(index) //nolint:errcheck
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() { //nolint
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx) //nolint:errcheck
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		host.Log().Errorf("skydex-market: operator UI server stopped: %v", err)
	}
}
