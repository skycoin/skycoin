// Network-related information for the GUI
package gui

import (
    "github.com/skycoin/skycoin/src/daemon"
    "net/http"
)

func connectionPage(w http.ResponseWriter, r *http.Request) {
    addr := r.FormValue("addr")
    if addr == "" {
        Error404(w)
        return
    }
    m := daemon.GetConnection(addr)
    if m == nil {
        Error404(w)
        return
    }
    if SendJSON(w, m) != nil {
        Error500(w)
    }
}

func connectionsPage(w http.ResponseWriter, r *http.Request) {
    m := daemon.GetConnections()
    if m == nil {
        Error404(w)
        return
    }
    if SendJSON(w, m) != nil {
        Error500(w)
    }
}

func RegisterNetworkHandlers(mux *http.ServeMux) {
    mux.HandleFunc("/api/network/connection", connectionPage)
    mux.HandleFunc("/api/network/connections", connectionsPage)
}
