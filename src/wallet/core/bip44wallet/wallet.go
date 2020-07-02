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

const (
	// WalletType represents the bip44 wallet type
	WalletType = "bip44"
	// DefaultAccountName is the default bip44 account name
	DefaultAccountName = "default"
)

var (
	// defaultWalletDecoder is the default bip44 wallet decoder
	defaultWalletDecoder = &JSONDecoder{}
)

var logger = logging.MustGetLogger("bip44wallet")

// Wallet manages keys using the original Skycoin deterministic
// keypair generator method.
type Wallet struct {
	//Meta wallet meta data
	wallet.Meta
	// accounts bip44 wallet accounts
	accountManager
	// decoder is used to encode/decode bip44 wallet to/from []byte
	decoder wallet.Decoder
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
	// reset reset all accounts' entries
	reset()
}

// ChainEntry represents an item on the bip44 wallet chain
type ChainEntry struct {
	Address cipher.Addresser
}

// NewWallet create a bip44 wallet with options
// TODO: encrypt the wallet if the options.Encrypt is true
// TODO: generate a default account when create a new wallet
// also, a default address will be generated
func NewWallet(filename, label, seed, seedPassphrase string, options ...wallet.Option) (*Wallet, error) {
	wlt := &Wallet{
		Meta: wallet.Meta{
			wallet.MetaFilename:       filename,
			wallet.MetaLabel:          label,
			wallet.MetaSeed:           seed,
			wallet.MetaSeedPassphrase: seedPassphrase,
			wallet.MetaEncrypted:      "false",
			wallet.MetaType:           WalletType,
			wallet.MetaVersion:        wallet.Version,
			wallet.MetaCoin:           string(wallet.CoinTypeSkycoin),
			wallet.MetaCryptoType:     string(crypto.DefaultCryptoType),
			wallet.MetaTimestamp:      strconv.FormatInt(time.Now().Unix(), 10),
		},
		accountManager: &bip44Accounts{},
		decoder:        defaultWalletDecoder,
	}

	advOpts := wallet.AdvancedOptions{}
	// applies options to wallet and AdvancedOptions
	for _, opt := range options {
		opt(wlt)
		opt(&advOpts)
	}

	if wlt.Bip44Coin() == nil {
		switch wlt.Coin() {
		case wallet.CoinTypeSkycoin:
			wlt.SetBip44Coin(bip44.CoinTypeSkycoin)
		case wallet.CoinTypeBitcoin:
			wlt.SetBip44Coin(bip44.CoinTypeBitcoin)
		default:
			return nil, errors.New("bip44 coin type not set")
		}
	}

	// validateMeta wallet before encrypting
	if err := validateMeta(wlt.Meta); err != nil {
		return nil, err
	}

	// generates a default account
	var accountName = DefaultAccountName
	if advOpts.DefaultBip44AccountName != "" {
		accountName = advOpts.DefaultBip44AccountName
	}

	_, err := wlt.NewAccount(accountName)
	if err != nil {
		return nil, fmt.Errorf("generate default account failed: %v", err)
	}

	// Generate addresses if options.GenrateN > 0
	generateN := advOpts.GenerateN
	if generateN == 0 {
		generateN = 1
	}

	if _, err := wlt.GenerateAddresses(generateN); err != nil {
		return nil, err
	}

	// Generate a default change address
	if _, err := wlt.GenerateAddresses(1, wallet.OptionChange(true)); err != nil {
		return nil, err
	}

	scanN := advOpts.ScanN
	// scans addresses if options.ScanN > 0
	if scanN > 0 {
		if advOpts.TF == nil {
			return nil, errors.New("missing transaction finder for scanning addresses")
		}

		if scanN > generateN {
			scanN = scanN - generateN
		}

		_, err := wlt.ScanAddresses(scanN, advOpts.TF)
		if err != nil {
			return nil, err
		}
	}

	// encrypts wallet if options.Encrypt is true
	if advOpts.Encrypt {
		if len(advOpts.Password) == 0 {
			return nil, errors.New("missing password for encrypting wallet")
		}

		if err := wlt.Lock(advOpts.Password); err != nil {
			return nil, err
		}
	}

	// validateMeta the wallet again after encrypted
	if err := validateMeta(wlt.Meta); err != nil {
		return nil, err
	}
	return wlt, nil
}

func validateMeta(m wallet.Meta) error {
	if m[wallet.MetaType] != WalletType {
		return errors.New("invalid wallet type")
	}

	if s := m[wallet.MetaBip44Coin]; s == "" {
		return errors.New("missing bip44 coin type")
	} else if _, err := strconv.ParseUint(s, 10, 32); err != nil {
		return fmt.Errorf("invalid bip44 coin type: %v", err)
	}

	if err := wallet.ValidateMeta(m); err != nil {
		return err
	}

	if s := m[wallet.MetaSeed]; s != "" {
		if err := bip39.ValidateMnemonic(s); err != nil {
			return err
		}
	}

	if err := wallet.ValidateMetaCryptoType(m); err != nil {
		return err
	}

	return wallet.ValidateMetaSeed(m)
}

