package server

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/skycoin/skycoin/src/skydex/db"
	"github.com/skycoin/skycoin/src/skydex/jobs"
	"github.com/skycoin/skycoin/src/skydex/protocol"
	"github.com/skycoin/skycoin/src/skydex/walletaddr"
)

// Decimal precision per asset. SKY droplets are 1e-6, but the Skycoin network
// only accepts transaction amounts to 3 decimals (MaxDropletPrecision), so
// on-chain SKY amounts (deposits, deliveries, refunds) are rounded to 3. The
// supported external payment coins (BTC/BCH/LTC/DOGE/DASH) use 8.
const (
	skyOnChainDecimals = 3
	paymentDecimals    = 8
)

// paymentTol is the amount-collision tolerance for external payment coins: half
// their smallest unit, so two amounts are "the same" only if they round to the
// same on-chain value.
const paymentTol = 0.5e-8

// Upper bounds on request amounts. Beyond these, base*10^decimals would exceed
// float64's exact-integer range (2^53 ≈ 9e15), so nonRound's small delta would
// be lost to rounding and uniqueness could silently break. The caps keep
// base*10^decimals ≤ 1e15 and sit far above any real trade (SKY total supply is
// ~1e8; a payment-coin price of 1e7 is already absurd).
const (
	maxAmount = 1e9
	maxPrice  = 1e7
)

// roundToDecimals rounds v to an asset's on-chain precision (SKY 6, payment coin
// 8), so stored amounts are always valid, representable values and the escrow
// deposit reconciles exactly with later delivery.
func roundToDecimals(v float64, decimals int) float64 {
	scale := math.Pow(10, float64(decimals))
	return math.Round(v*scale) / scale
}

// errNoUniqueAmount means the non-round amount space to a wallet is momentarily
// exhausted (too many concurrent pending payments of the same base amount).
var errNoUniqueAmount = errors.New("could not allocate a unique payment amount")

// allocUniqueAmount generates a non-round amount for base that is not already
// pending (per the exists check), retrying on collision. The caller must hold a
// lock spanning this call and the subsequent insert so the chosen amount can't
// be taken by a concurrent allocation.
func allocUniqueAmount(base float64, decimals int, exists func(float64) (bool, error)) (float64, error) { //nolint
	for i := 0; i < 100; i++ {
		amt := nonRound(base, decimals)
		dup, err := exists(amt)
		if err != nil {
			return 0, err
		}
		if !dup {
			return amt, nil
		}
	}
	return 0, errNoUniqueAmount
}

// success builds a success response, degrading to an internal error if the
// payload cannot be marshaled.
func success(id string, data any) protocol.Envelope {
	env, err := protocol.NewResponse(id, protocol.StatusSuccess, data)
	if err != nil {
		return protocol.ErrorResponse(id, protocol.CodeInternalError, "failed to encode response")
	}
	return env
}

// internal logs op with err and returns a generic internal-error response so
// implementation details never leak to clients.
func (s *Server) internal(id, op string, err error) protocol.Envelope {
	s.log.WithError(err).Errorf("skydex-market: %s failed", op)
	return protocol.ErrorResponse(id, protocol.CodeInternalError, "internal error")
}

// unfreeze best-effort returns a product to the market after a failed buy.
func (s *Server) unfreeze(productID string) {
	if err := s.db.UnfreezeProduct(productID); err != nil {
		s.log.WithError(err).Errorf("skydex-market: failed to unfreeze product %s", productID)
	}
}

