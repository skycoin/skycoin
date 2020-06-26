package collection

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
)

const (
	// WalletType represents the collection wallet type
	WalletType = "collection"
)

var defaultWalletDecoder = &JSONDecoder{}

// Wallet manages keys as an arbitrary collection.
// It has no defined keypair generator. The only way to add keys to the
// wallet is to explicitly add them.
// This wallet does not support address scanning or generation.
// This wallet does not use seeds.
type Wallet struct {
	wallet.Meta
	entries wallet.Entries
	decoder wallet.Decoder
}

// NewWallet creates a collection wallet
func NewWallet(filename, label string, options ...wallet.Option) (*Wallet, error) {
	var wlt = &Wallet{
		Meta: wallet.Meta{
			wallet.MetaFilename:   filename,
			wallet.MetaLabel:      label,
			wallet.MetaEncrypted:  "false",
			wallet.MetaType:       WalletType,
			wallet.MetaVersion:    wallet.Version,
			wallet.MetaCoin:       string(wallet.CoinTypeSkycoin),
			wallet.MetaCryptoType: string(crypto.DefaultCryptoType),
			wallet.MetaTimestamp:  strconv.FormatInt(time.Now().Unix(), 10),
		},
		entries: wallet.Entries{},
		decoder: defaultWalletDecoder,
	}

	advOpts := &wallet.AdvancedOptions{}
	for _, opt := range options {
		opt(wlt)
		opt(advOpts)
	}

	// validateMeta wallet before encryption
	if err := validateMeta(wlt.Meta); err != nil {
		return nil, err
	}

	if advOpts.GenerateN != 0 || advOpts.ScanN != 0 {
		return nil, wallet.NewError(fmt.Errorf("wallet scanning is not defined for %q wallet", WalletType))
	}

	if advOpts.Encrypt {
		if len(advOpts.Password) == 0 {
			return nil, wallet.ErrMissingPassword
		}

		if err := wlt.Lock(advOpts.Password); err != nil {
			return nil, err
		}
	} else {
		if len(advOpts.Password) > 0 {
			return nil, wallet.ErrMissingEncrypt
		}
	}

	if err := validateMeta(wlt.Meta); err != nil {
		return nil, err
	}

	return wlt, nil
}

func validateMeta(m wallet.Meta) error {
	if m[wallet.MetaType] != WalletType {
		return wallet.ErrInvalidWalletType
	}

	if m[wallet.MetaSeed] != "" {
		return wallet.NewError(fmt.Errorf("seed should not be provided for %q wallets", WalletType))
	}
	return wallet.ValidateMeta(m)
}

// SetDecoder sets the decoder
func (w *Wallet) SetDecoder(d wallet.Decoder) {
	w.decoder = d
}

// Serialize encode the wallet to byte slice
func (w Wallet) Serialize() ([]byte, error) {
	if w.decoder == nil {
		w.decoder = defaultWalletDecoder
	}

	return w.decoder.Encode(&w)
}

// Deserialize decodes wallet from byte slice
func (w *Wallet) Deserialize(data []byte) error {
	if w.decoder == nil {
		w.decoder = defaultWalletDecoder
	}

	wlt, err := w.decoder.Decode(data)
	if err != nil {
		return err
	}

	w2 := wlt.(*Wallet)
	w2.decoder = w.decoder
	*w = *w2
	return nil
}

// Lock encrypts the wallet secrets
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

	cryptoType := wlt.CryptoType()
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

	wlt.SetEncrypted(cryptoType, string(encSecret))

	wlt.Erase()
	w.Erase()
	w.copyFrom(wlt)
	return nil
}

