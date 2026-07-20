package server_test

import (
	"math"
	"net"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/skydex/db"
	"github.com/skycoin/skycoin/src/skydex/market"
	"github.com/skycoin/skycoin/src/skydex/protocol"
	"github.com/skycoin/skycoin/src/skydex/server"
)

// Valid mainnet payout addresses (SKY base58check, BTC P2PKH) — registration
// now validates addresses, so tests use real ones. Generated deterministically;
// see internal/skydex-market/walletaddr for the format rules.
const (
	skySeller = "FbmPhy5bhsMX8JNeAQdQ4W3DeCYFH97FBg"
	skyBuyer  = "23Lqf6WpmaiFdzr4g5gerqUBu8H7SyeDHJU"
	skyBuyer2 = "2joBpuNB6z3MReAjzQwGdxphxD5BmbrGV1k"
	skyOther  = "xPHfK6xn5LZvWAPB5EV6A62WyzQLywci1J"
	btcSeller = "19wiPpvWPxHEmcANKzwUjiRgABQk8wZzCp"
	btcBuyer  = "15rB1CxysqYZ6soyEoVpiUcz7Sb2yJwpS2"
	btcBuyer2 = "17qzzrRtvo6FYVHiN2xdtyrWXBEhujfAmN"
)

// newTestDB spins up a real (temp-file) SQLite database with migrations and
// default config applied.
func newTestDB(t *testing.T) *db.Database {
	t.Helper()
	path := filepath.Join(t.TempDir(), "market.db")
	database, err := db.New(path, "")
	if err != nil {
		t.Fatalf("db.New: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() }) //nolint
	if err := database.Migrate(); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := database.InitDefaultConfig(); err != nil {
		t.Fatalf("init config: %v", err)
	}
	// Start from a clean payment-coin slate so each test enables exactly what it
	// needs (InitDefaultConfig now enables BTC + LTC out of the box).
	for _, k := range []string{"explorer_btc_provider", "explorer_ltc_provider"} {
		if err := database.SetConfig(k, ""); err != nil {
			t.Fatal(err)
		}
	}
	// Make SKY a complete, enabled sell coin: the seeded row is disabled with no
	// escrow seed, and availability now requires node + seed + address. Set the
	// legacy escrow keys (SellCoinConfig bridges them for SKY) and enable it.
	if err := database.SetConfig("sky_wallet_seed", "test escrow seed words for the market"); err != nil {
		t.Fatal(err)
	}
	if err := database.SetConfig("wallet_sky", "sky-market-wallet"); err != nil {
		t.Fatal(err)
	}
	if err := database.SetSellCoinEnabled("SKY", true); err != nil {
		t.Fatal(err)
	}
	return database
}

// zeroCommission disables the SKY commission so a listing's expected deposit
// equals its amount, keeping tests that aren't about commission simple.
func zeroCommission(t *testing.T, database *db.Database) {
	t.Helper()
	if err := database.SetConfig("commission_rate_percent", "0"); err != nil {
		t.Fatal(err)
	}
	if err := database.SetConfig("commission_min_sky", "0"); err != nil {
		t.Fatal(err)
	}
}

// dialTestServer wires a client Conn to a Server.Serve loop over net.Pipe,
// authenticating the client as buyerPK.
func dialTestServer(t *testing.T, database *db.Database, buyerPK string) *market.Conn {
	t.Helper()
	srv := server.New(database, nil, protocol.DefaultPort)
	clientEnd, serverEnd := net.Pipe()
	go srv.Serve(serverEnd, buyerPK)
	c := market.NewConn(clientEnd)
	t.Cleanup(func() { _ = c.Close() }) //nolint
	return c
}

// TestTradeRoundTrip drives register -> get_currencies -> create_listing ->
// get_products -> buy_product -> get_orders over the framed transport against a
// real database, verifying both the protocol and the market business logic.
func TestTradeRoundTrip(t *testing.T) {
	database := newTestDB(t)

	const sellerPK = "0311111111111111111111111111111111111111111111111111111111111111aa"
	const buyerPK = "0322222222222222222222222222222222222222222222222222222222222222bb"

	// Enable BTC as a payment currency, and set the market escrow wallet.
	if err := database.SetConfig("explorer_btc_provider", "esplora"); err != nil {
		t.Fatal(err)
	}
	if err := database.SetConfig("wallet_sky", "sky-market-wallet"); err != nil {
		t.Fatal(err)
	}
	zeroCommission(t, database)

	// --- seller registers and lists a product ---
	seller := dialTestServer(t, database, sellerPK)

	mustSuccess(t, seller, protocol.TypeRegister, protocol.RegisterRequest{
		WalletSKY: skySeller, WalletBTC: btcSeller,
	})

	// get_currencies should now include BTC (explorer configured) but not LTC.
	var cur protocol.GetCurrenciesResponse
	bindSuccess(t, seller, protocol.TypeGetCurrencies, nil, &cur)
	if !slices.Contains(cur.Currencies, "BTC") || slices.Contains(cur.Currencies, "LTC") {
		t.Fatalf("available currencies = %v, want BTC present and LTC absent", cur.Currencies)
	}

	// A currency without a configured explorer must be rejected.
	rejected, err := seller.Do(protocol.TypeCreateListing, protocol.CreateListingRequest{
		Amount: 10, Price: 2, PaymentCurrency: "LTC",
	})
	if err != nil {
		t.Fatalf("create_listing(LTC) transport error: %v", err)
	}
	if !rejected.IsError() || errCode(t, rejected) != protocol.CodeCurrencyUnavailable {
		t.Fatalf("create_listing(LTC) = %+v, want CURRENCY_UNAVAILABLE", rejected)
	}

	var listing protocol.CreateListingResponse
	bindSuccess(t, seller, protocol.TypeCreateListing, protocol.CreateListingRequest{
		Amount: 10, Price: 2, PaymentCurrency: "BTC",
	}, &listing)
	// The seller deposits the exact round amount (identified by sender + window).
	if listing.ExpectedAmount != 10 || listing.MarketWallet != "sky-market-wallet" {
		t.Fatalf("unexpected listing response: %+v", listing)
	}

	// The Listing Checker job (not part of this phase) would normally promote a
	// confirmed listing into products. Simulate that so we can exercise buying.
	activeProduct := &db.Product{
		ID: "prod-1", SellerPubKey: sellerPK, Amount: 10, Price: 2,
		PaymentCurrency: "BTC", Status: "active",
	}
	if err := database.CreateProduct(activeProduct); err != nil {
		t.Fatalf("seed product: %v", err)
	}

	// --- buyer registers, sees the product, and buys it ---
	buyer := dialTestServer(t, database, buyerPK)
	mustSuccess(t, buyer, protocol.TypeRegister, protocol.RegisterRequest{
		WalletSKY: skyBuyer, WalletBTC: btcBuyer,
	})

	var products protocol.GetProductsResponse
	bindSuccess(t, buyer, protocol.TypeGetProducts, nil, &products)
	if len(products.Products) != 1 || products.Products[0].ID != "prod-1" {
		t.Fatalf("get_products = %+v, want the one seeded product", products.Products)
	}

	var buy protocol.BuyProductResponse
	bindSuccess(t, buyer, protocol.TypeBuyProduct, protocol.BuyProductRequest{ProductID: "prod-1"}, &buy)
	if buy.SellerWallet != btcSeller || buy.ExpectedPaymentAmount <= 2 {
		t.Fatalf("unexpected buy response: %+v", buy)
	}

	// The product is now frozen: a second buy must be rejected.
	second, err := buyer.Do(protocol.TypeBuyProduct, protocol.BuyProductRequest{ProductID: "prod-1"})
	if err != nil {
		t.Fatalf("second buy transport error: %v", err)
	}
	if !second.IsError() || errCode(t, second) != protocol.CodeProductUnavailable {
		t.Fatalf("second buy = %+v, want PRODUCT_UNAVAILABLE", second)
	}

	// The buyer's order shows up in get_orders as a pending_payment buy.
	var orders protocol.GetOrdersResponse
	bindSuccess(t, buyer, protocol.TypeGetOrders, nil, &orders)
	if len(orders.Orders) != 1 || orders.Orders[0].Type != "buy" || orders.Orders[0].Status != "pending_payment" {
		t.Fatalf("get_orders = %+v, want one pending_payment buy", orders.Orders)
	}
}

// TestCreateListing_SellCoin verifies the sell coin is validated and routed to
// the right escrow wallet: an unknown/disabled coin is rejected, and an enabled
// fibercoin lists against its own escrow address.
func TestCreateListing_SellCoin(t *testing.T) {
	database := newTestDB(t)
	if err := database.SetConfig("explorer_btc_provider", "esplora"); err != nil {
		t.Fatal(err)
	}
	if err := database.SetConfig("wallet_sky", "sky-market-wallet"); err != nil {
		t.Fatal(err)
	}
	zeroCommission(t, database)

	const sellerPK = "0311111111111111111111111111111111111111111111111111111111111111aa"
	seller := dialTestServer(t, database, sellerPK)
	mustSuccess(t, seller, protocol.TypeRegister, protocol.RegisterRequest{WalletSKY: skySeller, WalletBTC: btcSeller})

	// get_currencies advertises SKY as a sell coin (seeded, enabled by default).
	var cur protocol.GetCurrenciesResponse
	bindSuccess(t, seller, protocol.TypeGetCurrencies, nil, &cur)
	if !slices.Contains(cur.SellCoins, "SKY") {
		t.Fatalf("sell coins = %v, want SKY present", cur.SellCoins)
	}

	// An unknown sell coin is rejected.
	rej, err := seller.Do(protocol.TypeCreateListing, protocol.CreateListingRequest{
		SellCoin: "NOPE", Amount: 10, Price: 2, PaymentCurrency: "BTC",
	})
	if err != nil {
		t.Fatalf("transport: %v", err)
	}
	if !rej.IsError() || errCode(t, rej) != protocol.CodeCurrencyUnavailable {
		t.Fatalf("unknown sell coin = %+v, want CURRENCY_UNAVAILABLE", rej)
	}

	// Add and enable a fibercoin with its own escrow wallet.
	if err := database.UpsertSellCoin(&db.SellCoin{
		Symbol: "MDL", Name: "Mobile", NodeURL: "http://mdl:6420",
		WalletSeed: "mdl-seed", WalletAddr: "MDL-escrow", Confirmations: 1, Enabled: true,
	}); err != nil {
		t.Fatalf("add MDL: %v", err)
	}

	var listing protocol.CreateListingResponse
	bindSuccess(t, seller, protocol.TypeCreateListing, protocol.CreateListingRequest{
		SellCoin: "MDL", Amount: 10, Price: 2, PaymentCurrency: "BTC",
	}, &listing)
	if listing.SellCoin != "MDL" || listing.MarketWallet != "MDL-escrow" {
		t.Fatalf("MDL listing routed wrong: %+v", listing)
	}

	// Disabling MDL makes further MDL listings unavailable.
	if err := database.SetSellCoinEnabled("MDL", false); err != nil {
		t.Fatal(err)
	}
	rej2, err := seller.Do(protocol.TypeCreateListing, protocol.CreateListingRequest{
		SellCoin: "MDL", Amount: 10, Price: 2, PaymentCurrency: "BTC",
	})
	if err != nil {
		t.Fatalf("transport: %v", err)
	}
	if !rej2.IsError() || errCode(t, rej2) != protocol.CodeCurrencyUnavailable {
		t.Fatalf("disabled sell coin = %+v, want CURRENCY_UNAVAILABLE", rej2)
	}
}

// TestFiberPayment verifies a fibercoin↔fibercoin trade: a SKY listing priced in
// another fibercoin (FIB) charges the buyer a FIXED amount (no unique non-round
// amount), and a coin cannot be traded for itself.
func TestFiberPayment(t *testing.T) {
	database := newTestDB(t)
	if err := database.SetConfig("wallet_sky", "sky-market-wallet"); err != nil {
		t.Fatal(err)
	}
	zeroCommission(t, database)
	// Add FIB as a second sell coin (enabling it as a fibercoin payment option).
	if err := database.UpsertSellCoin(&db.SellCoin{
		Symbol: "FIB", Name: "Fiber", NodeURL: "http://fib:6420",
		WalletSeed: "fib-seed", WalletAddr: "FIB-escrow", Confirmations: 1, Enabled: true,
	}); err != nil {
		t.Fatalf("add FIB: %v", err)
	}

	const sellerPK = "0311111111111111111111111111111111111111111111111111111111111111aa"
	const buyerPK = "0322222222222222222222222222222222222222222222222222222222222222bb"
	seller := dialTestServer(t, database, sellerPK)
	mustSuccess(t, seller, protocol.TypeRegister, protocol.RegisterRequest{WalletSKY: skySeller})

	// A coin can't be traded for itself.
	self, err := seller.Do(protocol.TypeCreateListing, protocol.CreateListingRequest{
		SellCoin: "SKY", Amount: 10, Price: 100, PaymentCurrency: "SKY",
	})
	if err != nil {
		t.Fatalf("transport: %v", err)
	}
	if !self.IsError() || errCode(t, self) != protocol.CodeInvalidRequest {
		t.Fatalf("SKY/SKY listing = %+v, want INVALID_REQUEST", self)
	}

	// List SKY priced in FIB.
	var listing protocol.CreateListingResponse
	bindSuccess(t, seller, protocol.TypeCreateListing, protocol.CreateListingRequest{
		SellCoin: "SKY", Amount: 10, Price: 100, PaymentCurrency: "FIB",
	}, &listing)
	if listing.SellCoin != "SKY" || listing.MarketWallet != "sky-market-wallet" || listing.ExpectedAmount != 10 {
		t.Fatalf("unexpected SKY/FIB listing: %+v", listing)
	}

	// Simulate the confirmed product (the Listing Checker would create it).
	if err := database.CreateProduct(&db.Product{
		ID: "prod-fib", SellerPubKey: sellerPK, SellCoin: "SKY", Amount: 10, Price: 100,
		PaymentCurrency: "FIB", Status: "active",
	}); err != nil {
		t.Fatalf("seed product: %v", err)
	}

	buyer := dialTestServer(t, database, buyerPK)
	mustSuccess(t, buyer, protocol.TypeRegister, protocol.RegisterRequest{WalletSKY: skyBuyer})

	var buy protocol.BuyProductResponse
	bindSuccess(t, buyer, protocol.TypeBuyProduct, protocol.BuyProductRequest{ProductID: "prod-fib"}, &buy)
	// FIXED amount — exactly the price, not a non-round 100.00x. Paid to the
	// seller's shared Skycoin-family address.
	if buy.PaymentCurrency != "FIB" || buy.ExpectedPaymentAmount != 100 || buy.SellerWallet != skySeller {
		t.Fatalf("fiber buy = %+v, want fixed 100 FIB to %s", buy, skySeller)
	}
}

// TestOnePendingListingPerSeller verifies a seller sends the exact round deposit
// amount and can only have one pending listing at a time — a second is rejected
// while the first is still awaiting its deposit. This keeps a single deposit from
// the seller's address unambiguous within one listing window.
func TestOnePendingListingPerSeller(t *testing.T) {
	database := newTestDB(t)
	if err := database.SetConfig("explorer_btc_provider", "esplora"); err != nil {
		t.Fatal(err)
	}
	if err := database.SetConfig("wallet_sky", "sky-market-wallet"); err != nil {
		t.Fatal(err)
	}
	zeroCommission(t, database)

	const sellerPK = "0311111111111111111111111111111111111111111111111111111111111111aa"
	seller := dialTestServer(t, database, sellerPK)
	mustSuccess(t, seller, protocol.TypeRegister, protocol.RegisterRequest{WalletSKY: skySeller, WalletBTC: btcSeller})

	var a protocol.CreateListingResponse
	req := protocol.CreateListingRequest{Amount: 10, Price: 2, PaymentCurrency: "BTC"}
	bindSuccess(t, seller, protocol.TypeCreateListing, req, &a)
	if a.ExpectedAmount != 10 {
		t.Fatalf("deposit amount = %v, want the exact round amount 10", a.ExpectedAmount)
	}

	// A second pending listing from the same seller is rejected.
	second, err := seller.Do(protocol.TypeCreateListing, req)
	if err != nil {
		t.Fatalf("second create transport error: %v", err)
	}
	if !second.IsError() || errCode(t, second) != protocol.CodeInvalidRequest {
		t.Fatalf("second listing = %+v, want INVALID_REQUEST (one pending listing per seller)", second)
	}
}

// TestRegisterRejectsBadWallet verifies registration validates payout address
// format: a malformed BTC address (or SKY address) is rejected up front.
func TestRegisterRejectsBadWallet(t *testing.T) {
	database := newTestDB(t)
	const pk = "0311111111111111111111111111111111111111111111111111111111111111aa"
	c := dialTestServer(t, database, pk)

	// Valid SKY but a malformed BTC address.
	badBTC, err := c.Do(protocol.TypeRegister, protocol.RegisterRequest{
		WalletSKY: skySeller, WalletBTC: "bc1-not-a-real-address",
	})
	if err != nil {
		t.Fatalf("transport error: %v", err)
	}
	if !badBTC.IsError() || errCode(t, badBTC) != protocol.CodeInvalidWallet {
		t.Fatalf("register(bad BTC) = %+v, want INVALID_WALLET", badBTC)
	}

	// Malformed SKY address.
	badSKY, err := c.Do(protocol.TypeRegister, protocol.RegisterRequest{WalletSKY: "not-a-sky-address"})
	if err != nil {
		t.Fatalf("transport error: %v", err)
	}
	if !badSKY.IsError() || errCode(t, badSKY) != protocol.CodeInvalidWallet {
		t.Fatalf("register(bad SKY) = %+v, want INVALID_WALLET", badSKY)
	}

	// A fully valid registration still succeeds.
	mustSuccess(t, c, protocol.TypeRegister, protocol.RegisterRequest{WalletSKY: skySeller, WalletBTC: btcSeller})
}

// TestTradeSizeBounds verifies the operator's min/max trade-size limits are
// enforced on create_listing.
func TestTradeSizeBounds(t *testing.T) {
	database := newTestDB(t)
	if err := database.SetConfig("explorer_btc_provider", "esplora"); err != nil {
		t.Fatal(err)
	}
	if err := database.SetConfig("wallet_sky", "sky-market-wallet"); err != nil {
		t.Fatal(err)
	}
	zeroCommission(t, database)
	if err := database.SetConfig("min_trade_sky", "5"); err != nil {
		t.Fatal(err)
	}
	if err := database.SetConfig("max_trade_sky", "100"); err != nil {
		t.Fatal(err)
	}
	const sellerPK = "0311111111111111111111111111111111111111111111111111111111111111aa"
	seller := dialTestServer(t, database, sellerPK)
	mustSuccess(t, seller, protocol.TypeRegister, protocol.RegisterRequest{WalletSKY: skySeller, WalletBTC: btcSeller})

	for _, amt := range []float64{4, 101} { // below min, above max
		resp, err := seller.Do(protocol.TypeCreateListing, protocol.CreateListingRequest{Amount: amt, Price: 2, PaymentCurrency: "BTC"})
		if err != nil {
			t.Fatalf("create_listing(%v) transport: %v", amt, err)
		}
		if !resp.IsError() || errCode(t, resp) != protocol.CodeInvalidRequest {
			t.Fatalf("create_listing(%v) = %+v, want INVALID_REQUEST (out of bounds)", amt, resp)
		}
	}

	// An in-bounds amount is accepted.
	var ok protocol.CreateListingResponse
	bindSuccess(t, seller, protocol.TypeCreateListing, protocol.CreateListingRequest{Amount: 50, Price: 2, PaymentCurrency: "BTC"}, &ok)
	if ok.ExpectedAmount != 50 {
		t.Fatalf("in-bounds listing deposit = %v, want 50", ok.ExpectedAmount)
	}
}

// TestListingCommission verifies the SKY commission is folded into the seller's
// expected deposit (amount + commission) at the default 0.5% rate and its floor.
func TestListingCommission(t *testing.T) {
	database := newTestDB(t) // default commission: 0.5%, min 0.001, no cap
	if err := database.SetConfig("explorer_btc_provider", "esplora"); err != nil {
		t.Fatal(err)
	}
	if err := database.SetConfig("wallet_sky", "sky-market-wallet"); err != nil {
		t.Fatal(err)
	}
	const sellerPK = "0311111111111111111111111111111111111111111111111111111111111111aa"
	seller := dialTestServer(t, database, sellerPK)
	mustSuccess(t, seller, protocol.TypeRegister, protocol.RegisterRequest{WalletSKY: skySeller, WalletBTC: btcSeller})

	// 0.5% of 100 = 0.5 → seller deposits 100.5, buyer will receive 100.
	var big protocol.CreateListingResponse
	bindSuccess(t, seller, protocol.TypeCreateListing, protocol.CreateListingRequest{
		Amount: 100, Price: 2, PaymentCurrency: "BTC",
	}, &big)
	if math.Abs(big.ExpectedAmount-100.5) > 1e-9 {
		t.Fatalf("deposit = %v, want amount+commission = 100.5", big.ExpectedAmount)
	}
}

// TestCreateListingPrecision verifies SKY amounts are normalized to the network's
// 3-decimal precision (so the seller can deposit the exact amount) and that out-
// of-range amounts are rejected.
func TestCreateListingPrecision(t *testing.T) {
	database := newTestDB(t)
	if err := database.SetConfig("explorer_btc_provider", "esplora"); err != nil {
		t.Fatal(err)
	}
	if err := database.SetConfig("wallet_sky", "sky-market-wallet"); err != nil {
		t.Fatal(err)
	}
	zeroCommission(t, database)
	const sellerPK = "0311111111111111111111111111111111111111111111111111111111111111aa"
	seller := dialTestServer(t, database, sellerPK)
	mustSuccess(t, seller, protocol.TypeRegister, protocol.RegisterRequest{WalletSKY: skySeller, WalletBTC: btcSeller})

	// A high-precision Amount is rounded to 3 decimals; the deposit is that
	// exact round amount (no non-round delta).
	var l protocol.CreateListingResponse
	bindSuccess(t, seller, protocol.TypeCreateListing, protocol.CreateListingRequest{
		Amount: 10.123456789, Price: 2, PaymentCurrency: "BTC",
	}, &l)
	if math.Abs(l.ExpectedAmount-10.123) > 1e-9 {
		t.Fatalf("deposit %v, want normalized to 3 decimals (10.123)", l.ExpectedAmount)
	}

	// An amount above the safe float range is rejected.
	tooBig, err := seller.Do(protocol.TypeCreateListing, protocol.CreateListingRequest{
		Amount: 1e10, Price: 2, PaymentCurrency: "BTC",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !tooBig.IsError() || errCode(t, tooBig) != protocol.CodeInvalidRequest {
		t.Fatalf("too-large listing = %+v, want INVALID_REQUEST", tooBig)
	}

	// An amount that rounds to zero at SKY precision is rejected.
	tooSmall, err := seller.Do(protocol.TypeCreateListing, protocol.CreateListingRequest{
		Amount: 1e-9, Price: 2, PaymentCurrency: "BTC",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !tooSmall.IsError() || errCode(t, tooSmall) != protocol.CodeInvalidRequest {
		t.Fatalf("too-small listing = %+v, want INVALID_REQUEST", tooSmall)
	}
}

// TestUnknownMessageType verifies the dispatcher rejects unknown types cleanly.
// TestGetListings verifies a seller sees their own pending listings (with the
// deposit amount and market wallet) and their live status, and that another
// caller does not see them.
func TestGetListings(t *testing.T) {
	database := newTestDB(t)

	const sellerPK = "0311111111111111111111111111111111111111111111111111111111111111aa"
	const otherPK = "0322222222222222222222222222222222222222222222222222222222222222bb"

	if err := database.SetConfig("explorer_btc_provider", "esplora"); err != nil {
		t.Fatal(err)
	}
	if err := database.SetConfig("wallet_sky", "sky-market-wallet"); err != nil {
		t.Fatal(err)
	}
	zeroCommission(t, database)

	seller := dialTestServer(t, database, sellerPK)
	mustSuccess(t, seller, protocol.TypeRegister, protocol.RegisterRequest{
		WalletSKY: skySeller, WalletBTC: btcSeller,
	})

	var listing protocol.CreateListingResponse
	bindSuccess(t, seller, protocol.TypeCreateListing, protocol.CreateListingRequest{
		Amount: 10, Price: 2, PaymentCurrency: "BTC",
	}, &listing)

	// The seller sees the pending listing with the deposit amount + wallet.
	var mine protocol.GetListingsResponse
	bindSuccess(t, seller, protocol.TypeGetListings, nil, &mine)
	if len(mine.Listings) != 1 {
		t.Fatalf("get_listings = %+v, want one pending listing", mine.Listings)
	}
	l := mine.Listings[0]
	if l.Status != "pending" || l.ExpectedAmount != listing.ExpectedAmount || mine.MarketWallet != "sky-market-wallet" {
		t.Fatalf("listing view = %+v (wallet %q), want pending with the deposit amount", l, mine.MarketWallet)
	}

	// Once the deposit is confirmed, the listing shows as confirmed.
	if err := database.UpdatePendingListingStatus(l.ID, "confirmed", "deposit-tx"); err != nil {
		t.Fatal(err)
	}
	bindSuccess(t, seller, protocol.TypeGetListings, nil, &mine)
	if len(mine.Listings) != 1 || mine.Listings[0].Status != "confirmed" || mine.Listings[0].TxHash != "deposit-tx" {
		t.Fatalf("after confirm, get_listings = %+v, want one confirmed listing with tx", mine.Listings)
	}

	// A different caller does not see the seller's listings.
	other := dialTestServer(t, database, otherPK)
	mustSuccess(t, other, protocol.TypeRegister, protocol.RegisterRequest{WalletSKY: skyOther})
	var theirs protocol.GetListingsResponse
	bindSuccess(t, other, protocol.TypeGetListings, nil, &theirs)
	if len(theirs.Listings) != 0 {
		t.Fatalf("other caller get_listings = %+v, want none", theirs.Listings)
	}
}

// TestBuyerCancelBlocksRebuy verifies a buyer can cancel their in-flight buy
// (releasing the product) but is then blocked from buying that same product
// again, while another buyer still can.
func TestBuyerCancelBlocksRebuy(t *testing.T) {
	database := newTestDB(t)
	if err := database.SetConfig("explorer_btc_provider", "esplora"); err != nil {
		t.Fatal(err)
	}

	const sellerPK = "0311111111111111111111111111111111111111111111111111111111111111aa"
	const buyerPK = "0322222222222222222222222222222222222222222222222222222222222222bb"
	const buyer2PK = "0333333333333333333333333333333333333333333333333333333333333333cc"

	seller := dialTestServer(t, database, sellerPK)
	mustSuccess(t, seller, protocol.TypeRegister, protocol.RegisterRequest{WalletSKY: skySeller, WalletBTC: btcSeller})
	if err := database.CreateProduct(&db.Product{
		ID: "prod-1", SellerPubKey: sellerPK, Amount: 10, Price: 2, PaymentCurrency: "BTC", Status: "active",
	}); err != nil {
		t.Fatal(err)
	}

	buyer := dialTestServer(t, database, buyerPK)
	mustSuccess(t, buyer, protocol.TypeRegister, protocol.RegisterRequest{WalletSKY: skyBuyer, WalletBTC: btcBuyer})

	var buy protocol.BuyProductResponse
	bindSuccess(t, buyer, protocol.TypeBuyProduct, protocol.BuyProductRequest{ProductID: "prod-1"}, &buy)

	// Cancel the in-flight order; the product returns to active.
	mustSuccess(t, buyer, protocol.TypeCancelOrder, protocol.CancelOrderRequest{OrderID: buy.OrderID})
	if p, _ := database.GetProduct("prod-1"); p == nil || p.Status != "active" { //nolint
		t.Fatalf("product after cancel = %+v, want active", p)
	}
	// A voluntary cancel counts as a freeze violation (toward the ban system).
	if n, _ := database.CountRecentViolations(buyerPK, time.Now().UTC().Add(-time.Hour)); n != 1 { //nolint
		t.Fatalf("violations after cancel = %d, want 1", n)
	}

	// With the default buy-cancel limit of 2, one cancel does not block the buyer:
	// they can still buy the released product a second time.
	var buyAgain protocol.BuyProductResponse
	bindSuccess(t, buyer, protocol.TypeBuyProduct, protocol.BuyProductRequest{ProductID: "prod-1"}, &buyAgain)
	// Cancel again — this reaches the limit (2 cancels for this product).
	mustSuccess(t, buyer, protocol.TypeCancelOrder, protocol.CancelOrderRequest{OrderID: buyAgain.OrderID})
	if p, _ := database.GetProduct("prod-1"); p == nil || p.Status != "active" { //nolint
		t.Fatalf("product after second cancel = %+v, want active", p)
	}

	// Now the same buyer is blocked from buying it again.
	reb, err := buyer.Do(protocol.TypeBuyProduct, protocol.BuyProductRequest{ProductID: "prod-1"})
	if err != nil {
		t.Fatalf("rebuy transport error: %v", err)
	}
	if !reb.IsError() || errCode(t, reb) != protocol.CodeBuyerBlocked {
		t.Fatalf("rebuy = %+v, want BUYER_BLOCKED", reb)
	}

	// The operator clears the block; the buyer can buy the product again.
	if _, err := database.ClearBuyerProductBlock(buyerPK, "prod-1"); err != nil {
		t.Fatalf("clear block: %v", err)
	}
	var buyAfterClear protocol.BuyProductResponse
	bindSuccess(t, buyer, protocol.TypeBuyProduct, protocol.BuyProductRequest{ProductID: "prod-1"}, &buyAfterClear)
	if buyAfterClear.OrderID == "" {
		t.Fatal("buyer should be able to buy again after the operator clears the block")
	}
	// Reset for the "different buyer" check below: cancel and clear so the product
	// is active and unblocked again.
	mustSuccess(t, buyer, protocol.TypeCancelOrder, protocol.CancelOrderRequest{OrderID: buyAfterClear.OrderID})
	if _, err := database.ClearBuyerProductBlock(buyerPK, "prod-1"); err != nil {
		t.Fatalf("clear block (reset): %v", err)
	}

	// A different buyer still can.
	buyer2 := dialTestServer(t, database, buyer2PK)
	mustSuccess(t, buyer2, protocol.TypeRegister, protocol.RegisterRequest{WalletSKY: skyBuyer2, WalletBTC: btcBuyer2})
	var buy2 protocol.BuyProductResponse
	bindSuccess(t, buyer2, protocol.TypeBuyProduct, protocol.BuyProductRequest{ProductID: "prod-1"}, &buy2)
	if buy2.OrderID == "" {
		t.Fatal("second buyer should be able to buy the released product")
	}
}

// TestSellerCancelConfirmedOffer verifies a seller can cancel a confirmed offer
// (deactivating the product and marking the listing canceled for refund), but
// not once a buyer has frozen it.
func TestSellerCancelConfirmedOffer(t *testing.T) {
	database := newTestDB(t)
	if err := database.SetConfig("explorer_btc_provider", "esplora"); err != nil {
		t.Fatal(err)
	}
	if err := database.SetConfig("wallet_sky", "sky-market-wallet"); err != nil {
		t.Fatal(err)
	}
	zeroCommission(t, database)

	const sellerPK = "0311111111111111111111111111111111111111111111111111111111111111aa"
	const buyerPK = "0322222222222222222222222222222222222222222222222222222222222222bb"

	seller := dialTestServer(t, database, sellerPK)
	mustSuccess(t, seller, protocol.TypeRegister, protocol.RegisterRequest{WalletSKY: skySeller, WalletBTC: btcSeller})

	var listing protocol.CreateListingResponse
	bindSuccess(t, seller, protocol.TypeCreateListing, protocol.CreateListingRequest{
		Amount: 10, Price: 2, PaymentCurrency: "BTC",
	}, &listing)

	// Simulate the Listing Checker promoting the confirmed deposit into a product.
	if err := database.UpdatePendingListingStatus(listing.ListingID, "confirmed", "dep-tx"); err != nil {
		t.Fatal(err)
	}
	if err := database.CreateProduct(&db.Product{
		ID: "prod-1", ListingID: listing.ListingID, SellerPubKey: sellerPK,
		Amount: 10, Price: 2, PaymentCurrency: "BTC", Status: "active",
	}); err != nil {
		t.Fatal(err)
	}

	// Cancel the confirmed offer: product deactivated, listing canceled.
	mustSuccess(t, seller, protocol.TypeCancelListing, protocol.CancelListingRequest{ListingID: listing.ListingID})
	if p, _ := database.GetProduct("prod-1"); p == nil || p.Status != "cancelled" { //nolint
		t.Fatalf("product after seller cancel = %+v, want canceled", p)
	}
	if l, _ := database.GetPendingListing(listing.ListingID); l == nil || l.Status != "canceled" { //nolint
		t.Fatalf("listing after seller cancel = %+v, want canceled", l)
	}

	// A second, frozen offer cannot be canceled.
	var listing2 protocol.CreateListingResponse
	bindSuccess(t, seller, protocol.TypeCreateListing, protocol.CreateListingRequest{
		Amount: 5, Price: 1, PaymentCurrency: "BTC",
	}, &listing2)
	if err := database.UpdatePendingListingStatus(listing2.ListingID, "confirmed", "dep-tx-2"); err != nil {
		t.Fatal(err)
	}
	if err := database.CreateProduct(&db.Product{
		ID: "prod-2", ListingID: listing2.ListingID, SellerPubKey: sellerPK,
		Amount: 5, Price: 1, PaymentCurrency: "BTC", Status: "active",
	}); err != nil {
		t.Fatal(err)
	}
	if err := database.FreezeProduct("prod-2", buyerPK); err != nil {
		t.Fatal(err)
	}
	resp, err := seller.Do(protocol.TypeCancelListing, protocol.CancelListingRequest{ListingID: listing2.ListingID})
	if err != nil {
		t.Fatalf("cancel frozen transport error: %v", err)
	}
	if !resp.IsError() || errCode(t, resp) != protocol.CodeProductUnavailable {
		t.Fatalf("cancel of frozen offer = %+v, want PRODUCT_UNAVAILABLE", resp)
	}
}

func TestUnknownMessageType(t *testing.T) {
	database := newTestDB(t)
	c := dialTestServer(t, database, "03deadbeef")
	resp, err := c.Do("client.nonsense", nil)
	if err != nil {
		t.Fatalf("transport error: %v", err)
	}
	if !resp.IsError() || errCode(t, resp) != protocol.CodeInvalidRequest {
		t.Fatalf("resp = %+v, want INVALID_REQUEST", resp)
	}
}

// --- helpers ---

func mustSuccess(t *testing.T, c *market.Conn, msgType string, data any) protocol.Envelope {
	t.Helper()
	resp, err := c.Do(msgType, data)
	if err != nil {
		t.Fatalf("%s transport error: %v", msgType, err)
	}
	if resp.IsError() {
		t.Fatalf("%s failed: %s", msgType, errCode(t, resp))
	}
	return resp
}

func bindSuccess(t *testing.T, c *market.Conn, msgType string, data, out any) {
	t.Helper()
	resp := mustSuccess(t, c, msgType, data)
	if err := resp.Bind(out); err != nil {
		t.Fatalf("%s bind response: %v", msgType, err)
	}
}

func errCode(t *testing.T, env protocol.Envelope) string {
	t.Helper()
	var e protocol.ErrorData
	if err := env.Bind(&e); err != nil {
		t.Fatalf("bind error data: %v", err)
	}
	return e.Code
}