// handleRegister upserts the caller's wallet addresses, keyed by the
// authenticated connection public key (never a self-reported field).
func (s *Server) handleRegister(pk string, req protocol.Envelope) protocol.Envelope {
	var r protocol.RegisterRequest
	if err := req.Bind(&r); err != nil {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "invalid register payload")
	}
	if strings.TrimSpace(r.WalletSKY) == "" {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "wallet_sky is required")
	}
	// Collect the per-currency payout addresses: the fixed fields plus the dynamic
	// Wallets map (which may carry arbitrary operator-added coins). A value in
	// Wallets overrides a fixed field for the same coin.
	payouts := map[string]string{}
	for _, w := range []struct{ currency, addr string }{
		{"BTC", r.WalletBTC}, {"BCH", r.WalletBCH}, {"LTC", r.WalletLTC},
		{"USDT_ERC20", r.WalletUSDTERC20}, {"USDT_TRC20", r.WalletUSDTTRC20},
	} {
		if a := strings.TrimSpace(w.addr); a != "" {
			payouts[w.currency] = a
		}
	}
	for cur, addr := range r.Wallets {
		if a := strings.TrimSpace(addr); a != "" {
			payouts[strings.ToUpper(strings.TrimSpace(cur))] = a
		}
	}
	// Validate the Skycoin address and any payout address for a coin we know how
	// to check (SKY/BTC/LTC). walletaddr.Validate accepts a custom coin it has no
	// validator for, so a typo there is the operator/seller's responsibility.
	if err := walletaddr.Validate("SKY", r.WalletSKY); err != nil {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidWallet, "invalid SKY wallet address")
	}
	for cur, addr := range payouts {
		if err := walletaddr.Validate(cur, addr); err != nil {
			return protocol.ErrorResponse(req.ID, protocol.CodeInvalidWallet, "invalid "+cur+" wallet address")
		}
	}

	u := &db.User{PubKey: pk, WalletSKY: r.WalletSKY}
	if err := s.db.CreateUser(u); err != nil {
		return s.internal(req.ID, "register", err)
	}
	for cur, addr := range payouts {
		if err := s.db.SetUserWallet(pk, cur, addr); err != nil {
			return s.internal(req.ID, "register.wallet", err)
		}
	}
	return success(req.ID, protocol.MessageData{Message: "User registered successfully"})
}

// handleGetCurrencies returns the payment currencies this market accepts, i.e.
// the ones whose explorer the operator has configured.
func (s *Server) handleGetCurrencies(req protocol.Envelope) protocol.Envelope {
	cur, err := s.db.AvailableCurrencies()
	if err != nil {
		return s.internal(req.ID, "get_currencies", err)
	}
	if cur == nil {
		cur = []string{}
	}
	sell, err := s.db.AvailableSellCoins()
	if err != nil {
		return s.internal(req.ID, "get_currencies.sellcoins", err)
	}
	if sell == nil {
		sell = []string{}
	}
	return success(req.ID, protocol.GetCurrenciesResponse{
		Currencies: cur,
		SellCoins:  sell,
		MarketName: s.db.GetMarketName(),
	})
}

// handleGetProducts returns all active products.
func (s *Server) handleGetProducts(req protocol.Envelope) protocol.Envelope {
	products, err := s.db.GetActiveProducts()
	if err != nil {
		return s.internal(req.ID, "get_products", err)
	}
	views := make([]protocol.ProductView, 0, len(products))
	for _, p := range products {
		views = append(views, protocol.ProductView{
			ID:              p.ID,
			SellerPubKey:    p.SellerPubKey,
			SellCoin:        p.SellCoin,
			Amount:          p.Amount,
			Price:           p.Price,
			PaymentCurrency: p.PaymentCurrency,
			CreatedAt:       p.CreatedAt.Format(time.RFC3339),
		})
	}
	return success(req.ID, protocol.GetProductsResponse{Products: views})
}

