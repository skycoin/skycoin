// Package commands cmd/skycoin-web/commands/ui2.go c4-ui-wallet
package commands

import (
	"net/http"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/ui2"
)

// mountUI2 serves the server rendered wallet at /ui2/, beside the Angular
// one. Its API client calls mux's /api routes in process, so it reaches the
// node through the same proxy, node choice and dialer as the Angular wallet.
func mountUI2(mux *http.ServeMux) {
	client := api.HandlerClient(mux, "localhost")
	mux.Handle("/ui2/", http.StripPrefix("/ui2", ui2.Handler(client, ui2.Options{Coins: ui2.NewCoinHost(mux, "localhost")})))
}
