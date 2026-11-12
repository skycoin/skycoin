// Package commands cmd/skycoin-web/commands/ui2.go c4-ui-wallet
package commands

import (
	"io/fs"
	"net/http"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/skycoin-web/src/gui"
	"github.com/skycoin/skycoin/src/ui2"
)

// mountUIs serves the server rendered wallet at the root and the Angular
// one, built into guiFS, at /angular/. The pages' API client calls mux's
// /api routes in process, so it reaches the node through the same proxy,
// node choice and dialer as the Angular wallet.
func mountUIs(mux *http.ServeMux, guiFS fs.FS) {
	client := api.HandlerClient(mux, "localhost")
	mux.Handle("/", ui2.Handler(client, ui2.Options{Coins: ui2.NewCoinHost(mux, "localhost")}))
	mux.Handle("/ui2/", api.MovedPrefix("/ui2", ""))
	api.MountAngular(mux, gui.HandlerFS(guiFS), true)
}
