// Package db internal/db/listings.go
package db

import (
	"database/sql"
	"fmt"
	"time"
)

// CreatePendingListing creates a new pending listing (sell order awaiting SKY transfer).
func (d *Database) CreatePendingListing(listing *PendingListing) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if listing.SellCoin == "" {
		listing.SellCoin = "SKY"
	}

	listing.CreatedAt = time.Now().UTC()

	_, err := d.db.Exec(`
		INSERT INTO pending_listings (id, seller_pubkey, sell_coin, amount_sky, expected_amount_sky, price, payment_currency, status, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, listing.ID, listing.SellerPubKey, listing.SellCoin, listing.Amount, listing.ExpectedAmount,
		listing.Price, listing.PaymentCurrency, listing.Status, listing.ExpiresAt, listing.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create pending listing: %w", err)
	}

	return nil
}

// GetPendingListing retrieves a pending listing by ID.
func (d *Database) GetPendingListing(id string) (*PendingListing, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	listing := &PendingListing{}
	err := d.db.QueryRow(`
		SELECT id, seller_pubkey, COALESCE(sell_coin, 'SKY') AS sell_coin, amount_sky, expected_amount_sky, price, payment_currency, status, expires_at, created_at, confirmed_at, COALESCE(tx_hash, '') AS tx_hash, closed_at, returned_at, COALESCE(return_tx_hash, '') AS return_tx_hash
		FROM pending_listings
		WHERE id = ?
	`, id).Scan(&listing.ID, &listing.SellerPubKey, &listing.SellCoin, &listing.Amount, &listing.ExpectedAmount,
		&listing.Price, &listing.PaymentCurrency, &listing.Status, &listing.ExpiresAt,
		&listing.CreatedAt, &listing.ConfirmedAt, &listing.TxHash,
		&listing.ClosedAt, &listing.ReturnedAt, &listing.ReturnTxHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get pending listing: %w", err)
	}

	return listing, nil
}

// UpdatePendingListingStatus updates the status of a pending listing.
func (d *Database) UpdatePendingListingStatus(id, status string, txHash string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now().UTC()
	var confirmedAt, closedAt *time.Time
	switch status {
	case "confirmed":
		confirmedAt = &now
	case "expired", "canceled":
		// Stamp the moment the listing went terminal; the Return Scheduler
		// starts the refund-delay clock from here.
		closedAt = &now
	}

	// COALESCE so setting one timestamp never clears the other on a later
	// transition.
	_, err := d.db.Exec(`
		UPDATE pending_listings
		SET status = ?,
		    confirmed_at = COALESCE(?, confirmed_at),
		    closed_at = COALESCE(?, closed_at),
		    tx_hash = ?
		WHERE id = ?
	`, status, confirmedAt, closedAt, txHash, id)

	if err != nil {
		return fmt.Errorf("failed to update pending listing status: %w", err)
	}

	return nil
}

// GetExpiredPendingListings returns all pending listings that have expired.
func (d *Database) GetExpiredPendingListings() ([]*PendingListing, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	now := time.Now().UTC()
	rows, err := d.db.Query(`
		SELECT id, seller_pubkey, COALESCE(sell_coin, 'SKY') AS sell_coin, amount_sky, expected_amount_sky, price, payment_currency, status, expires_at, created_at, confirmed_at, COALESCE(tx_hash, '') AS tx_hash
		FROM pending_listings
		WHERE status = 'pending' AND expires_at <= ?
	`, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired pending listings: %w", err)
	}
	defer rows.Close() //nolint

	var listings []*PendingListing
	for rows.Next() {
		listing := &PendingListing{}
		err := rows.Scan(&listing.ID, &listing.SellerPubKey, &listing.SellCoin, &listing.Amount, &listing.ExpectedAmount,
			&listing.Price, &listing.PaymentCurrency, &listing.Status, &listing.ExpiresAt,
			&listing.CreatedAt, &listing.ConfirmedAt, &listing.TxHash)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pending listing: %w", err)
		}
		listings = append(listings, listing)
	}

	return listings, nil
}

// GetSellerPendingListings returns all pending listings for a specific seller.
func (d *Database) GetSellerPendingListings(sellerPubKey string) ([]*PendingListing, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query(`
		SELECT id, seller_pubkey, COALESCE(sell_coin, 'SKY') AS sell_coin, amount_sky, expected_amount_sky, price, payment_currency, status, expires_at, created_at, confirmed_at, COALESCE(tx_hash, '') AS tx_hash
		FROM pending_listings
		WHERE seller_pubkey = ? AND status IN ('pending', 'confirmed')
		ORDER BY created_at DESC
	`, sellerPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller pending listings: %w", err)
	}
	defer rows.Close() //nolint

	var listings []*PendingListing
	for rows.Next() {
		listing := &PendingListing{}
		err := rows.Scan(&listing.ID, &listing.SellerPubKey, &listing.SellCoin, &listing.Amount, &listing.ExpectedAmount,
			&listing.Price, &listing.PaymentCurrency, &listing.Status, &listing.ExpiresAt,
			&listing.CreatedAt, &listing.ConfirmedAt, &listing.TxHash)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pending listing: %w", err)
		}
		listings = append(listings, listing)
	}

	return listings, nil
}

// GetPendingListings returns pending listings that have not yet expired. The
// Listing Checker polls these to detect the seller's SKY deposit.
func (d *Database) GetPendingListings() ([]*PendingListing, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	now := time.Now().UTC()
	rows, err := d.db.Query(`
		SELECT id, seller_pubkey, COALESCE(sell_coin, 'SKY') AS sell_coin, amount_sky, expected_amount_sky, price, payment_currency, status, expires_at, created_at, confirmed_at, COALESCE(tx_hash, '') AS tx_hash
		FROM pending_listings
		WHERE status = 'pending' AND expires_at > ?
		ORDER BY created_at ASC
	`, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending listings: %w", err)
	}
	defer rows.Close() //nolint

	var listings []*PendingListing
	for rows.Next() {
		listing := &PendingListing{}
		if err := rows.Scan(&listing.ID, &listing.SellerPubKey, &listing.SellCoin, &listing.Amount, &listing.ExpectedAmount,
			&listing.Price, &listing.PaymentCurrency, &listing.Status, &listing.ExpiresAt,
			&listing.CreatedAt, &listing.ConfirmedAt, &listing.TxHash); err != nil {
			return nil, fmt.Errorf("failed to scan pending listing: %w", err)
		}
		listings = append(listings, listing)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate pending listings: %w", err)
	}

	return listings, nil
}

// GetReturnableListings returns terminal (expired/canceled) pending listings that
// have not yet been refunded — every one is eligible immediately. closed_at IS
// NOT NULL restricts to actually terminated listings; oldest deposit first. The
// Return Scheduler still verifies the on-chain deposit before sending anything
// back.
func (d *Database) GetReturnableListings() ([]*PendingListing, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query(`
		SELECT id, seller_pubkey, COALESCE(sell_coin, 'SKY') AS sell_coin, amount_sky, expected_amount_sky, price, payment_currency, status, expires_at, created_at, confirmed_at, COALESCE(tx_hash, '') AS tx_hash
		FROM pending_listings
		WHERE status IN ('expired', 'canceled', 'cancelled')
		  AND returned_at IS NULL
		  AND closed_at IS NOT NULL
		ORDER BY COALESCE(confirmed_at, created_at) ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get returnable listings: %w", err)
	}
	defer rows.Close() //nolint

	var listings []*PendingListing
	for rows.Next() {
		listing := &PendingListing{}
		if err := rows.Scan(&listing.ID, &listing.SellerPubKey, &listing.SellCoin, &listing.Amount, &listing.ExpectedAmount,
			&listing.Price, &listing.PaymentCurrency, &listing.Status, &listing.ExpiresAt,
			&listing.CreatedAt, &listing.ConfirmedAt, &listing.TxHash); err != nil {
			return nil, fmt.Errorf("failed to scan returnable listing: %w", err)
		}
		listings = append(listings, listing)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate returnable listings: %w", err)
	}

	return listings, nil
}

