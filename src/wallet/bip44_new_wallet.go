// This file is still under refactoring, once it gets done, we will replace the
// old bip44_wallet.go.

package wallet

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/SkycoinProject/skycoin/src/cipher"

	"github.com/SkycoinProject/skycoin/src/cipher/bip32"
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
)

const (
	// Bip44WalletVersion Bip44 wallet version
	Bip44WalletVersion = "0.4"
)

// Bip44WalletNew manages keys using the original Skycoin deterministic
// keypair generator method.
// With this generator, a single chain of addresses is created, each one dependent
// on the previous.
type Bip44WalletNew struct {
	Meta
	accounts accountManager
}

type accountManager interface {
	// Len returns the account number
	Len() uint32
	// New creates a new account, returns the account index, and error if any
	New(opts bip44AccountCreateOptions) (uint32, error)
	// NewAddresses generates addresses on account
	NewAddresses(index, chain, num uint32) ([]cipher.Addresser, error)
}

// Bip44WalletCreateOptions options for creating the bip44 wallet
type Bip44WalletCreateOptions struct {
	Filename       string
	Version        string
	Label          string
	Seed           string
	SeedPassphrase string
	CoinType       CoinType
}

// NewBip44WalletNew create a bip44 wallet base on options,
func NewBip44WalletNew(opts Bip44WalletCreateOptions) *Bip44WalletNew {
	wlt := &Bip44WalletNew{
		Meta: Meta{
			metaFilename:       opts.Filename,
			metaVersion:        Bip44WalletVersion,
			metaLabel:          opts.Label,
			metaSeed:           opts.Seed,
			metaSeedPassphrase: opts.SeedPassphrase,
			metaCoin:           string(opts.CoinType),
			metaTimestamp:      strconv.FormatInt(time.Now().Unix(), 10),
			metaEncrypted:      "false",
		},
	}

	ca := resolveCoinAdapter(opts.CoinType)
	wlt.Meta.setBip44Coin(ca.Bip44CoinType())

	return wlt
}

// NewAccount create a bip44 wallet account, returns account index and
// error if any.
func (w *Bip44WalletNew) NewAccount(name string) (uint32, error) {
	opts := bip44AccountCreateOptions{
		name:           name,
		seed:           w.Meta.Seed(),
		seedPassphrase: w.Meta.SeedPassphrase(),
		coinType:       w.Meta.Coin(),
	}

	return w.accounts.New(opts)
}

// NewAddresses creates addresses
func (w *Bip44WalletNew) NewAddresses(account, chain, n uint32) ([]cipher.Addresser, error) {
	return w.accounts.NewAddresses(account, chain, n)
}

func makeChainPubKeys(a *bip44.Account) (*bip32.PublicKey, *bip32.PublicKey, error) {
	external, err := a.NewPublicChildKey(0)
	if err != nil {
		return nil, nil, fmt.Errorf("create external chain public key failed: %v", err)
	}

	change, err := a.NewPublicChildKey(1)
	if err != nil {
		return nil, nil, fmt.Errorf("create change chain public key failed: %v", err)
	}
	return external, change, nil
}

// Serialize returns the JSON representation of the wallet
func (w *Bip44WalletNew) Serialize() ([]byte, error) {
	return json.MarshalIndent(w, "", "    ")
}

// Unserialize unserialize data to bip44 wallet
func Unserialize(data []byte) *Bip44WalletNew {
	return nil
}

// // ReadableBip44WalletNew readable bip44 wallet
// type ReadableBip44WalletNew struct {
// 	Meta     `json:"meta"`
// 	Accounts ReadableBip44Accounts `json:"accounts"`
// }
