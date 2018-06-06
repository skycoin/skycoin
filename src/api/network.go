package api

// APIs for network-related information

import (
	"net/http"
	"sort"

	daemon "github.com/skycoin/skycoin/src/daemon" //http,json helpers
	wh "github.com/skycoin/skycoin/src/util/http"  //http,json helpers
)

// Connection wrapper around daemon connection with info about block height added
type Connection struct {
	*daemon.Connection
	Height uint64 `json:"height"`
}

// Connections an array of connections
// Arrays must be wrapped in structs to avoid certain javascript exploits
type Connections struct {
	Connections []Connection `json:"connections"`
}

func connectionHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		addr := r.FormValue("addr")
		if addr == "" {
			wh.Error400(w, "addr is required")
			return
		}

		c := gateway.GetConnection(addr)
		if c == nil {
			wh.Error404(w, "")
			return
		}
		cnx := Connection{
			Connection: c,
			Height:     0,
		}
		bcp, err := gateway.GetBlockchainProgress()
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}
		for _, ph := range bcp.Peers {
			if ph.Address == c.Addr {
				cnx.Height = ph.Height
				break
			}
		}

		wh.SendJSONOr500(logger, w, cnx)
	}
}

func connectionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		dcnxs := gateway.GetConnections()
		bcp, err := gateway.GetBlockchainProgress()
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		peerHeights := bcp.Peers
		index := make(map[string]uint64, len(peerHeights))

		for i := 0; i < len(peerHeights); i++ {
			index[peerHeights[i].Address] = peerHeights[i].Height
		}

		cnxs := Connections{}
		for _, c := range dcnxs.Connections {
			cnx := Connection{
				Connection: c,
				Height:     index[c.Addr],
			}
			cnxs.Connections = append(cnxs.Connections, cnx)
		}
		wh.SendJSONOr500(logger, w, cnxs)
	}
}

func defaultConnectionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		conns := gateway.GetDefaultConnections()
		sort.Strings(conns)

		wh.SendJSONOr500(logger, w, conns)
	}
}

func trustConnectionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		conns := gateway.GetTrustConnections()
		sort.Strings(conns)

		wh.SendJSONOr500(logger, w, conns)
	}
}

func exchgConnectionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		conns := gateway.GetExchgConnection()
		sort.Strings(conns)

		wh.SendJSONOr500(logger, w, conns)
	}
}
