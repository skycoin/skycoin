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
func DropletPrecisionToDivisor(precision uint8) uint64 {
	if precision > droplet.Exponent {
		panic("precision must be <= droplet.Exponent")
	}

	n := droplet.Exponent - precision
	var i uint64 = 1
	for k := uint8(0); k < n; k++ {
		i = i * 10
	}
	return i
}

// DropletPrecisionCheck checks if an amount of coins is valid given decimal place restrictions
func DropletPrecisionCheck(precision uint8, amount uint64) error {
	if amount%DropletPrecisionToDivisor(precision) != 0 {
		return ErrInvalidDecimals
	}
	return nil
}