// handleCreateListing records a sell order awaiting the seller's SKY deposit to
// the market wallet, and returns the non-round deposit amount to transfer.
func (s *Server) handleCreateListing(pk string, req protocol.Envelope) protocol.Envelope {
	var r protocol.CreateListingRequest
	if err := req.Bind(&r); err != nil {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "invalid listing payload")
	}
	// Default the sell coin to SKY when the client omits it (backward compatible).
	sellCoin := strings.ToUpper(strings.TrimSpace(r.SellCoin))
	if sellCoin == "" {
		sellCoin = "SKY"
	}
	if r.Amount <= 0 || r.Price <= 0 {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "amount and price must be positive")
	}
	if r.Amount > maxAmount || r.Price > maxPrice {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "amount or price is too large")
	}
	// Normalize the sell-coin amount to Skycoin-family on-chain precision (3 dp),
	// the network's spendable precision, so the seller can deposit and the market
	// can later deliver the exact amount. (Price is normalized below, once we know
	// whether the payment currency is a fibercoin or an external coin.) Re-check
	// positivity: a value below half a unit rounds to zero.
	r.Amount = roundToDecimals(r.Amount, skyOnChainDecimals)
	if r.Amount <= 0 {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "amount is too small to represent")
	}

	// Enforce the operator's trade-size bounds (0 = unset).
	if minTrade, err := s.db.GetMinTradeSKY(); err == nil && minTrade > 0 && r.Amount < minTrade {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest,
			fmt.Sprintf("amount is below the minimum trade size of %s", strconv.FormatFloat(minTrade, 'f', -1, 64)))
	}
	if maxTrade, err := s.db.GetMaxTradeSKY(); err == nil && maxTrade > 0 && r.Amount > maxTrade {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest,
			fmt.Sprintf("amount is above the maximum trade size of %s", strconv.FormatFloat(maxTrade, 'f', -1, 64)))
	}

	// The seller must be registered (wallet addresses on file).
	if u, err := s.db.GetUser(pk); err != nil {
		return s.internal(req.ID, "create_listing.get_user", err)
	} else if u == nil {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "register your wallets first")
	}

	// The chosen sell coin must be enabled at this market.
	sellOK, err := s.db.IsSellCoinAvailable(sellCoin)
	if err != nil {
		return s.internal(req.ID, "create_listing.sellcoin", err)
	}
	if !sellOK {
		return protocol.ErrorResponse(req.ID, protocol.CodeCurrencyUnavailable, "sell coin not available at this market: "+sellCoin)
	}

	// Resolve the payment currency. It can be an external explorer coin (BTC/LTC/…)
	// or a fibercoin (incl. SKY) that the buyer pays peer-to-peer on its own chain.
	payCur := strings.ToUpper(strings.TrimSpace(r.PaymentCurrency))
	if payCur == sellCoin {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "sell coin and payment currency must differ")
	}
	payOK, err := s.db.IsPaymentAvailable(payCur)
	if err != nil {
		return s.internal(req.ID, "create_listing.currency", err)
	}
	if !payOK {
		return protocol.ErrorResponse(req.ID, protocol.CodeCurrencyUnavailable, "payment currency not available at this market: "+payCur)
	}
	// A fibercoin payment is a FIXED on-chain amount at 3-dp precision; an external
	// coin uses 8-dp (and later a unique non-round amount). Normalize the price to
	// whichever applies, then re-check positivity.
	isFiberPay, err := s.db.IsSellCoinAvailable(payCur)
	if err != nil {
		return s.internal(req.ID, "create_listing.fiberpay", err)
	}
	payDec := paymentDecimals
	if isFiberPay {
		payDec = skyOnChainDecimals
	}
	r.PaymentCurrency = payCur
	r.Price = roundToDecimals(r.Price, payDec)
	if r.Price <= 0 {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "price is too small to represent")
	}

	// The seller must have a payout address for the payment coin, or the buyer's
	// payment would have nowhere to go. (A fibercoin payment resolves to the
	// seller's Skycoin address, which is always set.)
	if payout, err := s.db.GetUserWallet(pk, payCur); err != nil {
		return s.internal(req.ID, "create_listing.payout", err)
	} else if strings.TrimSpace(payout) == "" {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest,
			"add a "+payCur+" payout address in Settings before selling for "+payCur)
	}

	// The deposit target is the escrow wallet of the chosen sell coin.
	marketWallet, err := s.db.SellCoinWallet(sellCoin)
	if err != nil {
		return s.internal(req.ID, "create_listing.wallet", err)
	}
	if marketWallet == "" {
		return protocol.ErrorResponse(req.ID, protocol.CodeInternalError, "escrow wallet not configured for "+sellCoin)
	}

	expiry := s.listingExpiry()

	// The seller sends the exact round amount; the deposit is identified by their
	// registered SKY address (the tx sender) within this listing's window. That
	// only works if the seller has no other pending listing whose window overlaps
	// — otherwise a single deposit from that address is ambiguous. So allow one
	// pending listing per seller at a time. Hold allocMu across the check and the
	// insert so two concurrent requests can't both pass.
	s.allocMu.Lock()
	defer s.allocMu.Unlock()
	hasPending, err := s.db.SellerHasPendingListing(pk)
	if err != nil {
		return s.internal(req.ID, "create_listing.pending_check", err)
	}
	if hasPending {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest,
			"you already have a pending listing; complete, cancel, or wait for it to expire before creating another")
	}

	// The market commission (in the sell coin) is paid by the seller on top of the
	// amount: the seller deposits amount + commission, the buyer receives amount,
	// and the market retains the commission. So the expected deposit is
	// amount+commission.
	rate, _ := s.db.GetCommissionRatePercent() //nolint:errcheck
	minC, _ := s.db.GetCommissionMinSKY()      //nolint:errcheck
	maxC, _ := s.db.GetCommissionMaxSKY()      //nolint:errcheck
	commission := jobs.Commission(r.Amount, rate, minC, maxC)
	expected := roundToDecimals(r.Amount+commission, skyOnChainDecimals)

	listing := &db.PendingListing{
		ID:              uuid.NewString(),
		SellerPubKey:    pk,
		SellCoin:        sellCoin,
		Amount:          r.Amount,
		ExpectedAmount:  expected, // amount + commission (the deposit)
		Price:           r.Price,
		PaymentCurrency: r.PaymentCurrency,
		Status:          "pending",
		ExpiresAt:       time.Now().UTC().Add(expiry),
	}
	if err := s.db.CreatePendingListing(listing); err != nil {
		return s.internal(req.ID, "create_listing", err)
	}

	return success(req.ID, protocol.CreateListingResponse{
		ListingID:      listing.ID,
		SellCoin:       listing.SellCoin,
		Amount:         listing.Amount,
		Commission:     commission,
		ExpectedAmount: listing.ExpectedAmount,
		MarketWallet:   marketWallet,
		ExpiresAt:      listing.ExpiresAt.Format(time.RFC3339),
	})
}

