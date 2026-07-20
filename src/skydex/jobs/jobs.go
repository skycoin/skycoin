package jobs

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/skydex/db"
)

// RequiredConfirmations is the number of blockchain confirmations a payment must
// reach before a trade completes (skydex-design.md §3).
const RequiredConfirmations = 2

// depositConfirmGrace widens a listing's deposit window on the upper end so a
// deposit the seller sent just before expiry, but which only got mined into a
// block shortly after, is still credited. The lower bound (listing creation) is
// exact — it blocks replay of an older, unrelated transaction.
const depositConfirmGrace = 10 * time.Minute

// paymentConfirmGrace is the equivalent for buyer payments. BTC/LTC block and
// confirmation times are far slower and more variable than SKY's, so the upper
// bound is generous; the lower bound (order creation) is what actually blocks
// replay of an old payment of a recycled amount.
const paymentConfirmGrace = 2 * time.Hour

// Commission returns the market's SKY commission for a sale of amountSKY:
// ratePercent% of the amount, rounded UP to the network's 3-decimal precision,
// then clamped to [minSKY, maxSKY] (maxSKY <= 0 means no cap). The seller pays it
// by depositing amountSKY + commission into escrow; the market delivers amountSKY
// to the buyer and retains the commission.
func Commission(amountSKY, ratePercent, minSKY, maxSKY float64) float64 {
	if amountSKY <= 0 {
		return 0
	}
	// Percentage of the amount, rounded up to 3 decimals so we never undercharge
	// and the result is a valid on-chain amount.
	c := math.Ceil(amountSKY*ratePercent/100*1000) / 1000
	if c < minSKY {
		c = minSKY
	}
	if maxSKY > 0 && c > maxSKY {
		c = maxSKY
	}
	// The floor/cap come from operator config; snap to 3 decimals to stay valid.
	return math.Round(c*1000) / 1000
}

// Runner owns and schedules the market's background jobs.
type Runner struct {
	db    *db.Database
	chain Chain
	log   logrus.FieldLogger
	confs int
}

// NewRunner builds a Runner. A nil chain defaults to NoopChain (no on-chain
// action), a nil log to a default logger.
func NewRunner(database *db.Database, chain Chain, log logrus.FieldLogger) *Runner {
	if chain == nil {
		chain = NoopChain{}
	}
	if log == nil {
		log = logrus.New()
	}
	return &Runner{db: database, chain: chain, log: log, confs: RequiredConfirmations}
}

// requiredConfirmations returns the operator-configured confirmation threshold,
// falling back to the built-in default when unset or invalid. Read per tick so a
// change takes effect without a restart.
func (r *Runner) requiredConfirmations() int {
	if c, err := r.db.GetConfirmationsRequired(); err == nil && c > 0 {
		return c
	}
	return r.confs
}

// Run starts every background job on its own interval and blocks until ctx is
// canceled. Intervals follow skydex-design.md §7.8.
func (r *Runner) Run(ctx context.Context) {
	go r.loop(ctx, 10*time.Second, "expiry", r.RunExpiry)
	go r.loop(ctx, 30*time.Second, "listing-check", r.RunListingCheck)
	go r.loop(ctx, 30*time.Second, "escrow-check", r.RunEscrowCheck)
	go r.loop(ctx, time.Minute, "return-scheduler", r.RunReturnScheduler)
	go r.loop(ctx, 5*time.Minute, "escrow-audit", r.RunEscrowAudit)
	go r.loop(ctx, time.Minute, "ban-manager", r.RunBanManager)
	go r.loop(ctx, time.Hour, "cleanup", r.RunCleanup)
	<-ctx.Done()
}

// loop runs fn every interval until ctx is canceled, logging any error.
func (r *Runner) loop(ctx context.Context, interval time.Duration, name string, fn func() error) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if err := fn(); err != nil {
				r.log.WithError(err).Warnf("skydex-market: job %s failed", name)
			}
		}
	}
}

