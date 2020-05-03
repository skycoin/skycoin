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
	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"

	"github.com/SkycoinProject/skycoin/src/cipher/bip32"
	"github.com/SkycoinProject/skycoin/src/cipher/bip39"
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
)

const walletType = "bip44"

var (
	// defaultBip44WalletDecoder is the default bip44 wallet decoder
	defaultBip44WalletDecoder = &JSONDecoder{}
)

var logger = logging.MustGetLogger("bip44wallet")

// Wallet manages keys using the original Skycoin deterministic
// keypair generator method.
type Wallet struct {
	//Meta wallet meta data
	wallet.Meta
	// accounts bip44 wallet accounts
	accounts accountManager
	// decoder is used to encode/decode bip44 wallet to/from []byte
	decoder Decoder
}

// accountManager is the interface that manages the bip44 wallet accounts.
type accountManager interface {
	// new creates a new account, returns the account index, and error, if any
	new(opts bip44AccountCreateOptions) (uint32, error)
	// newAddresses generates addresses on selected account
	newAddresses(account, chain, num uint32) ([]cipher.Addresser, error)
	// entries returns entries of specific chain of the selected account
	entries(account, chain uint32) (wallet.Entries, error)
	// entriesLen returns the entries length of specific chain of selected account
	entriesLen(account, chain uint32) (uint32, error)
	// entryAt  returns the entry of specific index
	entryAt(account, chain, index uint32) (wallet.Entry, error)
	// getEntry returns the entry of given address
	getEntry(account uint32, address cipher.Addresser) (wallet.Entry, bool, error)
	// len returns the account number
	len() uint32
	// clone returns a deep clone accounts manager
	clone() accountManager
	// syncSecrets checks if there are any addresses that do not have secrets associated with,
	// if yes, generate the secrets for those addresses
	syncSecrets(ss wallet.Secrets) error
	// packSecrets packs secrets
	packSecrets(ss wallet.Secrets)
	// unpackSecrets unpacks secrets
	unpackSecrets(ss wallet.Secrets) error
	// erase wipes secrets
	erase()
	account(account uint32) (*bip44Account, error)
	// all returns all accounts in wallet.Account format
	all() []wallet.Bip44Account
}

// Decoder is the interface that wraps the Encode and Decode methods.
// Encode method encodes the wallet to bytes, Decode method decodes bytes to bip44 wallet.
type Decoder interface {
	Encode(w Wallet) ([]byte, error)
	Decode(b []byte) (*Wallet, error)
}

// ChainEntry represents an item on the bip44 wallet chain
type ChainEntry struct {
	Address cipher.Addresser
}

// Options options for creating the bip44 wallet
type Options struct {
	Version       string
	Bip44CoinType *bip44.CoinType
	CryptoType    crypto.CryptoType
	WalletDecoder Decoder
}

// NewWallet create a bip44 wallet with options
func NewWallet(filename, label, seed, seedPassphrase string, coinType wallet.CoinType, opts *Options) (*Wallet, error) {
	wlt := &Wallet{
		Meta: wallet.Meta{
			wallet.MetaType:           walletType,
			wallet.MetaFilename:       filename,
			wallet.MetaLabel:          label,
			wallet.MetaSeed:           seed,
			wallet.MetaSeedPassphrase: seedPassphrase,
			wallet.MetaCoin:           string(coinType),
			wallet.MetaTimestamp:      strconv.FormatInt(time.Now().Unix(), 10),
			wallet.MetaEncrypted:      "false",
		},
		accounts: &bip44Accounts{},
		decoder:  opts.WalletDecoder,
	}

	if opts != nil {
		wlt.Meta[wallet.MetaVersion] = opts.Version
		wlt.Meta[wallet.MetaCryptoType] = string(opts.CryptoType)

		if coinType == "" {
			return nil, errors.New("missing coin type")
		}

		// Note: if opts.Bip44CoinType is nil, we will only set bip44 coin type for
		// skycoin and bitcoin. All other coins should explicitly set it, otherwise
		// an error will be reported.
		if opts.Bip44CoinType == nil {
			switch coinType {
			case wallet.CoinTypeSkycoin:
				wlt.Meta.SetBip44Coin(bip44.CoinTypeSkycoin)
			case wallet.CoinTypeBitcoin:
				wlt.Meta.SetBip44Coin(bip44.CoinTypeBitcoin)
			default:
				return nil, errors.New("missing bip44 coin type")
			}
		} else {
			wlt.Meta.SetBip44Coin(*opts.Bip44CoinType)
		}
	}

	if wlt.CryptoType() == "" {
		wlt.Meta[wallet.MetaCryptoType] = string(crypto.DefaultCryptoType)
	}

	if wlt.decoder == nil {
		wlt.decoder = defaultBip44WalletDecoder
	}

	if err := validateMeta(wlt.Meta); err != nil {
		return nil, err
	}
	return wlt, nil
}