// handleBuyProduct freezes a product for the buyer and opens a buy order,
// returning the non-round amount to pay the seller.
func (s *Server) handleBuyProduct(pk string, req protocol.Envelope) protocol.Envelope {
	var r protocol.BuyProductRequest
	if err := req.Bind(&r); err != nil || r.ProductID == "" {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "invalid buy payload")
	}

	banned, err := s.db.IsUserBanned(pk)
	if err != nil {
		return s.internal(req.ID, "buy.ban_check", err)
	}
	if banned {
		return protocol.ErrorResponse(req.ID, protocol.CodeUserBanned, "you are temporarily banned from buying")
	}

	// The buyer must be registered so we can deliver SKY to them later.
	if u, err := s.db.GetUser(pk); err != nil {
		return s.internal(req.ID, "buy.get_user", err)
	} else if u == nil {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "register your wallets first")
	}

	product, err := s.db.GetProduct(r.ProductID)
	if err != nil {
		return s.internal(req.ID, "buy.get_product", err)
	}
	if product == nil {
		return protocol.ErrorResponse(req.ID, protocol.CodeProductNotFound, "product not found")
	}
	if product.Status != "active" {
		return protocol.ErrorResponse(req.ID, protocol.CodeProductUnavailable, "this product is no longer available")
	}

	// A buyer who has canceled (or let expire) at least buy_cancel_limit orders
	// for this product may not buy it again until the operator clears the block.
	blocked, err := s.db.BuyerBlockedFromProduct(pk, product.ID, s.db.GetBuyCancelLimit())
	if err != nil {
		return s.internal(req.ID, "buy.block_check", err)
	}
	if blocked {
		return protocol.ErrorResponse(req.ID, protocol.CodeBuyerBlocked, "you can no longer buy this product")
	}

	// Resolve the seller's wallet for the product's payment currency.
	sellerWallet, err := s.db.GetUserWallet(product.SellerPubKey, product.PaymentCurrency)
	if err != nil {
		return s.internal(req.ID, "buy.seller_wallet", err)
	}
	if sellerWallet == "" {
		return protocol.ErrorResponse(req.ID, protocol.CodeInternalError, "seller wallet not configured")
	}

	// Freeze first: this atomically moves active->frozen and rejects a racing
	// second buyer (design §5 product freeze).
	if err := s.db.FreezeProduct(product.ID, pk); err != nil {
		return protocol.ErrorResponse(req.ID, protocol.CodeProductUnavailable, "this product is no longer available")
	}

	// Determine the exact amount the buyer must pay. For a FIBERCOIN payment the
	// buyer transfers a FIXED amount on that coin's chain; it is verified by sender
	// + window on the coin's node, so no unique non-round amount is needed (a fixed
	// 100 FIBER1, not 100.002). For an external explorer coin we allocate a unique
	// non-round amount so payments to a shared address are individually
	// identifiable. Either way we hold allocMu across the check and the insert.
	isFiberPay, err := s.db.IsSellCoinAvailable(product.PaymentCurrency)
	if err != nil {
		s.unfreeze(product.ID)
		return s.internal(req.ID, "buy.fiberpay", err)
	}
	s.allocMu.Lock()
	var expected float64
	if isFiberPay {
		expected = roundToDecimals(product.Price, skyOnChainDecimals)
		// A single buyer can't have two identical fixed payments to the same seller
		// in the same coin — the sender no longer disambiguates them. Different
		// buyers paying the same amount are fine (the sender tells them apart).
		dup, derr := s.db.BuyerHasPendingPayment(pk, sellerWallet, product.PaymentCurrency, expected, 0.5e-3)
		if derr != nil {
			s.allocMu.Unlock()
			s.unfreeze(product.ID)
			return s.internal(req.ID, "buy.dupcheck", derr)
		}
		if dup {
			s.allocMu.Unlock()
			s.unfreeze(product.ID)
			return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest,
				"you already have an identical pending payment to this seller; complete it first")
		}
	} else {
		expected, err = allocUniqueAmount(product.Price, paymentDecimals, func(a float64) (bool, error) {
			return s.db.OrderPaymentAmountExists(sellerWallet, product.PaymentCurrency, a, paymentTol)
		})
		if err != nil {
			s.allocMu.Unlock()
			s.unfreeze(product.ID)
			return s.internal(req.ID, "buy.alloc", err)
		}
	}
	order := &db.Order{
		ID:                    uuid.NewString(),
		ProductID:             product.ID,
		BuyerPubKey:           pk,
		SellCoin:              product.SellCoin,
		Amount:                product.Amount,
		Price:                 product.Price,
		PaymentCurrency:       product.PaymentCurrency,
		ExpectedPaymentAmount: expected,
		SellerWallet:          sellerWallet,
		Status:                "pending_payment",
		ExpiresAt:             time.Now().UTC().Add(s.orderExpiry()),
	}
	err = s.db.CreateOrder(order)
	s.allocMu.Unlock()
	if err != nil {
		s.unfreeze(product.ID)
		return s.internal(req.ID, "buy.create_order", err)
	}

	return success(req.ID, protocol.BuyProductResponse{
		OrderID:               order.ID,
		SellCoin:              order.SellCoin,
		Amount:                order.Amount,
		Price:                 order.Price,
		PaymentCurrency:       order.PaymentCurrency,
		ExpectedPaymentAmount: order.ExpectedPaymentAmount,
		SellerWallet:          order.SellerWallet,
		ExpiresAt:             order.ExpiresAt.Format(time.RFC3339),
	})
}

