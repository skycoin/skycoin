package transaction

import (
	"bytes"
	"errors"
	"sort"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/params"
	"github.com/SkycoinProject/skycoin/src/util/fee"
)

var (
	// ErrInsufficientBalance is returned if a wallet does not have enough balance for a spend
	ErrInsufficientBalance = NewError(errors.New("balance is not sufficient"))
	// ErrInsufficientHours is returned if a wallet does not have enough hours for a spend with requested hours
	ErrInsufficientHours = NewError(errors.New("hours are not sufficient"))
	// ErrZeroSpend is returned if a transaction is trying to spend 0 coins
	ErrZeroSpend = NewError(errors.New("zero spend amount"))
	// ErrNoUnspents is returned if a Create is called with no unspent outputs
	ErrNoUnspents = NewError(errors.New("no unspents to spend"))
)

// UxBalance is an intermediate representation of a UxOut for sorting and spend choosing
type UxBalance struct {
	Hash           cipher.SHA256
	BkSeq          uint64
	Time           uint64
	Address        cipher.Address
	Coins          uint64
	InitialHours   uint64
	Hours          uint64
	SrcTransaction cipher.SHA256
}

// NewUxBalances converts coin.UxArray to []UxBalance. headTime is required to calculate coin hours.
func NewUxBalances(uxa coin.UxArray, headTime uint64) ([]UxBalance, error) {
	uxb := make([]UxBalance, len(uxa))
	for i, ux := range uxa {
		b, err := NewUxBalance(headTime, ux)
		if err != nil {
			return nil, err
		}
		uxb[i] = b
	}

	return uxb, nil
}

// NewUxBalance converts coin.UxOut to UxBalance. headTime is required to calculate coin hours.
func NewUxBalance(headTime uint64, ux coin.UxOut) (UxBalance, error) {
	hours, err := ux.CoinHours(headTime)
	if err != nil {
		return UxBalance{}, err
	}

	return UxBalance{
		Hash:           ux.Hash(),
		BkSeq:          ux.Head.BkSeq,
		Time:           ux.Head.Time,
		Address:        ux.Body.Address,
		Coins:          ux.Body.Coins,
		InitialHours:   ux.Body.Hours,
		Hours:          hours,
		SrcTransaction: ux.Body.SrcTransaction,
	}, nil
}

func uxBalancesSub(a, b []UxBalance) []UxBalance {
	var x []UxBalance

	bMap := make(map[cipher.SHA256]struct{}, len(b))
	for _, i := range b {
		bMap[i.Hash] = struct{}{}
	}

	for _, i := range a {
		if _, ok := bMap[i.Hash]; !ok {
			x = append(x, i)
		}
	}

	return x
}

// ChooseSpendsMinimizeUxOuts chooses uxout spends to satisfy an amount, using the least number of uxouts
//     -- PRO: Allows more frequent spending, less waiting for confirmations, useful for exchanges.
//     -- PRO: When transaction is volume is higher, transactions are prioritized by fee/size. Minimizing uxouts minimizes size.
//     -- CON: Would make the unconfirmed pool grow larger.
// Users with high transaction frequency will want to use this so that they will not need to wait as frequently
// for unconfirmed spends to complete before sending more.
// Alternatively, or in addition to this, they should batch sends into single transactions.
func ChooseSpendsMinimizeUxOuts(uxa []UxBalance, coins, hours uint64) ([]UxBalance, error) {
	return ChooseSpends(uxa, coins, hours, sortSpendsCoinsHighToLow)
}

// sortSpendsCoinsHighToLow sorts uxout spends with highest balance to lowest
func sortSpendsCoinsHighToLow(uxa []UxBalance) {
	sort.Slice(uxa, makeCmpUxOutByCoins(uxa, func(a, b uint64) bool {
		return a > b
	}))
}

// ChooseSpendsMaximizeUxOuts chooses uxout spends to satisfy an amount, using the most number of uxouts
// See the pros and cons of ChooseSpendsMinimizeUxOuts.
// This should be the default mode, because this keeps the unconfirmed pool smaller which will allow
// the network to scale better.
func ChooseSpendsMaximizeUxOuts(uxa []UxBalance, coins, hours uint64) ([]UxBalance, error) {
	return ChooseSpends(uxa, coins, hours, sortSpendsCoinsLowToHigh)
}

// sortSpendsCoinsLowToHigh sorts uxout spends with lowest balance to highest
func sortSpendsCoinsLowToHigh(uxa []UxBalance) {
	sort.Slice(uxa, makeCmpUxOutByCoins(uxa, func(a, b uint64) bool {
		return a < b
	}))
}

// sortSpendsHoursLowToHigh sorts uxout spends with lowest hours to highest
func sortSpendsHoursLowToHigh(uxa []UxBalance) {
	sort.Slice(uxa, makeCmpUxOutByHours(uxa, func(a, b uint64) bool {
		return a < b
	}))
}

