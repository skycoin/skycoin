package readable

import "github.com/skycoin/skycoin/src/wallet"

// Balance has coins and hours
type Balance struct {
	Coins uint64 `json:"coins"`
	Hours uint64 `json:"hours"`
}

// NewBalance copies from wallet.Balance
func NewBalance(b wallet.Balance) Balance {
	return Balance{
		Coins: b.Coins,
		Hours: b.Hours,
	}
}

// BalancePair records the confirmed and predicted balance of an address
type BalancePair struct {
	Confirmed Balance `json:"confirmed"`
	Predicted Balance `json:"predicted"` // TODO rename "pending"
}

// NewBalancePair copies from wallet.BalancePair
func NewBalancePair(bp wallet.BalancePair) BalancePair {
	return BalancePair{
		Confirmed: NewBalance(bp.Confirmed),
		Predicted: NewBalance(bp.Predicted),
	}
}

// AddressBalances represents a map of address balances
type AddressBalances map[string]BalancePair

// NewAddressBalances copies from wallet.AddressBalances
func NewAddressBalances(wab wallet.AddressBalances) AddressBalances {
	ab := make(AddressBalances, len(wab))
	for k, v := range wab {
		ab[k] = NewBalancePair(v)
	}
	return ab
}

// WalletEntry the wallet entry struct
type WalletEntry struct {
	Address string `json:"address"`
	Public  string `json:"public_key"`
}

// WalletMeta the wallet meta struct
type WalletMeta struct {
	Coin       wallet.CoinType   `json:"coin"`
	Filename   string            `json:"filename"`
	Label      string            `json:"label"`
	Type       string            `json:"type"`
	Version    string            `json:"version"`
	CryptoType wallet.CryptoType `json:"crypto_type"`
	Timestamp  int64             `json:"timestamp"`
	Encrypted  bool              `json:"encrypted"`
}
