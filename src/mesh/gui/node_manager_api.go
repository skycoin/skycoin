package gui

import (
	//"errors"
	"fmt"
	//"log"
	"net/http"
	//"os"
	//"strconv"
	//"strings"

	//"github.com/skycoin/skycoin/src/cipher"

	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers

	"github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func testHandler(gateway *nodemanager.NodeManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		wh.Error400(w, fmt.Sprintf("Works!"))

		if addr := r.FormValue("addr"); addr == "" {
			wh.Error404(w)
		} else {
			//wh.SendOr404(w, gateway.GetConnection(addr))
			wh.Error404(w)
		}
	}
}

func RegisterNodeManagerHandlers(mux *http.ServeMux, gateway *nodemanager.NodeManager) {
	// Returns wallet info
	// GET Arguments:
	//      id - Wallet ID.

	//  Gets a wallet .  Will be assigned name if present.
	mux.HandleFunc("/test", testHandler(gateway))

}
