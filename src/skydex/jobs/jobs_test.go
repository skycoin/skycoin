package jobs_test

import (
	"math"
	"path/filepath"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/skydex/db"
	"github.com/skycoin/skycoin/src/skydex/jobs"
)

// fakeChain is a controllable Chain for tests.
type fakeChain struct {
	deposit bool
	confs   int
	balance float64
	sends   []sendCall
}

type sendCall struct {
	addr string
	amt  float64
}

func (f *fakeChain) DepositConfirmed(string, string, string, float64, time.Time, time.Time) (bool, string, error) {
	return f.deposit, "deposit-tx", nil
}
func (f *fakeChain) PaymentConfirmations(string, string, string, float64, time.Time, time.Time) (int, string, error) {
	return f.confs, "pay-tx", nil
}
func (f *fakeChain) SendCoin(_ /*coin*/, addr string, amt float64) (string, error) {
	f.sends = append(f.sends, sendCall{addr, amt})
	return "send-tx", nil
}
func (f *fakeChain) EscrowBalance(string, string) (float64, error) {
	return f.balance, nil
}

func newDB(t *testing.T) *db.Database {
	t.Helper()
	database, err := db.New(filepath.Join(t.TempDir(), "m.db"), "")
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
	return database
}

func mustUser(t *testing.T, d *db.Database, pk, skyWallet string) {
	t.Helper()
	if err := d.CreateUser(&db.User{PubKey: pk, WalletSKY: skyWallet}); err != nil {
		t.Fatalf("create user: %v", err)
	}
}

// TestExpiry expires an overdue listing and an overdue order, releasing the
// product and recording a freeze violation.
func TestExpiry(t *testing.T) {
	d := newDB(t)
	r := jobs.NewRunner(d, &fakeChain{}, nil)

	const seller, buyer = "03seller", "03buyer"
	mustUser(t, d, seller, "sky-seller")
	mustUser(t, d, buyer, "sky-buyer")

	// Overdue pending listing.
	if err := d.CreatePendingListing(&db.PendingListing{
		ID: "l1", SellerPubKey: seller, Amount: 5, ExpectedAmount: 5.01, Price: 1,
		PaymentCurrency: "BTC", Status: "pending", ExpiresAt: time.Now().UTC().Add(-time.Minute),
	}); err != nil {
		t.Fatal(err)
	}

	// Active product, frozen by the buyer, with an overdue order.
	if err := d.CreateProduct(&db.Product{
		ID: "p1", SellerPubKey: seller, Amount: 5, Price: 1, PaymentCurrency: "BTC", Status: "active",
	}); err != nil {
		t.Fatal(err)
	}
	if err := d.FreezeProduct("p1", buyer); err != nil {
		t.Fatal(err)
	}
	if err := d.CreateOrder(&db.Order{
		ID: "o1", ProductID: "p1", BuyerPubKey: buyer, Amount: 5, Price: 1,
		PaymentCurrency: "BTC", ExpectedPaymentAmount: 1.01, SellerWallet: "bc1-seller",
		Status: "pending_payment", ExpiresAt: time.Now().UTC().Add(-time.Minute),
	}); err != nil {
		t.Fatal(err)
	}

	if err := r.RunExpiry(); err != nil {
		t.Fatalf("RunExpiry: %v", err)
	}

	listing, _ := d.GetPendingListing("l1") //nolint
	if listing.Status != "expired" {
		t.Fatalf("listing status = %q, want expired", listing.Status)
	}
	order, _ := d.GetOrder("o1") //nolint
	if order.Status != "expired" {
		t.Fatalf("order status = %q, want expired", order.Status)
	}
	product, _ := d.GetProduct("p1") //nolint
	if product.Status != "active" {
		t.Fatalf("product status = %q, want active (released)", product.Status)
	}
	n, _ := d.CountRecentViolations(buyer, time.Now().UTC().Add(-time.Hour)) //nolint
	if n != 1 {
		t.Fatalf("violations = %d, want 1", n)
	}
}

