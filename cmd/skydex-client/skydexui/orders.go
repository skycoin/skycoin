// Package skydexui cmd/skydex-client/skydexui/orders.go c3-ui-dex
package skydexui

import (
	"net/http"

	"github.com/a-h/templ"

	"github.com/skycoin/skycoin/src/skydex/protocol"
)

// defaultConfirmations is the threshold shown when a market does not report
// its own, as Orders.jsx.
const defaultConfirmations = 2

// inFlight are the buy states whose progress is polled.
var inFlight = map[string]bool{"pending_payment": true, "paid": true, "confirmed": true}

// OrdersView is what the My Orders tab shows.
type OrdersView struct {
	Orders []protocol.OrderView
	// Live is the polled status of the buys in flight, by order id.
	Live  map[string]protocol.GetOrderStatusResponse
	Error string
}

func (v OrdersView) status(o protocol.OrderView) string {
	if s, ok := v.Live[o.ID]; ok && s.Status != "" {
		return s.Status
	}
	return o.Status
}

func (v OrdersView) cancelable(o protocol.OrderView) bool {
	return o.Type == "buy" && v.status(o) == "pending_payment"
}

func (h *handler) orders() OrdersView {
	v := OrdersView{Live: map[string]protocol.GetOrderStatusResponse{}}
	var resp protocol.GetOrdersResponse
	if err := h.m.Call(protocol.TypeGetOrders, nil, &resp); err != nil {
		v.Error = err.Error()
		return v
	}
	v.Orders = resp.Orders
	// The market only answers the buyer about an order's live status.
	for _, o := range v.Orders {
		if o.Type == "buy" && inFlight[o.Status] {
			var s protocol.GetOrderStatusResponse
			if h.m.Call(protocol.TypeGetOrderStatus, protocol.GetOrderStatusRequest{OrderID: o.ID}, &s) == nil {
				v.Live[o.ID] = s
			}
		}
	}
	return v
}

// ListingsView is what the My Listings tab shows.
type ListingsView struct {
	Listings []protocol.ListingView
	Error    string
}

func (h *handler) listings() ListingsView {
	var v ListingsView
	var resp protocol.GetListingsResponse
	if err := h.m.Call(protocol.TypeGetListings, nil, &resp); err != nil {
		v.Error = err.Error()
	}
	v.Listings = resp.Listings
	return v
}

// table serves a polled table, or sends a disconnected browser to the
// connect screen.
func (h *handler) table(c func() templ.Component) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if connected, _, _ := h.m.Status(); !connected {
			redirect(w, r, "")
			return
		}
		render(w, r, c())
	}
}

// cancel cancels an order or a listing and answers with the refreshed table
// and the market's notice, or with its error.
func (h *handler) cancel(tab string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := h.page(r, tab)
		if !p.Connected {
			redirect(w, r, "")
			return
		}
		var msg protocol.MessageData
		var err error
		if tab == "orders" {
			err = h.m.Call(protocol.TypeCancelOrder, protocol.CancelOrderRequest{OrderID: r.FormValue("id")}, &msg)
			if msg.Message == "" {
				msg.Message = "Order canceled."
			}
		} else {
			err = h.m.Call(protocol.TypeCancelListing, protocol.CancelListingRequest{ListingID: r.FormValue("id")}, &msg)
			if msg.Message == "" {
				msg.Message = "Offer canceled."
			}
		}
		f := &Flash{Kind: "info", Text: msg.Message}
		if err != nil {
			f = &Flash{Kind: "danger", Text: err.Error()}
		}
		if !isHTMX(r) {
			render(w, r, TabPage(p, h.tabBodyWith(p, f)))
			return
		}
		if tab == "orders" {
			render(w, r, cancelReply(f, OrdersTable(h.orders())))
		} else {
			render(w, r, cancelReply(f, ListingsTable(h.listings())))
		}
	}
}

func (h *handler) tabBodyWith(p *Page, f *Flash) templ.Component {
	if p.Tab == "orders" {
		return OrdersTab(h.orders(), f)
	}
	return ListingsTab(h.listings(), f)
}