// Unlock unlocks the encrypted wallet
func (w *Wallet) Unlock(password []byte) (wallet.Wallet, error) {
	if !w.IsEncrypted() {
		return nil, wallet.ErrWalletNotEncrypted
	}

	if len(password) == 0 {
		return nil, wallet.ErrMissingPassword
	}

	sstr := w.Secrets()
	if sstr == "" {
		return nil, errors.New("missing secrets")
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
	if err := cw.unpackSecrets(ss); err != nil {
		return nil, err
	}
	cw.SetDecrypted()
	return cw, nil
}

// packSecrets copies data from decrypted wallets into the secrets container
func (w *Wallet) packSecrets(ss wallet.Secrets) {
	// ss.Set(wallet.SecretSeed, w.Meta.Seed())
	// ss.Set(wallet.SecretLastSeed, w.Meta.LastSeed())

	// Saves entry secret keys in secrets
	for _, e := range w.entries {
		ss.Set(e.Address.String(), e.Secret.Hex())
	}
}

// UnpackSecrets copies data from decrypted secrets into the wallet
func (w *Wallet) unpackSecrets(ss wallet.Secrets) error {
	return w.entries.UnpackSecretKeys(ss)
}

// Fingerprint returns an empty string; fingerprints are only defined for
// wallets with a seed
func (w *Wallet) Fingerprint() string {
	return ""
}

// Clone clones the wallet a new wallet object
func (w *Wallet) Clone() wallet.Wallet {
	return &Wallet{
		Meta:    w.Meta.Clone(),
		entries: w.entries.Clone(),
		decoder: w.decoder,
	}
}

// copyFrom copies the src wallet by reallocating
func (w *Wallet) copyFrom(src *Wallet) {
	w.Meta = src.Meta.Clone()
	w.entries = src.entries.Clone()
}

// CopyFromRef copies the src wallet with a pointer dereference
func (w *Wallet) CopyFromRef(src wallet.Wallet) {
	*w = *(src.(*Wallet))
}

// Accounts is not defined for collection wallet
func (w *Wallet) Accounts() []wallet.Bip44Account {
	return nil
}

// Erase wipes secret fields in wallet
func (w *Wallet) Erase() {
	// w.Meta.EraseSeeds()
	w.entries.Erase()
}

// Validate validates the wallet
func (w *Wallet) Validate() error {
	if err := w.Meta.Validate(); err != nil {
		return err
	}

	if w.Type() != WalletType {
		return wallet.ErrInvalidWalletType
	}

	if s := w.Meta[wallet.MetaSeed]; s != "" {
		return errors.New("seed should not be in collection wallets")
	}

	if s := w.Meta[wallet.MetaLastSeed]; s != "" {
		return errors.New("lastSeed should not be in collection wallets")
	}
	return nil
}

// ScanAddresses is a no-op for "collection" wallets
func (w *Wallet) ScanAddresses(scanN uint64, tf wallet.TransactionsFinder) ([]cipher.Addresser, error) {
	return nil, wallet.NewError(errors.New("A collection wallet does not implement ScanAddresses"))
}

// GenerateAddresses is a no-op for "collection" wallets
func (w *Wallet) GenerateAddresses(num uint64, _ ...wallet.Option) ([]cipher.Addresser, error) {
	return nil, wallet.NewError(errors.New("A collection wallet does not implement GenerateAddresses"))
}

// GetAddresses returns all addresses in wallet
func (w *Wallet) GetAddresses(_ ...wallet.Option) ([]cipher.Addresser, error) {
	return w.entries.GetAddresses(), nil
}

// GetEntries returns a copy of all entries held by the wallet
func (w *Wallet) GetEntries(_ ...wallet.Option) (wallet.Entries, error) {
	return w.entries.Clone(), nil
}

// GetEntryAt returns entry at a given index in the entries array
func (w *Wallet) GetEntryAt(i int, _ ...wallet.Option) (wallet.Entry, error) {
	if i < 0 || i >= len(w.entries) {
		return wallet.Entry{}, fmt.Errorf("entry index %d is out of range", i)
	}
	return w.entries[i], nil
}

// GetEntry returns entry of given address
func (w *Wallet) GetEntry(a cipher.Addresser, _ ...wallet.Option) (wallet.Entry, error) {
	e, ok := w.entries.Get(a)
	if !ok {
		return wallet.Entry{}, wallet.ErrEntryNotFound
	}
	return e, nil
}

// HasEntry returns true if the wallet has an entry.Entry with a given cipher.Address.
func (w *Wallet) HasEntry(a cipher.Addresser, _ ...wallet.Option) (bool, error) {
	return w.entries.Has(a), nil
}

// EntriesLen returns the number of entries in the wallet
func (w *Wallet) EntriesLen(_ ...wallet.Option) (int, error) {
	return len(w.entries), nil
}

//// GenerateSkycoinAddresses is a no-op for "collection" wallets
//func (w *Wallet) GenerateSkycoinAddresses(num uint64) ([]cipher.Address, error) {
//	return nil, wallet.NewError(errors.New("A collection wallet does not implement GenerateSkycoinAddresses"))
//}
//
//// GetSkycoinAddresses returns all Skycoin addresses in wallet. The wallet's coin type must be Skycoin.
//func (w *Wallet) GetSkycoinAddresses() ([]cipher.Address, error) {
//	if w.Meta.Coin() != wallet.CoinTypeSkycoin {
//		return nil, errors.New("Wallet coin type is not skycoin")
//	}
//
//	return w.Entries.GetSkycoinAddresses(), nil
//}

// AddEntry adds a new entry to the wallet.
func (w *Wallet) AddEntry(e wallet.Entry) error {
	if w.IsEncrypted() {
		return wallet.ErrWalletEncrypted
	}

	if err := e.Verify(); err != nil {
		return err
	}

	for _, entry := range w.entries {
		if e.SkycoinAddress() == entry.SkycoinAddress() {
			return errors.New("wallet already contains entry with this address")
		}
	}

	w.entries = append(w.entries, e)
	return nil
}

// Loader implements the wallet.Loader interface
type Loader struct{}

// Load loads wallet from byte slice
func (l Loader) Load(data []byte) (wallet.Wallet, error) {
	w := &Wallet{}
	if err := w.Deserialize(data); err != nil {
		return nil, err
	}
	return w, nil
}

// Creator implements the wallet.Creator interface
type Creator struct{}

// Create implements the wallet.Creator interface
func (c Creator) Create(filename, label string, options wallet.Options) (wallet.Wallet, error) {
	return NewWallet(filename, label, convertOptions(options)...)
}

func convertOptions(options wallet.Options) []wallet.Option {
	var opts []wallet.Option
	if options.Coin != "" {
		opts = append(opts, wallet.OptionCoinType(options.Coin))
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

	return opts
}