func validateMeta(m wallet.Meta) error {
	if fn := m[wallet.MetaFilename]; fn == "" {
		return errors.New("filename not set")
	}

	if tm := m[wallet.MetaTimestamp]; tm != "" {
		_, err := strconv.ParseInt(tm, 10, 64)
		if err != nil {
			return errors.New("invalid timestamp")
		}
	}

	walletType, ok := m[wallet.MetaType]
	if !ok {
		return errors.New("type field not set")
	}

	if walletType != walletType {
		return errors.New("invalid wallet type")
	}

	if coinType := m[wallet.MetaCoin]; coinType == "" {
		return errors.New("coin field not set")
	}

	var isEncrypted bool
	if encStr, ok := m[wallet.MetaEncrypted]; ok {
		// validate the encrypted value
		var err error
		isEncrypted, err = strconv.ParseBool(encStr)
		if err != nil {
			return errors.New("encrypted field is not a valid bool")
		}
	}

	if isEncrypted {
		cryptoType, ok := m[wallet.MetaCryptoType]
		if !ok {
			return errors.New("crypto type field not set")
		}

		if _, err := crypto.GetCrypto(crypto.CryptoType(cryptoType)); err != nil {
			return errors.New("unknown crypto type")
		}

		if s := m[wallet.MetaSecrets]; s == "" {
			return errors.New("wallet is encrypted, but secrets field not set")
		}

		if s := m[wallet.MetaSeed]; s != "" {
			return errors.New("seed should not be visible in encrypted wallets")
		}
	} else {
		if s := m[wallet.MetaSecrets]; s != "" {
			return errors.New("secrets should not be in unencrypted wallets")
		}
	}

	// bip44 wallet seeds must be a valid bip39 mnemonic
	if s := m[wallet.MetaSeed]; s == "" {
		return errors.New("seed missing in unencrypted bip44 wallet")
	} else if err := bip39.ValidateMnemonic(s); err != nil {
		return err
	}

	if s := m[wallet.MetaBip44Coin]; s == "" {
		return errors.New("missing bip44 coin type")
	} else if _, err := strconv.ParseUint(s, 10, 32); err != nil {
		return fmt.Errorf("invalid bip44 coin type: %v", err)
	}

	return nil
}

// NewAccount create a bip44 wallet account, returns account index and
// error, if any.
func (w *Wallet) NewAccount(name string) (uint32, error) {
	return w.accounts.new(bip44AccountCreateOptions{
		name:           name,
		seed:           w.Seed(),
		seedPassphrase: w.SeedPassphrase(),
		coinType:       w.Coin(),
		bip44CoinType:  w.Bip44Coin(),
	})
}

// NewExternalAddresses generates addresses on external chain of selected account
func (w *Wallet) NewExternalAddresses(account, n uint32) ([]cipher.Addresser, error) {
	return w.accounts.newAddresses(account, bip44.ExternalChainIndex, n)
}

