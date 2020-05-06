package deterministic

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
)

const walletType = "deterministic"

var defaultWalletDecoder = &JSONDecoder{}

// Wallet manages keys using the original Skycoin deterministic
// keypair generator method.
// With this generator, a single chain of addresses is created, each one dependent
// on the previous.
type Wallet struct {
	wallet.Meta
	entries wallet.Entries
	decoder wallet.Decoder
}

func NewWallet(filename, label, seed string, options ...wallet.Option) (*Wallet, error) {
	var wlt = &Wallet{
		Meta: wallet.Meta{
			wallet.MetaFilename:   filename,
			wallet.MetaLabel:      label,
			wallet.MetaSeed:       seed,
			wallet.MetaLastSeed:   seed,
			wallet.MetaEncrypted:  "false",
			wallet.MetaType:       walletType,
			wallet.MetaVersion:    wallet.Version,
			wallet.MetaCoin:       string(wallet.CoinTypeSkycoin),
			wallet.MetaCryptoType: string(crypto.DefaultCryptoType),
			wallet.MetaTimestamp:  strconv.FormatInt(time.Now().Unix(), 10),
		},
		entries: wallet.Entries{},
		decoder: defaultWalletDecoder,
	}

	moreOpts := moreOptions{}
	for _, opt := range options {
		opt(wlt)
		opt(moreOpts)
	}

	// validateMeta wallet before encrypting
	if err := validateMeta(wlt.Meta); err != nil {
		return nil, err
	}

	generateN := moreOpts.GenerateN
	if generateN > 0 {
		_, err := wlt.Entries().GenerateAddresses(generateN)
		if err != nil {
			return nil, err
		}
	}

	scanN := moreOpts.ScanN
	if scanN > 0 {
		if moreOpts.TF == nil {
			return nil, errors.New("missing transaction finder for scanning addresses")
		}

		if scanN > generateN {
			scanN = scanN - generateN
		}

		_, err := wlt.ScanAddresses(scanN, moreOpts.TF)
		if err != nil {
			return nil, err
		}
	}

	// encrypts wallet if options.Encrypt is true
	if moreOpts.Encrypt {
		if len(moreOpts.Password) == 0 {
			return nil, errors.New("missing password for encrypting wallet")
		}

		if err := wlt.Lock(moreOpts.Password); err != nil {
			return nil, err
		}
	}

	// validateMeta again after encrypted
	if err := validateMeta(wlt.Meta); err != nil {
		return nil, err
	}

	return wlt, nil
}

func validateMeta(m wallet.Meta) error {
	if m[wallet.MetaType] != walletType {
		return errors.New("invalid wallet type")
	}

	return wallet.ValidateMeta(m)
}

func (w Wallet) Serialize() ([]byte, error) {
	if w.decoder == nil {
		w.decoder = defaultWalletDecoder
	}

	return w.decoder.Encode(&w)
}

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