// TestListingCheck promotes a pending listing whose deposit is confirmed.
func TestListingCheck(t *testing.T) {
	d := newDB(t)
	if err := d.SetConfig("wallet_sky", "market-wallet"); err != nil {
		t.Fatal(err)
	}
	mustUser(t, d, "03seller", "sky-seller")
	if err := d.CreatePendingListing(&db.PendingListing{
		ID: "l1", SellerPubKey: "03seller", Amount: 10, ExpectedAmount: 10.02, Price: 2,
		PaymentCurrency: "BTC", Status: "pending", ExpiresAt: time.Now().UTC().Add(time.Hour),
	}); err != nil {
		t.Fatal(err)
	}

	r := jobs.NewRunner(d, &fakeChain{deposit: true}, nil)
	if err := r.RunListingCheck(); err != nil {
		t.Fatalf("RunListingCheck: %v", err)
	}

	listing, _ := d.GetPendingListing("l1") //nolint
	if listing.Status != "confirmed" {
		t.Fatalf("listing status = %q, want confirmed", listing.Status)
	}
	products, _ := d.GetActiveProducts() //nolint
	if len(products) != 1 || products[0].Amount != 10 {
		t.Fatalf("expected one active product from the listing, got %+v", products)
	}
}

// TestEscrowCheck completes an order once its payment reaches the threshold and
// delivers SKY to the buyer.
func TestEscrowCheck(t *testing.T) {
	d := newDB(t)
	mustUser(t, d, "03seller", "sky-seller")
	mustUser(t, d, "03buyer", "sky-buyer")
	if err := d.CreateProduct(&db.Product{
		ID: "p1", SellerPubKey: "03seller", Amount: 7, Price: 2, PaymentCurrency: "BTC", Status: "frozen",
	}); err != nil {
		t.Fatal(err)
	}
	if err := d.CreateOrder(&db.Order{
		ID: "o1", ProductID: "p1", BuyerPubKey: "03buyer", Amount: 7, Price: 2,
		PaymentCurrency: "BTC", ExpectedPaymentAmount: 2.01, SellerWallet: "bc1-seller",
		Status: "pending_payment", ExpiresAt: time.Now().UTC().Add(time.Hour),
	}); err != nil {
		t.Fatal(err)
	}

	chain := &fakeChain{confs: 2}
	r := jobs.NewRunner(d, chain, nil)
	if err := r.RunEscrowCheck(); err != nil {
		t.Fatalf("RunEscrowCheck: %v", err)
	}

	order, _ := d.GetOrder("o1") //nolint
	if order.Status != "completed" {
		t.Fatalf("order status = %q, want completed", order.Status)
	}
	product, _ := d.GetProduct("p1") //nolint
	if product.Status != "sold" {
		t.Fatalf("product status = %q, want sold", product.Status)
	}
	if len(chain.sends) != 1 || chain.sends[0].addr != "sky-buyer" || chain.sends[0].amt != 7 {
		t.Fatalf("SendSKY calls = %+v, want one send of 7 to sky-buyer", chain.sends)
	}
	// Commission booked at the default 0.5%: 0.5% of 7 = 0.035 SKY.
	all, _ := d.GetAllOrders() //nolint
	if len(all) != 1 || math.Abs(all[0].Commission-0.035) > 1e-9 {
		t.Fatalf("commission = %v, want 0.035 SKY on the completed order", all)
	}
}

// TestCommissionSKY covers the SKY commission: rate% of the amount, rounded up
// to 3 decimals, clamped to [min, cap] (cap 0 = none).
func TestCommissionSKY(t *testing.T) {
	cases := []struct {
		amount, rate, min, cap, want float64
	}{
		{100, 0.5, 0.001, 0, 0.5},     // 0.5% of 100
		{10, 0.5, 0.001, 0, 0.05},     // 0.5% of 10
		{1, 0.5, 0.001, 0, 0.005},     // 0.5% of 1
		{0.1, 0.5, 0.001, 0, 0.001},   // 0.5% = 0.0005 -> floor 0.001
		{100000, 0.5, 0.001, 50, 50},  // 0.5% = 500 -> capped at 50
		{33.33, 0.5, 0.001, 0, 0.167}, // 0.166650 -> rounds up to 0.167
		{10, 0, 0, 0, 0},              // rate 0 -> nothing (no floor)
		{0, 0.5, 0.001, 0, 0},         // no amount -> 0 (not floored)
		{-5, 0.5, 0.001, 0, 0},        // guard against negatives
	}
	for _, c := range cases {
		if got := jobs.Commission(c.amount, c.rate, c.min, c.cap); math.Abs(got-c.want) > 1e-9 {
			t.Errorf("Commission(%v, %v, %v, %v) = %v, want %v", c.amount, c.rate, c.min, c.cap, got, c.want)
		}
	}
}

