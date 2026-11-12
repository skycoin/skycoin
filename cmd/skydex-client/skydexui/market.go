// Package skydexui cmd/skydex-client/skydexui/market.go c3-ui-dex
package skydexui

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"

	"github.com/skycoin/skycoin/src/skydex/protocol"
)

func (h *handler) products() MarketView {
	var v MarketView
	var resp protocol.GetProductsResponse
	if err := h.m.Call(protocol.TypeGetProducts, nil, &resp); err != nil {
		v.LoadErr = err.Error()
	}
	v.Products = resp.Products
	return v
}

func (h *handler) productGrid(w http.ResponseWriter, r *http.Request) {
	p := h.page(r, "market")
	if !p.Connected {
		redirect(w, r, "")
		return
	}
	render(w, r, ProductGrid(p, h.products(), false))
}

// needsRegister is the market telling an unknown user to register first,
// which the React client turns into a trip to Settings.
func needsRegister(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "register")
}

// marketAction runs a buy or a new listing. ok renders the instructions into
// the view on success.
func (h *handler) marketAction(w http.ResponseWriter, r *http.Request, call func() error, ok func(*MarketView) (target string, box templ.Component)) {
	p := h.page(r, "market")
	if !p.Connected {
		redirect(w, r, "")
		return
	}
	err := call()
	if err != nil && needsRegister(err) {
		redirect(w, r, "settings")
		return
	}
	v := h.products()
	if err != nil {
		v.Error = err.Error()
		if isHTMX(r) {
			render(w, r, marketError(v.Error, false))
			return
		}
		render(w, r, TabPage(p, MarketTab(p, v)))
		return
	}
	target, box := ok(&v)
	if !isHTMX(r) {
		render(w, r, TabPage(p, MarketTab(p, v)))
		return
	}
	w.Header().Set("HX-Retarget", target)
	w.Header().Set("HX-Reswap", "innerHTML")
	render(w, r, templ.Join(box, marketError("", true), ProductGrid(p, v, true)))
}

func (h *handler) buy(w http.ResponseWriter, r *http.Request) {
	var o protocol.BuyProductResponse
	h.marketAction(w, r, func() error {
		return h.m.Call(protocol.TypeBuyProduct, protocol.BuyProductRequest{ProductID: r.FormValue("product_id")}, &o)
	}, func(v *MarketView) (string, templ.Component) {
		v.Payment = &o
		return "#market-payment", PaymentBox(o)
	})
}

func (h *handler) createListing(w http.ResponseWriter, r *http.Request) {
	var d protocol.CreateListingResponse
	h.marketAction(w, r, func() error {
		amount, _ := strconv.ParseFloat(strings.TrimSpace(r.FormValue("amount")), 64) //nolint:errcheck
		price, _ := strconv.ParseFloat(strings.TrimSpace(r.FormValue("price")), 64)   //nolint:errcheck
		return h.m.Call(protocol.TypeCreateListing, protocol.CreateListingRequest{
			SellCoin:        r.FormValue("sell_coin"),
			Amount:          amount,
			Price:           price,
			PaymentCurrency: r.FormValue("payment_currency"),
		}, &d)
	}, func(v *MarketView) (string, templ.Component) {
		v.Deposit = &d
		// skydex.js closes and clears the form on this event.
		w.Header().Set("HX-Trigger", "listing-created")
		return "#market-deposit", DepositBox(d)
	})
}
