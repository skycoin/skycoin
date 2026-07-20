// Package db internal/db/users.go
package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// CreateUser creates or updates a user's wallet addresses.
func (d *Database) CreateUser(user *User) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	user.UpdatedAt = time.Now().UTC()

	_, err := d.db.Exec(`
		INSERT INTO users (pubkey, wallet_sky, wallet_btc, wallet_bch, wallet_ltc, wallet_usdt_erc20, wallet_usdt_trc20, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(pubkey) DO UPDATE SET
			wallet_sky = excluded.wallet_sky,
			wallet_btc = excluded.wallet_btc,
			wallet_bch = excluded.wallet_bch,
			wallet_ltc = excluded.wallet_ltc,
			wallet_usdt_erc20 = excluded.wallet_usdt_erc20,
			wallet_usdt_trc20 = excluded.wallet_usdt_trc20,
			updated_at = excluded.updated_at
	`, user.PubKey, user.WalletSKY, user.WalletBTC, user.WalletBCH, user.WalletLTC,
		user.WalletUSDT_ERC20, user.WalletUSDT_TRC20, user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create/update user: %w", err)
	}

	return nil
}

// GetUser retrieves a user by their public key.
func (d *Database) GetUser(pubkey string) (*User, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	user := &User{}
	err := d.db.QueryRow(`
		SELECT pubkey, wallet_sky, wallet_btc, wallet_bch, wallet_ltc, wallet_usdt_erc20, wallet_usdt_trc20, updated_at
		FROM users
		WHERE pubkey = ?
	`, pubkey).Scan(&user.PubKey, &user.WalletSKY, &user.WalletBTC, &user.WalletBCH,
		&user.WalletLTC, &user.WalletUSDT_ERC20, &user.WalletUSDT_TRC20, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserWallet returns a user's address for a currency. SKY and every Skycoin
// fibercoin resolve to the user's single Skycoin-family address (they share the
// format) — that's where a buyer receives the purchased coin and a seller
// receives refunds. Every other (external) payment coin has its own payout
// address in user_wallets, so arbitrary operator-added coins are supported.
func (d *Database) GetUserWallet(pubkey, currency string) (string, error) {
	cur := strings.ToUpper(strings.TrimSpace(currency))

	if cur != "SKY" {
		// A configured sell coin (fibercoin) also resolves to the Skycoin address.
		sc, err := d.GetSellCoin(cur)
		if err != nil {
			return "", err
		}
		if sc == nil {
			// External payment coin: per-currency payout address.
			return d.getUserPayout(pubkey, cur)
		}
	}

	user, err := d.GetUser(pubkey)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("user not found: %s", pubkey)
	}
	return user.WalletSKY, nil
}

// getUserPayout returns the user's payout address for an external payment coin.
// It prefers the per-currency user_wallets entry and falls back to the legacy
// fixed columns for the built-in coins, so existing BTC/LTC registrations keep
// working unchanged.
func (d *Database) getUserPayout(pubkey, currency string) (string, error) {
	cur := strings.ToUpper(strings.TrimSpace(currency))

	d.mu.RLock()
	var addr string
	err := d.db.QueryRow(`SELECT address FROM user_wallets WHERE pubkey = ? AND currency = ?`, pubkey, cur).Scan(&addr)
	d.mu.RUnlock()
	if err == nil {
		return addr, nil
	}
	if err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to get payout wallet: %w", err)
	}

	// Back-compat: the built-in coins still have fixed columns on the users row.
	u, err := d.GetUser(pubkey)
	if err != nil || u == nil {
		return "", err
	}
	switch cur {
	case "BTC":
		return u.WalletBTC, nil
	case "BCH":
		return u.WalletBCH, nil
	case "LTC":
		return u.WalletLTC, nil
	case "USDT_ERC20":
		return u.WalletUSDT_ERC20, nil
	case "USDT_TRC20":
		return u.WalletUSDT_TRC20, nil
	}
	return "", nil
}

// SetUserWallet stores (or, for an empty address, clears) a user's external
// payout address for a currency.
func (d *Database) SetUserWallet(pubkey, currency, address string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	cur := strings.ToUpper(strings.TrimSpace(currency))
	addr := strings.TrimSpace(address)
	if cur == "" {
		return fmt.Errorf("currency is required")
	}
	if addr == "" {
		_, err := d.db.Exec(`DELETE FROM user_wallets WHERE pubkey = ? AND currency = ?`, pubkey, cur)
		if err != nil {
			return fmt.Errorf("failed to clear payout wallet: %w", err)
		}
		return nil
	}
	_, err := d.db.Exec(`
		INSERT INTO user_wallets (pubkey, currency, address) VALUES (?, ?, ?)
		ON CONFLICT(pubkey, currency) DO UPDATE SET address = excluded.address
	`, pubkey, cur, addr)
	if err != nil {
		return fmt.Errorf("failed to set payout wallet: %w", err)
	}
	return nil
}
