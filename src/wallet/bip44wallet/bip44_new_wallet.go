// This file is still under refactoring, once it gets done, we will replace the
// old bip44_wallet.go.

package bip44wallet

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/util/logging"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
	"github.com/SkycoinProject/skycoin/src/wallet/entry"
	"github.com/SkycoinProject/skycoin/src/wallet/meta"
	"github.com/SkycoinProject/skycoin/src/wallet/secrets"

	"github.com/SkycoinProject/skycoin/src/cipher/bip32"
	"github.com/SkycoinProject/skycoin/src/cipher/bip39"
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
)

const walletType = "bip44"

var (
	// defaultBip44WalletDecoder is the default bip44 wallet decoder
	defaultBip44WalletDecoder = &Bip44WalletJSONDecoder{}
)

var logger = logging.MustGetLogger("bip44wallet")

// Bip44WalletNew manages keys using the original Skycoin deterministic
// keypair generator method.
type Bip44WalletNew struct {
	//Meta wallet meta data
	meta.Meta
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
	newAddresses(account, chain, num uint32) ([]cipher.Addresser, error)
	// entries reutrns entries of specific chain of the selected account
	entries(account, chain uint32) (entry.Entries, error)
	// entriesLen returns the entries length of specific chain of selected account
	entriesLen(account, chain uint32) (uint32, error)
	// entryAt  returns the entry of specific index
	entryAt(account, chain, index uint32) (entry.Entry, error)
	// len returns the account number
	len() uint32
	// clone returns a deep clone accounts manager
	clone() accountManager
	// packSecrets packs secrets
	packSecrets(ss secrets.Secrets)
	// unpackSecrets unpacks secrets
	unpackSecrets(ss secrets.Secrets) error
	// erase wipes secrets
	erase()
}

// Bip44WalletDecoder is the interface that wraps the Encode and Decode methods.
// Encode method encodes the wallet to bytes, Decode method decodes bytes to bip44 wallet.
type Bip44WalletDecoder interface {
	Encode(w *Bip44WalletNew) ([]byte, error)
	Decode(b []byte) (*Bip44WalletNew, error)
}

// ChainEntry represents an item on the bip44 wallet chain
type ChainEntry struct {
	Address cipher.Addresser
}

// Bip44WalletCreateOptions options for creating the bip44 wallet
type Bip44WalletCreateOptions struct {
	Filename       string
	Version        string
	Label          string
	Seed           string
	SeedPassphrase string
	CoinType       meta.CoinType
	CryptoType     crypto.CryptoType
	WalletDecoder  Bip44WalletDecoder
}

// NewBip44WalletNew create a bip44 wallet with options
func NewBip44WalletNew(opts Bip44WalletCreateOptions) (*Bip44WalletNew, error) {
	wlt := &Bip44WalletNew{
		Meta: meta.Meta{
			meta.MetaType:           walletType,
			meta.MetaFilename:       opts.Filename,
			meta.MetaVersion:        opts.Version,
			meta.MetaLabel:          opts.Label,
			meta.MetaSeed:           opts.Seed,
			meta.MetaSeedPassphrase: opts.SeedPassphrase,
			meta.MetaCoin:           string(opts.CoinType),
			meta.MetaTimestamp:      strconv.FormatInt(time.Now().Unix(), 10),
			meta.MetaCryptoType:     string(opts.CryptoType),
			meta.MetaEncrypted:      "false",
		},
		accounts: &bip44Accounts{},
		decoder:  opts.WalletDecoder,
	}

	if wlt.CryptoType() == "" {
		wlt.Meta[meta.MetaCryptoType] = string(crypto.DefaultCryptoType)
	}

	if wlt.decoder == nil {
		wlt.decoder = defaultBip44WalletDecoder
	}

	bip44CoinType := resolveCoinAdapter(opts.CoinType).Bip44CoinType()
	wlt.Meta.SetBip44Coin(bip44CoinType)

	if err := bip44MetaValidate(wlt.Meta); err != nil {
		return nil, err
	}
	return wlt, nil
}