// TestEscrowCheckNoopDefers verifies that without a chain backend the order is
// confirmed but not completed (delivery deferred).
func TestEscrowCheckNoopDefers(t *testing.T) {
	d := newDB(t)
	mustUser(t, d, "03seller", "sky-seller")
	mustUser(t, d, "03buyer", "sky-buyer")
	if err := d.CreateProduct(&db.Product{
		ID: "p1", SellerPubKey: "03seller", Amount: 7, Price: 2, PaymentCurrency: "BTC", Status: "frozen",
	}); err != nil {
		t.Fatal(err)
	}
	if err := d.CreateOrder(&db.Order{
		ID: "o1", ProductID: "p1", BuyerPubKey: "03buyer", Amount: 7, Price: 2,
		PaymentCurrency: "BTC", ExpectedPaymentAmount: 2.01, SellerWallet: "bc1-seller",
		Status: "pending_payment", ExpiresAt: time.Now().UTC().Add(time.Hour),
	}); err != nil {
		t.Fatal(err)
	}

	// PaymentConfirmations reaches threshold, but NoopChain.SendSKY fails.
	r := jobs.NewRunner(d, noopButConfirmed{}, nil)
	if err := r.RunEscrowCheck(); err != nil {
		t.Fatalf("RunEscrowCheck: %v", err)
	}
	order, _ := d.GetOrder("o1") //nolint
	if order.Status != "confirmed" {
		t.Fatalf("order status = %q, want confirmed (delivery deferred)", order.Status)
	}
}

// noopButConfirmed confirms payment but cannot send (like a misconfigured chain).
type noopButConfirmed struct{ jobs.NoopChain }

func (noopButConfirmed) PaymentConfirmations(string, string, string, float64, time.Time, time.Time) (int, string, error) {
	return 2, "pay-tx", nil
}

// TestReturnScheduler refunds escrowed SKY for a canceled listing whose deposit
// landed — immediately (no delay) — and records the refund exactly once.
func TestReturnScheduler(t *testing.T) {
	d := newDB(t)
	if err := d.SetConfig("wallet_sky", "market-wallet"); err != nil {
		t.Fatal(err)
	}
	mustUser(t, d, "03seller", "sky-seller")

	// A listing the seller deposited for, then canceled.
	if err := d.CreatePendingListing(&db.PendingListing{
		ID: "l1", SellerPubKey: "03seller", Amount: 5, ExpectedAmount: 5.013, Price: 1,
		PaymentCurrency: "BTC", Status: "pending", ExpiresAt: time.Now().UTC().Add(time.Hour),
	}); err != nil {
		t.Fatal(err)
	}
	if err := d.UpdatePendingListingStatus("l1", "canceled", ""); err != nil {
		t.Fatal(err)
	}

	chain := &fakeChain{deposit: true}
	r := jobs.NewRunner(d, chain, nil)

	// The refund is due immediately once the listing is canceled.
	if err := r.RunReturnScheduler(); err != nil {
		t.Fatalf("RunReturnScheduler (due): %v", err)
	}
	if len(chain.sends) != 1 || chain.sends[0].addr != "sky-seller" || chain.sends[0].amt != 5.013 {
		t.Fatalf("SendSKY calls = %+v, want one refund of 5.013 to sky-seller", chain.sends)
	}
	listing, _ := d.GetPendingListing("l1") //nolint
	if listing.Status != "returned" {
		t.Fatalf("listing status = %q, want returned", listing.Status)
	}
	if listing.ReturnTxHash != "send-tx" {
		t.Fatalf("return_tx_hash = %q, want send-tx", listing.ReturnTxHash)
	}

	// Idempotent: a second pass must not refund again.
	if err := r.RunReturnScheduler(); err != nil {
		t.Fatalf("RunReturnScheduler (second pass): %v", err)
	}
	if len(chain.sends) != 1 {
		t.Fatalf("refund must happen at most once, got %+v", chain.sends)
	}
}