// RunExpiry expires overdue listings and orders. An expired buy order releases
// its product back to the market and records a freeze violation against the
// buyer (skydex-design.md §5, §7.8/3).
func (r *Runner) RunExpiry() error {
	listings, err := r.db.GetExpiredPendingListings()
	if err != nil {
		return err
	}
	for _, l := range listings {
		if err := r.db.UpdatePendingListingStatus(l.ID, "expired", ""); err != nil {
			r.log.WithError(err).Warnf("expiry: failed to expire listing %s", l.ID)
		}
	}

	orders, err := r.db.GetExpiredOrders()
	if err != nil {
		return err
	}
	for _, o := range orders {
		if err := r.db.MarkOrderExpired(o.ID); err != nil {
			r.log.WithError(err).Warnf("expiry: failed to expire order %s", o.ID)
			continue
		}
		if err := r.db.UnfreezeProduct(o.ProductID); err != nil {
			r.log.WithError(err).Warnf("expiry: failed to unfreeze product %s", o.ProductID)
		}
		if err := r.db.CreateFreezeViolation(o.BuyerPubKey, o.ID); err != nil {
			r.log.WithError(err).Warnf("expiry: failed to record violation for %s", o.BuyerPubKey)
		}
	}
	return nil
}

// RunListingCheck promotes pending listings whose SKY deposit has settled in the
// market wallet into active products (skydex-design.md §7.8/2).
func (r *Runner) RunListingCheck() error {
	listings, err := r.db.GetPendingListings()
	if err != nil {
		return err
	}
	for _, l := range listings {
		// Each listing's escrow wallet is the deposit address of its sell coin.
		marketWallet, err := r.db.SellCoinWallet(l.SellCoin)
		if err != nil || marketWallet == "" {
			r.log.Warnf("listing-check: sell coin %s has no escrow wallet configured; skipping listing %s", l.SellCoin, l.ID)
			continue
		}
		// The deposit is identified by the seller's registered Skycoin-family
		// address (the tx sender) within the listing's window — not by a fiddly
		// non-round amount. Without an address on file the deposit can't be verified.
		sellerAddr, err := r.db.GetUserWallet(l.SellerPubKey, l.SellCoin)
		if err != nil || sellerAddr == "" {
			r.log.Warnf("listing-check: seller %s has no sell-coin wallet; cannot verify deposit for listing %s", l.SellerPubKey, l.ID)
			continue
		}
		confirmed, txHash, err := r.chain.DepositConfirmed(l.SellCoin, marketWallet, sellerAddr, l.ExpectedAmount, l.CreatedAt, l.ExpiresAt.Add(depositConfirmGrace))
		if err != nil {
			r.log.WithError(err).Warnf("listing-check: deposit check failed for %s", l.ID)
			continue
		}
		if !confirmed {
			continue
		}
		if err := r.db.UpdatePendingListingStatus(l.ID, "confirmed", txHash); err != nil {
			r.log.WithError(err).Warnf("listing-check: failed to confirm listing %s", l.ID)
			continue
		}
		product := &db.Product{
			ID:              uuid.NewString(),
			ListingID:       l.ID,
			SellerPubKey:    l.SellerPubKey,
			SellCoin:        l.SellCoin,
			Amount:          l.Amount,
			Price:           l.Price,
			PaymentCurrency: l.PaymentCurrency,
			Status:          "active",
		}
		if err := r.db.CreateProduct(product); err != nil {
			r.log.WithError(err).Warnf("listing-check: failed to create product for listing %s", l.ID)
		}
	}
	return nil
}

