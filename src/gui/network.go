package gui

// Network-related information for the GUI
import (
	"net/http"

	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)

func connectionHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if addr := r.FormValue("addr"); addr == "" {
			wh.Error404(w)
		} else {
			wh.SendOr404(w, gateway.GetConnection(addr))
		}
	}
}

func connectionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wh.SendOr404(w, gateway.GetConnections())
	}
}

func defaultConnectionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wh.SendOr404(w, gateway.GetDefaultConnections())
	}
}

func trustConnectionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wh.SendOr404(w, gateway.GetTrustConnections())
	}
}

func exchgConnectionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wh.SendOr404(w, gateway.GetExchgConnection())
	}
}
