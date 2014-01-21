// Network-related information for the GUI
package gui

import (
    "github.com/skycoin/skycoin/src/daemon"
    "net/http"
)

type Connection struct {
    addr          string
    last_sent     int64
    last_received int64
}

type Connections struct {
    connections []Connection
}

func connectionsPage(w http.ResponseWriter, r *http.Request) {
    if daemon.Pool == nil {
        // the daemon is not initialized
        Error404(w)
        return
    }

    conns := make([]Connection, len(daemon.Pool.Pool))
    for _, v := range daemon.Pool.Pool {
        conns = append(conns, Connection{
            addr:          v.Addr(),
            last_sent:     v.LastSent.Unix(),
            last_received: v.LastReceived.Unix(),
        })
    }

    if SendJSON(w, &Connections{conns}) != nil {
        Error500(w)
    }
}

func RegisterNetworkHandlers(mux *http.ServeMux) {
    mux.HandleFunc("/api/network/connections", connectionsPage)
}