// TestReturnSchedulerNoDeposit verifies escrow is NOT refunded when no matching
// deposit is found on-chain (guards against draining the wallet for a listing
// the seller never funded).
func TestReturnSchedulerNoDeposit(t *testing.T) {
	d := newDB(t)
	if err := d.SetConfig("wallet_sky", "market-wallet"); err != nil {
		t.Fatal(err)
	}
	mustUser(t, d, "03seller", "sky-seller")
	if err := d.CreatePendingListing(&db.PendingListing{
		ID: "l1", SellerPubKey: "03seller", Amount: 5, ExpectedAmount: 5.013, Price: 1,
		PaymentCurrency: "BTC", Status: "pending", ExpiresAt: time.Now().UTC().Add(time.Hour),
	}); err != nil {
		t.Fatal(err)
	}
	if err := d.UpdatePendingListingStatus("l1", "expired", ""); err != nil {
		t.Fatal(err)
	}

	// deposit:false => DepositConfirmed reports the SKY never arrived.
	chain := &fakeChain{deposit: false}
	r := jobs.NewRunner(d, chain, nil)
	if err := r.RunReturnScheduler(); err != nil {
		t.Fatalf("RunReturnScheduler: %v", err)
	}
	if len(chain.sends) != 0 {
		t.Fatalf("must not refund without an on-chain deposit, got %+v", chain.sends)
	}
	listing, _ := d.GetPendingListing("l1") //nolint
	if listing.Status != "expired" {
		t.Fatalf("listing status = %q, want expired (unchanged)", listing.Status)
	}
}

// TestBanManager bans a user over the violation limit and lifts an expired ban.
func TestBanManager(t *testing.T) {
	d := newDB(t)
	mustUser(t, d, "03seller", "sky")
	mustUser(t, d, "03offender", "sky")

	// A product + order to satisfy the freeze_violations foreign keys.
	if err := d.CreateProduct(&db.Product{
		ID: "p1", SellerPubKey: "03seller", Amount: 1, Price: 1, PaymentCurrency: "BTC", Status: "active",
	}); err != nil {
		t.Fatal(err)
	}
	if err := d.CreateOrder(&db.Order{
		ID: "o1", ProductID: "p1", BuyerPubKey: "03offender", Amount: 1, Price: 1,
		PaymentCurrency: "BTC", ExpectedPaymentAmount: 1.01, SellerWallet: "w",
		Status: "expired", ExpiresAt: time.Now().UTC(),
	}); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		if err := d.CreateFreezeViolation("03offender", "o1"); err != nil {
			t.Fatal(err)
		}
	}

	// An already-expired ban that should be lifted.
	mustUser(t, d, "03old", "sky")
	if err := d.CreateBan("03old", 3, time.Now().UTC().Add(-time.Hour)); err != nil {
		t.Fatal(err)
	}

	r := jobs.NewRunner(d, &fakeChain{}, nil)
	if err := r.RunBanManager(); err != nil {
		t.Fatalf("RunBanManager: %v", err)
	}

	banned, _ := d.IsUserBanned("03offender") //nolint
	if !banned {
		t.Fatal("offender should be banned after 3 violations")
	}
	stillBanned, _ := d.IsUserBanned("03old") //nolint
	if stillBanned {
		t.Fatal("expired ban should have been lifted")
	}
}

// TestCleanup deletes completed orders past the retention window.
func TestCleanup(t *testing.T) {
	d := newDB(t)
	if err := d.SetConfig("cleanup_days", "0"); err != nil { // retention 0 => delete once completed
		t.Fatal(err)
	}
	mustUser(t, d, "03seller", "sky")
	mustUser(t, d, "03buyer", "sky")
	if err := d.CreateProduct(&db.Product{
		ID: "p1", SellerPubKey: "03seller", Amount: 1, Price: 1, PaymentCurrency: "BTC", Status: "sold",
	}); err != nil {
		t.Fatal(err)
	}
	if err := d.CreateOrder(&db.Order{
		ID: "o1", ProductID: "p1", BuyerPubKey: "03buyer", Amount: 1, Price: 1,
		PaymentCurrency: "BTC", ExpectedPaymentAmount: 1.01, SellerWallet: "w",
		Status: "completed", ExpiresAt: time.Now().UTC(),
	}); err != nil {
		t.Fatal(err)
	}
	if err := d.MarkOrderCompleted("o1"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Millisecond) // ensure completed_at is strictly in the past

	r := jobs.NewRunner(d, &fakeChain{}, nil)
	if err := r.RunCleanup(); err != nil {
		t.Fatalf("RunCleanup: %v", err)
	}
	if o, _ := d.GetOrder("o1"); o != nil { //nolint
		t.Fatalf("completed order should have been cleaned up, got %+v", o)
	}
}

