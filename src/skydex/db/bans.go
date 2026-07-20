// Package db internal/db/bans.go
package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CreateFreezeViolation records a freeze violation for a buyer.
func (d *Database) CreateFreezeViolation(buyerPubKey, orderID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	violationID := uuid.New().String()
	now := time.Now().UTC()

	_, err := d.db.Exec(`
		INSERT INTO freeze_violations (id, buyer_pubkey, order_id, created_at)
		VALUES (?, ?, ?, ?)
	`, violationID, buyerPubKey, orderID, now)

	if err != nil {
		return fmt.Errorf("failed to create freeze violation: %w", err)
	}

	return nil
}

// CountRecentViolations counts the number of violations for a buyer within a time window.
func (d *Database) CountRecentViolations(buyerPubKey string, since time.Time) (int, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var count int
	err := d.db.QueryRow(`
		SELECT COUNT(*) FROM freeze_violations
		WHERE buyer_pubkey = ? AND created_at >= ?
	`, buyerPubKey, since).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to count violations: %w", err)
	}

	return count, nil
}

// GetBan retrieves a ban record for a user.
func (d *Database) GetBan(pubkey string) (*Ban, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	ban := &Ban{}
	err := d.db.QueryRow(`
		SELECT pubkey, violations, ban_until, created_at
		FROM bans
		WHERE pubkey = ?
	`, pubkey).Scan(&ban.PubKey, &ban.Violations, &ban.BanUntil, &ban.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not banned
		}
		return nil, fmt.Errorf("failed to get ban: %w", err)
	}

	return ban, nil
}

// IsUserBanned checks if a user is currently banned.
func (d *Database) IsUserBanned(pubkey string) (bool, error) {
	ban, err := d.GetBan(pubkey)
	if err != nil {
		return false, err
	}
	if ban == nil {
		return false, nil
	}

	// Check if ban is still active
	return time.Now().UTC().Before(ban.BanUntil), nil
}

// CreateBan creates or updates a ban record for a user.
func (d *Database) CreateBan(pubkey string, violations int, banUntil time.Time) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now().UTC()
	_, err := d.db.Exec(`
		INSERT INTO bans (pubkey, violations, ban_until, created_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(pubkey) DO UPDATE SET
			violations = excluded.violations,
			ban_until = excluded.ban_until,
			created_at = excluded.created_at
	`, pubkey, violations, banUntil, now)

	if err != nil {
		return fmt.Errorf("failed to create ban: %w", err)
	}

	return nil
}

// DeleteBan removes a ban record for a user.
func (d *Database) DeleteBan(pubkey string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec(`DELETE FROM bans WHERE pubkey = ?`, pubkey)
	if err != nil {
		return fmt.Errorf("failed to delete ban: %w", err)
	}

	return nil
}

// GetExpiredBans returns all bans that have expired (ban_until <= now).
func (d *Database) GetExpiredBans() ([]*Ban, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	now := time.Now().UTC()
	rows, err := d.db.Query(`
		SELECT pubkey, violations, ban_until, created_at
		FROM bans
		WHERE ban_until <= ?
	`, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired bans: %w", err)
	}
	defer rows.Close() //nolint

	var bans []*Ban
	for rows.Next() {
		ban := &Ban{}
		err := rows.Scan(&ban.PubKey, &ban.Violations, &ban.BanUntil, &ban.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ban: %w", err)
		}
		bans = append(bans, ban)
	}

	return bans, nil
}

// DeleteOldViolations removes freeze violations older than the specified duration.
func (d *Database) DeleteOldViolations(age time.Duration) (int64, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	cutoff := time.Now().UTC().Add(-age)
	result, err := d.db.Exec(`
		DELETE FROM freeze_violations
		WHERE created_at <= ?
	`, cutoff)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old violations: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// GetActiveBans returns bans that are still in effect (ban_until in the future),
// soonest-to-expire first. Used by the market operator UI.
func (d *Database) GetActiveBans() ([]*Ban, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	now := time.Now().UTC()
	rows, err := d.db.Query(`
		SELECT pubkey, violations, ban_until, created_at
		FROM bans
		WHERE ban_until > ?
		ORDER BY ban_until ASC
	`, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get active bans: %w", err)
	}
	defer rows.Close() //nolint

	var bans []*Ban
	for rows.Next() {
		ban := &Ban{}
		if err := rows.Scan(&ban.PubKey, &ban.Violations, &ban.BanUntil, &ban.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan ban: %w", err)
		}
		bans = append(bans, ban)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate bans: %w", err)
	}

	return bans, nil
}

// GetBannablePubKeys returns buyer public keys that have accrued at least limit
// freeze violations since the given time and are not already banned.
func (d *Database) GetBannablePubKeys(since time.Time, limit int) ([]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query(`
		SELECT fv.buyer_pubkey
		FROM freeze_violations fv
		WHERE fv.created_at >= ?
		GROUP BY fv.buyer_pubkey
		HAVING COUNT(*) >= ?
		   AND fv.buyer_pubkey NOT IN (SELECT pubkey FROM bans WHERE ban_until > ?)
	`, since, limit, time.Now().UTC())
	if err != nil {
		return nil, fmt.Errorf("failed to get bannable pubkeys: %w", err)
	}
	defer rows.Close() //nolint

	var pks []string
	for rows.Next() {
		var pk string
		if err := rows.Scan(&pk); err != nil {
			return nil, fmt.Errorf("failed to scan bannable pubkey: %w", err)
		}
		pks = append(pks, pk)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate bannable pubkeys: %w", err)
	}

	return pks, nil
}
