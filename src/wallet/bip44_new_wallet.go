// This file is still under refactoring, once it gets done, we will replace the
// old bip44_wallet.go.

package wallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/util/mathutil"

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
	Accounts []*bip44Account
}

// Bip44WalletCreateOptions options for creating the bip44 wallet
type Bip44WalletCreateOptions struct {
	Filename       string
	Version        string
	Label          string
	Seed           string
	SeedPassphrase string
	Coin           CoinType
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
			metaCoin:           string(opts.Coin),
			metaTimestamp:      strconv.FormatInt(time.Now().Unix(), 10),
			metaEncrypted:      "false",
		},
	}

	bip44Coin := registeredCoinAdapters.get(opts.Coin).Bip44CoinType()
	wlt.Meta.setBip44Coin(bip44Coin)

	return wlt
}

// NewAccount create a bip44 wallet account, returns account index and
// error if any.
func (w *Bip44WalletNew) NewAccount(name string) (uint32, error) {
	accountIndex, err := w.nextAccountIndex()
	if err != nil {
		return 0, err
	}

	opts := bip44AccountCreateOptions{
		name:           name,
		index:          accountIndex,
		seed:           w.Meta.Seed(),
		seedPassphrase: w.Meta.SeedPassphrase(),
		coinType:       w.Meta.Coin(),
	}

	ba, err := newBip44Account(opts)
	if err != nil {
		return 0, err
	}

	w.Accounts = append(w.Accounts, ba)

	return ba.Index, nil
}

func (w *Bip44WalletNew) nextAccountIndex() (uint32, error) {
	if _, err := mathutil.AddUint32(uint32(len(w.Accounts)), 1); err != nil {
		return 0, errors.New("Maximum bip44 account number reached")
	}

	return uint32(len(w.Accounts)), nil
}

// NewAddresses creates addresses
func (w *Bip44WalletNew) NewAddresses(account, chain, n uint32) ([]cipher.Addresser, error) {
	a, err := w.account(account)
	if err != nil {
		return nil, err
	}
	return a.newAddresses(chain, n)
}

// account returns the wallet account
func (w *Bip44WalletNew) account(index uint32) (*bip44Account, error) {
	if index >= uint32(len(w.Accounts)) {
		return nil, fmt.Errorf("account of index %d does not exist", index)
	}
	if a := w.Accounts[index]; a != nil {
		return a, nil
	}

	return nil, fmt.Errorf("account  of index %d does not exist", index)
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
	// rw := ReadableBip44WalletNew{
	// 	Meta: w.Meta.clone(),
	// }
	// rw.Accounts = make([]ReadableBip44Account, len(w.accounts))
	// for i, a := range w.accounts {
	// 	rw.Accounts[i] = newReadableBip44Account(a)
	// }
	// return json.MarshalIndent(rw, "", "    ")
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
