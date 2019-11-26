package params

import (
	"errors"

	"github.com/SkycoinProject/skycoin/src/util/droplet"
)

const (
	// MinBurnFactor minimum value for BurnFactor
	MinBurnFactor uint32 = 2
	// MinTransactionSize minimum value for MaxTransactionSize
	MinTransactionSize uint32 = 1024
)

var (
	// ErrInvalidBurnFactor BurnFactor value is out of range
	ErrInvalidBurnFactor = errors.New("BurnFactor value is out of range")
	// ErrInvalidMaxTransactionSize MaxTransactionSize value is out of range
	ErrInvalidMaxTransactionSize = errors.New("MaxTransactionSize value is out of range")
	// ErrInvalidMaxDropletPrecision MaxDropletPrecision value is out of range
	ErrInvalidMaxDropletPrecision = errors.New("MaxDropletPrecision value is out of range")
)

// VerifyTxn are parameters for verifying a transaction
type VerifyTxn struct {
	// BurnFactor inverse fraction of coinhours that must be burned
	BurnFactor uint32
	// MaxTransactionSize maximum size of a transaction in bytes
	MaxTransactionSize uint32
	// MaxDropletPrecision maximum decimal precision of droplets
	MaxDropletPrecision uint8
}

// MaxDropletDivisor return the modulus divisor used when checking droplet precision rules
func (v VerifyTxn) MaxDropletDivisor() uint64 {
	return DropletPrecisionToDivisor(v.MaxDropletPrecision)
}

// Validate validates the configured parameters
func (v VerifyTxn) Validate() error {
	if v.BurnFactor < MinBurnFactor {
		return ErrInvalidBurnFactor
	}

	if v.MaxTransactionSize < MinTransactionSize {
		return ErrInvalidMaxTransactionSize
	}

	if v.MaxDropletPrecision > droplet.Exponent {
		return ErrInvalidMaxDropletPrecision
	}

	return nil
}