// NewChangeAddresses generates addresses on change chain of selected account
func (w *Wallet) NewChangeAddresses(account, n uint32) ([]cipher.Addresser, error) {
	return w.accounts.newAddresses(account, bip44.ChangeChainIndex, n)
}

// ExternalEntries returns the entries on external chain
func (w *Wallet) ExternalEntries(account uint32) (wallet.Entries, error) {
	return w.accounts.entries(account, bip44.ExternalChainIndex)
}

// ChangeEntries returns the entries on change chain
func (w *Wallet) ChangeEntries(account uint32) (wallet.Entries, error) {
	return w.accounts.entries(account, bip44.ChangeChainIndex)
}

// ExternalEntriesLen returns the external chain entries length of selected account
func (w *Wallet) ExternalEntriesLen(account uint32) (uint32, error) {
	return w.accounts.entriesLen(account, bip44.ExternalChainIndex)
}

// ChangeEntriesLen returns the change chain entries length of selected account
func (w *Wallet) ChangeEntriesLen(account uint32) (uint32, error) {
	return w.accounts.entriesLen(account, bip44.ChangeChainIndex)
}

// ExternalEntryAt returns the entry at the given index on external chain of selected account
func (w *Wallet) ExternalEntryAt(account, i uint32) (wallet.Entry, error) {
	return w.accounts.entryAt(account, bip44.ExternalChainIndex, i)
}

// ChangeEntryAt returns the entry at the given index on change chain of selected account
func (w *Wallet) ChangeEntryAt(account, i uint32) (wallet.Entry, error) {
	return w.accounts.entryAt(account, bip44.ChangeChainIndex, i)
}

// GetEntry returns the entry of given address on selected account
func (w *Wallet) GetEntry(account uint32, address cipher.Addresser) (wallet.Entry, bool, error) {
	return w.accounts.getEntry(account, address)
}

// Serialize encodes the bip44 wallet to []byte
func (w Wallet) Serialize() ([]byte, error) {
	if w.decoder == nil {
		w.decoder = defaultBip44WalletDecoder
	}
	return w.decoder.Encode(w)
}