func bip44MetaValidate(m meta.Meta) error {
	if fn := m[meta.MetaFilename]; fn == "" {
		return errors.New("Filename not set")
	}

	if tm := m[meta.MetaTimestamp]; tm != "" {
		_, err := strconv.ParseInt(tm, 10, 64)
		if err != nil {
			return errors.New("Invalid timestamp")
		}
	}

	walletType, ok := m[meta.MetaType]
	if !ok {
		return errors.New("Type field not set")
	}

	if walletType != walletType {
		return errors.New("Invalid wallet type")
	}

	if coinType := m[meta.MetaCoin]; coinType == "" {
		return errors.New("Coin field not set")
	}

	var isEncrypted bool
	if encStr, ok := m[meta.MetaEncrypted]; ok {
		// validate the encrypted value
		var err error
		isEncrypted, err = strconv.ParseBool(encStr)
		if err != nil {
			return errors.New("Encrypted field is not a valid bool")
		}
	}

	if isEncrypted {
		cryptoType, ok := m[meta.MetaCryptoType]
		if !ok {
			return errors.New("Crypto type field not set")
		}

		if _, err := crypto.GetCrypto(crypto.CryptoType(cryptoType)); err != nil {
			return errors.New("Unknown crypto type")
		}

		if s := m[meta.MetaSecrets]; s == "" {
			return errors.New("Wallet is encrypted, but secrets field not set")
		}

		if s := m[meta.MetaSeed]; s != "" {
			return errors.New("Seed should not be visible in encrypted wallets")
		}
	} else {
		if s := m[meta.MetaSecrets]; s != "" {
			return errors.New("wallet.Secrets should not be in unencrypted wallets")
		}
	}

	// bip44 wallet seeds must be a valid bip39 mnemonic
	if s := m[meta.MetaSeed]; s == "" {
		return errors.New("Seed missing in unencrypted bip44 wallet")
	} else if err := bip39.ValidateMnemonic(s); err != nil {
		return err
	}

	if s := m[meta.MetaBip44Coin]; s == "" {
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
		seed:           w.Seed(),
		seedPassphrase: w.SeedPassphrase(),
		coinType:       meta.CoinType(w.Coin()),
	}

	return w.accounts.new(opts)
}

// NewExternalAddresses generates addresses on external chain of selected account
func (w *Bip44WalletNew) NewExternalAddresses(account, n uint32) ([]cipher.Addresser, error) {
	return w.accounts.newAddresses(account, bip44.ExternalChainIndex, n)
}

// NewChangeAddresses generates addresses on change chain of selected account
func (w *Bip44WalletNew) NewChangeAddresses(account, n uint32) ([]cipher.Addresser, error) {
	return w.accounts.newAddresses(account, bip44.ChangeChainIndex, n)
}

// ExternalEntries returns the entries on external external chain
func (w *Bip44WalletNew) ExternalEntries(account uint32) (entry.Entries, error) {
	return w.accounts.entries(account, bip44.ExternalChainIndex)
}

// ChangeEntries returns the entries on external external chain
func (w *Bip44WalletNew) ChangeEntries(account uint32) (entry.Entries, error) {
	return w.accounts.entries(account, bip44.ChangeChainIndex)
}

// ExternalEntriesLen returns the external chain entries length of selected account
func (w *Bip44WalletNew) ExternalEntriesLen(account uint32) (uint32, error) {
	return w.accounts.entriesLen(account, bip44.ExternalChainIndex)
}

// ChangeEntriesLen returns the change chain entries length of selected account
func (w *Bip44WalletNew) ChangeEntriesLen(account uint32) (uint32, error) {
	return w.accounts.entriesLen(account, bip44.ChangeChainIndex)
}

// ExternalEntryAt returns the entry at the given index on external chain of selected account
func (w *Bip44WalletNew) ExternalEntryAt(account, i uint32) (entry.Entry, error) {
	return w.accounts.entryAt(account, bip44.ExternalChainIndex, i)
}

