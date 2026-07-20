// Package db internal/db/sellcoins.go — the sell-side coin registry.
//
// The market sells SKY and any operator-added Skycoin fibercoins. Each is a row
// in the sell_coins table with its own fullnode URL and escrow hot wallet
// (seed + address) and confirmation depth. This mirrors the per-currency
// explorer config on the buy side: operators add/enable/disable/delete coins,
// and a listing/product/order records which coin it trades.
package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// ListSellCoins returns every configured sell coin, SKY first then alphabetical.
func (d *Database) ListSellCoins() ([]*SellCoin, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query(`
		SELECT symbol, name, node_url, wallet_seed, wallet_addr, confirmations, enabled, updated_at
		FROM sell_coins
		ORDER BY (symbol = 'SKY') DESC, symbol ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list sell coins: %w", err)
	}
	defer rows.Close() //nolint

	var out []*SellCoin
	for rows.Next() {
		sc, err := scanSellCoin(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, sc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate sell coins: %w", err)
	}
	return out, nil
}

// GetSellCoin returns one sell coin by symbol, or nil if it does not exist.
func (d *Database) GetSellCoin(symbol string) (*SellCoin, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	row := d.db.QueryRow(`
		SELECT symbol, name, node_url, wallet_seed, wallet_addr, confirmations, enabled, updated_at
		FROM sell_coins WHERE symbol = ?
	`, strings.ToUpper(strings.TrimSpace(symbol)))
	sc, err := scanSellCoin(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return sc, nil
}

// scanner abstracts *sql.Row and *sql.Rows so scanSellCoin serves both.
type scanner interface {
	Scan(dest ...any) error
}

func scanSellCoin(s scanner) (*SellCoin, error) {
	sc := &SellCoin{}
	var enabled int
	if err := s.Scan(&sc.Symbol, &sc.Name, &sc.NodeURL, &sc.WalletSeed, &sc.WalletAddr,
		&sc.Confirmations, &enabled, &sc.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("failed to scan sell coin: %w", err)
	}
	sc.Enabled = enabled != 0
	return sc, nil
}

// UpsertSellCoin inserts or updates a sell coin, keyed by symbol (upper-cased).
func (d *Database) UpsertSellCoin(sc *SellCoin) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	symbol := strings.ToUpper(strings.TrimSpace(sc.Symbol))
	if symbol == "" {
		return fmt.Errorf("sell coin symbol is required")
	}
	confs := sc.Confirmations
	if confs < 1 {
		confs = 1
	}
	enabled := 0
	if sc.Enabled {
		enabled = 1
	}
	_, err := d.db.Exec(`
		INSERT INTO sell_coins (symbol, name, node_url, wallet_seed, wallet_addr, confirmations, enabled, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(symbol) DO UPDATE SET
			name = excluded.name,
			node_url = excluded.node_url,
			wallet_seed = excluded.wallet_seed,
			wallet_addr = excluded.wallet_addr,
			confirmations = excluded.confirmations,
			enabled = excluded.enabled,
			updated_at = excluded.updated_at
	`, symbol, strings.TrimSpace(sc.Name), strings.TrimSpace(sc.NodeURL),
		strings.TrimSpace(sc.WalletSeed), strings.TrimSpace(sc.WalletAddr),
		confs, enabled, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("failed to upsert sell coin %s: %w", symbol, err)
	}
	return nil
}

// DeleteSellCoin removes a sell coin. Existing listings/orders keep their
// recorded sell_coin symbol; they simply can no longer settle if the coin's
// escrow config is gone, so operators should disable rather than delete a coin
// with in-flight trades.
func (d *Database) DeleteSellCoin(symbol string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(`DELETE FROM sell_coins WHERE symbol = ?`, strings.ToUpper(strings.TrimSpace(symbol)))
	if err != nil {
		return fmt.Errorf("failed to delete sell coin %s: %w", symbol, err)
	}
	return nil
}

// SetSellCoinEnabled toggles a sell coin's availability without touching its
// escrow config.
func (d *Database) SetSellCoinEnabled(symbol string, enabled bool) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	e := 0
	if enabled {
		e = 1
	}
	_, err := d.db.Exec(`UPDATE sell_coins SET enabled = ?, updated_at = ? WHERE symbol = ?`,
		e, time.Now().UTC(), strings.ToUpper(strings.TrimSpace(symbol)))
	if err != nil {
		return fmt.Errorf("failed to set sell coin enabled %s: %w", symbol, err)
	}
	return nil
}

// IsSellCoinAvailable reports whether symbol is an enabled sell coin with a
// complete escrow config (node + seed + address). An enabled-but-incomplete coin
// is never treated as available, so a trade can't be created against a coin that
// could not settle.
func (d *Database) IsSellCoinAvailable(symbol string) (bool, error) {
	nodeURL, seed, wallet, _, enabled, err := d.SellCoinConfig(symbol)
	if err != nil {
		return false, err
	}
	return enabled && SellCoinComplete(nodeURL, seed, wallet), nil
}

// AvailableSellCoins returns the symbols of enabled, fully-configured sell coins,
// SKY first.
func (d *Database) AvailableSellCoins() ([]string, error) {
	coins, err := d.ListSellCoins()
	if err != nil {
		return nil, err
	}
	var out []string
	for _, c := range coins {
		if !c.Enabled {
			continue
		}
		nodeURL, seed, wallet, _, _, err := d.SellCoinConfig(c.Symbol)
		if err != nil {
			return nil, err
		}
		if SellCoinComplete(nodeURL, seed, wallet) {
			out = append(out, c.Symbol)
		}
	}
	return out, nil
}

// SellCoinConfig returns the escrow-node config for symbol, satisfying
// chain.SellCoinConfigStore. For SKY, empty row fields fall back to the legacy
// single-coin config keys (sky_fullnode_url / sky_wallet_seed / wallet_sky) so a
// pre-fibercoin deployment keeps working unchanged.
func (d *Database) SellCoinConfig(symbol string) (nodeURL, seed, wallet string, confs int, enabled bool, err error) {
	sym := strings.ToUpper(strings.TrimSpace(symbol))
	sc, err := d.GetSellCoin(sym)
	if err != nil {
		return "", "", "", 0, false, err
	}
	if sc == nil {
		return "", "", "", 0, false, nil // unknown coin => unavailable
	}
	nodeURL, seed, wallet, confs, enabled = sc.NodeURL, sc.WalletSeed, sc.WalletAddr, sc.Confirmations, sc.Enabled
	if sym == "SKY" {
		if nodeURL == "" {
			nodeURL, _ = d.configOrEmpty("sky_fullnode_url") //nolint:errcheck
		}
		if seed == "" {
			seed, _ = d.configOrEmpty("sky_wallet_seed") //nolint:errcheck
		}
		if wallet == "" {
			wallet, _ = d.configOrEmpty("wallet_sky") //nolint:errcheck
		}
		if strings.TrimSpace(nodeURL) == "" {
			nodeURL = "https://node.skycoin.com" // public Skycoin fullnode default
		}
	}
	if confs < 1 {
		confs = 1
	}
	return nodeURL, seed, wallet, confs, enabled, nil
}

// SellCoinComplete reports whether a sell coin has everything it needs to settle
// trades: a fullnode URL, an escrow seed and an escrow address. A coin missing
// any of these can be saved as a draft but must not be enabled — it could never
// verify a deposit or sign a delivery/refund.
func SellCoinComplete(nodeURL, seed, wallet string) bool {
	return strings.TrimSpace(nodeURL) != "" &&
		strings.TrimSpace(seed) != "" &&
		strings.TrimSpace(wallet) != ""
}

// SellCoinWallet returns the escrow (deposit) address for symbol — where sellers
// send the coin. For SKY it falls back to the legacy wallet_sky key.
func (d *Database) SellCoinWallet(symbol string) (string, error) {
	_, _, wallet, _, _, err := d.SellCoinConfig(symbol)
	return wallet, err
}

// IsPaymentAvailable reports whether currency may be used as a listing's payment
// currency: either an enabled external explorer coin (BTC/LTC/…) or an enabled
// sell coin (a fibercoin the buyer pays peer-to-peer, verified on its own node).
func (d *Database) IsPaymentAvailable(currency string) (bool, error) {
	if ok, err := d.IsCurrencyAvailable(currency); err != nil || ok {
		return ok, err
	}
	return d.IsSellCoinAvailable(currency)
}
