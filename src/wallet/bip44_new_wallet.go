// This file is still under refactoring, once it gets done, we will replace the
// old bip44_wallet.go.

package wallet

import (
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
type Bip44WalletNew struct {
	Meta
	accounts accountManager
	decoder  WalletDecoder
}

type accountManager interface {
	// Len returns the account number
	Len() uint32
	// New creates a new account, returns the account index, and error, if any
	New(opts bip44AccountCreateOptions) (uint32, error)
	// NewAddresses generates addresses on selected account
	NewAddresses(index, chain, num uint32) ([]cipher.Addresser, error)
	// ToReadable converts the bip44 accounts to readable accounts with JSON tags
	ToReadable() ReadableBip44Accounts
}

// WalletDecoder is the interface that wraps the Encode and Decode methods.
//
// Encode method encodes the wallet to bytes, Decode method decodes
type WalletDecoder interface {
	Encode(w *Bip44WalletNew) ([]byte, error)
	Decode(b []byte) (*Bip44WalletNew, error)
}

// Bip44WalletCreateOptions options for creating the bip44 wallet
type Bip44WalletCreateOptions struct {
	Filename       string
	Version        string
	Label          string
	Seed           string
	SeedPassphrase string
	CoinType       CoinType
	WalletDecoder  WalletDecoder
}

// NewBip44WalletNew create a bip44 wallet base on options
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
		accounts: &bip44Accounts{},
		decoder:  opts.WalletDecoder,
	}

	if wlt.decoder == nil {
		// TODO:
		// wlt.decoder = defaultWalletDecoder
	}

	bip44CoinType := resolveCoinAdapter(opts.CoinType).Bip44CoinType()
	wlt.Meta.setBip44Coin(bip44CoinType)

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

// Serialize serializes the bip44 wallet to bytes
func (w *Bip44WalletNew) Serialize() ([]byte, error) {
	return w.decoder.Encode(w)
}

// Unserialize decode the bytes into
func (w *Bip44WalletNew) Unserialize(b []byte) error {
	if w.decoder == nil {
		// TODO: use default wallet decoder if wallet's decoder is nil
		// w.decoder = defaultWalletDecoder
	}
	toW, err := w.decoder.Decode(b)
	if err != nil {
		return err
	}

	toW.decoder = w.decoder
	*w = *toW
	return nil
}
