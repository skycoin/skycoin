package api

// APIs for network-related information

import (
	"net/http"
	"sort"

	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/readable"
	wh "github.com/skycoin/skycoin/src/util/http"
)

// connectionHandler returns a specific connection
// URI: /api/v1/network/connections
// Method: GET
// Args:
//	addr - An IP:Port string
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

		c, err := gateway.GetConnection(addr)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		if c == nil {
			wh.Error404(w, "")
			return
		}

		wh.SendJSONOr500(logger, w, readable.NewConnection(c))
	}
}

// Connections wraps []Connection
type Connections struct {
	Connections []readable.Connection `json:"connections"`
}

// NewConnections copies []daemon.Connection to a struct with json tags
func NewConnections(dconns []daemon.Connection) Connections {
	conns := make([]readable.Connection, len(dconns))
	for i, dc := range dconns {
		conns[i] = readable.NewConnection(&dc)
	}

	return Connections{
		Connections: conns,
	}
}

// connectionsHandler returns all outgoing connections
// URI: /api/v1/network/connections
// Method: GET
// Args: type [optional] either outgoing or incoming
func connectionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		var conns []daemon.Connection
		var err error

		typ := r.FormValue("type")
		switch typ {
		case "":
			conns, err = gateway.GetConnections()
		case "outgoing":
			conns, err = gateway.GetOutgoingConnections()
		case "incoming":
			conns, err = gateway.GetIncomingConnections()
		default:
			wh.Error400(w, "invalid type")
			return
		}

		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, NewConnections(conns))
	}
}

// defaultConnectionsHandler returns the list of default hardcoded bootstrap addresses.
// They are not necessarily connected to.
// URI: /api/v1/network/defaultConnections
// Method: GET
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

// trustConnectionsHandler returns all trusted connections
// In the default configuration, these will be a subset of the default hardcoded bootstrap addresses
// URI: /api/v1/network/trust
// Method: GET
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

// exchgConnectionsHandler returns all connections found through peer exchange
// URI: /api/v1/network/exchange
// Method: GET
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