// handleCancelListing cancels the caller's own sell offer. It works while the
// listing is still pending (no deposit yet — nothing to return) and after it has
// become an active product, provided no buyer has selected (frozen) it. For a
// confirmed listing the derived product is deactivated and the escrowed SKY is
// returned by the Return Scheduler (measured from when the SKY was deposited: if
// that was already longer than the return delay ago, it goes back on the next
// scheduler tick, otherwise once the delay elapses).
func (s *Server) handleCancelListing(pk string, req protocol.Envelope) protocol.Envelope {
	var r protocol.CancelListingRequest
	if err := req.Bind(&r); err != nil || r.ListingID == "" {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "invalid cancel payload")
	}

	listing, err := s.db.GetPendingListing(r.ListingID)
	if err != nil {
		return s.internal(req.ID, "cancel.get_listing", err)
	}
	// Treat "not yours" the same as "not found" so listing IDs aren't probeable.
	if listing == nil || listing.SellerPubKey != pk {
		return protocol.ErrorResponse(req.ID, protocol.CodeListingNotFound, "listing not found")
	}

	switch listing.Status {
	case "pending":
		// No deposit confirmed yet: just cancel. The Return Scheduler still
		// checks the chain in case a deposit lands late.
		if err := s.db.UpdatePendingListingStatus(listing.ID, "canceled", ""); err != nil {
			return s.internal(req.ID, "cancel.update", err)
		}
		return success(req.ID, protocol.MessageData{Message: "Listing canceled."})

	case "confirmed":
		// The deposit settled and became a product. It can only be canceled
		// while still active — once a buyer freezes or buys it, the seller is
		// committed to the trade.
		product, err := s.db.GetProductByListingID(listing.ID)
		if err != nil {
			return s.internal(req.ID, "cancel.get_product", err)
		}
		if product != nil && product.Status != "active" {
			return protocol.ErrorResponse(req.ID, protocol.CodeProductUnavailable, "a buyer has already selected this offer; it can no longer be canceled")
		}
		if product != nil {
			if err := s.db.MarkProductCancelled(product.ID); err != nil {
				return s.internal(req.ID, "cancel.deactivate_product", err)
			}
		}
		if err := s.db.UpdatePendingListingStatus(listing.ID, "canceled", listing.TxHash); err != nil {
			return s.internal(req.ID, "cancel.update", err)
		}
		return success(req.ID, protocol.MessageData{Message: "Offer canceled. Your SKY will be returned shortly."})

	default:
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "listing can no longer be canceled")
	}
}