// SetDecoder sets the wallet decoder
func (w *Wallet) SetDecoder(d wallet.Decoder) {
	w.decoder = d
}

// NewAccount create a bip44 wallet account, returns account index and
// error, if any.
func (w *Wallet) NewAccount(name string) (uint32, error) {
	return w.accountManager.new(bip44AccountCreateOptions{
		name:           name,
		seed:           w.Seed(),
		seedPassphrase: w.SeedPassphrase(),
		coinType:       w.Coin(),
		bip44CoinType:  w.Bip44Coin(),
	})
}

// newExternalAddresses generates addresses on external chain of selected account
func (w *Wallet) newExternalAddresses(account, n uint32) ([]cipher.Addresser, error) {
	return w.newAddresses(account, bip44.ExternalChainIndex, n)
}

// NewChangeAddresses generates addresses on change chain of selected account
func (w *Wallet) newChangeAddresses(account, n uint32) ([]cipher.Addresser, error) {
	return w.newAddresses(account, bip44.ChangeChainIndex, n)
}

// externalEntries returns the entries on external chain
func (w *Wallet) externalEntries(account uint32) (wallet.Entries, error) {
	return w.entries(account, bip44.ExternalChainIndex)
}

//
//// ChangeEntries returns the entries on change chain
//func (w *Wallet) ChangeEntries(account uint32) (wallet.Entries, error) {
//	return w.accounts.entries(account, bip44.ChangeChainIndex)
//}
//
//// ExternalEntriesLen returns the external chain entries length of selected account
//func (w *Wallet) ExternalEntriesLen(account uint32) (uint32, error) {
//	return w.accounts.entriesLen(account, bip44.ExternalChainIndex)
//}
//
//// ChangeEntriesLen returns the change chain entries length of selected account
//func (w *Wallet) ChangeEntriesLen(account uint32) (uint32, error) {
//	return w.accounts.entriesLen(account, bip44.ChangeChainIndex)
//}
//
//// ExternalEntryAt returns the entry at the given index on external chain of selected account
//func (w *Wallet) ExternalEntryAt(account, i uint32) (wallet.Entry, error) {
//	return w.accounts.entryAt(account, bip44.ExternalChainIndex, i)
//}
//
//// ChangeEntryAt returns the entry at the given index on change chain of selected account
//func (w *Wallet) ChangeEntryAt(account, i uint32) (wallet.Entry, error) {
//	return w.accounts.entryAt(account, bip44.ChangeChainIndex, i)
//}
//
//// GetEntry returns the entry of given address on selected account
//func (w *Wallet) GetEntry(account uint32, address cipher.Addresser) (wallet.Entry, bool, error) {
//	return w.accounts.getEntry(account, address)
//}

// Serialize encodes the bip44 wallet to []byte
func (w Wallet) Serialize() ([]byte, error) {
	if w.decoder == nil {
		w.decoder = defaultWalletDecoder
	}
	return w.decoder.Encode(&w)
}

