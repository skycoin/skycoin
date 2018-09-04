package api

// APIs for network-related information

import (
	"net/http"
	"sort"

	"github.com/skycoin/skycoin/src/readable"
	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)

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

func connectionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		dcnxs, err := gateway.GetConnections()
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, readable.NewConnections(dcnxs))
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
