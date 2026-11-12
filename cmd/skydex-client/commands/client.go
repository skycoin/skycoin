// Package commands cmd/skydex-client/commands/client.go c4-app-dex
//
// Package commands is the skydex-client app: the SkyDEX trading UI.
//
// It serves the embedded single-page trading client and a small local control
// API, and dials a market on demand through a MarketDialer. It carries no
// skywire dependency: over skywire the thin skywire app supplies an appnet
// dialer; the standalone command below supplies a plain TCP dialer.
package commands

import (
	"context"
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/0magnet/logrus"

	"github.com/skycoin/skycoin/cmd/skydex-client/skydexui"
)

//go:embed static
var uiFS embed.FS

// Config parameterizes a client run.
type Config struct {
	// UIAddr is the localhost address the trading UI is served on.
	UIAddr string
	// DefaultPK is the market reference pre-filled in the UI's connect screen.
	// The client never connects automatically; the user confirms and clicks
	// Connect.
	DefaultPK string
}

// Run serves the trading UI and its local control API on cfg.UIAddr, dialing
// the market on demand through dialer, until ctx is canceled.
func Run(ctx context.Context, cfg Config, dialer MarketDialer, log logrus.FieldLogger) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sess := newSession(dialer, cfg.DefaultPK)
	defer sess.close()

	log.Infof("skydex-client: trading UI available at http://%s", cfg.UIAddr)
	serveUI(ctx, log, sess, cfg.UIAddr) // blocks until ctx is canceled

	return nil
}

// serveUI serves the embedded SPA plus the local control API on addr until ctx
// is canceled. Any path that isn't a real asset falls back to index.html so
// client-side routes resolve.
func serveUI(ctx context.Context, log logrus.FieldLogger, sess *session, addr string) {
	handler, err := uiHandler(sess)
	if err != nil {
		log.Errorf("skydex-client: failed to open embedded UI: %v", err)
		return
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() { //nolint
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx) //nolint:errcheck
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Errorf("skydex-client: UI server stopped: %v", err)
	}
}

// reactPrefix is where the React client is served beside the server
// rendered one. Its build expects this base (vite.config.js).
const reactPrefix = "/angular"

// uiHandler is the server rendered client, the React client at
// reactPrefix and the local control API. Writes from another site are
// refused on all of them, so no other page can place or cancel orders
// through the API.
func uiHandler(sess *session) (http.Handler, error) {
	distFS, err := fs.Sub(uiFS, "static")
	if err != nil {
		return nil, err
	}
	fileServer := http.StripPrefix(reactPrefix, http.FileServer(http.FS(distFS)))

	mux := http.NewServeMux()
	registerAPI(mux, sess)
	mux.Handle("/", skydexui.Handler(uiMarket{sess}))
	mux.HandleFunc(reactPrefix+"/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(path.Clean(strings.TrimPrefix(r.URL.Path, reactPrefix)), "/")
		if st, statErr := fs.Stat(distFS, name); statErr != nil || st.IsDir() {
			r.URL.Path = reactPrefix + "/"
		}
		fileServer.ServeHTTP(w, r)
	})
	// /ui2/ held the server rendered client while it sat beside the React one.
	mux.HandleFunc("/ui2/", func(w http.ResponseWriter, r *http.Request) {
		// Clean, so /ui2//host cannot become a redirect to another site.
		to := path.Clean("/" + strings.TrimPrefix(r.URL.Path, "/ui2"))
		if r.URL.RawQuery != "" {
			to += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, to, http.StatusMovedPermanently) //nolint:gosec // to is a cleaned path on this host
	})

	return skydexui.SameOriginWrites(mux), nil
}
