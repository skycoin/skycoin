// Package params src/params/droplet.go c1-coin-chain
package params

import (
	"errors"

	"github.com/skycoin/skycoin/src/util/droplet"
)

var (
	// ErrInvalidDecimals is returned by DropletPrecisionCheck if a coin amount has an invalid number of decimal places
	ErrInvalidDecimals = errors.New("invalid amount, too many decimal places")
)

// DropletPrecisionToDivisor converts number of allowed decimal places to the modulus divisor used when checking droplet precision rules
func DropletPrecisionToDivisor(precision uint8) (uint64, error) {
	if precision > droplet.Exponent {
		return 0, ErrInvalidMaxDropletPrecision
	}

	n := droplet.Exponent - precision
	var i uint64 = 1
	for k := uint8(0); k < n; k++ {
		i = i * 10
	}
	return i, nil
}

// DropletPrecisionCheck checks if an amount of coins is valid given decimal place restrictions
func DropletPrecisionCheck(precision uint8, amount uint64) error {
	divisor, err := DropletPrecisionToDivisor(precision)
	if err != nil {
		return err
	}
	if amount%divisor != 0 {
		return ErrInvalidDecimals
	}
	return nil
}
