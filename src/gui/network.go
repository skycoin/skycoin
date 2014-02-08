// Network-related information for the GUI
package gui

import (
    "github.com/skycoin/skycoin/src/daemon"
    "net/http"
)

func connectionHandler(rpc *daemon.RPC) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if addr := r.FormValue("addr"); addr == "" {
            Error404(w)
        } else {
            SendOr404(w, rpc.GetConnection(addr))
        }
    }
}

func connectionsHandler(rpc *daemon.RPC) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        SendOr404(w, rpc.GetConnections())
    }
}

func RegisterNetworkHandlers(mux *http.ServeMux, rpc *daemon.RPC) {
    mux.HandleFunc("/api/network/connection", connectionHandler(rpc))
    mux.HandleFunc("/api/network/connections", connectionsHandler(rpc))
}
