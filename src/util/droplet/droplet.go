package droplet

import (
	"errors"
	"fmt"
	"math"

	logging "github.com/op/go-logging"
	"github.com/shopspring/decimal"
)

const (
	// Exponent is the number of decimal places held by droplets
	Exponent = 6
	// Multiplier is how much to multiply coins by to get droplets
	Multiplier = 1e6
)

var (
	// ErrNegativeValue is returned if a balance string is a negative number
	ErrNegativeValue = errors.New("Droplet string conversion failed: Negative balance")
	// ErrTooManyDecimals is returned if a balance string has more than 6 decimal places
	ErrTooManyDecimals = errors.New("Droplet string conversion failed: Too many decimal places")
	// ErrTooLarge is returned if a balance string is greater than math.MaxInt64
	ErrTooLarge = errors.New("Droplet string conversion failed: Value is too large")

	logger     = logging.MustGetLogger("convert")
	maxDecimal decimal.Decimal
)

func init() {
	max, err := decimal.NewFromString(fmt.Sprint(math.MaxInt64))
	if err != nil {
		panic(err)
	}

	maxDecimal = max
}

// FromString converts a skycoin balance string with decimal places to uint64 droplets.
// For example, "123.000456" becomes 123000456
func FromString(b string) (uint64, error) {
	d, err := decimal.NewFromString(b)
	if err != nil {
		return 0, err
	}

	// Values must be zero or positive
	if d.Sign() == -1 {
		return 0, ErrNegativeValue
	}

	// Skycoins have a maximum of 6 decimal places
	if d.Exponent() < -Exponent {
		return 0, ErrTooManyDecimals
	}

	// Multiply the coin balance by 1e6 to obtain droplets amount
	e := d.Mul(decimal.New(1, Exponent))

	// Check that there are no decimal places remaining. This error should not
	// occur, because of the earlier check of Exponent()
	if e.Exponent() < 0 {
		logger.Critical("Balance still has decimals after converting to droplets: %s", b)
		return 0, ErrTooManyDecimals
	}

	// Values greater than math.MaxInt64 will overflow after conversion to int64
	// using decimal.IntPart()
	if e.GreaterThan(maxDecimal) {
		return 0, ErrTooLarge
	}

	return uint64(e.IntPart()), nil
}

// ToString converts droplets to a skycoin balance fixed-point decimal string.
// String will always have a decimal precision of droplet.Exponent (6).
// For example, 123000456 becomes "123.000456" and
// 123000000 becomes "123.000000".
func ToString(n uint64) (string, error) {
	if n > math.MaxInt64 {
		return "", ErrTooLarge
	}

	d := decimal.New(int64(n), -Exponent)

	return d.StringFixed(Exponent), nil
}