// TestEscrowAudit sums the market's outstanding SKY obligations (products
// awaiting delivery + confirmed deposits awaiting refund) and confirms the audit
// runs cleanly whether the on-chain balance covers them or falls short.
func TestEscrowAudit(t *testing.T) {
	d := newDB(t)
	if err := d.SetConfig("wallet_sky", "market-wallet"); err != nil {
		t.Fatal(err)
	}
	mustUser(t, d, "03seller", "sky-seller")

	// Active (10) + frozen (5) products are SKY awaiting delivery.
	if err := d.CreateProduct(&db.Product{ID: "p-active", SellerPubKey: "03seller", Amount: 10, Price: 1, PaymentCurrency: "BTC", Status: "active"}); err != nil {
		t.Fatal(err)
	}
	if err := d.CreateProduct(&db.Product{ID: "p-frozen", SellerPubKey: "03seller", Amount: 5, Price: 1, PaymentCurrency: "BTC", Status: "active"}); err != nil {
		t.Fatal(err)
	}
	if err := d.FreezeProduct("p-frozen", "03buyer"); err != nil {
		t.Fatal(err)
	}
	// A sold product (7) has already been delivered — excluded.
	if err := d.CreateProduct(&db.Product{ID: "p-sold", SellerPubKey: "03seller", Amount: 7, Price: 1, PaymentCurrency: "BTC", Status: "sold"}); err != nil {
		t.Fatal(err)
	}

	// A canceled listing whose deposit arrived (confirmed) but isn't refunded: +3.
	if err := d.CreatePendingListing(&db.PendingListing{ID: "l-refund", SellerPubKey: "03seller", Amount: 3, ExpectedAmount: 3, Price: 1, PaymentCurrency: "BTC", Status: "pending", ExpiresAt: time.Now().UTC().Add(time.Hour)}); err != nil {
		t.Fatal(err)
	}
	if err := d.UpdatePendingListingStatus("l-refund", "confirmed", "dep-tx"); err != nil {
		t.Fatal(err)
	}
	if err := d.UpdatePendingListingStatus("l-refund", "canceled", "dep-tx"); err != nil {
		t.Fatal(err)
	}
	// An expired listing whose deposit never arrived (confirmed_at NULL) — excluded.
	if err := d.CreatePendingListing(&db.PendingListing{ID: "l-nodep", SellerPubKey: "03seller", Amount: 99, ExpectedAmount: 99, Price: 1, PaymentCurrency: "BTC", Status: "pending", ExpiresAt: time.Now().UTC().Add(-time.Minute)}); err != nil {
		t.Fatal(err)
	}
	if err := d.UpdatePendingListingStatus("l-nodep", "expired", ""); err != nil {
		t.Fatal(err)
	}

	// Expected outstanding obligations: 10 + 5 + 3 = 18.
	got, err := d.OutstandingEscrow("SKY")
	if err != nil {
		t.Fatalf("OutstandingEscrow: %v", err)
	}
	if math.Abs(got-18) > 1e-9 {
		t.Fatalf("OutstandingEscrow = %v, want 18", got)
	}

	// Healthy (balance covers obligations), shortfall (balance below), and no
	// backend all complete without error.
	for _, tc := range []struct {
		name  string
		chain jobs.Chain
	}{
		{"healthy", &fakeChain{balance: 20}},
		{"shortfall", &fakeChain{balance: 1}},
		{"no-backend", jobs.NoopChain{}},
	} {
		if err := jobs.NewRunner(d, tc.chain, nil).RunEscrowAudit(); err != nil {
			t.Fatalf("RunEscrowAudit(%s): %v", tc.name, err)
		}
	}
}
