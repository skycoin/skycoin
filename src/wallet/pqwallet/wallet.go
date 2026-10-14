// Package pqwallet implements post-quantum wallets. Key i is the ML-DSA-44 seed
// pqsig.DeriveSecKey(seed, i), and its address is a version 1 Skycoin address.
package pqwallet

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/cipher/pqsig"
	"github.com/skycoin/skycoin/src/wallet"
)

// WalletType is the post-quantum wallet type
const WalletType = wallet.WalletTypePQ

// ErrCoinType is returned for a coin other than Skycoin, which is the only coin with
// post-quantum addresses
var ErrCoinType = wallet.NewError(errors.New("post-quantum wallets only hold skycoin addresses"))

var defaultWalletDecoder = &JSONDecoder{}

func init() {
	if err := wallet.RegisterCreator(WalletType, &Creator{}); err != nil {
		panic(err)
	}

	if err := wallet.RegisterLoader(WalletType, &Loader{}); err != nil {
		panic(err)
	}
}

// Wallet holds post-quantum keys derived from one seed
type Wallet struct {
	wallet.Meta
	entries wallet.Entries
	decoder wallet.Decoder
}

// NewWallet creates a post-quantum wallet
func NewWallet(filename, label, seed string, options ...wallet.Option) (*Wallet, error) {
	if label == "" {
		return nil, wallet.ErrMissingLabel
	}

	if seed == "" {
		return nil, wallet.ErrMissingSeed
	}

	wlt := &Wallet{
		Meta: wallet.Meta{
			wallet.MetaFilename:   filename,
			wallet.MetaLabel:      label,
			wallet.MetaSeed:       seed,
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

	if err := validateMeta(wlt.Meta); err != nil {
		return nil, err
	}

	generateN := advOpts.GenerateN
	if generateN > 0 {
		if _, err := wlt.GenerateAddresses(wallet.OptionGenerateN(generateN)); err != nil {
			return nil, err
		}
	}

	scanN := advOpts.ScanN
	if scanN > 0 {
		if advOpts.TF == nil {
			return nil, errors.New("missing transaction finder for scanning addresses")
		}

		if scanN > generateN {
			scanN -= generateN
		}

		if _, err := wlt.ScanAddresses(scanN, advOpts.TF); err != nil {
			return nil, err
		}
	}

	if advOpts.Encrypt {
		if wlt.IsTemp() {
			return nil, wallet.ErrEncryptTempWallet
		}

		if len(advOpts.Password) == 0 {
			return nil, wallet.ErrMissingPassword
		}

		if err := wlt.Lock(advOpts.Password); err != nil {
			return nil, err
		}
	} else if len(advOpts.Password) > 0 {
		return nil, wallet.ErrMissingEncrypt
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

	if m.Coin() != wallet.CoinTypeSkycoin {
		return ErrCoinType
	}

	if err := wallet.ValidateMeta(m); err != nil {
		return err
	}

	if err := wallet.ValidateMetaCryptoType(m); err != nil {
		return err
	}

	return wallet.ValidateMetaSeed(m)
}

// newEntry derives the entry at index from seed
func newEntry(seed string, index uint32) (wallet.Entry, error) {
	sk := pqsig.DeriveSecKey([]byte(seed), index)
	pk, err := sk.PubKey()
	if err != nil {
		return wallet.Entry{}, err
	}

	e := wallet.Entry{
		Address:     cipher.AddressFromPQPubKey(pk),
		ChildNumber: index,
	}
	copy(e.Secret[:], sk)
	return e, nil
}

// SetDecoder sets the decoder
func (w *Wallet) SetDecoder(d wallet.Decoder) {
	w.decoder = d
}

// Serialize encodes the wallet into bytes
func (w Wallet) Serialize() ([]byte, error) {
	if w.decoder == nil {
		w.decoder = defaultWalletDecoder
	}

	return w.decoder.Encode(&w)
}

// Deserialize decodes the wallet from bytes
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

// Lock encrypts the sensitive data of the wallet
func (w *Wallet) Lock(password []byte) error {
	if w.IsTemp() {
		return wallet.ErrEncryptTempWallet
	}

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
		cryptoType = crypto.DefaultCryptoType
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

// Unlock decrypts the wallet and returns a decrypted copy
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

func (w *Wallet) packSecrets(ss wallet.Secrets) {
	ss.Set(wallet.SecretSeed, w.Seed())
	for _, e := range w.entries {
		ss.Set(e.Address.String(), e.Secret.Hex())
	}
}

func (w *Wallet) unpackSecrets(ss wallet.Secrets) error {
	seed, ok := ss.Get(wallet.SecretSeed)
	if !ok {
		return errors.New("seed doesn't exist in wallet")
	}
	w.SetSeed(seed)

	return w.entries.UnpackSecretKeys(ss)
}

// Fingerprint returns a unique ID for the wallet, made of its type and first address
func (w *Wallet) Fingerprint() string {
	addr := ""
	if len(w.entries) == 0 {
		if !w.IsEncrypted() {
			if e, err := newEntry(w.Meta.Seed(), 0); err == nil {
				addr = e.Address.String()
			}
		}
	} else {
		addr = w.entries[0].Address.String()
	}
	return fmt.Sprintf("%s-%s", w.Type(), addr)
}

// Clone returns a copy of the wallet
func (w *Wallet) Clone() wallet.Wallet {
	return &Wallet{
		Meta:    w.Meta.Clone(),
		entries: w.entries.Clone(),
		decoder: w.decoder,
	}
}

func (w *Wallet) copyFrom(src *Wallet) {
	w.Meta = src.Meta.Clone()
	w.entries = src.entries.Clone()
}

// CopyFromRef copies the src wallet with a pointer dereference
func (w *Wallet) CopyFromRef(src wallet.Wallet) {
	*w = *(src.(*Wallet))
}

// Accounts is not implemented, it is an interface for bip44 wallet.
func (w *Wallet) Accounts() []wallet.Bip44Account {
	return nil
}

// Erase wipes secret fields in wallet
func (w *Wallet) Erase() {
	w.Meta.EraseSeeds()
	w.entries.Erase()
}

// Validate validates the wallet
func (w *Wallet) Validate() error {
	if err := w.Meta.Validate(); err != nil {
		return err
	}

	if !wallet.IsValidWalletType(w.Meta.Type()) {
		return wallet.ErrInvalidWalletType
	}

	if !w.IsEncrypted() && w.Seed() == "" {
		return errors.New("seed missing in unencrypted post-quantum wallet")
	}
	return nil
}

// ScanAddresses scans ahead N addresses, truncating up to the highest address with any transaction history.
func (w *Wallet) ScanAddresses(scanN uint64, tf wallet.TransactionsFinder) ([]cipher.Addresser, error) {
	if w.IsEncrypted() {
		return nil, wallet.ErrWalletEncrypted
	}

	if scanN == 0 {
		return nil, nil
	}

	w2 := w.Clone().(*Wallet)

	addrs, err := w2.GenerateAddresses(wallet.OptionGenerateN(scanN))
	if err != nil {
		return nil, err
	}

	active, err := tf.AddressesActivity(addrs)
	if err != nil {
		return nil, err
	}

	var keepNum int
	for i := len(active) - 1; i >= 0; i-- {
		if active[i] {
			keepNum = i + 1
			break
		}
	}

	w.entries = w2.entries[:len(w.entries)+keepNum]
	return addrs[:keepNum], nil
}

// GenerateAddresses generates N addresses
func (w *Wallet) GenerateAddresses(options ...wallet.Option) ([]cipher.Addresser, error) {
	if w.Meta.IsEncrypted() {
		return nil, wallet.ErrWalletEncrypted
	}

	num := wallet.GetGenerateNFromOptions(options...)
	if num == 0 {
		return nil, nil
	}

	start := len(w.entries)
	if uint64(start)+num > uint64(^uint32(0)) {
		return nil, errors.New("too many addresses")
	}

	addrs := make([]cipher.Addresser, 0, num)
	for i := uint64(0); i < num; i++ {
		e, err := newEntry(w.Meta.Seed(), uint32(uint64(start)+i)) //nolint:gosec // bounded above
		if err != nil {
			return nil, err
		}
		w.entries = append(w.entries, e)
		addrs = append(addrs, e.Address)
	}
	return addrs, nil
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

// HasEntry returns true if the wallet has an Entry with a given cipher.Address.
func (w *Wallet) HasEntry(a cipher.Addresser, _ ...wallet.Option) (bool, error) {
	return w.entries.Has(a), nil
}

// EntriesLen returns the number of entries in the wallet
func (w *Wallet) EntriesLen(_ ...wallet.Option) (int, error) {
	return len(w.entries), nil
}

// Loader implements the wallet.Loader interface
type Loader struct{}

// Load loads a post-quantum wallet from bytes
func (l Loader) Load(data []byte) (wallet.Wallet, error) {
	w := &Wallet{}
	if err := w.Deserialize(data); err != nil {
		return nil, err
	}

	return w, nil
}

// Creator implements the wallet.Creator interface
type Creator struct{}

// Create creates a post-quantum wallet
func (c Creator) Create(filename, label, seed string, options wallet.Options) (wallet.Wallet, error) {
	return NewWallet(filename, label, seed, convertOptions(options)...)
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

	if options.GenerateN > 0 {
		opts = append(opts, wallet.OptionGenerateN(options.GenerateN))
	}

	if options.ScanN > 0 {
		opts = append(opts, wallet.OptionScanN(options.ScanN))
		opts = append(opts, wallet.OptionTransactionsFinder(options.TF))
	}

	if options.Temp {
		opts = append(opts, wallet.OptionTemp(true))
	}

	return opts
}
