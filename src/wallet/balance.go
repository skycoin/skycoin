package wallet

import (
	"log"

	"github.com/skycoin/skycoin/src/coin"
)

/*
Do not show balances or outputs that have not cleared yet
- should only allow spends against outputs that are on head
*/

type BalancePair struct {
	Confirmed Balance `json:"confirmed"`
	Predicted Balance `json:"predicted"` //do "pending"
}

type Balance struct {
	Coins uint64 `json:"coins"`
	Hours uint64 `json:"hours"`
}

func NewBalance(coins, hours uint64) Balance {
	return Balance{
		Coins: coins,
		Hours: hours,
	}
}

func NewBalanceFromUxOut(headTime uint64, ux *coin.UxOut) Balance {
	return Balance{
		Coins: ux.Body.Coins,
		Hours: ux.CoinHours(headTime),
	}
}

// Deprecate
func (self Balance) Add(other Balance) Balance {
	return Balance{
		Coins: self.Coins + other.Coins,
		Hours: self.Hours + other.Hours,
	}
}

// Subtracts other from self and returns the new Balance.  Will panic if
// other is greater than balance, because Coins and Hours are unsigned.
// Deprecate
func (self Balance) Sub(other Balance) Balance {
	if other.Coins > self.Coins || other.Hours > self.Hours {
		log.Panic("Cannot subtract balances, second balance is too large")
	}
	return Balance{
		Coins: self.Coins - other.Coins,
		Hours: self.Hours - other.Hours,
	}
}

// Deprecate
func (self Balance) Equals(other Balance) bool {
	return self.Coins == other.Coins && self.Hours == other.Hours
}

// Deprecate
func (self Balance) IsZero() bool {
	return self.Coins == 0 && self.Hours == 0
}
