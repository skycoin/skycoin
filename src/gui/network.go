// Network-related information for the GUI
package gui

import (
    "github.com/skycoin/skycoin/src/daemon"
    "net/http"
)

func connectionHandler(gateway *daemon.Gateway) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if addr := r.FormValue("addr"); addr == "" {
            Error404(w)
        } else {
            SendOr404(w, gateway.GetConnection(addr))
        }
    }
}

func connectionsHandler(gateway *daemon.Gateway) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        SendOr404(w, gateway.GetConnections())
    }
}

func RegisterNetworkHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
    mux.HandleFunc("/api/network/connection", connectionHandler(gateway))
    mux.HandleFunc("/api/network/connections", connectionsHandler(gateway))
}