// RunEscrowCheck advances in-flight orders: it detects the buyer's payment,
// tracks confirmations, and on reaching the threshold delivers the escrowed SKY
// to the buyer and completes the trade (skydex-design.md §7.8/1).
func (r *Runner) RunEscrowCheck() error {
	orders, err := r.db.GetInFlightOrders()
	if err != nil {
		return err
	}
	for _, o := range orders {
		// The buyer's registered payment address (may be empty) is a soft
		// corroborating signal. The payment must land within the order's window
		// (created .. expiry + grace), which blocks replay of an old payment of a
		// recycled unique amount.
		buyerPay, _ := r.db.GetUserWallet(o.BuyerPubKey, o.PaymentCurrency) //nolint:errcheck
		confs, txHash, err := r.chain.PaymentConfirmations(o.PaymentCurrency, o.SellerWallet, buyerPay, o.ExpectedPaymentAmount, o.CreatedAt, o.ExpiresAt.Add(paymentConfirmGrace))
		if err != nil {
			r.log.WithError(err).Warnf("escrow-check: payment check failed for order %s", o.ID)
			continue
		}
		if confs <= 0 {
			continue // payment not observed yet
		}
		if o.Status == "pending_payment" {
			if err := r.db.MarkOrderPaid(o.ID, txHash); err != nil {
				r.log.WithError(err).Warnf("escrow-check: failed to mark order %s paid", o.ID)
				continue
			}
		}
		if err := r.db.SetOrderConfirmations(o.ID, confs); err != nil {
			r.log.WithError(err).Warnf("escrow-check: failed to update confirmations for %s", o.ID)
		}
		if confs < r.requiredConfirmations() {
			continue // not enough confirmations yet
		}
		if o.Status != "confirmed" {
			if err := r.db.MarkOrderConfirmed(o.ID); err != nil {
				r.log.WithError(err).Warnf("escrow-check: failed to confirm order %s", o.ID)
				continue
			}
		}
		buyerAddr, err := r.db.GetUserWallet(o.BuyerPubKey, o.SellCoin)
		if err != nil || buyerAddr == "" {
			r.log.Warnf("escrow-check: buyer %s has no sell-coin wallet; deferring delivery", o.BuyerPubKey)
			continue
		}
		if _, err := r.chain.SendCoin(o.SellCoin, buyerAddr, o.Amount); err != nil {
			// e.g. NoopChain (no backend) — leave the order 'confirmed' for a retry.
			r.log.WithError(err).Warnf("escrow-check: %s delivery deferred for order %s", o.SellCoin, o.ID)
			continue
		}
		if err := r.db.MarkOrderCompleted(o.ID); err != nil {
			r.log.WithError(err).Warnf("escrow-check: failed to complete order %s", o.ID)
			continue
		}
		// Book the sell-coin commission retained on the completed sale (for reporting;
		// it was physically collected as amount+commission at deposit and only the
		// amount was delivered). Recompute from the current config.
		rate, _ := r.db.GetCommissionRatePercent() //nolint:errcheck
		minC, _ := r.db.GetCommissionMinSKY()      //nolint:errcheck
		maxC, _ := r.db.GetCommissionMaxSKY()      //nolint:errcheck
		if c := Commission(o.Amount, rate, minC, maxC); c > 0 {
			if err := r.db.SetOrderCommission(o.ID, c); err != nil {
				r.log.WithError(err).Warnf("escrow-check: failed to record commission for order %s", o.ID)
			}
		}
		if err := r.db.MarkProductSold(o.ProductID); err != nil {
			r.log.WithError(err).Warnf("escrow-check: failed to mark product %s sold", o.ProductID)
		}
	}
	return nil
}

// RunReturnScheduler refunds escrowed SKY to sellers whose offer was canceled or
// expired after they had already deposited into the market wallet. Refunds are
// immediate — as soon as a listing goes terminal it is eligible. A refund is only
// sent once the deposit is still verifiable on-chain; if the escrow can't yet
// fund the refund transaction (coin-hours still accruing) SendSKY fails and it is
// retried next tick. The refund is recorded so it happens at most once
// (design §7.8/4).
func (r *Runner) RunReturnScheduler() error {
	listings, err := r.db.GetReturnableListings()
	if err != nil {
		return err
	}
	for _, l := range listings {
		marketWallet, err := r.db.SellCoinWallet(l.SellCoin)
		if err != nil || marketWallet == "" {
			r.log.Warnf("return-scheduler: sell coin %s has no escrow wallet; deferring refund of listing %s", l.SellCoin, l.ID)
			continue
		}
		sellerAddr, err := r.db.GetUserWallet(l.SellerPubKey, l.SellCoin)
		if err != nil || sellerAddr == "" {
			r.log.Warnf("return-scheduler: seller %s has no sell-coin wallet; deferring refund of listing %s", l.SellerPubKey, l.ID)
			continue
		}
		// Only refund coin that actually arrived in the market wallet, from this
		// seller, within the listing's deposit window. Without this on-chain check
		// a seller could create-and-cancel listings to drain the escrow wallet for
		// deposits they never made.
		confirmed, _, err := r.chain.DepositConfirmed(l.SellCoin, marketWallet, sellerAddr, l.ExpectedAmount, l.CreatedAt, l.ExpiresAt.Add(depositConfirmGrace))
		if err != nil {
			r.log.WithError(err).Warnf("return-scheduler: deposit check failed for listing %s", l.ID)
			continue
		}
		if !confirmed {
			continue // nothing escrowed for this listing
		}
		txHash, err := r.chain.SendCoin(l.SellCoin, sellerAddr, l.ExpectedAmount)
		if err != nil {
			// e.g. NoopChain (no backend) — leave it and retry next tick.
			r.log.WithError(err).Warnf("return-scheduler: %s refund deferred for listing %s", l.SellCoin, l.ID)
			continue
		}
		if err := r.db.MarkListingReturned(l.ID, txHash); err != nil {
			r.log.WithError(err).Warnf("return-scheduler: failed to record refund for listing %s", l.ID)
		}
	}
	return nil
}

