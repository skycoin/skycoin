// Network-related information for the GUI
package gui

import (
    "github.com/skycoin/skycoin/src/daemon"
    "net/http"
)

func connectionHandler(rpc *daemon.RPC) HTTPHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        addr := r.FormValue("addr")
        if addr == "" {
            Error404(w)
            return
        }
        m := rpc.GetConnection(addr)
        if m == nil {
            Error404(w)
            return
        }
        if SendJSON(w, m) != nil {
            Error500(w)
        }
    }
}

func connectionsHandler(rpc *daemon.RPC) HTTPHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        m := rpc.GetConnections()
        if m == nil {
            Error404(w)
            return
        }
        if SendJSON(w, m) != nil {
            Error500(w)
        }
    }
}

func RegisterNetworkHandlers(mux *http.ServeMux, rpc *daemon.RPC) {
    mux.HandleFunc("/api/network/connection", connectionHandler(rpc))
    mux.HandleFunc("/api/network/connections", connectionsHandler(rpc))
}
