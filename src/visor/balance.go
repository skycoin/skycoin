package visor

import (
    "github.com/skycoin/skycoin/src/coin"
    "log"
)

type BalancePair struct {
    Confirmed Balance `json:"confirmed"`
    Predicted Balance `json:"predicted"`
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

func (self Balance) Add(other Balance) Balance {
    return Balance{
        Coins: self.Coins + other.Coins,
        Hours: self.Hours + other.Hours,
    }
}

// Subtracts other from self and returns the new Balance.  Will panic if
// other is greater than balance, because Coins and Hours are unsigned.
func (self Balance) Sub(other Balance) Balance {
    if other.Coins > self.Coins || other.Hours > self.Hours {
        log.Panic("Cannot subtract balances, second balance is too large")
    }
    return Balance{
        Coins: self.Coins - other.Coins,
        Hours: self.Hours - other.Hours,
    }
}

func (self Balance) Equals(other Balance) bool {
    return self.Coins == other.Coins && self.Hours == other.Hours
}

func (self Balance) IsZero() bool {
    return self.Coins == 0 && self.Hours == 0
}
