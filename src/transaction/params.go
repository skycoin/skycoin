package transaction

import (
	"errors"

	"github.com/shopspring/decimal"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
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
	// ErrInsufficientBalance is returned if a wallet does not have enough balance for a spend
	ErrInsufficientBalance = NewError(errors.New("balance is not sufficient"))
	// ErrInsufficientHours is returned if a wallet does not have enough hours for a spend with requested hours
	ErrInsufficientHours = NewError(errors.New("hours are not sufficient"))
	// ErrZeroSpend is returned if a transaction is trying to spend 0 coins
	ErrZeroSpend = NewError(errors.New("zero spend amount"))
	// ErrInvalidHoursSelectionMode for invalid HoursSelection mode values
	ErrInvalidHoursSelectionMode = NewError(errors.New("invalid hours selection mode"))
	// ErrInvalidHoursSelectionType for invalid HoursSelection type values
	ErrInvalidHoursSelectionType = NewError(errors.New("invalid hours selection type"))
	// ErrNoUnspents is returned if a wallet has no unspents to spend
	ErrNoUnspents = NewError(errors.New("no unspents to spend"))
	// ErrNullChangeAddress ChangeAddress must not be the null address
	ErrNullChangeAddress = NewError(errors.New("ChangeAddress must not be the null address"))
	// ErrMissingTo To is required
	ErrMissingTo = NewError(errors.New("To is required"))
	// ErrZeroCoinsTo To.Coins must not be zero
	ErrZeroCoinsTo = NewError(errors.New("To.Coins must not be zero"))
	// ErrNullAddressTo To.Address must not be the null address
	ErrNullAddressTo = NewError(errors.New("To.Address must not be the null address"))
	// ErrDuplicateTo To contains duplicate values
	ErrDuplicateTo = NewError(errors.New("To contains duplicate values"))
	// ErrZeroToHoursAuto To.Hours must be zero for auto type hours selection
	ErrZeroToHoursAuto = NewError(errors.New("To.Hours must be zero for auto type hours selection"))
	// ErrMissingModeAuto HoursSelection.Mode is required for auto type hours selection
	ErrMissingModeAuto = NewError(errors.New("HoursSelection.Mode is required for auto type hours selection"))
	// ErrInvalidHoursSelMode Invalid HoursSelection.Mode
	ErrInvalidHoursSelMode = NewError(errors.New("Invalid HoursSelection.Mode"))
	// ErrInvalidModeManual HoursSelection.Mode cannot be used for manual type hours selection
	ErrInvalidModeManual = NewError(errors.New("HoursSelection.Mode cannot be used for manual type hours selection"))
	// ErrInvalidHoursSelType Invalid HoursSelection.Type
	ErrInvalidHoursSelType = NewError(errors.New("Invalid HoursSelection.Type"))
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
		return ErrMissingTo
	}

	for _, to := range c.To {
		if to.Coins == 0 {
			return ErrZeroCoinsTo
		}

		if to.Address.Null() {
			return ErrNullAddressTo
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
		return ErrDuplicateTo
	}

	switch c.HoursSelection.Type {
	case HoursSelectionTypeAuto:
		for _, to := range c.To {
			if to.Hours != 0 {
				return ErrZeroToHoursAuto
			}
		}

		switch c.HoursSelection.Mode {
		case HoursSelectionModeShare:
		case "":
			return ErrMissingModeAuto
		default:
			return ErrInvalidHoursSelMode
		}

	case HoursSelectionTypeManual:
		if c.HoursSelection.Mode != "" {
			return ErrInvalidModeManual
		}

	default:
		return ErrInvalidHoursSelType
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
