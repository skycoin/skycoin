// Package params src/params/init.go c1-coin-chain
package params

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/skycoin/skycoin/src/util/droplet"
)

// envErr holds what init found wrong with the USER_* overrides, reported by CheckEnv.
var envErr error

func init() {
	envErr = errors.Join(
		loadUserBurnFactor(),
		loadUserMaxTransactionSize(),
		loadUserMaxDecimals(),
		sanityCheck(),
	)
}

// CheckEnv reports an invalid USER_BURN_FACTOR, USER_MAX_TXN_SIZE or
// USER_MAX_DECIMALS. An invalid value leaves the default in place.
func CheckEnv() error {
	return envErr
}

func sanityCheck() error {
	if err := UserVerifyTxn.Validate(); err != nil {
		return fmt.Errorf("user verify txn params: %w", err)
	}
	if err := MainNetDistribution.Validate(); err != nil {
		return fmt.Errorf("mainnet distribution: %w", err)
	}
	return nil
}

func loadUserBurnFactor() error {
	xs := os.Getenv("USER_BURN_FACTOR")
	if xs == "" {
		return nil
	}

	x, err := strconv.ParseUint(xs, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid USER_BURN_FACTOR %q: %w", xs, err)
	}

	if x < uint64(MinBurnFactor) {
		return fmt.Errorf("USER_BURN_FACTOR must be >= %d", MinBurnFactor)
	}

	UserVerifyTxn.BurnFactor = uint32(x)
	return nil
}

func loadUserMaxTransactionSize() error {
	xs := os.Getenv("USER_MAX_TXN_SIZE")
	if xs == "" {
		return nil
	}

	x, err := strconv.ParseUint(xs, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid USER_MAX_TXN_SIZE %q: %w", xs, err)
	}

	if x < uint64(MinTransactionSize) {
		return fmt.Errorf("USER_MAX_TXN_SIZE must be >= %d", MinTransactionSize)
	}

	UserVerifyTxn.MaxTransactionSize = uint32(x)
	return nil
}

func loadUserMaxDecimals() error {
	xs := os.Getenv("USER_MAX_DECIMALS")
	if xs == "" {
		return nil
	}

	x, err := strconv.ParseUint(xs, 10, 8)
	if err != nil {
		return fmt.Errorf("invalid USER_MAX_DECIMALS %q: %w", xs, err)
	}

	if x > uint64(droplet.Exponent) {
		return fmt.Errorf("USER_MAX_DECIMALS must be <= %d", droplet.Exponent)
	}

	UserVerifyTxn.MaxDropletPrecision = uint8(x)
	return nil
}
