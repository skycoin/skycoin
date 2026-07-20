// Package db internal/db/orders.go
package db

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateOrder creates a new buy order.
func (d *Database) CreateOrder(order *Order) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if order.SellCoin == "" {
		order.SellCoin = "SKY"
	}

	order.CreatedAt = time.Now().UTC()

	_, err := d.db.Exec(`
		INSERT INTO orders (id, product_id, buyer_pubkey, sell_coin, amount_sky, price, payment_currency, expected_payment_amount, seller_wallet, status, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, order.ID, order.ProductID, order.BuyerPubKey, order.SellCoin, order.Amount, order.Price,
		order.PaymentCurrency, order.ExpectedPaymentAmount, order.SellerWallet,
		order.Status, order.ExpiresAt, order.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

// GetOrder retrieves an order by ID.
func (d *Database) GetOrder(id string) (*Order, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	order := &Order{}
	err := d.db.QueryRow(`
		SELECT id, product_id, buyer_pubkey, COALESCE(sell_coin, 'SKY') AS sell_coin, amount_sky, price, payment_currency, expected_payment_amount, seller_wallet, status, expires_at, created_at, paid_at, COALESCE(payment_tx_hash, '') AS payment_tx_hash, confirmations, completed_at
		FROM orders
		WHERE id = ?
	`, id).Scan(&order.ID, &order.ProductID, &order.BuyerPubKey, &order.SellCoin, &order.Amount, &order.Price,
		&order.PaymentCurrency, &order.ExpectedPaymentAmount, &order.SellerWallet,
		&order.Status, &order.ExpiresAt, &order.CreatedAt, &order.PaidAt,
		&order.PaymentTxHash, &order.Confirmations, &order.CompletedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return order, nil
}

// GetBuyerOrders returns all orders for a specific buyer.
func (d *Database) GetBuyerOrders(buyerPubKey string) ([]*Order, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query(`
		SELECT id, product_id, buyer_pubkey, COALESCE(sell_coin, 'SKY') AS sell_coin, amount_sky, price, payment_currency, expected_payment_amount, seller_wallet, status, expires_at, created_at, paid_at, COALESCE(payment_tx_hash, '') AS payment_tx_hash, confirmations, completed_at
		FROM orders
		WHERE buyer_pubkey = ?
		ORDER BY created_at DESC
	`, buyerPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get buyer orders: %w", err)
	}
	defer rows.Close() //nolint

	var orders []*Order
	for rows.Next() {
		order := &Order{}
		err := rows.Scan(&order.ID, &order.ProductID, &order.BuyerPubKey, &order.SellCoin, &order.Amount, &order.Price,
			&order.PaymentCurrency, &order.ExpectedPaymentAmount, &order.SellerWallet,
			&order.Status, &order.ExpiresAt, &order.CreatedAt, &order.PaidAt,
			&order.PaymentTxHash, &order.Confirmations, &order.CompletedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// GetSellerOrders returns all orders where the seller matches the given pubkey.
func (d *Database) GetSellerOrders(sellerPubKey string) ([]*Order, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query(`
		SELECT o.id, o.product_id, o.buyer_pubkey, COALESCE(o.sell_coin, 'SKY') AS sell_coin, o.amount_sky, o.price, o.payment_currency, o.expected_payment_amount, o.seller_wallet, o.status, o.expires_at, o.created_at, o.paid_at, COALESCE(o.payment_tx_hash, '') AS payment_tx_hash, o.confirmations, o.completed_at
		FROM orders o
		INNER JOIN products p ON o.product_id = p.id
		WHERE p.seller_pubkey = ?
		ORDER BY o.created_at DESC
	`, sellerPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller orders: %w", err)
	}
	defer rows.Close() //nolint

	var orders []*Order
	for rows.Next() {
		order := &Order{}
		err := rows.Scan(&order.ID, &order.ProductID, &order.BuyerPubKey, &order.SellCoin, &order.Amount, &order.Price,
			&order.PaymentCurrency, &order.ExpectedPaymentAmount, &order.SellerWallet,
			&order.Status, &order.ExpiresAt, &order.CreatedAt, &order.PaidAt,
			&order.PaymentTxHash, &order.Confirmations, &order.CompletedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// UpdateOrderStatus updates the status of an order.
func (d *Database) UpdateOrderStatus(id, status string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(`
		UPDATE orders
		SET status = ?
		WHERE id = ?
	`, status, id)

	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

// MarkOrderPaid marks an order as paid with the transaction hash.
func (d *Database) MarkOrderPaid(orderID, txHash string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now().UTC()
	_, err := d.db.Exec(`
		UPDATE orders
		SET status = 'paid', paid_at = ?, payment_tx_hash = ?
		WHERE id = ?
	`, now, txHash, orderID)

	if err != nil {
		return fmt.Errorf("failed to mark order as paid: %w", err)
	}

	return nil
}

// IncrementOrderConfirmations increments the confirmation count for an order.
func (d *Database) IncrementOrderConfirmations(orderID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(`
		UPDATE orders
		SET confirmations = confirmations + 1
		WHERE id = ?
	`, orderID)

	if err != nil {
		return fmt.Errorf("failed to increment order confirmations: %w", err)
	}

	return nil
}

// MarkOrderConfirmed marks an order as confirmed (2+ confirmations received).
func (d *Database) MarkOrderConfirmed(orderID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(`
		UPDATE orders
		SET status = 'confirmed'
		WHERE id = ?
	`, orderID)

	if err != nil {
		return fmt.Errorf("failed to mark order as confirmed: %w", err)
	}

	return nil
}

// MarkOrderCompleted marks an order as completed (SKY transferred to buyer).
func (d *Database) MarkOrderCompleted(orderID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now().UTC()
	_, err := d.db.Exec(`
		UPDATE orders
		SET status = 'completed', completed_at = ?
		WHERE id = ?
	`, now, orderID)

	if err != nil {
		return fmt.Errorf("failed to mark order as completed: %w", err)
	}

	return nil
}

// MarkOrderCanceled marks an in-flight order as canceled by the buyer.
func (d *Database) MarkOrderCanceled(orderID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(`
		UPDATE orders
		SET status = 'canceled'
		WHERE id = ?
	`, orderID)

	if err != nil {
		return fmt.Errorf("failed to mark order as canceled: %w", err)
	}

	return nil
}

// CountBuyerProductCancels returns how many orders this buyer has canceled or
// let expire for this product. Cleared blocks ('cancelled_cleared', set when the
// operator unblocks a buyer) are not counted.
func (d *Database) CountBuyerProductCancels(buyerPubKey, productID string) (int, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var n int
	err := d.db.QueryRow(`
		SELECT COUNT(*) FROM orders
		WHERE buyer_pubkey = ? AND product_id = ?
		  AND status IN ('canceled', 'cancelled', 'expired')
	`, buyerPubKey, productID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("failed to count buyer cancels: %w", err)
	}
	return n, nil
}

// BuyerBlockedFromProduct reports whether the buyer has reached the cancel limit
// for this product — i.e. has canceled or let expire at least `limit` orders for
// it — and so may not buy that same product again until the operator clears it.
func (d *Database) BuyerBlockedFromProduct(buyerPubKey, productID string, limit int) (bool, error) {
	n, err := d.CountBuyerProductCancels(buyerPubKey, productID)
	if err != nil {
		return false, err
	}
	return n >= limit, nil
}

// ListBuyerProductBlocks returns every (buyer, product) pair that has reached the
// given cancel limit, so the operator can review and clear them. Ordered by the
// most cancellations first.
func (d *Database) ListBuyerProductBlocks(limit int) ([]*BuyerProductBlock, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query(`
		SELECT buyer_pubkey, product_id, COUNT(*) AS cancels
		FROM orders
		WHERE status IN ('canceled', 'cancelled', 'expired')
		GROUP BY buyer_pubkey, product_id
		HAVING cancels >= ?
		ORDER BY cancels DESC
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list buyer blocks: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	var blocks []*BuyerProductBlock
	for rows.Next() {
		b := &BuyerProductBlock{}
		if err := rows.Scan(&b.BuyerPubKey, &b.ProductID, &b.Cancels); err != nil {
			return nil, fmt.Errorf("failed to scan buyer block: %w", err)
		}
		blocks = append(blocks, b)
	}
	return blocks, rows.Err()
}