// Deserialize decodes the []byte to a bip44 wallet
func (w *Wallet) Deserialize(b []byte) error {
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
func (w Wallet) IsEncrypted() bool {
	return w.Meta.IsEncrypted()
}

// Lock encrypts the wallet if it is unencrypted, return false
// if it is already encrypted.
func (w *Wallet) Lock(password []byte) error {
	if len(password) == 0 {
		return errors.New("missing password when locking bip44 wallet")
	}

	if w.IsEncrypted() {
		return errors.New("wallet is already encrypted")
	}

	wlt := w.Clone().(*Wallet)

	ss := make(wallet.Secrets)
	defer func() {
		ss.Erase()
		wlt.Erase()
	}()

	wlt.packSecrets(ss)

	sb, err := ss.Serialize()
	if err != nil {
		return err
	}

	cryptoType := wlt.Meta.CryptoType()
	if cryptoType == "" {
		return errors.New("crypto type field not set")
	}
	cryptor, err := crypto.GetCrypto(cryptoType)
	if err != nil {
		return err
	}

	encSecret, err := cryptor.Encrypt(sb, password)
	if err != nil {
		return err
	}

	// Sets wallet as encrypted, updates secret field.
	wlt.SetEncrypted(cryptoType, string(encSecret))

	// Wipes the secret fields in wlt
	wlt.Erase()

	// Wipes the secret fields in w
	w.Erase()

	w.copyFrom(wlt)
	return nil
}

// Unlock decrypt the wallet
func (w *Wallet) Unlock(password []byte) (wallet.Wallet, error) {
	if !w.IsEncrypted() {
		return nil, errors.New("wallet is not encrypted")
	}

	if len(password) == 0 {
		return nil, errors.New("missing password")
	}

	sstr := w.Secrets()
	if sstr == "" {
		return nil, errors.New("wallet.Secrets missing from wallet")
	}

	ct := w.CryptoType()
	if ct == "" {
		return nil, errors.New("missing crypto type")
	}

	cryptor, err := crypto.GetCrypto(ct)
	if err != nil {
		return nil, err
	}

	sb, err := cryptor.Decrypt([]byte(sstr), password)
	if err != nil {
		return nil, errors.New("invalid password")
	}

	defer func() {
		// Wipes the data from secrets bytes buffer
		for i := range sb {
			sb[i] = 0
		}
	}()

	ss := make(wallet.Secrets)
	defer ss.Erase()
	if err := ss.Deserialize(sb); err != nil {
		return nil, err
	}

	cw := w.Clone().(*Wallet)

	initSSLen := len(ss)
	// fills secrets for those new generated addresses
	if err := cw.syncSecrets(ss); err != nil {
		return nil, err
	}

	if len(ss) > initSSLen {
		// new secrets generated, update the secrets field of the locked wallet
		sb, err := ss.Serialize()
		if err != nil {
			return nil, err
		}

		encSecret, err := cryptor.Encrypt(sb, password)
		if err != nil {
			return nil, err
		}

		// Sets wallet as encrypted, updates secret field.
		w.SetEncrypted(ct, string(encSecret))
	}

	if err := cw.unpackSecrets(ss); err != nil {
		return nil, err
	}
	cw.SetDecrypted()
	return cw, nil
}

func (w Wallet) Fingerprint() string {
	addr := ""
	entries, err := w.ExternalEntries(0)
	if err != nil {
		logger.WithError(err).Panic("Fingerprint get external entries failed")
		return ""
	}

	if len(entries) == 0 {
		if !w.IsEncrypted() {
			addrs, err := w.NewExternalAddresses(0, 1)
			if err != nil {
				logger.WithError(err).Panic("Fingerprint failed to generate initial entry for empty wallet")
			}
			addr = addrs[0].String()
		}
	} else {
		addr = entries[0].Address.String()
	}
	return fmt.Sprintf("%s-%s", w.Type(), addr)
}

// Clone deep clone of the bip44 wallet
func (w Wallet) Clone() wallet.Wallet {
	return &Wallet{
		Meta:     w.Meta.Clone(),
		accounts: w.accounts.clone(),
		decoder:  w.decoder,
	}
}

func (w *Wallet) CopyFrom(src wallet.Wallet) {
	w.copyFrom(src.(*Wallet))
}

// CopyFromRef copies the src wallet with a pointer dereference
func (w *Wallet) CopyFromRef(src wallet.Wallet) {
	*w = *(src.(*Wallet))
}

func (w *Wallet) Accounts() []wallet.Bip44Account {
	return w.accounts.all()
}

// Entries provides entries service to access the external chain of given account
func (w *Wallet) Entries(options ...wallet.Option) wallet.EntriesService {
	eos := &wallet.Bip44EntriesOptions{}
	for _, opt := range options {
		opt(eos)
	}

	aes := &accountEntriesService{}
	a, err := w.accounts.account(eos.Account)
	if err != nil {
		aes.err = err
		return aes
	}

	if eos.Change {
		aes.chain = uint32(1)
	}

	return aes
}

func (w *Wallet) copyFrom(wlt *Wallet) {
	w.Meta = wlt.Meta.Clone()
	w.accounts = wlt.accounts.clone()
	w.decoder = wlt.decoder
}

// Erase wipes all sensitive data
func (w *Wallet) Erase() {
	w.SetSeed("")
	w.SetSeedPassphrase("")
	w.accounts.erase()
}

// immutableMeta records the meta keys of a wallet that should not be modified
// once after they are initialized.
//func immutableMeta() map[string]struct{} {
//	empty := struct{}{}
//	return map[string]struct{}{
//		wallet.MetaFilename:       empty,
//		wallet.MetaCoin:           empty,
//		wallet.MetaType:           empty,
//		wallet.MetaCryptoType:     empty,
//		wallet.MetaSeed:           empty,
//		wallet.MetaSeedPassphrase: empty,
//	}
//}

//func secretsMeta() map[string]struct{} {
//	empty := struct{}{}
//	return map[string]struct{}{
//		wallet.MetaSeed:           empty,
//		wallet.MetaSeedPassphrase: empty,
//		wallet.MetaSecrets:        empty,
//		wallet.MetaEncrypted:      empty,
//	}
//}
//
//// WalletDiff records the wallet differences
//type WalletDiff struct {
//	Meta     wallet.Meta
//	Accounts []AccountDiff
//}
//
//// AccountDiff records the account differences
//type AccountDiff struct {
//	NewExternalAddressNum int
//	NewChangeAddressNum   int
//}

// DiffNoneSecrets gets the differences of none secrets between wallets
//
// Note: immutable meta like the wallet filename, coin type, wallet type, etc.
// will be filter out, they won't be recognized as changes.
//func (w *Wallet) DiffNoneSecrets(wlt *Wallet) (*WalletDiff, error) {
//	diff := &WalletDiff{
//		Meta:     make(wallet.Meta),
//		Accounts: make([]AccountDiff, w.accounts.len()),
//	}
//
//	im := immutableMeta()
//	sm := secretsMeta()
//
//	// check the meta change
//	for k, v := range wlt.Meta {
//		// filter out the immutable meta data
//		if _, ok := im[k]; ok {
//			continue
//		}
//
//		// filter out the secrets meta
//		if _, ok := sm[k]; ok {
//			continue
//		}
//
//		if w.Meta[k] != v {
//			diff.Meta[k] = v
//		}
//	}
//
//	accountsDiff := w.accounts.diff(wlt.accounts)
//	for i, adf := range accountsDiff {
//		diff.Accounts[i].NewExternalAddressNum = int(adf.chainsDiff[bip44.ExternalChainIndex])
//		diff.Accounts[i].NewChangeAddressNum = int(adf.chainsDiff[bip44.ChangeChainIndex])
//	}
//
//	return diff, nil
//}
//
// CommitDiffs applies the wallet differences
//
// Immutable meta data will be filter out
// Secrets meta data will be committed
//func (w *Wallet) CommitDiffs(diff *WalletDiff) error {
//	w2 := w.Clone()
//	im := immutableMeta()
//
//	// filter out the immutable meta data
//	for k, v := range diff.Meta {
//		if _, ok := im[k]; ok {
//			continue
//		}
//		w2.Meta[k] = v
//	}
//
//	for i, a := range diff.Accounts {
//		if a.NewExternalAddressNum > 0 {
//			_, err := w2.NewExternalAddresses(uint32(i), uint32(a.NewExternalAddressNum))
//			if err != nil {
//				return err
//			}
//		}
//
//		if a.NewChangeAddressNum > 0 {
//			_, err := w2.NewChangeAddresses(uint32(i), uint32(a.NewChangeAddressNum))
//			if err != nil {
//				return err
//			}
//		}
//	}
//	*w = w2
//	return nil
//}

// syncSecrets synchronize the secrets with all addresses, ensure that
// each address has the secret key stored in the secrets
func (w Wallet) syncSecrets(ss wallet.Secrets) error {
	return w.accounts.syncSecrets(ss)
}

// packSecrets saves all sensitive data to the secrets map.
func (w Wallet) packSecrets(ss wallet.Secrets) {
	ss.Set(wallet.SecretSeed, w.Meta.Seed())
	ss.Set(wallet.SecretSeedPassphrase, w.Meta.SeedPassphrase())
	w.accounts.packSecrets(ss)
}

func (w *Wallet) unpackSecrets(ss wallet.Secrets) error {
	seed, ok := ss.Get(wallet.SecretSeed)
	if !ok {
		return errors.New("Seed does not exist in secrets")
	}
	w.Meta.SetSeed(seed)

	passphrase, _ := ss.Get(wallet.SecretSeedPassphrase)
	w.Meta.SetSeedPassphrase(passphrase)

	return w.accounts.unpackSecrets(ss)
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
