// Package db internal/db/config.go
package db

import (
	"database/sql"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/skydex/protocol"
)

// ExplorerKeyPrefix returns the market_config key prefix for a currency's
// explorer settings, e.g. "explorer_btc". The full keys are <prefix>_provider,
// <prefix>_url and <prefix>_key.
func ExplorerKeyPrefix(currency string) string {
	return "explorer_" + strings.ToLower(currency)
}

// configOrEmpty reads a config value, returning "" (not an error) when the key
// has never been set.
func (d *Database) configOrEmpty(key string) (string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var v string
	err := d.db.QueryRow(`SELECT value FROM market_config WHERE key = ?`, key).Scan(&v)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get config %s: %w", key, err)
	}
	return v, nil
}

// ExplorerConfig returns the operator's configured explorer provider, optional
// base-URL override and API key for a currency. An empty provider means the
// currency has no explorer configured (and is therefore unavailable). Satisfies
// chain.ExplorerConfigStore.
func (d *Database) ExplorerConfig(currency string) (provider, baseURL, apiKey string, err error) {
	p := ExplorerKeyPrefix(currency)
	if provider, err = d.configOrEmpty(p + "_provider"); err != nil {
		return "", "", "", err
	}
	if baseURL, err = d.configOrEmpty(p + "_url"); err != nil {
		return "", "", "", err
	}
	if apiKey, err = d.configOrEmpty(p + "_key"); err != nil {
		return "", "", "", err
	}
	return provider, baseURL, apiKey, nil
}

// IsCurrencyAvailable reports whether currency has an explorer provider
// configured. Any coin symbol is allowed — the Esplora adapter is universal — so
// availability is purely "has the operator configured an explorer for it".
func (d *Database) IsCurrencyAvailable(currency string) (bool, error) {
	provider, _, _, err := d.ExplorerConfig(currency)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(provider) != "", nil
}