func makeCmpUxOutByCoins(uxa []UxBalance, coinsCmp func(a, b uint64) bool) func(i, j int) bool {
	// Sort by:
	// coins highest or lowest depending on coinsCmp
	//  hours lowest
	//   oldest first
	//    tie break with hash comparison
	return func(i, j int) bool {
		a := uxa[i]
		b := uxa[j]

		if a.Coins == b.Coins {
			if a.Hours == b.Hours {
				if a.BkSeq == b.BkSeq {
					return cmpUxBalanceByUxID(a, b)
				}
				return a.BkSeq < b.BkSeq
			}
			return a.Hours < b.Hours
		}
		return coinsCmp(a.Coins, b.Coins)
	}
}

func makeCmpUxOutByHours(uxa []UxBalance, hoursCmp func(a, b uint64) bool) func(i, j int) bool {
	// Sort by:
	// hours highest or lowest depending on hoursCmp
	//  coins lowest
	//   oldest first
	//    tie break with hash comparison
	return func(i, j int) bool {
		a := uxa[i]
		b := uxa[j]

		if a.Hours == b.Hours {
			if a.Coins == b.Coins {
				if a.BkSeq == b.BkSeq {
					return cmpUxBalanceByUxID(a, b)
				}
				return a.BkSeq < b.BkSeq
			}
			return a.Coins < b.Coins
		}
		return hoursCmp(a.Hours, b.Hours)
	}
}

func cmpUxBalanceByUxID(a, b UxBalance) bool {
	cmp := bytes.Compare(a.Hash[:], b.Hash[:])
	if cmp == 0 {
		logger.Panic("Duplicate UxOut when sorting")
	}
	return cmp < 0
}

// ChooseSpends chooses uxouts from a list of uxouts.
// It first chooses the uxout with the most number of coins that has nonzero coinhours.
// It then chooses uxouts with zero coinhours, ordered by sortStrategy
// It then chooses remaining uxouts with nonzero coinhours, ordered by sortStrategy
func ChooseSpends(uxa []UxBalance, coins, hours uint64, sortStrategy func([]UxBalance)) ([]UxBalance, error) {
	if coins == 0 {
		return nil, ErrZeroSpend
	}

	if len(uxa) == 0 {
		return nil, ErrNoUnspents
	}

	for _, ux := range uxa {
		if ux.Coins == 0 {
			logger.Panic("UxOut coins are 0, can't spend")
			return nil, errors.New("UxOut coins are 0, can't spend")
		}
	}

	// Split UxBalances into those with and without hours
	var nonzero, zero []UxBalance
	for _, ux := range uxa {
		if ux.Hours == 0 {
			zero = append(zero, ux)
		} else {
			nonzero = append(nonzero, ux)
		}
	}

	// Abort if there are no uxouts with non-zero coinhours, they can't be spent yet
	if len(nonzero) == 0 {
		return nil, fee.ErrTxnNoFee
	}

	// Sort uxouts with hours lowest to highest and coins highest to lowest
	sortSpendsCoinsHighToLow(nonzero)

	var haveCoins uint64
	var haveHours uint64
	var spending []UxBalance

	// Use the first nonzero output. This output will have the least hours possible
	firstNonzero := nonzero[0]
	if firstNonzero.Hours == 0 {
		logger.Panic("balance has zero hours unexpectedly")
		return nil, errors.New("balance has zero hours unexpectedly")
	}

	nonzero = nonzero[1:]

	spending = append(spending, firstNonzero)

	haveCoins += firstNonzero.Coins
	haveHours += firstNonzero.Hours

	if haveCoins >= coins && fee.RemainingHours(haveHours, params.UserVerifyTxn.BurnFactor) >= hours {
		return spending, nil
	}

	// Sort uxouts without hours according to the sorting strategy
	sortStrategy(zero)

	for _, ux := range zero {
		spending = append(spending, ux)

		haveCoins += ux.Coins
		haveHours += ux.Hours

		if haveCoins >= coins {
			break
		}
	}

	if haveCoins >= coins && fee.RemainingHours(haveHours, params.UserVerifyTxn.BurnFactor) >= hours {
		return spending, nil
	}

	// Sort remaining uxouts with hours according to the sorting strategy
	sortStrategy(nonzero)

	for _, ux := range nonzero {
		spending = append(spending, ux)

		haveCoins += ux.Coins
		haveHours += ux.Hours

		if haveCoins >= coins && fee.RemainingHours(haveHours, params.UserVerifyTxn.BurnFactor) >= hours {
			return spending, nil
		}
	}

	if haveCoins < coins {
		return nil, ErrInsufficientBalance
	}

	return nil, ErrInsufficientHours
}