// Deserialize decodes the []byte to a bip44 wallet
func (w *Wallet) Deserialize(b []byte) error {
	if w.decoder == nil {
		w.decoder = defaultWalletDecoder
	}
	toW, err := w.decoder.Decode(b)
	if err != nil {
		return err
	}
	toW2 := toW.(*Wallet)

	toW2.decoder = w.decoder
	*w = *toW2
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
		return wallet.ErrMissingPassword
	}

	if w.IsEncrypted() {
		return wallet.ErrWalletEncrypted
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
		return nil, wallet.ErrWalletNotEncrypted
	}

	if len(password) == 0 {
		return nil, wallet.ErrMissingPassword
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
		return nil, wallet.ErrInvalidPassword
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

// Fingerprint returns a unique ID fingerprint for this wallet, composed of its wallet type and initial address
func (w Wallet) Fingerprint() string {
	addr := ""
	entries, err := w.externalEntries(0)
	if err != nil {
		logger.WithError(err).Panic("Fingerprint get external entries failed")
		return ""
	}

	if len(entries) == 0 {
		if !w.IsEncrypted() {
			addrs, err := w.newExternalAddresses(0, 1)
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
		Meta:           w.Meta.Clone(),
		accountManager: w.accountManager.clone(),
		decoder:        w.decoder,
	}
}

func (w *Wallet) CopyFrom(src wallet.Wallet) {
	w.copyFrom(src.(*Wallet))
}

func (w *Wallet) copyFrom(wlt *Wallet) {
	w.Meta = wlt.Meta.Clone()
	w.accountManager = wlt.accountManager.clone()
	w.decoder = wlt.decoder
}

// CopyFromRef copies the src wallet with a pointer dereference
func (w *Wallet) CopyFromRef(src wallet.Wallet) {
	*w = *(src.(*Wallet))
}

// Accounts returns the list of accounts
func (w *Wallet) Accounts() []wallet.Bip44Account {
	return w.accountManager.all()
}

// GetEntries provides entries service to access the external chain of given account
func (w *Wallet) GetEntries(options ...wallet.Option) (wallet.Entries, error) {
	opts := getBip44Options(options...)
	return w.entries(opts.Account, opts.Change)
}

// Erase wipes all sensitive data
func (w *Wallet) Erase() {
	w.SetSeed("")
	w.SetSeedPassphrase("")
	w.accountManager.erase()
}

// syncSecrets synchronize the secrets with all addresses, ensure that
// each address has the secret key stored in the secrets
func (w Wallet) syncSecrets(ss wallet.Secrets) error {
	return w.accountManager.syncSecrets(ss)
}

// packSecrets saves all sensitive data to the secrets map.
func (w Wallet) packSecrets(ss wallet.Secrets) {
	ss.Set(wallet.SecretSeed, w.Meta.Seed())
	ss.Set(wallet.SecretSeedPassphrase, w.Meta.SeedPassphrase())
	w.accountManager.packSecrets(ss)
}

func (w *Wallet) unpackSecrets(ss wallet.Secrets) error {
	seed, ok := ss.Get(wallet.SecretSeed)
	if !ok {
		return errors.New("seed does not exist in secrets")
	}
	w.Meta.SetSeed(seed)

	passphrase, _ := ss.Get(wallet.SecretSeedPassphrase)
	w.Meta.SetSeedPassphrase(passphrase)

	return w.accountManager.unpackSecrets(ss)
}

func getBip44Options(options ...wallet.Option) *wallet.Bip44EntriesOptions {
	v := &wallet.Bip44EntriesOptions{}
	for _, opt := range options {
		opt(v)
	}

	return v
}

// ScanAddresses scans both the external and change addresses to find addresses with
// transactions.
// Only external addresses will be returned.
func (w *Wallet) ScanAddresses(scanN uint64, tf wallet.TransactionsFinder) ([]cipher.Addresser, error) {
	if scanN == 0 {
		return nil, nil
	}

	w2 := w.Clone().(*Wallet)

	accounts := w2.Accounts()
	scanAddresses := func(account, chain uint32) ([]cipher.Addresser, int, int, error) {
		nExistingAddrs, err := w2.entriesLen(account, chain)
		if err != nil {
			return nil, 0, 0, err
		}

		// generates the addresses  to scan
		addrs, err := w2.accountManager.newAddresses(account, chain, uint32(scanN))
		if err != nil {
			return nil, 0, 0, err
		}

		// finds if these addresses had any activity
		active, err := tf.AddressesActivity(addrs)
		if err != nil {
			return nil, 0, 0, err
		}

		// checks activity from the last one until we find the address that has activity
		var keepNum uint64
		for i := len(active) - 1; i >= 0; i-- {
			if active[i] {
				keepNum = uint64(i + 1)
				break
			}
		}

		return addrs[:keepNum], int(nExistingAddrs), int(keepNum), nil
	}

	// [accounts][chains] array
	generateAddresses := make([][]uint32, len(accounts))

	// only external addresses will be returned
	var retAddrs []cipher.Addresser

	for i, a := range accounts {
		addrs, initLen, keepNum, err := scanAddresses(a.Index, bip44.ExternalChainIndex)
		if err != nil {
			return nil, err
		}

		retAddrs = append(retAddrs, addrs...)

		generateAddresses[i] = append(generateAddresses[i], uint32(initLen+keepNum))

		_, initLen, keepNum, err = scanAddresses(a.Index, bip44.ChangeChainIndex)
		if err != nil {
			return nil, err
		}

		generateAddresses[i] = append(generateAddresses[i], uint32(initLen+keepNum))
	}

	w2.reset()
	for i, a := range accounts {
		// generate addresses on external chains
		for _, c := range []uint32{bip44.ExternalChainIndex, bip44.ChangeChainIndex} {
			_, err := w2.newAddresses(a.Index, c, generateAddresses[i][c])
			if err != nil {
				return nil, err
			}
		}
	}

	*w = *w2

	return retAddrs, nil
}

// GetAddresses returns all addresses on selected account and chain,
// if no options ware provided, addresses on external chain of account 0 will be returned.
func (w *Wallet) GetAddresses(options ...wallet.Option) ([]cipher.Addresser, error) {
	opts := getBip44Options(options...)
	entries, err := w.entries(opts.Account, opts.Change)
	if err != nil {
		return nil, err
	}

	var addrs []cipher.Addresser
	for _, e := range entries {
		addrs = append(addrs, e.Address)
	}

	return addrs, nil
}

// GenerateAddresses generates addresses on selected account and chain,
// if no options are provided, addresses will be generated on external chain of account 0.
func (w *Wallet) GenerateAddresses(num uint64, options ...wallet.Option) ([]cipher.Addresser, error) {
	opts := getBip44Options(options...)

	return w.newAddresses(opts.Account, opts.Change, uint32(num))
}

// GetEntryAt returns the entry at specific index
func (w *Wallet) GetEntryAt(i int, options ...wallet.Option) (wallet.Entry, error) {
	opts := getBip44Options(options...)
	return w.entryAt(opts.Account, opts.Change, uint32(i))
}

// GetEntry returns the entry of given address on selected account and chain,
// if no options are provided, check the external chain of account 0.
func (w *Wallet) GetEntry(addr cipher.Addresser, options ...wallet.Option) (wallet.Entry, error) {
	opts := getBip44Options(options...)
	e, ok, err := w.getEntry(opts.Account, addr)
	if err != nil {
		return wallet.Entry{}, err
	}

	if !ok {
		return wallet.Entry{}, wallet.ErrEntryNotFound
	}

	return e, nil
}

// HasEntry checks whether the entry of given address exists on selected account and chain,
// if no options are provided, check the external chain of account 0.
func (w *Wallet) HasEntry(addr cipher.Addresser, options ...wallet.Option) (bool, error) {
	opts := getBip44Options(options...)
	_, ok, err := w.getEntry(opts.Account, addr)
	if err != nil {
		return false, err
	}

	return ok, nil
}

// EntriesLen returns the entries length of selected account and chain,
// if no options are provided, entries length of external chain on account 0 will
// be returned.
func (w *Wallet) EntriesLen(options ...wallet.Option) (int, error) {
	opts := getBip44Options(options...)
	l, err := w.entriesLen(opts.Account, opts.Change)
	return int(l), err
}

func (w *Wallet) reset() {
	w.accountManager.reset()
}

// PeekChangeAddress returns the last entry address on change chain if
// no transactions are found, otherwise, return with a new address.
func (w *Wallet) PeekChangeAddress(tf wallet.TransactionsFinder, options ...wallet.Option) (cipher.Addresser, error) {
	options = append(options, wallet.OptionChange(true))
	entries, err := w.GetEntries(options...)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		// generate a new address and return
		addrs, err := w.GenerateAddresses(1, options...)
		if err != nil {
			return nil, err
		}

		return addrs[0], nil
	}

	// checks if the last entry address has transactions
	addr := entries[len(entries)-1].Address
	oks, err := tf.AddressesActivity([]cipher.Addresser{addr})
	if err != nil {
		return nil, err
	}

	if oks[0] == false {
		return addr, nil
	}

	// generate a new address and return it
	addrs, err := w.GenerateAddresses(1, options...)
	if err != nil {
		return nil, err
	}
	return addrs[0], nil
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

// Loader implements the wallet.Loader interface
type Loader struct{}

// Load implements the Load method of the wallet.Loader interface
func (l Loader) Load(data []byte) (wallet.Wallet, error) {
	w := &Wallet{}
	if err := w.Deserialize(data); err != nil {
		return nil, err
	}

	return w, nil
}

// Creator implements the wallet.Creator interface
type Creator struct{}

// Create implements the Create method of wallet.Creator interface
func (c Creator) Create(filename, label, seed string, options wallet.Options) (wallet.Wallet, error) {
	opts := convertOptions(options)
	return NewWallet(
		filename,
		label,
		seed,
		options.SeedPassphrase,
		opts...)
}

// convertOptions collects the cared fields from wallet.Options
// and converts them to an wallet.Option slice
func convertOptions(options wallet.Options) []wallet.Option {
	var opts []wallet.Option

	if options.Coin != "" {
		opts = append(opts, wallet.OptionCoinType(options.Coin))
	}

	if options.Bip44Coin != nil {
		opts = append(opts, wallet.OptionBip44Coin(options.Bip44Coin))
	}

	if options.CryptoType != "" {
		opts = append(opts, wallet.OptionCryptoType(options.CryptoType))
	}

	if options.Decoder != nil {
		opts = append(opts, wallet.OptionDecoder(options.Decoder))
	}

	if options.Encrypt {
		opts = append(opts, wallet.OptionEncrypt(true))
		opts = append(opts, wallet.OptionPassword(options.Password))
	}

	if options.GenerateN > 0 {
		opts = append(opts, wallet.OptionGenerateN(options.GenerateN))
	}

	if options.ScanN > 0 {
		opts = append(opts, wallet.OptionScanN(options.ScanN))
		opts = append(opts, wallet.OptionTransactionsFinder(options.TF))
	}

	return opts
}
