/*
Package bip44 implements the bip44 spec https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
*/
package bip44

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher/bip32"
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

// Purpose constants for BIP43 purpose field
const (
	// PurposeBIP44 is the purpose for BIP44 wallets (P2PKH)
	PurposeBIP44 uint32 = 44
	// PurposeBIP84 is the purpose for BIP84 wallets (P2WPKH native segwit)
	PurposeBIP84 uint32 = 84
)

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

// NewCoin creates a bip32 node at the `coin_type` level of a bip44 path (m/44'/coin_type')
func NewCoin(seed []byte, coinType CoinType) (*Coin, error) {
	return NewCoinWithPurpose(seed, coinType, PurposeBIP44)
}

// NewCoinWithPurpose creates a bip32 node at the `coin_type` level with a specific BIP43 purpose.
// Use PurposeBIP44 (44) for standard BIP44 P2PKH addresses.
// Use PurposeBIP84 (84) for BIP84 native segwit (P2WPKH) addresses.
func NewCoinWithPurpose(seed []byte, coinType CoinType, purpose uint32) (*Coin, error) {
	if uint32(coinType) >= bip32.FirstHardenedChild {
		return nil, ErrInvalidCoinType
	}

	path := fmt.Sprintf("m/%d'/%d'", purpose, coinType)
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
	if a.PrivateKey != nil {
		key := a.PrivateKey.Clone()
		na.PrivateKey = &key
	}
	return na
}
