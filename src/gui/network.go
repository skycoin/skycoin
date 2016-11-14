// Network-related information for the GUI
package gui

import (
	"net/http"

	"github.com/skycoin/skycoin/src/daemon"
	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)

func connectionHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if addr := r.FormValue("addr"); addr == "" {
			wh.Error404(w)
		} else {
			wh.SendOr404(w, gateway.GetConnection(addr))
		}
	}
}

func connectionsHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wh.SendOr404(w, gateway.GetConnections())
	}
}

func defaultConnectionsHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wh.SendOr404(w, gateway.GetDefaultConnections())
	}
}

func RegisterNetworkHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
	mux.HandleFunc("/network/connection", connectionHandler(gateway))
	mux.HandleFunc("/network/connections", connectionsHandler(gateway))
	mux.HandleFunc("/network/defaultConnections", defaultConnectionsHandler(gateway))
}
