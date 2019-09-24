package transaction

import (
	"errors"

	"github.com/shopspring/decimal"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
)

// Error wraps transaction creation-related errors.
// It wraps errors caused by user input, but not errors caused by programmer input or internal issues.
type Error struct {
	error
}

// NewError creates an Error
func NewError(err error) error {
	if err == nil {
		return nil
	}
	return Error{err}
}

const (
	// HoursSelectionTypeManual is used to specify manual hours selection in advanced spend
	HoursSelectionTypeManual = "manual"
	// HoursSelectionTypeAuto is used to specify automatic hours selection in advanced spend
	HoursSelectionTypeAuto = "auto"

	// HoursSelectionModeShare will distribute coin hours equally amongst destinations
	HoursSelectionModeShare = "share"
)

var (
	// ErrNullChangeAddress ChangeAddress must not be the null address
	ErrNullChangeAddress = NewError(errors.New("ChangeAddress must not be the null address"))
	// ErrMissingReceivers To is required
	ErrMissingReceivers = NewError(errors.New("To is required"))
	// ErrZeroCoinsReceiver To.Coins must not be zero
	ErrZeroCoinsReceiver = NewError(errors.New("To.Coins must not be zero"))
	// ErrNullAddressReceiver To.Address must not be the null address
	ErrNullAddressReceiver = NewError(errors.New("To.Address must not be the null address"))
	// ErrDuplicateReceiver To contains duplicate values
	ErrDuplicateReceiver = NewError(errors.New("To contains duplicate values"))
	// ErrReceiverZeroHoursAuto To.Hours must be zero for auto type hours selection
	ErrReceiverZeroHoursAuto = NewError(errors.New("To.Hours must be zero for auto type hours selection"))
	// ErrMissingHoursSelectionModeAuto HoursSelection.Mode is required for auto type hours selection
	ErrMissingHoursSelectionModeAuto = NewError(errors.New("HoursSelection.Mode is required for auto type hours selection"))
	// ErrInvalidHoursSelelectionMode Invalid HoursSelection.Mode
	ErrInvalidHoursSelelectionMode = NewError(errors.New("Invalid HoursSelection.Mode"))
	// ErrInvalidHoursSelectionModeManual HoursSelection.Mode cannot be used for manual type hours selection
	ErrInvalidHoursSelectionModeManual = NewError(errors.New("HoursSelection.Mode cannot be used for manual type hours selection"))
	// ErrInvalidHoursSelectionType Invalid HoursSelection.Type
	ErrInvalidHoursSelectionType = NewError(errors.New("Invalid HoursSelection.Type"))
	// ErrMissingShareFactor HoursSelection.ShareFactor must be set for share mode
	ErrMissingShareFactor = NewError(errors.New("HoursSelection.ShareFactor must be set for share mode"))
	// ErrInvalidShareFactor HoursSelection.ShareFactor can only be used for share mode
	ErrInvalidShareFactor = NewError(errors.New("HoursSelection.ShareFactor can only be used for share mode"))
	// ErrShareFactorOutOfRange HoursSelection.ShareFactor must be >= 0 and <= 1
	ErrShareFactorOutOfRange = NewError(errors.New("HoursSelection.ShareFactor must be >= 0 and <= 1"))
)

// HoursSelection defines options for hours distribution
type HoursSelection struct {
	Type        string
	Mode        string
	ShareFactor *decimal.Decimal
}

// Params defines control parameters for transaction construction
type Params struct {
	HoursSelection HoursSelection
	To             []coin.TransactionOutput
	ChangeAddress  *cipher.Address
}

// Validate validates Params
func (c Params) Validate() error {
	if c.ChangeAddress != nil && c.ChangeAddress.Null() {
		return ErrNullChangeAddress
	}

	if len(c.To) == 0 {
		return ErrMissingReceivers
	}

	for _, to := range c.To {
		if to.Coins == 0 {
			return ErrZeroCoinsReceiver
		}

		if to.Address.Null() {
			return ErrNullAddressReceiver
		}
	}

	// Check for duplicate outputs, a transaction can't have outputs with
	// the same (address, coins, hours)
	// Auto mode would distribute hours to the outputs and could hypothetically
	// avoid assigning duplicate hours in many cases, but the complexity for doing
	// so is very high, so also reject duplicate (address, coins) for auto mode.
	outputs := make(map[coin.TransactionOutput]struct{}, len(c.To))
	for _, to := range c.To {
		outputs[to] = struct{}{}
	}

	if len(outputs) != len(c.To) {
		return ErrDuplicateReceiver
	}

	switch c.HoursSelection.Type {
	case HoursSelectionTypeAuto:
		for _, to := range c.To {
			if to.Hours != 0 {
				return ErrReceiverZeroHoursAuto
			}
		}

		switch c.HoursSelection.Mode {
		case HoursSelectionModeShare:
		case "":
			return ErrMissingHoursSelectionModeAuto
		default:
			return ErrInvalidHoursSelelectionMode
		}

	case HoursSelectionTypeManual:
		if c.HoursSelection.Mode != "" {
			return ErrInvalidHoursSelectionModeManual
		}

	default:
		return ErrInvalidHoursSelectionType
	}

	if c.HoursSelection.ShareFactor == nil {
		if c.HoursSelection.Mode == HoursSelectionModeShare {
			return ErrMissingShareFactor
		}
	} else {
		if c.HoursSelection.Mode != HoursSelectionModeShare {
			return ErrInvalidShareFactor
		}

		zero := decimal.New(0, 0)
		one := decimal.New(1, 0)
		if c.HoursSelection.ShareFactor.LessThan(zero) || c.HoursSelection.ShareFactor.GreaterThan(one) {
			return ErrShareFactorOutOfRange
		}
	}

	return nil
}