// ChangeEntryAt returns the entry at the given index on change chain of selected account
func (w *Bip44WalletNew) ChangeEntryAt(account, i uint32) (entry.Entry, error) {
	return w.accounts.entryAt(account, bip44.ChangeChainIndex, i)
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

// Lock encrypts the wallet if it is unencrypted, return false
// if it is already encrypted.
func (w *Bip44WalletNew) Lock(password []byte) error {
	if len(password) == 0 {
		return errors.New("Missing password when locking bip44 wallet")
	}

	if w.IsEncrypted() {
		return errors.New("Wallet is already encrypted")
	}

	wlt := w.Clone()

	ss := make(secrets.Secrets)
	defer func() {
		ss.Erase()
		wlt.erase()
	}()

	wlt.packSecrets(ss)

	sb, err := ss.Serialize()
	if err != nil {
		return err
	}

	cryptoType := wlt.Meta.CryptoType()
	if cryptoType == "" {
		return errors.New("Crypto type field not set")
	}
	crypto, err := crypto.GetCrypto(cryptoType)
	if err != nil {
		return err
	}

	encSecret, err := crypto.Encrypt(sb, password)
	if err != nil {
		return err
	}

	// Sets wallet as encrypted, updates secret field.
	wlt.SetEncrypted(cryptoType, string(encSecret))

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
		return nil, errors.New("Wallet is not encrypted")
	}

	if len(password) == 0 {
		return nil, errors.New("Missing password")
	}

	sstr := w.Secrets()
	if sstr == "" {
		return nil, errors.New("wallet.Secrets missing from wallet")
	}

	ct := w.CryptoType()
	if ct == "" {
		return nil, errors.New("Missing crypto type")
	}

	crypto, err := crypto.GetCrypto(ct)
	if err != nil {
		return nil, err
	}

	sb, err := crypto.Decrypt([]byte(sstr), password)
	if err != nil {
		return nil, errors.New("Invalid password")
	}

	defer func() {
		// Wipes the dat from secrets bytes buffer
		for i := range sb {
			sb[i] = 0
		}
	}()

	ss := make(secrets.Secrets)
	defer ss.Erase()
	if err := ss.Deserialize(sb); err != nil {
		return nil, err
	}

	cw := w.Clone()
	if err := cw.unpackSecrets(ss); err != nil {
		return nil, err
	}
	cw.SetDecrypted()

	return &cw, nil
}

// Clone deep clone of the bip44 wallet
func (w Bip44WalletNew) Clone() Bip44WalletNew {
	nw := Bip44WalletNew{
		Meta:     w.Meta.Clone(),
		accounts: w.accounts.clone(),
		decoder:  w.decoder,
	}

	return nw
}

func (w *Bip44WalletNew) copyFrom(wlt *Bip44WalletNew) {
	w.Meta = wlt.Meta.Clone()
	w.accounts = wlt.accounts.clone()
	w.decoder = wlt.decoder
}

func (w *Bip44WalletNew) erase() {
	w.SetSeed("")
	w.SetSeedPassphrase("")
	w.accounts.erase()
}

// packSecrets saves all sensitive data to the secrets map.
func (w Bip44WalletNew) packSecrets(ss secrets.Secrets) {
	ss.Set(secrets.SecretSeed, w.Meta.Seed())
	ss.Set(secrets.SecretSeedPassphrase, w.Meta.SeedPassphrase())
	w.accounts.packSecrets(ss)
}

func (w *Bip44WalletNew) unpackSecrets(ss secrets.Secrets) error {
	seed, ok := ss.Get(secrets.SecretSeed)
	if !ok {
		return errors.New("Seed does not exist in secrets")
	}
	w.Meta.SetSeed(seed)

	passphrase, _ := ss.Get(secrets.SecretSeedPassphrase)
	w.Meta.SetSeedPassphrase(passphrase)

	w.accounts.unpackSecrets(ss)
	return nil
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

// TODO:
// - Integrate the bip44 wallet to wallet system