// ConfiguredPaymentCurrencies returns the symbols of every coin that has an
// explorer provider configured (built-in BTC/LTC or any operator-added coin),
// extracted from the explorer_<coin>_provider config keys.
func (d *Database) ConfiguredPaymentCurrencies() ([]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query(`
		SELECT key FROM market_config
		WHERE key LIKE 'explorer\_%\_provider' ESCAPE '\' AND TRIM(value) <> ''
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list configured currencies: %w", err)
	}
	defer rows.Close() //nolint

	var out []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, fmt.Errorf("failed to scan currency key: %w", err)
		}
		coin := strings.TrimSuffix(strings.TrimPrefix(key, "explorer_"), "_provider")
		if coin != "" {
			out = append(out, strings.ToUpper(coin))
		}
	}
	return out, rows.Err()
}

// AvailableCurrencies returns the payment currencies that have an explorer
// configured — the canonical coins first (in display order), then any custom
// coins alphabetically.
func (d *Database) AvailableCurrencies() ([]string, error) {
	configured, err := d.ConfiguredPaymentCurrencies()
	if err != nil {
		return nil, err
	}
	set := make(map[string]bool, len(configured))
	for _, c := range configured {
		set[c] = true
	}
	var out []string
	for _, c := range protocol.PaymentCurrencies {
		if set[c] {
			out = append(out, c)
			delete(set, c)
		}
	}
	rest := make([]string, 0, len(set))
	for c := range set {
		rest = append(rest, c)
	}
	slices.Sort(rest)
	return append(out, rest...), nil
}

// GetConfig retrieves a configuration value by key.
func (d *Database) GetConfig(key string) (string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var value string
	err := d.db.QueryRow(`
		SELECT value FROM market_config WHERE key = ?
	`, key).Scan(&value)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("config key not found: %s", key)
		}
		return "", fmt.Errorf("failed to get config: %w", err)
	}

	return value, nil
}

// GetConfigInt retrieves a configuration value as an integer.
func (d *Database) GetConfigInt(key string) (int, error) {
	value, err := d.GetConfig(key)
	if err != nil {
		return 0, err
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("config value is not an integer: %s = %s", key, value)
	}

	return result, nil
}

// GetConfigFloat retrieves a configuration value as a float.
func (d *Database) GetConfigFloat(key string) (float64, error) {
	value, err := d.GetConfig(key)
	if err != nil {
		return 0, err
	}

	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("config value is not a float: %s = %s", key, value)
	}

	return result, nil
}

// SetConfig sets a configuration value.
func (d *Database) SetConfig(key, value string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now().UTC()
	_, err := d.db.Exec(`
		INSERT INTO market_config (key, value, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			updated_at = excluded.updated_at
	`, key, value, now)

	if err != nil {
		return fmt.Errorf("failed to set config: %w", err)
	}

	return nil
}

// GetAllConfig returns all configuration key-value pairs.
func (d *Database) GetAllConfig() (map[string]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query(`
		SELECT key, value FROM market_config
		ORDER BY key
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all config: %w", err)
	}
	defer rows.Close() //nolint

	config := make(map[string]string)
	for rows.Next() {
		var key, value string
		err := rows.Scan(&key, &value)
		if err != nil {
			return nil, fmt.Errorf("failed to scan config: %w", err)
		}
		config[key] = value
	}

	return config, nil
}

// GetMarketWallet returns the market's SKY wallet address for escrow.
func (d *Database) GetMarketWallet() (string, error) {
	return d.GetConfig("wallet_sky")
}

// GetMarketName returns the operator-set display name for this market, or an
// empty string if none is set.
func (d *Database) GetMarketName() string {
	name, err := d.GetConfig("market_name")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(name)
}

// GetFreezeViolationsLimit returns the number of violations before a user is banned.
func (d *Database) GetFreezeViolationsLimit() (int, error) {
	return d.GetConfigInt("freeze_violations_limit")
}

// GetBanDurationDays returns the ban duration in days.
func (d *Database) GetBanDurationDays() (int, error) {
	return d.GetConfigInt("ban_duration_days")
}

// BuyCancelLimitFloor is the smallest per-product buy-cancel limit allowed: a
// buyer always gets at least this many cancellations before being blocked.
const BuyCancelLimitFloor = 2

// GetBuyCancelLimit returns how many times a buyer may cancel (or let expire) an
// order for the same product before being blocked from re-buying it. It never
// returns less than BuyCancelLimitFloor, so a misconfiguration can't ban a buyer
// on their first cancellation.
func (d *Database) GetBuyCancelLimit() int {
	n, err := d.GetConfigInt("buy_cancel_limit")
	if err != nil || n < BuyCancelLimitFloor {
		return BuyCancelLimitFloor
	}
	return n
}

// GetListingExpiryMinutes returns how many minutes a pending listing is valid for.
func (d *Database) GetListingExpiryMinutes() (int, error) {
	return d.GetConfigInt("listing_expiry_minutes")
}

// GetOrderExpiryMinutes returns how many minutes a pending order is valid for.
func (d *Database) GetOrderExpiryMinutes() (int, error) {
	return d.GetConfigInt("order_expiry_minutes")
}

// GetConfirmationsRequired returns the blockchain confirmation depth a payment
// must reach before a trade completes (default 2).
func (d *Database) GetConfirmationsRequired() (int, error) {
	return d.GetConfigInt("confirmations_required")
}

// GetMinTradeSKY returns the minimum SKY amount a listing may sell (0 = no
// minimum beyond the round-to-zero check).
func (d *Database) GetMinTradeSKY() (float64, error) {
	return d.GetConfigFloat("min_trade_sky")
}

// GetMaxTradeSKY returns the maximum SKY amount a listing may sell (0 = no cap
// beyond the built-in float-safety ceiling).
func (d *Database) GetMaxTradeSKY() (float64, error) {
	return d.GetConfigFloat("max_trade_sky")
}

// GetCleanupDays returns how many days after completion orders are deleted.
func (d *Database) GetCleanupDays() (int, error) {
	return d.GetConfigInt("cleanup_days")
}

// GetCommissionRatePercent returns the market commission rate as a percent of the
// SKY sold (default 0.5 = 0.5%).
func (d *Database) GetCommissionRatePercent() (float64, error) {
	return d.GetConfigFloat("commission_rate_percent")
}

// GetCommissionMinSKY returns the minimum SKY commission per sale (floor), so
// tiny trades still pay something (default 0.001).
func (d *Database) GetCommissionMinSKY() (float64, error) {
	return d.GetConfigFloat("commission_min_sky")
}

// GetCommissionMaxSKY returns the maximum SKY commission per sale (cap). Zero (or
// negative) means no cap.
func (d *Database) GetCommissionMaxSKY() (float64, error) {
	return d.GetConfigFloat("commission_max_sky")
}