// MarkListingReturned records that escrowed SKY for a terminal listing has been
// refunded to the seller, storing the refund transaction hash. Idempotent: it
// only stamps listings that have not already been marked returned, so a
// duplicate call cannot double-account a refund.
func (d *Database) MarkListingReturned(id, txHash string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(`
		UPDATE pending_listings
		SET status = 'returned', returned_at = ?, return_tx_hash = ?
		WHERE id = ? AND returned_at IS NULL
	`, time.Now().UTC(), txHash, id)
	if err != nil {
		return fmt.Errorf("failed to mark listing returned: %w", err)
	}
	return nil
}

// DeleteOldPendingListings removes terminal (expired/confirmed/canceled/
// returned) pending listings older than age. Returns the number deleted.
func (d *Database) DeleteOldPendingListings(age time.Duration) (int64, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	cutoff := time.Now().UTC().Add(-age)
	res, err := d.db.Exec(`
		DELETE FROM pending_listings
		WHERE status IN ('expired', 'confirmed', 'canceled', 'cancelled', 'returned') AND created_at < ?
	`, cutoff)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old pending listings: %w", err)
	}
	n, _ := res.RowsAffected() //nolint:errcheck
	return n, nil
}

// OutstandingEscrow returns the total of a single sell coin the market should
// currently be holding in that coin's escrow wallet: the amount for every product
// still awaiting delivery (active or frozen) plus the deposit for every terminal
// listing whose coin arrived on-chain (confirmed_at set) but has not yet been
// refunded. Delivered ('sold' products) and already-refunded listings are
// excluded — their coin has left escrow. Used by the per-coin escrow-audit job.
func (d *Database) OutstandingEscrow(coin string) (float64, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var total float64
	err := d.db.QueryRow(`
		SELECT
		  COALESCE((SELECT SUM(amount_sky) FROM products
		            WHERE status IN ('active','frozen') AND sell_coin = ?), 0)
		+ COALESCE((SELECT SUM(expected_amount_sky) FROM pending_listings
		            WHERE status IN ('expired','canceled','cancelled')
		              AND returned_at IS NULL
		              AND closed_at IS NOT NULL
		              AND confirmed_at IS NOT NULL
		              AND sell_coin = ?), 0)
	`, coin, coin).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to sum outstanding escrow: %w", err)
	}
	return total, nil
}

// SellerHasPendingListing reports whether the seller already has a listing that
// is still awaiting its SKY deposit (status 'pending' and not yet expired). A
// deposit is identified by the seller's SKY address within the listing's window,
// so two concurrently-pending listings from one seller would be ambiguous; the
// create-listing handler uses this to allow only one pending listing per seller.
func (d *Database) SellerHasPendingListing(sellerPubKey string) (bool, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var n int
	err := d.db.QueryRow(`
		SELECT COUNT(*) FROM pending_listings
		WHERE seller_pubkey = ? AND status = 'pending' AND expires_at > ?
	`, sellerPubKey, time.Now().UTC()).Scan(&n)
	if err != nil {
		return false, fmt.Errorf("failed to check seller pending listing: %w", err)
	}
	return n > 0, nil
}