// ClearBuyerProductBlock lifts a buyer's block for one product by marking their
// canceled/expired orders for it as 'cancelled_cleared', which no longer counts
// toward the cancel limit. Returns the number of orders cleared.
func (d *Database) ClearBuyerProductBlock(buyerPubKey, productID string) (int64, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	res, err := d.db.Exec(`
		UPDATE orders
		SET status = 'cancelled_cleared'
		WHERE buyer_pubkey = ? AND product_id = ?
		  AND status IN ('canceled', 'cancelled', 'expired')
	`, buyerPubKey, productID)
	if err != nil {
		return 0, fmt.Errorf("failed to clear buyer block: %w", err)
	}
	n, _ := res.RowsAffected() //nolint:errcheck
	return n, nil
}

// MarkOrderExpired marks an order as expired.
func (d *Database) MarkOrderExpired(orderID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(`
		UPDATE orders
		SET status = 'expired'
		WHERE id = ?
	`, orderID)

	if err != nil {
		return fmt.Errorf("failed to mark order as expired: %w", err)
	}

	return nil
}

// GetExpiredOrders returns all orders that have expired (pending_payment status + past expires_at).
func (d *Database) GetExpiredOrders() ([]*Order, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	now := time.Now().UTC()
	rows, err := d.db.Query(`
		SELECT id, product_id, buyer_pubkey, COALESCE(sell_coin, 'SKY') AS sell_coin, amount_sky, price, payment_currency, expected_payment_amount, seller_wallet, status, expires_at, created_at, paid_at, COALESCE(payment_tx_hash, '') AS payment_tx_hash, confirmations, completed_at
		FROM orders
		WHERE status = 'pending_payment' AND expires_at <= ?
	`, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired orders: %w", err)
	}
	defer rows.Close() //nolint

	var orders []*Order
	for rows.Next() {
		order := &Order{}
		err := rows.Scan(&order.ID, &order.ProductID, &order.BuyerPubKey, &order.SellCoin, &order.Amount, &order.Price,
			&order.PaymentCurrency, &order.ExpectedPaymentAmount, &order.SellerWallet,
			&order.Status, &order.ExpiresAt, &order.CreatedAt, &order.PaidAt,
			&order.PaymentTxHash, &order.Confirmations, &order.CompletedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// GetPendingPaymentOrders returns all orders waiting for payment confirmation.
func (d *Database) GetPendingPaymentOrders() ([]*Order, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query(`
		SELECT id, product_id, buyer_pubkey, COALESCE(sell_coin, 'SKY') AS sell_coin, amount_sky, price, payment_currency, expected_payment_amount, seller_wallet, status, expires_at, created_at, paid_at, COALESCE(payment_tx_hash, '') AS payment_tx_hash, confirmations, completed_at
		FROM orders
		WHERE status = 'pending_payment' AND expires_at > ?
		ORDER BY created_at ASC
	`, time.Now().UTC())
	if err != nil {
		return nil, fmt.Errorf("failed to get pending payment orders: %w", err)
	}
	defer rows.Close() //nolint

	var orders []*Order
	for rows.Next() {
		order := &Order{}
		err := rows.Scan(&order.ID, &order.ProductID, &order.BuyerPubKey, &order.SellCoin, &order.Amount, &order.Price,
			&order.PaymentCurrency, &order.ExpectedPaymentAmount, &order.SellerWallet,
			&order.Status, &order.ExpiresAt, &order.CreatedAt, &order.PaidAt,
			&order.PaymentTxHash, &order.Confirmations, &order.CompletedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// GetCompletedOrdersOlderThan returns all completed orders older than the specified duration.
func (d *Database) GetCompletedOrdersOlderThan(age time.Duration) ([]*Order, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	cutoff := time.Now().UTC().Add(-age)
	rows, err := d.db.Query(`
		SELECT id, product_id, buyer_pubkey, COALESCE(sell_coin, 'SKY') AS sell_coin, amount_sky, price, payment_currency, expected_payment_amount, seller_wallet, status, expires_at, created_at, paid_at, COALESCE(payment_tx_hash, '') AS payment_tx_hash, confirmations, completed_at
		FROM orders
		WHERE status = 'completed' AND completed_at <= ?
	`, cutoff)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed orders: %w", err)
	}
	defer rows.Close() //nolint

	var orders []*Order
	for rows.Next() {
		order := &Order{}
		err := rows.Scan(&order.ID, &order.ProductID, &order.BuyerPubKey, &order.SellCoin, &order.Amount, &order.Price,
			&order.PaymentCurrency, &order.ExpectedPaymentAmount, &order.SellerWallet,
			&order.Status, &order.ExpiresAt, &order.CreatedAt, &order.PaidAt,
			&order.PaymentTxHash, &order.Confirmations, &order.CompletedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// DeleteOrder deletes an order by ID.
func (d *Database) DeleteOrder(orderID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(`DELETE FROM orders WHERE id = ?`, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	return nil
}

// GetAllOrders returns every order regardless of status, newest first.
// Used by the market operator UI for monitoring.
func (d *Database) GetAllOrders() ([]*Order, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query(`
		SELECT id, product_id, buyer_pubkey, COALESCE(sell_coin, 'SKY') AS sell_coin, amount_sky, price, payment_currency, expected_payment_amount, seller_wallet, status, expires_at, created_at, paid_at, COALESCE(payment_tx_hash, '') AS payment_tx_hash, confirmations, completed_at, commission_sky
		FROM orders
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders: %w", err)
	}
	defer rows.Close() //nolint

	var orders []*Order
	for rows.Next() {
		order := &Order{}
		err := rows.Scan(&order.ID, &order.ProductID, &order.BuyerPubKey, &order.SellCoin, &order.Amount, &order.Price,
			&order.PaymentCurrency, &order.ExpectedPaymentAmount, &order.SellerWallet,
			&order.Status, &order.ExpiresAt, &order.CreatedAt, &order.PaidAt,
			&order.PaymentTxHash, &order.Confirmations, &order.CompletedAt, &order.Commission)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate orders: %w", err)
	}

	return orders, nil
}

// SetOrderCommission records the SKY commission the market retained on a sale.
func (d *Database) SetOrderCommission(orderID string, sky float64) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(`UPDATE orders SET commission_sky = ? WHERE id = ?`, sky, orderID)
	if err != nil {
		return fmt.Errorf("failed to set order commission: %w", err)
	}
	return nil
}

// SetOrderConfirmations sets the current confirmation count for an order.
func (d *Database) SetOrderConfirmations(orderID string, confirmations int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(`
		UPDATE orders SET confirmations = ? WHERE id = ?
	`, confirmations, orderID)
	if err != nil {
		return fmt.Errorf("failed to set order confirmations: %w", err)
	}
	return nil
}

// GetInFlightOrders returns orders that are actively being processed by the
// escrow checker: awaiting payment (not yet expired), paid, or confirmed.
func (d *Database) GetInFlightOrders() ([]*Order, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	now := time.Now().UTC()
	rows, err := d.db.Query(`
		SELECT id, product_id, buyer_pubkey, COALESCE(sell_coin, 'SKY') AS sell_coin, amount_sky, price, payment_currency, expected_payment_amount, seller_wallet, status, expires_at, created_at, paid_at, COALESCE(payment_tx_hash, '') AS payment_tx_hash, confirmations, completed_at
		FROM orders
		WHERE status IN ('paid', 'confirmed')
		   OR (status = 'pending_payment' AND expires_at > ?)
		ORDER BY created_at ASC
	`, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get in-flight orders: %w", err)
	}
	defer rows.Close() //nolint

	var orders []*Order
	for rows.Next() {
		order := &Order{}
		if err := rows.Scan(&order.ID, &order.ProductID, &order.BuyerPubKey, &order.SellCoin, &order.Amount, &order.Price,
			&order.PaymentCurrency, &order.ExpectedPaymentAmount, &order.SellerWallet,
			&order.Status, &order.ExpiresAt, &order.CreatedAt, &order.PaidAt,
			&order.PaymentTxHash, &order.Confirmations, &order.CompletedAt); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate in-flight orders: %w", err)
	}

	return orders, nil
}

// BuyerHasPendingPayment reports whether the given buyer already has a pending
// payment of ~amount (within tol) to the same seller wallet in the same currency.
// Fibercoin payments use a FIXED amount, so two payments of the same amount from
// the SAME buyer to the same seller would be ambiguous (the sender no longer
// disambiguates them); the buy handler uses this to reject the second one. Across
// different buyers the same fixed amount is fine — the sender tells them apart.
func (d *Database) BuyerHasPendingPayment(buyerPubKey, sellerWallet, currency string, amount, tol float64) (bool, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var n int
	err := d.db.QueryRow(`
		SELECT COUNT(*) FROM orders
		WHERE status = 'pending_payment' AND buyer_pubkey = ? AND seller_wallet = ? AND payment_currency = ?
		  AND ABS(expected_payment_amount - ?) < ?
	`, buyerPubKey, sellerWallet, currency, amount, tol).Scan(&n)
	if err != nil {
		return false, fmt.Errorf("failed to check buyer pending payment: %w", err)
	}
	return n > 0, nil
}

// OrderPaymentAmountExists reports whether an order awaiting payment already
// expects a payment within tol of amount to the same seller wallet in the same
// currency. Used to keep every pending payment's non-round amount unique per
// recipient wallet.
func (d *Database) OrderPaymentAmountExists(sellerWallet, currency string, amount, tol float64) (bool, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var n int
	err := d.db.QueryRow(`
		SELECT COUNT(*) FROM orders
		WHERE status = 'pending_payment' AND seller_wallet = ? AND payment_currency = ?
		  AND ABS(expected_payment_amount - ?) < ?
	`, sellerWallet, currency, amount, tol).Scan(&n)
	if err != nil {
		return false, fmt.Errorf("failed to check pending payment amount: %w", err)
	}
	return n > 0, nil
}
