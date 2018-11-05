package params

import (
	"errors"

	"github.com/skycoin/skycoin/src/util/droplet"
)

var (
	// ErrInvalidDecimals is returned by DropletPrecisionCheck if a coin amount has an invalid number of decimal places
	ErrInvalidDecimals = errors.New("invalid amount, too many decimal places")

	// maxDropletDivisor represents the modulus divisor when checking droplet precision rules.
	// It is computed from MaxDropletPrecision in init()
	maxDropletDivisor uint64
)

// MaxDropletDivisor represents the modulus divisor when checking droplet precision rules.
func MaxDropletDivisor() uint64 {
	// The value is wrapped in a getter to make it immutable to external packages
	return maxDropletDivisor
}

// DropletPrecisionCheck checks if an amount of coins is valid given decimal place restrictions
func DropletPrecisionCheck(amount uint64) error {
	if amount%maxDropletDivisor != 0 {
		return ErrInvalidDecimals
	}
	return nil
}

func calculateDivisor(precision uint64) uint64 {
	if precision > droplet.Exponent {
		panic("precision must be <= droplet.Exponent")
	}

	n := droplet.Exponent - precision
	var i uint64 = 1
	for k := uint64(0); k < n; k++ {
		i = i * 10
	}
	return i
}
