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
	// new creates a new account, returns the account index, and error, if any
	new(opts bip44AccountCreateOptions) (uint32, error)
	// newAddresses generates addresses on selected account
	newAddresses(index, chain, num uint32) ([]cipher.Addresser, error)
	// len returns the account number
	len() uint32
	// clone returns a deep clone accounts manager
	clone() accountManager
	// packSecrets packs secrets
	packSecrets(ss Secrets)
	// unpackSecrets unpacks secrets
	unpackSecrets(ss Secrets) error
	// erase wipes secrets
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
	CryptoType     CryptoType
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
			metaCryptoType:     string(opts.CryptoType),
			metaEncrypted:      "false",
		},
		accounts: &bip44Accounts{},
		decoder:  opts.WalletDecoder,
	}

	if wlt.Meta.CryptoType() == "" {
		wlt.Meta[metaCryptoType] = string(DefaultCryptoType)
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
		return errors.New("Filename not set")
	}

	if tm := m[metaTimestamp]; tm != "" {
		_, err := strconv.ParseInt(tm, 10, 64)
		if err != nil {
			return errors.New("Invalid timestamp")
		}
	}

	walletType, ok := m[metaType]
	if !ok {
		return errors.New("Type field not set")
	}

	if walletType != WalletTypeBip44 {
		return ErrInvalidWalletType
	}

	if coinType := m[metaCoin]; coinType == "" {
		return errors.New("Coin field not set")
	}

	var isEncrypted bool
	if encStr, ok := m[metaEncrypted]; ok {
		// validate the encrypted value
		var err error
		isEncrypted, err = strconv.ParseBool(encStr)
		if err != nil {
			return errors.New("Encrypted field is not a valid bool")
		}
	}

	if isEncrypted {
		cryptoType, ok := m[metaCryptoType]
		if !ok {
			return errors.New("Crypto type field not set")
		}

		if _, err := getCrypto(CryptoType(cryptoType)); err != nil {
			return errors.New("Unknown crypto type")
		}

		if s := m[metaSecrets]; s == "" {
			return errors.New("Wallet is encrypted, but secrets field not set")
		}

		if s := m[metaSeed]; s != "" {
			return errors.New("Seed should not be visible in encrypted wallets")
		}
	} else {
		if s := m[metaSecrets]; s != "" {
			return errors.New("Secrets should not be in unencrypted wallets")
		}
	}

	// bip44 wallet seeds must be a valid bip39 mnemonic
	if s := m[metaSeed]; s == "" {
		return errors.New("Seed missing in unencrypted bip44 wallet")
	} else if err := bip39.ValidateMnemonic(s); err != nil {
		return err
	}

	if s := m[metaBip44Coin]; s == "" {
		return errors.New("Bip44Coin missing")
	} else if _, err := strconv.ParseUint(s, 10, 32); err != nil {
		return fmt.Errorf("Bip44Coin invalid: %v", err)
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
		return nil, nil, fmt.Errorf("Create external chain public key failed: %v", err)
	}

	change, err := a.NewPublicChildKey(1)
	if err != nil {
		return nil, nil, fmt.Errorf("Create change chain public key failed: %v", err)
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

// IsEncrypted returns whether the wallet is encrypted
func (w Bip44WalletNew) IsEncrypted() bool {
	return w.Meta.IsEncrypted()
}

// clone deep clone of the bip44 wallet
func (w Bip44WalletNew) clone() Bip44WalletNew {
	nw := Bip44WalletNew{
		Meta:     w.Meta.clone(),
		accounts: w.accounts.clone(),
		decoder:  w.decoder,
	}

	return nw
}

func (w *Bip44WalletNew) copyFrom(wlt *Bip44WalletNew) {
	w.Meta = wlt.Meta.clone()
	w.accounts = wlt.accounts.clone()
	w.decoder = wlt.decoder
}

func (w *Bip44WalletNew) erase() {
	w.setSeed("")
	w.setSeedPassphrase("")
	w.accounts.erase()
}

// Lock encrypts the wallet if it is unencrypted, return false
// if it is already encrypted.
func (w *Bip44WalletNew) Lock(password []byte) error {
	if len(password) == 0 {
		return ErrMissingPassword
	}

	if w.IsEncrypted() {
		return ErrWalletEncrypted
	}

	wlt := w.clone()

	ss := make(Secrets)
	defer func() {
		ss.erase()
		wlt.erase()
	}()

	wlt.packSecrets(ss)

	sb, err := ss.serialize()
	if err != nil {
		return err
	}

	cryptoType := wlt.Meta.CryptoType()
	if cryptoType == "" {
		return errors.New("Crypto type field not set")
	}
	crypto, err := getCrypto(cryptoType)
	if err != nil {
		return err
	}

	encSecret, err := crypto.Encrypt(sb, password)
	if err != nil {
		return err
	}

	// Sets wallet as encrypted, updates secret field.
	wlt.SetEncrypted(cryptoType, string(encSecret))

	// Update wallet to the latest version, which indicates encryption support
	wlt.SetVersion(Version)

	// Wipes the secret fields in wlt
	wlt.erase()

	// Wipes the secret fields in w
	w.erase()

	w.copyFrom(&wlt)
	return nil
}

// Unlock decrypt the wallet
func (w *Bip44WalletNew) Unlock(password []byte) (*Bip44WalletNew, error) {
	if !w.IsEncrypted() {
		return nil, ErrWalletNotEncrypted
	}

	if len(password) == 0 {
		return nil, ErrMissingPassword
	}

	sstr := w.Secrets()
	if sstr == "" {
		return nil, errors.New("Secrets missing from wallet")
	}

	ct := w.CryptoType()
	if ct == "" {
		return nil, errors.New("Missing crypto type")
	}

	crypto, err := getCrypto(ct)
	if err != nil {
		return nil, err
	}

	sb, err := crypto.Decrypt([]byte(sstr), password)
	if err != nil {
		return nil, ErrInvalidPassword
	}

	defer func() {
		// Wipes the dat from secrets bytes buffer
		for i := range sb {
			sb[i] = 0
		}
	}()

	ss := make(Secrets)
	defer ss.erase()
	if err := ss.deserialize(sb); err != nil {
		return nil, err
	}

	cw := w.clone()
	if err := cw.unpackSecrets(ss); err != nil {
		return nil, err
	}
	cw.SetDecrypted()

	return &cw, nil
}

// packSecrets saves all sensitive data to the secrets map.
func (w Bip44WalletNew) packSecrets(ss Secrets) {
	ss.set(secretSeed, w.Meta.Seed())
	ss.set(secretSeedPassphrase, w.Meta.SeedPassphrase())
	w.accounts.packSecrets(ss)
}

func (w *Bip44WalletNew) unpackSecrets(ss Secrets) error {
	seed, ok := ss.get(secretSeed)
	if !ok {
		return errors.New("Seed does not exist in secrets")
	}
	w.Meta.setSeed(seed)

	passphrase, _ := ss.get(secretSeedPassphrase)
	w.Meta.setSeedPassphrase(passphrase)

	w.accounts.unpackSecrets(ss)
	return nil
}
