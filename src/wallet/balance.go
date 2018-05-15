package wallet

import (
	"github.com/skycoin/skycoin/src/coin"
)

/*
Do not show balances or outputs that have not cleared yet
- should only allow spends against outputs that are on head
*/

// BalancePair records the confirmed and predicted balance
type BalancePair struct {
	Confirmed Balance `json:"confirmed"`
	Predicted Balance `json:"predicted"` //do "pending"
}

// AddressBalance represents a map of address balances
type AddressBalance map[string]BalancePair

// Balance is consisted of Coins and Hours
type Balance struct {
	Coins uint64 `json:"coins"`
	Hours uint64 `json:"hours"`
}

// NewBalance creates balance
func NewBalance(coins, hours uint64) Balance {
	return Balance{
		Coins: coins,
		Hours: hours,
	}
}

// NewBalanceFromUxOut creates Balance from UxOut
func NewBalanceFromUxOut(headTime uint64, ux *coin.UxOut) (Balance, error) {
	hours, err := ux.CoinHours(headTime)
	if err != nil {
		return Balance{}, err
	}

	return Balance{
		Coins: ux.Body.Coins,
		Hours: hours,
	}, nil
}

// Add adds two Balances
func (bal Balance) Add(other Balance) (Balance, error) {
	coins, err := coin.AddUint64(bal.Coins, other.Coins)
	if err != nil {
		return Balance{}, err
	}

	hours, err := coin.AddUint64(bal.Hours, other.Hours)
	if err != nil {
		return Balance{}, err
	}

	return Balance{
		Coins: coins,
		Hours: hours,
	}, nil
}

// Sub subtracts other from self and returns the new Balance.  Will panic if
// other is greater than balance, because Coins and Hours are unsigned.
func (bal Balance) Sub(other Balance) Balance {
	if other.Coins > bal.Coins || other.Hours > bal.Hours {
		logger.Panic("Cannot subtract balances, second balance is too large")
	}
	return Balance{
		Coins: bal.Coins - other.Coins,
		Hours: bal.Hours - other.Hours,
	}
}

// Equals compares two Balances
func (bal Balance) Equals(other Balance) bool {
	return bal.Coins == other.Coins && bal.Hours == other.Hours
}

// IsZero returns true if the Balance is empty (both coins and hours)
func (bal Balance) IsZero() bool {
	return bal.Coins == 0 && bal.Hours == 0
}
