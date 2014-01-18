// Network-related information for the GUI
package gui

import (
    "net/http"
)

type PeersMessage struct {
    Count int
    Peers []string
}

func peersPage(w http.ResponseWriter, r *http.Request) {
    m := PeersMessage{}
    SendJSON(w, m)
}

func RegisterNetworkHandlers(mux *http.ServeMux) {
    mux.HandleFunc("/peers", peersPage)
}
