package skydexui

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/skycoin/skycoin/src/skydex/protocol"
)

// terminal are the order states of a finished trade.
var terminal = map[string]bool{"completed": true, "canceled": true, "cancelled": true, "expired": true} //nolint:misspell // both spellings occur

// Trade is an entry of the browser's trade history, kept in localStorage
// under exchange:history as History.jsx keeps it. Fields this client does
// not know are kept as they are.
type Trade map[string]any

func (t Trade) str(k string) string {
	s, _ := t[k].(string) //nolint:errcheck
	return s
}

func (t Trade) num(k string) float64 {
	f, _ := t[k].(float64) //nolint:errcheck
	return f
}

// mergeTrades folds the finished orders into saved, as mergeTradeHistory in
// src/api.js: new orders are added, known ones refreshed, and orders the
// market has since purged are kept. The result is newest first.
func mergeTrades(saved []Trade, orders []protocol.OrderView, marketPK, marketName string, now time.Time) []Trade {
	byID := map[string]Trade{}
	var ids []string
	for _, t := range saved {
		id := t.str("id")
		if id == "" {
			continue
		}
		if _, ok := byID[id]; !ok {
			ids = append(ids, id)
		}
		byID[id] = t
	}
	for _, o := range orders {
		if o.ID == "" || !terminal[o.Status] {
			continue
		}
		prev, ok := byID[o.ID]
		if !ok {
			prev = Trade{}
			ids = append(ids, o.ID)
		}
		next := Trade{}
		for k, v := range prev {
			next[k] = v
		}
		b, _ := json.Marshal(o)      //nolint:errcheck
		_ = json.Unmarshal(b, &next) //nolint:errcheck
		next["market_pk"] = firstNonEmpty(marketPK, prev.str("market_pk"))
		next["market_name"] = firstNonEmpty(marketName, prev.str("market_name"))
		if _, ok := prev["saved_at"]; !ok {
			next["saved_at"] = now.UnixMilli()
		}
		byID[o.ID] = next
	}
	out := make([]Trade, 0, len(ids))
	for _, id := range ids {
		out = append(out, byID[id])
	}
	sort.SliceStable(out, func(i, j int) bool {
		return created(out[i]).After(created(out[j]))
	})
	return out
}

func created(t Trade) time.Time {
	c, _ := time.Parse(time.RFC3339Nano, t.str("created_at")) //nolint:errcheck
	return c
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

// history renders the browser's trade history merged with the market's
// finished orders. skydex.js stores the merged list back.
func (h *handler) history(w http.ResponseWriter, r *http.Request) {
	var saved []Trade
	_ = json.Unmarshal([]byte(r.FormValue("saved")), &saved) //nolint:errcheck
	connected, ref, name := h.m.Status()
	var orders protocol.GetOrdersResponse
	if connected && h.m.Call(protocol.TypeGetOrders, nil, &orders) != nil {
		// Market unreachable: still show what the browser kept.
		orders.Orders = nil
	}
	render(w, r, HistoryBody(mergeTrades(saved, orders.Orders, ref, name, time.Now())))
}
