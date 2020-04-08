/*
Package bip44 implements the bip44 spec https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
*/
package bip44

import (
	"errors"
	"fmt"

	"github.com/SkycoinProject/skycoin/src/cipher/bip32"
)

// Bip44's bip32 path: m / purpose' / coin_type' / account' / change / address_index

var (
	// ErrInvalidCoinType coin_type is >= 0x80000000
	ErrInvalidCoinType = errors.New("invalid bip44 coin type")

	// ErrInvalidAccount account is >= 0x80000000
	ErrInvalidAccount = errors.New("bip44 account index must be < 0x80000000")
)

// CoinType is the coin_type part of the bip44 path
type CoinType uint32

const (
	// CoinTypeBitcoin is the coin_type for Bitcoin
	CoinTypeBitcoin CoinType = 0
	// CoinTypeBitcoinTestnet is the coin_type for Skycoin
	CoinTypeBitcoinTestnet CoinType = 1
	// CoinTypeSkycoin is the coin_type for Skycoin
	CoinTypeSkycoin CoinType = 8000

	// ExternalChainIndex is the index of the external chain
	ExternalChainIndex uint32 = 0
	// ChangeChainIndex is the index of the change chain
	ChangeChainIndex uint32 = 1
)

// Coin is a bip32 node at the `coin_type` level of a bip44 path
type Coin struct {
	*bip32.PrivateKey
}

// NewCoin creates a bip32 node at the `coin_type` level of a bip44 path
func NewCoin(seed []byte, coinType CoinType) (*Coin, error) {
	if uint32(coinType) >= bip32.FirstHardenedChild {
		return nil, ErrInvalidCoinType
	}

	path := fmt.Sprintf("m/44'/%d'", coinType)
	pk, err := bip32.NewPrivateKeyFromPath(seed, path)
	if err != nil {
		return nil, err
	}

	return &Coin{
		pk,
	}, nil
}

// Account creates a bip32 node at the `account'` level of the bip44 path.
// The account number should be as it would appear in the path string, without
// the apostrophe that indicates hardening
func (c *Coin) Account(account uint32) (*Account, error) {
	if account >= bip32.FirstHardenedChild {
		return nil, ErrInvalidAccount
	}

	pk, err := c.NewPrivateChildKey(account + bip32.FirstHardenedChild)
	if err != nil {
		return nil, err
	}

	return &Account{
		pk,
	}, nil
}

// Account is a bip32 node at the `account` level of a bip44 path
type Account struct {
	*bip32.PrivateKey
}

// External returns the external chain node, to be used for receiving coins
func (a *Account) External() (*bip32.PrivateKey, error) {
	return a.NewPrivateChildKey(ExternalChainIndex)
}

// Change returns the change chain node, to be used for change addresses
func (a *Account) Change() (*bip32.PrivateKey, error) {
	return a.NewPrivateChildKey(ChangeChainIndex)
}

// Clone clones the account
func (a *Account) Clone() Account {
	na := Account{}
	key := a.PrivateKey.Clone()
	na.PrivateKey = &key
	return na
}
