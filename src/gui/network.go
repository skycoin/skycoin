package gui

// Network-related information for the GUI
import (
	"net/http"
	"sort"

	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)

// ConnectionStatus structs
type ConnectionStatus struct {
	Connection string `json:"Connection"`
	IsAlive    bool   `json:"Status"`
}

// ConnectionsHealth struct
type ConnectionsHealth struct {
	Count        int                `json:"count"`
	TotalAlive   int                `json:"total_alive"`
	TotalOffline int                `json:"total_offline"`
	Connections  []ConnectionStatus `json:IsAlive`
}

func defaultStatus(gateway Gatewayer) ConnectionsHealth {

	connsDefault := gateway.GetDefaultConnections()
	sort.Strings(connsDefault)
	connsAll := gateway.GetConnections().Connections

	countDefault, totalAlive := len(connsDefault), 0
	totalOffline := countDefault

	var connections []ConnectionStatus
	connsMap := make(map[string]*ConnectionStatus, countDefault)
	for _, conn := range connsDefault {

		status := ConnectionStatus{
			Connection: conn,
			IsAlive:    false,
		}
		connections = append(connections, status)
		connsMap[conn] = &status
	}

	for _, conn := range connsAll {
		if status, isDefault := connsMap[conn.Addr]; isDefault {
			if !status.IsAlive {
				status.IsAlive = true
				totalAlive++
				totalOffline--
			}
		}
	}

	resp := ConnectionsHealth{
		Count:        countDefault,
		TotalAlive:   totalAlive,
		TotalOffline: totalOffline,
		Connections:  connections,
	}

	return resp
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
			wh.Error404(w)
			return
		}

		wh.SendJSONOr500(logger, w, c)
	}
}

func connectionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		wh.SendJSONOr500(logger, w, gateway.GetConnections())
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

func defaultStatusHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}
		resp := defaultStatus(gateway)

		wh.SendJSONOr500(logger, w, &resp)

	}
}