// handleCancelOrder cancels the caller's own in-flight buy order before payment
// is observed, releasing the product back to the market. The buyer is then
// blocked from re-buying that same product (BuyerBlockedFromProduct).
func (s *Server) handleCancelOrder(pk string, req protocol.Envelope) protocol.Envelope {
	var r protocol.CancelOrderRequest
	if err := req.Bind(&r); err != nil || r.OrderID == "" {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "invalid cancel payload")
	}

	order, err := s.db.GetOrder(r.OrderID)
	if err != nil {
		return s.internal(req.ID, "cancel_order.get", err)
	}
	// Only the buyer on the order may cancel it; hide others as not-found.
	if order == nil || order.BuyerPubKey != pk {
		return protocol.ErrorResponse(req.ID, protocol.CodeOrderNotFound, "order not found")
	}
	// Only cancelable before payment is observed. Once paid/confirmed the
	// payment is in flight and the trade must settle.
	if order.Status != "pending_payment" {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "this order can no longer be canceled")
	}

	if err := s.db.MarkOrderCanceled(order.ID); err != nil {
		return s.internal(req.ID, "cancel_order.update", err)
	}
	// Release the product so other buyers can purchase it.
	s.unfreeze(order.ProductID)
	// A voluntary cancel counts as a freeze violation, same as letting the order
	// expire: a buyer who repeatedly freezes then backs out is banned by the Ban
	// Manager (design §5). Best-effort — the cancel itself already succeeded.
	if err := s.db.CreateFreezeViolation(order.BuyerPubKey, order.ID); err != nil {
		s.log.WithError(err).Errorf("skydex-market: failed to record freeze violation for %s", order.BuyerPubKey)
	}

	// Tell the buyer where they stand against the per-product cancel limit.
	limit := s.db.GetBuyCancelLimit()
	msg := "Order canceled."
	if n, err := s.db.CountBuyerProductCancels(order.BuyerPubKey, order.ProductID); err == nil {
		if n >= limit {
			msg = "Order canceled. You can no longer buy this product (cancellation limit reached)."
		} else if remaining := limit - n; remaining == 1 {
			msg = "Order canceled. You can buy this product one more time before you're blocked from it."
		} else {
			msg = fmt.Sprintf("Order canceled. You can buy this product %d more times before you're blocked from it.", remaining)
		}
	}
	return success(req.ID, protocol.MessageData{Message: msg})
}