func (w *Wallet) Lock(password []byte) error {
	if len(password) == 0 {
		return errors.New("missing password when locking deterministic wallet")
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

	cryptoType := wlt.CryptoType()
	if cryptoType == "" {
		return errors.New("crypto type  field not set")
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

func (w *Wallet) Unlock(password []byte) (wallet.Wallet, error) {
	if !w.IsEncrypted() {
		return nil, errors.New("wallet is not encrypted")
	}

	if len(password) == 0 {
		return nil, errors.New("missing password")
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
	if err := cw.unpackSecrets(ss); err != nil {
		return nil, err
	}
	cw.SetDecrypted()
	return cw, nil
}

// packSecrets copies data from decrypted wallets into the wallet container
func (w *Wallet) packSecrets(ss wallet.Secrets) {
	ss.Set(wallet.SecretSeed, w.Seed())
	ss.Set(wallet.SecretLastSeed, w.LastSeed())

	// Saves entry secret keys in wallet
	for _, e := range w.entries {
		ss.Set(e.Address.String(), e.Secret.Hex())
	}
}

// unpackSecrets copies data from decrypted wallet into the wallet
func (w *Wallet) unpackSecrets(ss wallet.Secrets) error {
	seed, ok := ss.Get(wallet.SecretSeed)
	if !ok {
		return errors.New("seed doesn't exist in wallet")
	}
	w.SetSeed(seed)

	lastSeed, ok := ss.Get(wallet.SecretLastSeed)
	if !ok {
		return errors.New("lastSeed doesn't exist in wallet")
	}
	w.SetLastSeed(lastSeed)

	return w.entries.UnpackSecretKeys(ss)
}

// Fingerprint returns a unique ID fingerprint for this wallet, composed of its initial address
// and wallet type
func (w *Wallet) Fingerprint() string {
	addr := ""
	if len(w.entries) == 0 {
		if !w.IsEncrypted() {
			_, pk, _ := cipher.MustDeterministicKeyPairIterator([]byte(w.Meta.Seed()))
			addr = wallet.AddressConstructor(w.Meta)(pk).String()
		}
	} else {
		addr = w.entries[0].Address.String()
	}
	return fmt.Sprintf("%s-%s", w.Type(), addr)
}

// Clone clones the wallet a new wallet object
func (w *Wallet) Clone() wallet.Wallet {
	return &Wallet{
		Meta:    w.Meta.Clone(),
		entries: w.entries.Clone(),
	}
}

// copyFrom copies the src wallet to w
func (w *Wallet) copyFrom(src *Wallet) {
	w.Meta = src.Meta.Clone()
	w.entries = src.entries.Clone()
}

// CopyFromRef copies the src wallet with a pointer dereference
func (w *Wallet) CopyFromRef(src wallet.Wallet) {
	*w = *(src.(*Wallet))
}

func (w *Wallet) Accounts() []wallet.Bip44Account {
	return nil
}

func (w *Wallet) Entries(options ...wallet.Option) wallet.EntriesService {
	panic("implement me")
}

// Erase wipes secret fields in wallet
func (w *Wallet) Erase() {
	w.Meta.EraseSeeds()
	w.entries.Erase()
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

	nExistingAddrs := uint64(len(w2.entries))

	// Generate the addresses to scan
	addrs, err := w2.GenerateSkycoinAddresses(scanN)
	if err != nil {
		return nil, err
	}

	// Find if these addresses had any activity
	active, err := tf.AddressesActivity(addrs)
	if err != nil {
		return nil, err
	}

	// Check activity from the last one until we find the address that has activity
	var keepNum uint64
	for i := len(active) - 1; i >= 0; i-- {
		if active[i] {
			keepNum = uint64(i + 1)
			break
		}
	}

	// Regenerate addresses up to nExistingAddrs + nAddAddrs.
	// This is necessary to keep the lastSeed updated.
	w2.reset()
	if _, err := w2.GenerateSkycoinAddresses(nExistingAddrs + keepNum); err != nil {
		return nil, err
	}

	*w = *w2

	return nil, nil
}

// GenerateAddresses generates addresses
func (w *Wallet) GenerateAddresses(num uint64) ([]cipher.Addresser, error) {
	if w.Meta.IsEncrypted() {
		return nil, wallet.ErrWalletEncrypted
	}

	if num == 0 {
		return nil, nil
	}

	var seckeys []cipher.SecKey
	var seed []byte
	if len(w.entries) == 0 {
		seed, seckeys = cipher.MustGenerateDeterministicKeyPairsSeed([]byte(w.Meta.Seed()), int(num))
	} else {
		sd, err := hex.DecodeString(w.Meta.LastSeed())
		if err != nil {
			return nil, fmt.Errorf("decode hex seed failed: %v", err)
		}
		seed, seckeys = cipher.MustGenerateDeterministicKeyPairsSeed(sd, int(num))
	}

	w.Meta.SetLastSeed(hex.EncodeToString(seed))

	addrs := make([]cipher.Addresser, len(seckeys))
	makeAddress := wallet.AddressConstructor(w.Meta)
	for i, s := range seckeys {
		p := cipher.MustPubKeyFromSecKey(s)
		a := makeAddress(p)
		addrs[i] = a
		w.entries = append(w.entries, wallet.Entry{
			Address: a,
			Secret:  s,
			Public:  p,
		})
	}
	return addrs, nil
}

// Validate validates the wallet
func (w *Wallet) Validate() error {
	if err := w.Meta.Validate(); err != nil {
		return err
	}

	walletType := w.Meta.Type()
	if !wallet.IsValidWalletType(walletType) {
		return wallet.ErrInvalidWalletType
	}

	if !w.IsEncrypted() {
		if s := w.Seed(); s == "" {
			return errors.New("seed missing in unencrypted deterministic wallet")
		}

		if s := w.LastSeed(); s == "" {
			return errors.New("lastSeed missing in unencrypted deterministic wallet")
		}
	}
	return nil
}

// GetAddresses returns all addresses in wallet
func (w *Wallet) GetAddresses() []cipher.Addresser {
	return w.entries.GetAddresses()
}

// GetSkycoinAddresses returns all Skycoin addresses in wallet. The wallet's coin type must be Skycoin.
func (w *Wallet) GetSkycoinAddresses() ([]cipher.Address, error) {
	if w.Meta.Coin() != wallet.CoinTypeSkycoin {
		return nil, errors.New("Wallet coin type is not skycoin")
	}

	return w.entries.GetSkycoinAddresses(), nil
}

// GetEntries returns a copy of all entries held by the wallet
func (w *Wallet) GetEntries() wallet.Entries {
	return w.entries.Clone()
}

// EntriesLen returns the number of entries in the wallet
func (w *Wallet) EntriesLen() int {
	return len(w.entries)
}

// GetEntryAt returns entry at a given index in the entries array
func (w *Wallet) GetEntryAt(i int) wallet.Entry {
	return w.entries[i]
}

// GetEntry returns entry of given address
func (w *Wallet) GetEntry(a cipher.Address) (wallet.Entry, bool) {
	return w.entries.Get(a)
}

// HasEntry returns true if the wallet has an Entry with a given cipher.Address.
func (w *Wallet) HasEntry(a cipher.Address) bool {
	return w.entries.Has(a)
}

// GenerateSkycoinAddresses generates Skycoin addresses. If the wallet's coin type is not Skycoin, returns an error
func (w *Wallet) GenerateSkycoinAddresses(num uint64) ([]cipher.Address, error) {
	if w.Meta.Coin() != wallet.CoinTypeSkycoin {
		return nil, errors.New("GenerateSkycoinAddresses called for non-skycoin wallet")
	}

	addrs, err := w.GenerateAddresses(num)
	if err != nil {
		return nil, err
	}

	skyAddrs := make([]cipher.Address, len(addrs))
	for i, a := range addrs {
		skyAddrs[i] = a.(cipher.Address)
	}

	return skyAddrs, nil
}

// reset resets the wallet entries and move the lastSeed to origin
func (w *Wallet) reset() {
	w.entries = wallet.Entries{}
	w.Meta.SetLastSeed(w.Meta.Seed())
}

// Loader implements the wallet.Loader interface
type Loader struct{}

func (l Loader) Load(data []byte) (wallet.Wallet, error) {
	w := &Wallet{}
	if err := w.Deserialize(data); err != nil {
		return nil, err
	}

	return w, nil
}

func (l Loader) Type() string {
	return walletType
}

// Creator implements the wallet.Creator and wallet.Typer interface
type Creator struct{}

func (c Creator) Create(filename, label, seed string, options wallet.Options) (wallet.Wallet, error) {
	return NewWallet(
		filename,
		label,
		seed,
		convertOptions(options)...)
}

func (c Creator) Type() string {
	return walletType
}

func convertOptions(options wallet.Options) []wallet.Option {
	var opts []wallet.Option

	if options.Version != "" {
		opts = append(opts, Version(options.Version))
	}

	if options.Coin != "" {
		opts = append(opts, CoinType(options.Coin))
	}

	if options.CryptoType != "" {
		opts = append(opts, CryptoType(options.CryptoType))
	}

	if options.Decoder != nil {
		opts = append(opts, Decoder(options.Decoder))
	}

	if options.Encrypt {
		opts = append(opts, Encrypt(true))
		opts = append(opts, Password(options.Password))
	}

	if options.GenerateN > 0 {
		opts = append(opts, GenerateN(options.GenerateN))
	}

	if options.ScanN > 0 {
		opts = append(opts, ScanN(options.ScanN))
		opts = append(opts, TransactionsFinder(options.TF))
	}

	return opts
}

// entriesService wraps the wallet.Entries to implement the wallet.EntriesService interface
type entriesService struct {
	*Wallet
}

// Entries returns all entries in the wallet
func (es *entriesService) Entries() (wallet.Entries, error) {
	return es.entries.Clone(), nil
}

// GetEntryAt returns entry of the given index
func (es entriesService) GetEntryAt(i int) (wallet.Entry, error) {
	if i >= len(es.entries) {
		return wallet.Entry{}, fmt.Errorf("entry index %d is out of range", i)
	}

	return es.entries[i], nil
}

// GetEntry returns the entry by address
func (es entriesService) GetEntry(addresser cipher.Addresser) (wallet.Entry, error) {
	for _, e := range es.entries {
		if e.Address == addresser {
			return e, nil
		}
	}
	return wallet.Entry{}, wallet.ErrEntryNotFound
}

// HasEntry checks whether the address belongs to the wallet
func (es entriesService) HasEntry(addresser cipher.Addresser) (bool, error) {
	for _, e := range es.entries {
		if e.Address == addresser {
			return true, nil
		}
	}
	return false, nil
}

// Len returns the entries length
func (es entriesService) Len() (int, error) {
	return len(es.entries), nil
}

// GetAddresses returns all addresses in wallet
func (es entriesService) GetAddresses() ([]cipher.Addresser, error) {
	addrs := make([]cipher.Addresser, len(es.entries))
	for i, e := range es.entries {
		addrs[i] = e.Address
	}

	return addrs, nil
}

// GenerateAddresses generates addresses
func (es *entriesService) GenerateAddresses(num uint64) ([]cipher.Addresser, error) {
	return es.Wallet.GenerateAddresses(num)
}