// RunEscrowAudit cross-checks the escrow wallet's live on-chain SKY balance
// against the market's DB view of what it should be holding (products awaiting
// delivery + deposits awaiting refund). A shortfall — on-chain below outstanding
// obligations — means a delivery or refund could fail and is logged as an error;
// otherwise the healthy state (with any headroom) is logged. It only reads, never
// moves funds. Skipped when no chain backend is configured.
func (r *Runner) RunEscrowAudit() error {
	coins, err := r.db.AvailableSellCoins()
	if err != nil {
		return err
	}
	// Tolerance of half a milli-coin (below the network's 0.001 precision) absorbs
	// float rounding without masking a real deficit.
	const tol = 0.0005
	for _, coin := range coins {
		marketWallet, err := r.db.SellCoinWallet(coin)
		if err != nil || marketWallet == "" {
			continue // coin not fully configured yet — nothing to audit
		}
		onChain, err := r.chain.EscrowBalance(coin, marketWallet)
		if err != nil {
			if errors.Is(err, ErrNoChain) {
				continue // no backend configured — nothing to audit against
			}
			r.log.WithError(err).Warnf("escrow-audit: failed to read %s escrow balance", coin)
			continue
		}
		expected, err := r.db.OutstandingEscrow(coin)
		if err != nil {
			return err
		}
		if onChain+tol < expected {
			r.log.Errorf("escrow-audit: SHORTFALL — on-chain %s balance %.6f is below outstanding obligations %.6f (deficit %.6f); deliveries or refunds may fail",
				coin, onChain, expected, expected-onChain)
		} else {
			r.log.Infof("escrow-audit: %s healthy — on-chain %.6f covers %.6f of obligations (headroom %.6f)",
				coin, onChain, expected, onChain-expected)
		}
	}
	return nil
}

// RunBanManager releases expired bans and bans buyers who have exceeded the
// freeze-violation limit within the last 24 hours (skydex-design.md §5, §7.8/6).
func (r *Runner) RunBanManager() error {
	expired, err := r.db.GetExpiredBans()
	if err != nil {
		return err
	}
	for _, b := range expired {
		if err := r.db.DeleteBan(b.PubKey); err != nil {
			r.log.WithError(err).Warnf("ban-manager: failed to lift ban on %s", b.PubKey)
		}
	}

	limit, err := r.db.GetFreezeViolationsLimit()
	if err != nil {
		return err
	}
	days, err := r.db.GetBanDurationDays()
	if err != nil {
		return err
	}
	since := time.Now().UTC().Add(-24 * time.Hour)
	pks, err := r.db.GetBannablePubKeys(since, limit)
	if err != nil {
		return err
	}
	for _, pk := range pks {
		count, err := r.db.CountRecentViolations(pk, since)
		if err != nil {
			r.log.WithError(err).Warnf("ban-manager: failed to count violations for %s", pk)
			continue
		}
		banUntil := time.Now().UTC().Add(time.Duration(days) * 24 * time.Hour)
		if err := r.db.CreateBan(pk, count, banUntil); err != nil {
			r.log.WithError(err).Warnf("ban-manager: failed to ban %s", pk)
		}
	}
	return nil
}

// RunCleanup deletes data no longer needed for privacy: completed orders and
// terminal listings older than the configured retention, and old violations
// (skydex-design.md §7.8/5).
func (r *Runner) RunCleanup() error {
	days, err := r.db.GetCleanupDays()
	if err != nil {
		return err
	}
	age := time.Duration(days) * 24 * time.Hour

	old, err := r.db.GetCompletedOrdersOlderThan(age)
	if err != nil {
		return err
	}
	for _, o := range old {
		if err := r.db.DeleteOrder(o.ID); err != nil {
			r.log.WithError(err).Warnf("cleanup: failed to delete order %s", o.ID)
		}
	}

	if _, err := r.db.DeleteOldViolations(14 * 24 * time.Hour); err != nil {
		r.log.WithError(err).Warn("cleanup: failed to delete old violations")
	}
	if _, err := r.db.DeleteOldPendingListings(age); err != nil {
		r.log.WithError(err).Warn("cleanup: failed to delete old pending listings")
	}
	return nil
}