// handleGetOrders returns the caller's buy and sell orders.
func (s *Server) handleGetOrders(pk string, req protocol.Envelope) protocol.Envelope {
	buys, err := s.db.GetBuyerOrders(pk)
	if err != nil {
		return s.internal(req.ID, "get_orders.buyer", err)
	}
	sells, err := s.db.GetSellerOrders(pk)
	if err != nil {
		return s.internal(req.ID, "get_orders.seller", err)
	}

	views := make([]protocol.OrderView, 0, len(buys)+len(sells))
	for _, o := range buys {
		views = append(views, orderView(o, "buy"))
	}
	for _, o := range sells {
		views = append(views, orderView(o, "sell"))
	}
	return success(req.ID, protocol.GetOrdersResponse{Orders: views})
}

// handleGetListings returns the caller's own pending sell listings so the seller
// can track each one's lifecycle: pending (awaiting the SKY deposit) -> confirmed
// (deposit detected on-chain; the listing is now a purchasable product). It also
// returns the market wallet so the UI can restate the deposit target for
// still-pending listings.
func (s *Server) handleGetListings(pk string, req protocol.Envelope) protocol.Envelope {
	listings, err := s.db.GetSellerPendingListings(pk)
	if err != nil {
		return s.internal(req.ID, "get_listings", err)
	}
	wallet, _ := s.db.SellCoinWallet("SKY") //nolint:errcheck //SKY default for back-compat
	views := make([]protocol.ListingView, 0, len(listings))
	for _, l := range listings {
		coinWallet, _ := s.db.SellCoinWallet(l.SellCoin) //nolint:errcheck
		v := protocol.ListingView{
			ID:              l.ID,
			SellCoin:        l.SellCoin,
			Amount:          l.Amount,
			ExpectedAmount:  l.ExpectedAmount,
			Price:           l.Price,
			PaymentCurrency: l.PaymentCurrency,
			Status:          l.Status,
			MarketWallet:    coinWallet,
			ExpiresAt:       l.ExpiresAt.Format(time.RFC3339),
			CreatedAt:       l.CreatedAt.Format(time.RFC3339),
			TxHash:          l.TxHash,
		}
		if l.ConfirmedAt != nil {
			v.ConfirmedAt = l.ConfirmedAt.Format(time.RFC3339)
		}
		// A pending listing (no deposit yet) is always cancelable. A confirmed
		// one is only cancelable while its product is still active — once a buyer
		// freezes/buys it the seller is committed, so surface that to the client
		// instead of showing a Cancel button that the market would reject.
		switch l.Status {
		case "pending":
			v.Cancelable = true
		case "confirmed":
			if product, err := s.db.GetProductByListingID(l.ID); err == nil && product != nil {
				v.ProductStatus = product.Status
				v.Cancelable = product.Status == "active"
			} else {
				// No product row yet (just confirmed) — treat as still cancelable.
				v.Cancelable = true
			}
		}
		views = append(views, v)
	}
	return success(req.ID, protocol.GetListingsResponse{Listings: views, MarketWallet: wallet})
}

