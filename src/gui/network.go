// Network-related information for the GUI
package gui

import (
    "github.com/skycoin/skycoin/src/daemon"
    "net/http"
)

type Connection struct {
    Addr         string `json:"address"`
    LastSent     int64  `json:"last_sent"`
    LastReceived int64  `json:"last_received"`
}

type Connections struct {
    Connections []Connection `json:"connections"`
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
            Addr:         v.Addr(),
            LastSent:     v.LastSent.Unix(),
            LastReceived: v.LastReceived.Unix(),
        })
    }

    if SendJSON(w, &Connections{conns}) != nil {
        Error500(w)
    }
}

func RegisterNetworkHandlers(mux *http.ServeMux) {
    mux.HandleFunc("/api/network/connections", connectionsPage)
}
