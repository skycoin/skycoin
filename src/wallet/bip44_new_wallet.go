// This file is still under refactoring, once it gets done, we will replace the
// old bip44_wallet.go.

package wallet

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/SkycoinProject/skycoin/src/cipher"

	"github.com/SkycoinProject/skycoin/src/cipher/bip32"
	"github.com/SkycoinProject/skycoin/src/cipher/bip39"
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
)

var (
	// defaultBip44WalletDecoder is the default bip44 wallet decoder
	defaultBip44WalletDecoder = &Bip44WalletJSONDecoder{}
)

// Bip44WalletNew manages keys using the original Skycoin deterministic
// keypair generator method.
type Bip44WalletNew struct {
	// Meta wallet meta data
	Meta
	// accounts bip44 wallet accounts
	accounts accountManager
	// decoder is used to encode/decode bip44 wallet to/from []byte
	decoder Bip44WalletDecoder
}

// accountManager is the interface that manages the bip44 wallet accounts.
type accountManager interface {
	// New creates a new account, returns the account index, and error, if any
	new(opts bip44AccountCreateOptions) (uint32, error)
	// NewAddresses generates addresses on selected account
	newAddresses(index, chain, num uint32) ([]cipher.Addresser, error)
	// Len returns the account number
	len() uint32
	// Clone returns a deep clone accounts manager
	clone() accountManager
	// PackSecrets packs secrets
	packSecrets(ss Secrets)
	// UnpackSecrets unpacks secrets
	unpackSecrets(ss Secrets) error
	// Erase erase secrets
	erase()
}

// Bip44WalletDecoder is the interface that wraps the Encode and Decode methods.
//
// Encode method encodes the wallet to bytes, Decode method decodes bytes to bip44 wallet.
type Bip44WalletDecoder interface {
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
	WalletDecoder  Bip44WalletDecoder
}

// NewBip44WalletNew create a bip44 wallet with options
func NewBip44WalletNew(opts Bip44WalletCreateOptions) (*Bip44WalletNew, error) {
	wlt := &Bip44WalletNew{
		Meta: Meta{
			metaType:           WalletTypeBip44,
			metaFilename:       opts.Filename,
			metaVersion:        Version,
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
		wlt.decoder = defaultBip44WalletDecoder
	}

	bip44CoinType := resolveCoinAdapter(opts.CoinType).Bip44CoinType()
	wlt.Meta.setBip44Coin(bip44CoinType)

	if err := bip44MetaValidate(wlt.Meta); err != nil {
		return nil, err
	}
	return wlt, nil
}

func bip44MetaValidate(m Meta) error {
	if fn := m[metaFilename]; fn == "" {
		return errors.New("filename not set")
	}

	if tm := m[metaTimestamp]; tm != "" {
		_, err := strconv.ParseInt(tm, 10, 64)
		if err != nil {
			return errors.New("invalid timestamp")
		}
	}

	walletType, ok := m[metaType]
	if !ok {
		return errors.New("type field not set")
	}

	if walletType != WalletTypeBip44 {
		return ErrInvalidWalletType
	}

	if coinType := m[metaCoin]; coinType == "" {
		return errors.New("coin field not set")
	}

	// var isEncrypted bool
	// if encStr, ok := m[metaEncrypted]; ok {
	// 	// validate the encrypted value
	// 	var err error
	// 	isEncrypted, err = strconv.ParseBool(encStr)
	// 	if err != nil {
	// 		return errors.New("encrypted field is not a valid bool")
	// 	}
	// }

	// if isEncrypted {
	// cryptoType, ok := m[metaCryptoType]
	// if !ok {
	// 	return errors.New("crypto type field not set")
	// }

	// if _, err := getCrypto(CryptoType(cryptoType)); err != nil {
	// 	return errors.New("unknown crypto type")
	// }

	// if s := m[metaSecrets]; s == "" {
	// 	return errors.New("wallet is encrypted, but secrets field not set")
	// }

	// if s := m[metaSeed]; s != "" {
	// 	return errors.New("seed should not be visible in encrypted wallets")
	// }

	// if s := m[metaLastSeed]; s != "" {
	// 	return errors.New("lastSeed should not be visible in encrypted wallets")
	// }
	// } else {
	// if s := m[metaSecrets]; s != "" {
	// 	return errors.New("secrets should not be in unencrypted wallets")
	// }

	// bip44 wallet seeds must be a valid bip39 mnemonic
	if s := m[metaSeed]; s == "" {
		return errors.New("seed missing in unencrypted bip44 wallet")
	} else if err := bip39.ValidateMnemonic(s); err != nil {
		return err
	}
	// }

	if s := m[metaBip44Coin]; s == "" {
		return errors.New("bip44Coin missing")
	} else if _, err := strconv.ParseUint(s, 10, 32); err != nil {
		return fmt.Errorf("bip44Coin invalid: %v", err)
	}

	return nil
}

// NewAccount create a bip44 wallet account, returns account index and
// error, if any.
func (w *Bip44WalletNew) NewAccount(name string) (uint32, error) {
	opts := bip44AccountCreateOptions{
		name:           name,
		seed:           w.Meta.Seed(),
		seedPassphrase: w.Meta.SeedPassphrase(),
		coinType:       w.Meta.Coin(),
	}

	return w.accounts.new(opts)
}

// NewAddresses creates addresses
func (w *Bip44WalletNew) NewAddresses(account, chain, n uint32) ([]cipher.Addresser, error) {
	return w.accounts.newAddresses(account, chain, n)
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

// Serialize encodes the bip44 wallet to []byte
func (w *Bip44WalletNew) Serialize() ([]byte, error) {
	if w.decoder == nil {
		w.decoder = defaultBip44WalletDecoder
	}
	return w.decoder.Encode(w)
}

// Deserialize decodes the []byte to a bip44 wallet
func (w *Bip44WalletNew) Deserialize(b []byte) error {
	if w.decoder == nil {
		w.decoder = defaultBip44WalletDecoder
	}
	toW, err := w.decoder.Decode(b)
	if err != nil {
		return err
	}

	toW.decoder = w.decoder
	*w = *toW
	return nil
}