// handleGetOrderStatus returns the live status of one of the caller's orders.
func (s *Server) handleGetOrderStatus(pk string, req protocol.Envelope) protocol.Envelope {
	var r protocol.GetOrderStatusRequest
	if err := req.Bind(&r); err != nil || r.OrderID == "" {
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "invalid order-status payload")
	}
	order, err := s.db.GetOrder(r.OrderID)
	if err != nil {
		return s.internal(req.ID, "get_order_status", err)
	}
	// Only the buyer on the order may query it; hide others as not-found.
	if order == nil || order.BuyerPubKey != pk {
		return protocol.ErrorResponse(req.ID, protocol.CodeOrderNotFound, "order not found")
	}

	requiredConfs := jobs.RequiredConfirmations
	if c, err := s.db.GetConfirmationsRequired(); err == nil && c > 0 {
		requiredConfs = c
	}
	resp := protocol.GetOrderStatusResponse{
		OrderID:               order.ID,
		Status:                order.Status,
		Confirmations:         order.Confirmations,
		RequiredConfirmations: requiredConfs,
		PaymentTxHash:         order.PaymentTxHash,
	}
	if order.PaidAt != nil {
		resp.PaidAt = order.PaidAt.Format(time.RFC3339)
	}
	return success(req.ID, resp)
}

// orderView projects a db.Order into the client-facing view with a buy/sell tag.
func orderView(o *db.Order, kind string) protocol.OrderView {
	return protocol.OrderView{
		ID:                    o.ID,
		Type:                  kind,
		ProductID:             o.ProductID,
		SellCoin:              o.SellCoin,
		Amount:                o.Amount,
		Price:                 o.Price,
		PaymentCurrency:       o.PaymentCurrency,
		ExpectedPaymentAmount: o.ExpectedPaymentAmount,
		SellerWallet:          o.SellerWallet,
		Status:                o.Status,
		ExpiresAt:             o.ExpiresAt.Format(time.RFC3339),
		CreatedAt:             o.CreatedAt.Format(time.RFC3339),
	}
}

func (s *Server) listingExpiry() time.Duration {
	if m, err := s.db.GetListingExpiryMinutes(); err == nil && m > 0 {
		return time.Duration(m) * time.Minute
	}
	return 15 * time.Minute
}

func (s *Server) orderExpiry() time.Duration {
	if m, err := s.db.GetOrderExpiryMinutes(); err == nil && m > 0 {
		return time.Duration(m) * time.Minute
	}
	return 15 * time.Minute
}

// nonRound returns base plus a small random delta at the asset's precision, so
// the on-chain amount is unique and precisely identifiable in the recipient
// wallet (design §1 non-round-number validation). decimals is the asset's
// precision (6 for SKY, 8 for the payment coins); the result is rounded to it so
// it is always a valid on-chain amount.
func nonRound(base float64, decimals int) float64 {
	scale := math.Pow(10, float64(decimals))
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return math.Round(base*scale) / scale
	}
	// 1..1000 least-significant units of delta.
	n := binary.BigEndian.Uint32(b[:])%1000 + 1
	return math.Round(base*scale+float64(n)) / scale
}
