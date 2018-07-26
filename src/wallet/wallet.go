package wallet

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"encoding/hex"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"

	"github.com/shopspring/decimal"

	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/util/logging"
)

// Error wraps wallet related errors
type Error struct {
	error
}

// NewError creates an Error
func NewError(err error) error {
	if err == nil {
		return nil
	}
	return Error{err}
}

var (
	// Version represents the current wallet version
	Version = "0.2"

	logger = logging.MustGetLogger("wallet")

	// ErrInsufficientBalance is returned if a wallet does not have enough balance for a spend
	ErrInsufficientBalance = NewError(errors.New("balance is not sufficient"))
	// ErrInsufficientHours is returned if a wallet does not have enough hours for a spend with requested hours
	ErrInsufficientHours = NewError(errors.New("hours are not sufficient"))
	// ErrZeroSpend is returned if a transaction is trying to spend 0 coins
	ErrZeroSpend = NewError(errors.New("zero spend amount"))
	// ErrSpendingUnconfirmed is returned if caller attempts to spend unconfirmed outputs
	ErrSpendingUnconfirmed = NewError(errors.New("please spend after your pending transaction is confirmed"))
	// ErrInvalidEncryptedField is returned if a wallet's Meta.encrypted value is invalid.
	ErrInvalidEncryptedField = NewError(errors.New(`encrypted field value is not valid, must be "true", "false" or ""`))
	// ErrWalletEncrypted is returned when trying to generate addresses or sign tx in encrypted wallet
	ErrWalletEncrypted = NewError(errors.New("wallet is encrypted"))
	// ErrWalletNotEncrypted is returned when trying to decrypt unencrypted wallet
	ErrWalletNotEncrypted = NewError(errors.New("wallet is not encrypted"))
	// ErrMissingPassword is returned when trying to create wallet with encryption, but password is not provided.
	ErrMissingPassword = NewError(errors.New("missing password"))
	// ErrMissingEncrypt is returned when trying to create wallet with password, but options.Encrypt is not set.
	ErrMissingEncrypt = NewError(errors.New("missing encrypt"))
	// ErrInvalidPassword is returned if decrypts secrets failed
	ErrInvalidPassword = NewError(errors.New("invalid password"))
	// ErrMissingSeed is returned when trying to create wallet without a seed
	ErrMissingSeed = NewError(errors.New("missing seed"))
	// ErrMissingAuthenticated is returned if try to decrypt a scrypt chacha20poly1305 encrypted wallet, and find no authenticated metadata.
	ErrMissingAuthenticated = NewError(errors.New("missing authenticated metadata"))
	// ErrWrongCryptoType is returned when decrypting wallet with wrong crypto method
	ErrWrongCryptoType = NewError(errors.New("wrong crypto type"))
	// ErrWalletNotExist is returned if a wallet does not exist
	ErrWalletNotExist = NewError(errors.New("wallet doesn't exist"))
	// ErrSeedUsed is returned if a wallet already exists with the same seed
	ErrSeedUsed = NewError(errors.New("a wallet already exists with this seed"))
	// ErrWalletAPIDisabled is returned when trying to do wallet actions while the EnableWalletAPI option is false
	ErrWalletAPIDisabled = NewError(errors.New("wallet api is disabled"))
	// ErrSeedAPIDisabled is returned when trying to get seed of wallet while the EnableWalletAPI or EnableSeedAPI is false
	ErrSeedAPIDisabled = NewError(errors.New("wallet seed api is disabled"))
	// ErrWalletNameConflict represents the wallet name conflict error
	ErrWalletNameConflict = NewError(errors.New("wallet name would conflict with existing wallet, renaming"))
	// ErrInvalidHoursSelectionMode for invalid HoursSelection mode values
	ErrInvalidHoursSelectionMode = NewError(errors.New("invalid hours selection mode"))
	// ErrInvalidHoursSelectionType for invalid HoursSelection type values
	ErrInvalidHoursSelectionType = NewError(errors.New("invalid hours selection type"))
	// ErrUnknownAddress is returned if an address is not found in a wallet
	ErrUnknownAddress = NewError(errors.New("address not found in wallet"))
	// ErrUnknownUxOut is returned if a uxout is not owned by any address in a wallet
	ErrUnknownUxOut = NewError(errors.New("uxout is not owned by any address in the wallet"))
	// ErrNoUnspents is returned if a wallet has no unspents to spend
	ErrNoUnspents = NewError(errors.New("no unspents to spend"))
)

const (
	// WalletExt  wallet file extension
	WalletExt = "wlt"

	// WalletTimestampFormat  wallet timestamp layout
	WalletTimestampFormat = "2006_01_02"

	// CoinTypeSkycoin skycoin type
	CoinTypeSkycoin CoinType = "skycoin"
	// CoinTypeBitcoin bitcoin type
	CoinTypeBitcoin CoinType = "bitcoin"
)

// wallet meta fields
const (
	metaVersion    = "version"    // wallet version
	metaFilename   = "filename"   // wallet file name
	metaLabel      = "label"      // wallet label
	metaTm         = "tm"         // the timestamp when creating the wallet
	metaType       = "type"       // wallet type
	metaCoin       = "coin"       // coin type
	metaEncrypted  = "encrypted"  // whether the wallet is encrypted
	metaCryptoType = "cryptoType" // encrytion/decryption type
	metaSeed       = "seed"       // wallet seed
	metaLastSeed   = "lastSeed"   // seed for generating next address
	metaSecrets    = "secrets"    // secrets which records the encrypted seeds and secrets of address entries
)

// CoinType represents the wallet coin type
type CoinType string

// Options options that could be used when creating a wallet
type Options struct {
	Coin       CoinType   // coin type, skycoin, bitcoin, etc.
	Label      string     // wallet label.
	Seed       string     // wallet seed.
	Encrypt    bool       // whether the wallet need to be encrypted.
	Password   []byte     // password that would be used for encryption, and would only be used when 'Encrypt' is true.
	CryptoType CryptoType // wallet encryption type, scrypt-chacha20poly1305 or sha256-xor.
	ScanN      uint64     // number of addresses that're going to be scanned
}

const (
	// HoursSelectionTypeManual is used to specify manual hours selection in advanced spend
	HoursSelectionTypeManual = "manual"
	// HoursSelectionTypeAuto is used to specify automatic hours selection in advanced spend
	HoursSelectionTypeAuto = "auto"

	// HoursSelectionModeShare will distribute coin hours equally amongst destinations
	HoursSelectionModeShare = "share"
)

// HoursSelection defines options for hours distribution
type HoursSelection struct {
	Type        string
	Mode        string
	ShareFactor *decimal.Decimal
}

// CreateTransactionWalletParams defines a wallet to spend from and optionally which addresses in the wallet
type CreateTransactionWalletParams struct {
	ID        string
	UxOuts    []cipher.SHA256
	Addresses []cipher.Address
	Password  []byte
}

// CreateTransactionParams defines control parameters for transaction construction
type CreateTransactionParams struct {
	IgnoreUnconfirmed bool
	HoursSelection    HoursSelection
	Wallet            CreateTransactionWalletParams
	ChangeAddress     *cipher.Address
	To                []coin.TransactionOutput
}

// Validate validates CreateTransactionParams
func (c CreateTransactionParams) Validate() error {
	if c.ChangeAddress != nil && c.ChangeAddress.Null() {
		return NewError(errors.New("ChangeAddress must not be the null address"))
	}

	if len(c.To) == 0 {
		return NewError(errors.New("To is required"))
	}

	for _, to := range c.To {
		if to.Coins == 0 {
			return NewError(errors.New("To.Coins must not be zero"))
		}

		if to.Address.Null() {
			return NewError(errors.New("To.Address must not be the null address"))
		}
	}

	// Check for duplicate outputs, a transaction can't have outputs with
	// the same (address, coins, hours)
	// Auto mode would distribute hours to the outputs and could hypothetically
	// avoid assigning duplicate hours in many cases, but the complexity for doing
	// so is very high, so also reject duplicate (address, coins) for auto mode.
	outputs := make(map[coin.TransactionOutput]struct{}, len(c.To))
	for _, to := range c.To {
		outputs[to] = struct{}{}
	}

	if len(outputs) != len(c.To) {
		return NewError(errors.New("To contains duplicate values"))
	}

	if c.Wallet.ID == "" {
		return NewError(errors.New("Wallet.ID is required"))
	}

	addressMap := make(map[cipher.Address]struct{}, len(c.Wallet.Addresses))
	for _, a := range c.Wallet.Addresses {
		if a.Null() {
			return NewError(errors.New("Wallet.Addresses must not contain the null address"))
		}

		addressMap[a] = struct{}{}
	}

	if len(addressMap) != len(c.Wallet.Addresses) {
		return NewError(errors.New("Wallet.Addresses contains duplicate values"))
	}

	switch c.HoursSelection.Type {
	case HoursSelectionTypeAuto:
		for _, to := range c.To {
			if to.Hours != 0 {
				return NewError(errors.New("To.Hours must be zero for auto type hours selection"))
			}
		}

		switch c.HoursSelection.Mode {
		case HoursSelectionModeShare:
		case "":
			return NewError(errors.New("HoursSelection.Mode is required for auto type hours selection"))
		default:
			return NewError(errors.New("Invalid HoursSelection.Mode"))
		}

	case HoursSelectionTypeManual:
		if c.HoursSelection.Mode != "" {
			return NewError(errors.New("HoursSelection.Mode cannot be used for manual type hours selection"))
		}

	default:
		return NewError(errors.New("Invalid HoursSelection.Type"))
	}

	if c.HoursSelection.ShareFactor == nil {
		if c.HoursSelection.Mode == HoursSelectionModeShare {
			return NewError(errors.New("HoursSelection.ShareFactor must be set for share mode"))
		}
	} else {
		if c.HoursSelection.Mode != HoursSelectionModeShare {
			return NewError(errors.New("HoursSelection.ShareFactor can only be used for share mode"))
		}

		zero := decimal.New(0, 0)
		one := decimal.New(1, 0)
		if c.HoursSelection.ShareFactor.LessThan(zero) || c.HoursSelection.ShareFactor.GreaterThan(one) {
			return NewError(errors.New("HoursSelection.ShareFactor must be >= 0 and <= 1"))
		}
	}

	if len(c.Wallet.UxOuts) != 0 && len(c.Wallet.Addresses) != 0 {
		return NewError(errors.New("Wallet.UxOuts and Wallet.Addresses cannot be combined"))
	}

	// Check for duplicate spending uxouts
	uxouts := make(map[cipher.SHA256]struct{}, len(c.Wallet.UxOuts))
	for _, o := range c.Wallet.UxOuts {
		uxouts[o] = struct{}{}
	}

	if len(uxouts) != len(c.Wallet.UxOuts) {
		return NewError(errors.New("Wallet.UxOuts contains duplicate values"))
	}

	return nil
}

// newWalletFilename check for collisions and retry if failure
func newWalletFilename() string {
	timestamp := time.Now().Format(WalletTimestampFormat)
	// should read in wallet files and make sure does not exist
	padding := hex.EncodeToString((cipher.RandByte(2)))
	return fmt.Sprintf("%s_%s.%s", timestamp, padding, WalletExt)
}

// Wallet is consisted of meta and entries.
// Meta field records items that are not deterministic, like
// filename, lable, wallet type, secrets, etc.
// Entries field stores the address entries that are deterministically generated
// from seed.
// For wallet encryption
type Wallet struct {
	Meta    map[string]string
	Entries []Entry
}

// newWallet creates a wallet instance with given name and options.
func newWallet(wltName string, opts Options, bg BalanceGetter) (*Wallet, error) {
	if opts.Seed == "" {
		return nil, ErrMissingSeed
	}

	coin := opts.Coin
	if coin == "" {
		coin = CoinTypeSkycoin
	}

	w := &Wallet{
		Meta: map[string]string{
			metaFilename:   wltName,
			metaVersion:    Version,
			metaLabel:      opts.Label,
			metaSeed:       opts.Seed,
			metaLastSeed:   opts.Seed,
			metaTm:         fmt.Sprintf("%v", time.Now().Unix()),
			metaType:       "deterministic",
			metaCoin:       string(coin),
			metaEncrypted:  "false",
			metaCryptoType: "",
			metaSecrets:    "",
		},
	}

	// Create a default wallet
	_, err := w.GenerateAddresses(1)
	if err != nil {
		return nil, err
	}

	if opts.ScanN > 0 {
		// Scan for addresses with balances
		if bg != nil {
			if err := w.ScanAddresses(opts.ScanN-1, bg); err != nil {
				return nil, err
			}
		}
	}

	// Checks if the wallet need to encrypt
	if !opts.Encrypt {
		if len(opts.Password) != 0 {
			return nil, ErrMissingEncrypt
		}
		return w, nil
	}

	// Checks if the password is provided
	if len(opts.Password) == 0 {
		return nil, ErrMissingPassword
	}

	// Checks crypto type
	if _, err := getCrypto(opts.CryptoType); err != nil {
		return nil, err
	}

	// Encrypt the wallet
	if err := w.Lock(opts.Password, opts.CryptoType); err != nil {
		return nil, err
	}

	// Validate the wallet
	if err := w.Validate(); err != nil {
		return nil, err
	}

	return w, nil
}

// NewWallet creates wallet without scanning addresses
func NewWallet(wltName string, opts Options) (*Wallet, error) {
	if opts.ScanN != 0 {
		return nil, errors.New("scan number must be 0")
	}
	return newWallet(wltName, opts, nil)
}

// NewWalletScanAhead creates wallet and scan ahead N addresses
func NewWalletScanAhead(wltName string, opts Options, bg BalanceGetter) (*Wallet, error) {
	return newWallet(wltName, opts, bg)
}

// Lock encrypts the wallet with the given password and specific crypto type
func (w *Wallet) Lock(password []byte, cryptoType CryptoType) error {
	if len(password) == 0 {
		return ErrMissingPassword
	}

	if w.IsEncrypted() {
		return ErrWalletEncrypted
	}

	wlt := w.clone()

	// Records seeds in secrets
	ss := make(secrets)
	defer func() {
		// Wipes all unencrypted sensitive data
		ss.erase()
		wlt.erase()
	}()

	ss.set(secretSeed, wlt.seed())
	ss.set(secretLastSeed, wlt.lastSeed())

	// Saves address's secret keys in secrets
	for _, e := range wlt.Entries {
		ss.set(e.Address.String(), e.Secret.Hex())
	}

	sb, err := ss.serialize()
	if err != nil {
		return err
	}

	crypto, err := getCrypto(cryptoType)
	if err != nil {
		return err
	}

	// Encrypts the secrets
	encSecret, err := crypto.Encrypt(sb, password)
	if err != nil {
		return err
	}

	// Sets the crypto type
	wlt.setCryptoType(cryptoType)

	// Updates the secrets data in wallet
	wlt.setSecrets(string(encSecret))

	// Sets wallet as encrypted
	wlt.setEncrypted(true)

	// Sets the wallet version
	wlt.setVersion(Version)

	// Wipes unencrypted sensitive data
	wlt.erase()

	// Wipes the secret fields in w
	w.erase()

	// Replace the original wallet with new encrypted wallet
	w.copyFrom(wlt)
	return nil
}

// Unlock decrypts the wallet into a temporary decrypted copy of the wallet
// Returns error if the decryption fails
// The temporary decrypted wallet should be erased from memory when done.
func (w *Wallet) Unlock(password []byte) (*Wallet, error) {
	if !w.IsEncrypted() {
		return nil, ErrWalletNotEncrypted
	}

	if len(password) == 0 {
		return nil, ErrMissingPassword
	}

	wlt := w.clone()

	// Gets the secrets string
	sstr := wlt.secrets()
	if sstr == "" {
		return nil, errors.New("secrets doesn't exsit")
	}

	ct := w.cryptoType()
	if ct == "" {
		return nil, errors.New("missing crypto type")
	}

	// Gets the crypto
	crypto, err := getCrypto(ct)
	if err != nil {
		return nil, err
	}

	// Decrypts the secrets
	sb, err := crypto.Decrypt([]byte(sstr), password)
	if err != nil {
		return nil, ErrInvalidPassword
	}

	// Deserialize into secrets
	ss := make(secrets)
	defer ss.erase()
	if err := ss.deserialize(sb); err != nil {
		return nil, err
	}

	seed, ok := ss.get(secretSeed)
	if !ok {
		return nil, errors.New("seed doesn't exist in secrets")
	}
	wlt.setSeed(seed)

	lastSeed, ok := ss.get(secretLastSeed)
	if !ok {
		return nil, errors.New("lastSeed doesn't exist in secrets")
	}
	wlt.setLastSeed(lastSeed)

	// Gets addresses related secrets
	for i, e := range wlt.Entries {
		sstr, ok := ss.get(e.Address.String())
		if !ok {
			return nil, fmt.Errorf("secret of address %s doesn't exist in secrets", e.Address)
		}
		s, err := hex.DecodeString(sstr)
		if err != nil {
			return nil, fmt.Errorf("decode secret hex string failed: %v", err)
		}

		copy(wlt.Entries[i].Secret[:], s[:])
	}

	wlt.setEncrypted(false)
	wlt.setSecrets("")
	wlt.setCryptoType("")
	return wlt, nil
}

// copyFrom copies the src wallet to w
func (w *Wallet) copyFrom(src *Wallet) {
	// Clear the original info first
	w.Meta = make(map[string]string)
	w.Entries = w.Entries[:0]

	// Copies the meta
	for k, v := range src.Meta {
		w.Meta[k] = v
	}

	// Copies the address entries
	for _, e := range src.Entries {
		w.Entries = append(w.Entries, e)
	}
}

// erase wipes secret fields in wallet
func (w *Wallet) erase() {
	// Wipes the seed and last seed
	w.setSeed("")
	w.setLastSeed("")

	// Wipes private keys in entries
	for i := range w.Entries {
		for j := range w.Entries[i].Secret {
			w.Entries[i].Secret[j] = 0
		}

		w.Entries[i].Secret = cipher.SecKey{}
	}
}

// GuardUpdate executes a function within the context of a read-wirte managed decrypted wallet.
// Returns ErrWalletNotEncrypted if wallet is not encrypted.
func (w *Wallet) GuardUpdate(password []byte, fn func(w *Wallet) error) error {
	if !w.IsEncrypted() {
		return ErrWalletNotEncrypted
	}

	if len(password) == 0 {
		return ErrMissingPassword
	}

	cryptoType := w.cryptoType()
	wlt, err := w.Unlock(password)
	if err != nil {
		return err
	}

	defer wlt.erase()

	if err := fn(wlt); err != nil {
		return err
	}

	if err := wlt.Lock(password, cryptoType); err != nil {
		return err
	}

	*w = *wlt
	// Wipes all sensitive data
	w.erase()
	return nil
}

// GuardView executes a function within the context of a read-only managed decrypted wallet.
// Returns ErrWalletNotEncrypted if wallet is not encrypted.
func (w *Wallet) GuardView(password []byte, f func(w *Wallet) error) error {
	if !w.IsEncrypted() {
		return ErrWalletNotEncrypted
	}

	if len(password) == 0 {
		return ErrMissingPassword
	}

	wlt, err := w.Unlock(password)
	if err != nil {
		return err
	}

	defer wlt.erase()

	return f(wlt)
}

// Load loads wallet from a given file
func Load(wltFile string) (*Wallet, error) {
	if _, err := os.Stat(wltFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("load wallet file failed, wallet %s doesn't exist", wltFile)
	}

	r := &ReadableWallet{}
	if err := r.Load(wltFile); err != nil {
		return nil, err
	}

	// update filename meta info with the real filename
	r.Meta["filename"] = filepath.Base(wltFile)
	return r.ToWallet()
}

// Save saves the wallet to given dir
func (w *Wallet) Save(dir string) error {
	r := NewReadableWallet(w)
	return r.Save(filepath.Join(dir, w.Filename()))
}

// removeBackupFiles removes any *.wlt.bak files whom have version 0.1 and *.wlt matched in the given directory
func removeBackupFiles(dir string) error {
	fs, err := filterDir(dir, ".wlt")
	if err != nil {
		return err
	}

	// Creates the .wlt file map
	fm := make(map[string]struct{})
	for _, f := range fs {
		fm[f] = struct{}{}
	}

	// Filters all .wlt.bak files in the directory
	bakFs, err := filterDir(dir, ".wlt.bak")
	if err != nil {
		return err
	}

	// Removes the .wlt.bak file that has .wlt matched.
	for _, bf := range bakFs {
		f := strings.TrimRight(bf, ".bak")
		if _, ok := fm[f]; ok {
			// Load and check the wallet version
			w, err := Load(f)
			if err != nil {
				return err
			}

			if w.Version() == "0.1" {
				if err := os.Remove(bf); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func filterDir(dir string, suffix string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), suffix) {
			res = append(res, filepath.Join(dir, f.Name()))
		}
	}
	return res, nil
}

// reset resets the wallet entries and move the lastSeed to origin
func (w *Wallet) reset() {
	w.Entries = []Entry{}
	w.setLastSeed(w.seed())
}

// Validate validates the wallet
func (w *Wallet) Validate() error {
	if _, ok := w.Meta[metaFilename]; !ok {
		return errors.New("filename not set")
	}
	if _, ok := w.Meta[metaSeed]; !ok {
		return errors.New("seed field not set")
	}

	walletType, ok := w.Meta[metaType]
	if !ok {
		return errors.New("type field not set")
	}
	if walletType != "deterministic" {
		return errors.New("wallet type invalid")
	}

	if _, ok := w.Meta[metaCoin]; !ok {
		return errors.New("coin field not set")
	}

	if encStr, ok := w.Meta[metaEncrypted]; ok {
		// validate the encrypted value
		isEncrypted, err := strconv.ParseBool(encStr)
		if err != nil {
			return fmt.Errorf("invalid encrypted value: %v", err)
		}

		// checks if the secrets field is empty
		if isEncrypted {
			if _, ok := w.Meta[metaCryptoType]; !ok {
				return errors.New("crypto type field not set")
			}

			if _, ok := w.Meta[metaSecrets]; !ok {
				return errors.New("wallet is encrypted, but secrets field not set")
			}
		}
	}

	return nil
}

// Type gets the wallet type
func (w *Wallet) Type() string {
	return w.Meta[metaType]
}

// Version gets the wallet version
func (w *Wallet) Version() string {
	return w.Meta[metaVersion]
}

func (w *Wallet) setVersion(v string) {
	w.Meta[metaVersion] = v
}

// Filename gets the wallet filename
func (w *Wallet) Filename() string {
	return w.Meta[metaFilename]
}

// setFilename sets the wallet filename
func (w *Wallet) setFilename(fn string) {
	w.Meta[metaFilename] = fn
}

// Label gets the wallet label
func (w *Wallet) Label() string {
	return w.Meta[metaLabel]
}

// setLabel sets the wallet label
func (w *Wallet) setLabel(label string) {
	w.Meta[metaLabel] = label
}

// lastSeed returns the last seed
func (w *Wallet) lastSeed() string {
	return w.Meta[metaLastSeed]
}

func (w *Wallet) setLastSeed(lseed string) {
	w.Meta[metaLastSeed] = lseed
}

func (w *Wallet) seed() string {
	return w.Meta[metaSeed]
}

func (w *Wallet) setSeed(seed string) {
	w.Meta[metaSeed] = seed
}

func (w *Wallet) setEncrypted(encrypt bool) {
	w.Meta[metaEncrypted] = strconv.FormatBool(encrypt)
}

// IsEncrypted checks whether the wallet is encrypted.
func (w *Wallet) IsEncrypted() bool {
	encStr, ok := w.Meta[metaEncrypted]
	if !ok {
		return false
	}

	b, err := strconv.ParseBool(encStr)
	if err != nil {
		// This can not happen, the meta.encrypted value is either set by
		// setEncrypted() method or converted in ReadableWallet.toWallet().
		// toWallet() method will throw error if the meta.encrypted string is invalid.
		logger.Warning("parse wallet.meta.encrypted string failed: %v", err)
		return false
	}
	return b
}

func (w *Wallet) setCryptoType(tp CryptoType) {
	w.Meta[metaCryptoType] = string(tp)
}

func (w *Wallet) cryptoType() CryptoType {
	return CryptoType(w.Meta[metaCryptoType])
}

func (w *Wallet) secrets() string {
	return w.Meta[metaSecrets]
}

func (w *Wallet) setSecrets(s string) {
	w.Meta[metaSecrets] = s
}

// GenerateAddresses generates addresses
func (w *Wallet) GenerateAddresses(num uint64) ([]cipher.Address, error) {
	if num == 0 {
		return nil, nil
	}

	if w.IsEncrypted() {
		return nil, ErrWalletEncrypted
	}

	var seckeys []cipher.SecKey
	var seed []byte
	if len(w.Entries) == 0 {
		seed, seckeys = cipher.GenerateDeterministicKeyPairsSeed([]byte(w.seed()), int(num))
	} else {
		sd, err := hex.DecodeString(w.lastSeed())
		if err != nil {
			return nil, fmt.Errorf("decode hex seed failed: %v", err)
		}
		seed, seckeys = cipher.GenerateDeterministicKeyPairsSeed(sd, int(num))
	}

	w.setLastSeed(hex.EncodeToString(seed))

	addrs := make([]cipher.Address, len(seckeys))
	for i, s := range seckeys {
		p := cipher.PubKeyFromSecKey(s)
		a := cipher.AddressFromPubKey(p)
		addrs[i] = a
		w.Entries = append(w.Entries, Entry{
			Address: a,
			Secret:  s,
			Public:  p,
		})
	}
	return addrs, nil
}

// ScanAddresses scans ahead N addresses to find one with none-zero coins.
func (w *Wallet) ScanAddresses(scanN uint64, bg BalanceGetter) error {
	if w.IsEncrypted() {
		return ErrWalletEncrypted
	}

	if scanN <= 0 {
		return nil
	}

	nExistingAddrs := uint64(len(w.Entries))

	// Generate the addresses to scan
	addrs, err := w.GenerateAddresses(scanN)
	if err != nil {
		return err
	}

	// Get these addresses' balances
	bals, err := bg.GetBalanceOfAddrs(addrs)
	if err != nil {
		return err
	}

	// Check balance from the last one until we find the address that has coins
	var keepNum uint64
	for i := len(bals) - 1; i >= 0; i-- {
		if bals[i].Confirmed.Coins > 0 || bals[i].Predicted.Coins > 0 {
			keepNum = uint64(i + 1)
			break
		}
	}

	// Regenerate addresses up to keepNum.
	// This is necessary to keep the lastSeed updated.
	if keepNum != uint64(len(bals)) {
		w.reset()
		w.GenerateAddresses(nExistingAddrs + keepNum)
	}

	return nil
}

// GetAddresses returns all addresses in wallet
func (w *Wallet) GetAddresses() []cipher.Address {
	addrs := make([]cipher.Address, len(w.Entries))
	for i, e := range w.Entries {
		addrs[i] = e.Address
	}
	return addrs
}

// GetEntry returns entry of given address
func (w *Wallet) GetEntry(a cipher.Address) (Entry, bool) {
	for _, e := range w.Entries {
		if e.Address == a {
			return e, true
		}
	}
	return Entry{}, false
}

// AddEntry adds new entry
func (w *Wallet) AddEntry(entry Entry) error {
	// dup check
	for _, e := range w.Entries {
		if e.Address == entry.Address {
			return errors.New("duplicate address entry")
		}
	}

	w.Entries = append(w.Entries, entry)
	return nil
}

// clone returns the clone of self
func (w *Wallet) clone() *Wallet {
	wlt := Wallet{Meta: make(map[string]string)}
	for k, v := range w.Meta {
		wlt.Meta[k] = v
	}

	for _, e := range w.Entries {
		wlt.Entries = append(wlt.Entries, e)
	}

	return &wlt
}

// Validator validate if the wallet be able to create spending transaction
type Validator interface {
	// checks if any of the given addresses has unconfirmed spending transactions
	HasUnconfirmedSpendTx(addr []cipher.Address) (bool, error)
}

// CreateAndSignTransaction Creates a Transaction
// spending coins and hours from wallet
func (w *Wallet) CreateAndSignTransaction(auxs coin.AddressUxOuts, headTime, coins uint64, dest cipher.Address) (*coin.Transaction, error) {
	if w.IsEncrypted() {
		return nil, ErrWalletEncrypted
	}

	entriesMap := make(map[cipher.Address]Entry)
	for a := range auxs {
		e, ok := w.GetEntry(a)
		// Check that auxs does not contain addresses that are not known to this wallet
		if !ok {
			return nil, ErrUnknownAddress
		}
		entriesMap[e.Address] = e
	}

	// Determine which unspents to spend.
	// Use the MaximizeUxOuts strategy, this will keep the uxout pool smaller
	uxa := auxs.Flatten()
	uxb, err := NewUxBalances(headTime, uxa)
	if err != nil {
		return nil, err
	}

	spends, err := ChooseSpendsMaximizeUxOuts(uxb, coins, 0)
	if err != nil {
		return nil, err
	}

	// Add these unspents as tx inputs
	var txn coin.Transaction
	toSign := make([]cipher.SecKey, len(spends))
	spending := Balance{Coins: 0, Hours: 0}
	for i, au := range spends {
		entry, ok := entriesMap[au.Address]
		if !ok {
			return nil, NewError(fmt.Errorf("address %v does not exist in wallet %v", au.Address, w.Filename()))
		}

		txn.PushInput(au.Hash)

		toSign[i] = entry.Secret

		spending.Coins += au.Coins
		spending.Hours += au.Hours
	}

	if spending.Hours == 0 {
		return nil, fee.ErrTxnNoFee
	}

	// Calculate coin hour allocation
	changeCoins := spending.Coins - coins
	haveChange := changeCoins > 0
	changeHours, addrHours, outputHours := DistributeSpendHours(spending.Hours, 1, haveChange)

	logger.Infof("wallet.CreateAndSignTransaction: spending.Hours=%d, fee.VerifyTransactionFeeForHours(%d, %d)", spending.Hours, outputHours, spending.Hours-outputHours)
	if err := fee.VerifyTransactionFeeForHours(outputHours, spending.Hours-outputHours); err != nil {
		logger.Warningf("wallet.CreateAndSignTransaction: fee.VerifyTransactionFeeForHours failed: %v", err)
		return nil, err
	}

	if haveChange {
		changeAddr := spends[0].Address
		txn.PushOutput(changeAddr, changeCoins, changeHours)
	}

	txn.PushOutput(dest, coins, addrHours[0])

	txn.SignInputs(toSign)
	txn.UpdateHeader()

	return &txn, nil
}

// CreateAndSignTransactionAdvanced creates and signs a transaction based upon CreateTransactionParams.
// Set the password as nil if the wallet is not encrypted, otherwise the password must be provided.
// NOTE: Caller must ensure that auxs correspond to params.Wallet.Addresses and params.Wallet.UxOuts options
func (w *Wallet) CreateAndSignTransactionAdvanced(params CreateTransactionParams, auxs coin.AddressUxOuts, headTime uint64) (*coin.Transaction, []UxBalance, error) {
	if err := params.Validate(); err != nil {
		return nil, nil, err
	}

	if params.Wallet.ID != w.Filename() {
		return nil, nil, NewError(errors.New("params.Wallet.ID does not match wallet"))
	}

	if w.IsEncrypted() {
		return nil, nil, ErrWalletEncrypted
	}

	entriesMap := make(map[cipher.Address]Entry)
	for a := range auxs {
		// Check that auxs does not contain addresses that are not known to this wallet
		e, ok := w.GetEntry(a)
		if !ok {
			return nil, nil, ErrUnknownAddress
		}
		entriesMap[e.Address] = e
	}

	txn := &coin.Transaction{}

	// Determine which unspents to spend
	uxa := auxs.Flatten()

	uxb, err := NewUxBalances(headTime, uxa)
	if err != nil {
		return nil, nil, err
	}

	// Reverse lookup set to recover the inputs
	uxbMap := make(map[cipher.SHA256]UxBalance, len(uxb))
	for _, u := range uxb {
		if _, ok := uxbMap[u.Hash]; ok {
			return nil, nil, errors.New("Duplicate UxBalance in array")
		}
		uxbMap[u.Hash] = u
	}

	// calculate total coins and minimum hours to send
	var totalOutCoins uint64
	var requestedHours uint64
	for _, to := range params.To {
		totalOutCoins, err = coin.AddUint64(totalOutCoins, to.Coins)
		if err != nil {
			return nil, nil, NewError(fmt.Errorf("total output coins error: %v", err))
		}

		requestedHours, err = coin.AddUint64(requestedHours, to.Hours)
		if err != nil {
			return nil, nil, NewError(fmt.Errorf("total output hours error: %v", err))
		}
	}

	// Use the MinimizeUxOuts strategy, to use least possible uxouts
	// this will allow more frequent spending
	// we don't need to check whether we have sufficient balance beforehand as ChooseSpends already checks that
	spends, err := ChooseSpendsMinimizeUxOuts(uxb, totalOutCoins, requestedHours)
	if err != nil {
		return nil, nil, err
	}

	// calculate total coins and hours in spends
	var totalInputCoins uint64
	var totalInputHours uint64
	toSign := make([]cipher.SecKey, len(spends))
	for i, spend := range spends {
		totalInputCoins, err = coin.AddUint64(totalInputCoins, spend.Coins)
		if err != nil {
			return nil, nil, err
		}

		totalInputHours, err = coin.AddUint64(totalInputHours, spend.Hours)
		if err != nil {
			return nil, nil, err
		}

		entry, ok := entriesMap[spend.Address]
		if !ok {
			return nil, nil, fmt.Errorf("spend address %s not found in entriesMap", spend.Address.String())
		}

		toSign[i] = entry.Secret
		txn.PushInput(spend.Hash)
	}

	feeHours := fee.RequiredFee(totalInputHours)
	if feeHours == 0 {
		return nil, nil, fee.ErrTxnNoFee
	}
	remainingHours := totalInputHours - feeHours

	switch params.HoursSelection.Type {
	case HoursSelectionTypeManual:
		for _, out := range params.To {
			txn.Out = append(txn.Out, out)
		}

	case HoursSelectionTypeAuto:
		var addrHours []uint64

		switch params.HoursSelection.Mode {
		case HoursSelectionModeShare:
			// multiply remaining hours after fee burn with share factor
			hours, err := coin.Uint64ToInt64(remainingHours)
			if err != nil {
				return nil, nil, err
			}

			allocatedHoursInt := params.HoursSelection.ShareFactor.Mul(decimal.New(hours, 0)).IntPart()
			allocatedHours, err := coin.Int64ToUint64(allocatedHoursInt)
			if err != nil {
				return nil, nil, err
			}

			toCoins := make([]uint64, len(params.To))
			for i, to := range params.To {
				toCoins[i] = to.Coins
			}

			addrHours, err = DistributeCoinHoursProportional(toCoins, allocatedHours)
			if err != nil {
				return nil, nil, err
			}
		default:
			return nil, nil, ErrInvalidHoursSelectionType
		}

		for i, out := range params.To {
			out.Hours = addrHours[i]
			txn.Out = append(txn.Out, out)
		}

	default:
		return nil, nil, ErrInvalidHoursSelectionMode
	}

	totalOutHours, err := txn.OutputHours()
	if err != nil {
		return nil, nil, err
	}

	// Make sure we have enough coins and coin hours
	// If we don't, and we called ChooseSpends, then ChooseSpends has a bug, as it should have returned this error already
	if totalOutCoins > totalInputCoins {
		logger.WithError(ErrInsufficientBalance).Error("Insufficient coins after choosing spends, this should not occur")
		return nil, nil, ErrInsufficientBalance
	}

	if totalOutHours > remainingHours {
		logger.WithError(fee.ErrTxnInsufficientCoinHours).Error("Insufficient hours after choosing spends or distributing hours, this should not occur")
		return nil, nil, fee.ErrTxnInsufficientCoinHours
	}

	// create change output
	changeCoins := totalInputCoins - totalOutCoins
	changeHours := remainingHours - totalOutHours

	// If there are no change coins but there are change hours, try to add another
	// input to save the change hours.
	// This chooses an available input with the least number of coin hours;
	// if the extra coin hour fee incurred by this additional input is less than
	// the remaining coin hours, the input is added.
	if changeCoins == 0 && changeHours > 0 {
		// Find the output with the least coin hours
		// If size of the fee for this output is less than the changeHours, add it
		// Update changeCoins and changeHours
		z := uxBalancesSub(uxb, spends)
		sortSpendsHoursLowToHigh(z)
		if len(z) > 0 {
			extra := z[0]

			// Calculate the new hours being spent
			newTotalHours, err := coin.AddUint64(totalInputHours, extra.Hours)
			if err != nil {
				return nil, nil, err
			}

			// Calculate the new fee for this new amount of hours
			newFee := fee.RequiredFee(newTotalHours)
			if newFee < feeHours {
				err := errors.New("updated fee after adding extra input for change is unexpectedly less than it was initially")
				logger.WithError(err).Error()
				return nil, nil, err
			}

			// If the cost of adding this extra input is less than the amount of change hours we
			// can save, use the input
			additionalFee := newFee - feeHours
			if additionalFee < changeHours {
				changeCoins = extra.Coins

				if extra.Hours < additionalFee {
					err := errors.New("calculated additional fee is unexpectedly higher than the extra input's hours")
					logger.WithError(err).Error()
					return nil, nil, err
				}

				additionalHours := extra.Hours - additionalFee
				changeHours, err = coin.AddUint64(changeHours, additionalHours)
				if err != nil {
					return nil, nil, err
				}

				entry, ok := entriesMap[extra.Address]
				if !ok {
					return nil, nil, fmt.Errorf("extra spend address %s not found in entriesMap", extra.Address.String())
				}

				toSign = append(toSign, entry.Secret)
				txn.PushInput(extra.Hash)
			}
		}
	}

	// With auto share mode, if there are leftover hours and change couldn't be force-added,
	// recalculate that share ratio at 100%
	if changeCoins == 0 && changeHours > 0 && params.HoursSelection.Type == HoursSelectionTypeAuto && params.HoursSelection.Mode == HoursSelectionModeShare {
		oneDecimal := decimal.New(1, 0)
		if params.HoursSelection.ShareFactor.Equal(oneDecimal) {
			return nil, nil, errors.New("share factor is 1.0 but changeHours > 0 unexpectedly")
		}
		params.HoursSelection.ShareFactor = &oneDecimal
		return w.CreateAndSignTransactionAdvanced(params, auxs, headTime)
	}

	if changeCoins > 0 {
		var changeAddress cipher.Address
		if params.ChangeAddress != nil {
			changeAddress = *params.ChangeAddress
		} else {
			// Choose a change address from the unspent outputs
			// Sort spends by address, comparing bytes, and use the first
			// This provides deterministic change address selection from a set of unspent outputs
			if len(spends) == 0 {
				return nil, nil, errors.New("spends is unexpectedly empty when choosing an automatic change address")
			}

			addressBytes := make([][]byte, len(spends))
			for i, s := range spends {
				addressBytes[i] = s.Address.Bytes()
			}

			sort.Slice(addressBytes, func(i, j int) bool {
				return bytes.Compare(addressBytes[i], addressBytes[j]) < 0
			})

			var err error
			changeAddress, err = cipher.AddressFromBytes(addressBytes[0])
			if err != nil {
				logger.Critical().Errorf("cipher.AddressFromBytes failed for change address converted to bytes: %v", err)
				return nil, nil, err
			}
		}

		txn.PushOutput(changeAddress, changeCoins, changeHours)
	}

	txn.SignInputs(toSign)
	txn.UpdateHeader()

	inputs := make([]UxBalance, len(txn.In))
	for i, h := range txn.In {
		uxBalance, ok := uxbMap[h]
		if !ok {
			return nil, nil, errors.New("Created transaction's input is not in the UxBalanceSet, this should not occur")
		}
		inputs[i] = uxBalance
	}

	if err := verifyCreatedTransactionInvariants(params, txn, inputs); err != nil {
		logger.Critical().Errorf("CreateAndSignTransactionAdvanced created transaction that violates invariants, aborting: %v", err)
		return nil, nil, fmt.Errorf("Created transaction that violates invariants, this is a bug: %v", err)
	}

	return txn, inputs, nil
}

// verifyCreatedTransactionInvariants checks that the transaction that was created matches expectations.
// Does not call visor verification methods because that causes import cycle.
// daemon.Gateway checks that the transaction passes additional visor verification methods.
func verifyCreatedTransactionInvariants(params CreateTransactionParams, txn *coin.Transaction, inputs []UxBalance) error {
	for _, o := range txn.Out {
		// No outputs should be sent to the null address
		if o.Address.Null() {
			return errors.New("Output address is null")
		}

		if o.Coins == 0 {
			return errors.New("Output coins is 0")
		}
	}

	if len(txn.Out) != len(params.To) && len(txn.Out) != len(params.To)+1 {
		return errors.New("Transaction has unexpected number of outputs")
	}

	for i, o := range txn.Out[:len(params.To)] {
		if o.Address != params.To[i].Address {
			return errors.New("Output address does not match requested address")
		}

		if o.Coins != params.To[i].Coins {
			return errors.New("Output coins does not match requested coins")
		}

		if params.To[i].Hours != 0 && o.Hours != params.To[i].Hours {
			return errors.New("Output hours does not match requested hours")
		}
	}

	if len(txn.Sigs) != len(txn.In) {
		return errors.New("Number of signatures does not match number of inputs")
	}

	if len(txn.In) != len(inputs) {
		return errors.New("Number of UxOut inputs does not match number of transaction inputs")
	}

	for i, h := range txn.In {
		if inputs[i].Hash != h {
			return errors.New("Transaction input hash does not match UxOut inputs hash")
		}
	}

	inputsMap := make(map[cipher.SHA256]struct{}, len(inputs))

	for _, i := range inputs {
		if i.Hours < i.InitialHours {
			return errors.New("Calculated input hours are unexpectedly less than the initial hours")
		}

		if i.SrcTransaction.Null() {
			return errors.New("Input's source transaction is a null hash")
		}

		if i.Hash.Null() {
			return errors.New("Input's hash is a null hash")
		}

		if _, ok := inputsMap[i.Hash]; ok {
			return errors.New("Duplicate input in array")
		}

		inputsMap[i.Hash] = struct{}{}
	}

	var inputHours uint64
	for _, i := range inputs {
		var err error
		inputHours, err = coin.AddUint64(inputHours, i.Hours)
		if err != nil {
			return err
		}
	}

	var outputHours uint64
	for _, i := range txn.Out {
		var err error
		outputHours, err = coin.AddUint64(outputHours, i.Hours)
		if err != nil {
			return err
		}
	}

	if inputHours < outputHours {
		return errors.New("Total input hours is less than the output hours")
	}

	if inputHours-outputHours < fee.RequiredFee(inputHours) {
		return errors.New("Transaction will not satisy required fee")
	}

	return nil
}

// DistributeSpendHours calculates how many coin hours to transfer to the change address and how
// many to transfer to each of the other destination addresses.
// Input hours are split by BurnFactor (rounded down) to meet the fee requirement.
// The remaining hours are split in half, one half goes to the change address
// and the other half goes to the destination addresses.
// If the remaining hours are an odd number, the change address gets the extra hour.
// If the amount assigned to the destination addresses is not perfectly divisible by the
// number of destination addresses, the extra hours are distributed to some of these addresses.
// Returns the number of hours to send to the change address,
// an array of length nAddrs with the hours to give to each destination address,
// and a sum of these values.
func DistributeSpendHours(inputHours, nAddrs uint64, haveChange bool) (uint64, []uint64, uint64) {
	feeHours := fee.RequiredFee(inputHours)
	remainingHours := inputHours - feeHours

	var changeHours uint64
	if haveChange {
		// Split the remaining hours between the change output and the other outputs
		changeHours = remainingHours / 2

		// If remainingHours is an odd number, give the extra hour to the change output
		if remainingHours%2 == 1 {
			changeHours++
		}
	}

	// Distribute the remaining hours equally amongst the destination outputs
	remainingAddrHours := remainingHours - changeHours
	addrHoursShare := remainingAddrHours / nAddrs

	// Due to integer division, extra coin hours might remain after dividing by len(toAddrs)
	// Allocate these extra hours to the toAddrs
	addrHours := make([]uint64, nAddrs)
	for i := range addrHours {
		addrHours[i] = addrHoursShare
	}

	extraHours := remainingAddrHours - (addrHoursShare * nAddrs)
	i := 0
	for extraHours > 0 {
		addrHours[i] = addrHours[i] + 1
		i++
		extraHours--
	}

	// Assert that the hour calculation is correct
	var spendHours uint64
	for _, h := range addrHours {
		spendHours += h
	}
	spendHours += changeHours
	if spendHours != remainingHours {
		logger.Panicf("spendHours != remainingHours (%d != %d), calculation error", spendHours, remainingHours)
	}

	return changeHours, addrHours, spendHours
}

// DistributeCoinHoursProportional distributes hours amongst coins proportional to the coins amount
func DistributeCoinHoursProportional(coins []uint64, hours uint64) ([]uint64, error) {
	if len(coins) == 0 {
		return nil, errors.New("DistributeCoinHoursProportional coins array must not be empty")
	}

	coinsInt := make([]*big.Int, len(coins))

	var total uint64
	for i, c := range coins {
		if c == 0 {
			return nil, errors.New("DistributeCoinHoursProportional coins array has a zero value")
		}

		var err error
		total, err = coin.AddUint64(total, c)
		if err != nil {
			return nil, err
		}

		cInt64, err := coin.Uint64ToInt64(c)
		if err != nil {
			return nil, err
		}

		coinsInt[i] = big.NewInt(cInt64)
	}

	totalInt64, err := coin.Uint64ToInt64(total)
	if err != nil {
		return nil, err
	}
	totalInt := big.NewInt(totalInt64)

	hoursInt64, err := coin.Uint64ToInt64(hours)
	if err != nil {
		return nil, err
	}
	hoursInt := big.NewInt(hoursInt64)

	var assignedHours uint64
	addrHours := make([]uint64, len(coins))
	for i, c := range coinsInt {
		// Scale the ratio of coins to total coins proportionally by calculating
		// (coins * totalHours) / totalCoins
		// The remainder is truncated, remaining hours are appended after this
		num := &big.Int{}
		num.Mul(c, hoursInt)

		fracInt := big.Int{}
		fracInt.Div(num, totalInt)

		if !fracInt.IsUint64() {
			return nil, errors.New("DistributeCoinHoursProportional calculated fractional hours is not representable as a uint64")
		}

		fracHours := fracInt.Uint64()

		addrHours[i] = fracHours
		assignedHours, err = coin.AddUint64(assignedHours, fracHours)
		if err != nil {
			return nil, err
		}
	}

	if hours < assignedHours {
		return nil, errors.New("DistributeCoinHoursProportional assigned hours exceeding input hours, this is a bug")
	}

	remainingHours := hours - assignedHours

	if remainingHours > uint64(len(coins)) {
		return nil, errors.New("DistributeCoinHoursProportional remaining hours exceed len(coins), this is a bug")
	}

	// For remaining hours lost due to fractional cutoff when scaling,
	// first provide at least 1 coin hour to coins that were assigned 0.
	i := 0
	for remainingHours > 0 && i < len(coins) {
		if addrHours[i] == 0 {
			addrHours[i] = 1
			remainingHours--
		}
		i++
	}

	// Then, assign the extra coin hours
	i = 0
	for remainingHours > 0 {
		addrHours[i] = addrHours[i] + 1
		remainingHours--
		i++
	}

	return addrHours, nil
}

// UxBalance is an intermediate representation of a UxOut for sorting and spend choosing
type UxBalance struct {
	Hash           cipher.SHA256
	BkSeq          uint64
	Time           uint64
	Address        cipher.Address
	Coins          uint64
	InitialHours   uint64
	Hours          uint64
	SrcTransaction cipher.SHA256
}

// NewUxBalances converts coin.UxArray to []UxBalance. headTime is required to calculate coin hours.
func NewUxBalances(headTime uint64, uxa coin.UxArray) ([]UxBalance, error) {
	uxb := make([]UxBalance, len(uxa))
	for i, ux := range uxa {
		b, err := NewUxBalance(headTime, ux)
		if err != nil {
			return nil, err
		}
		uxb[i] = b
	}

	return uxb, nil
}

// NewUxBalance converts coin.UxOut to UxBalance. headTime is required to calculate coin hours.
func NewUxBalance(headTime uint64, ux coin.UxOut) (UxBalance, error) {
	hours, err := ux.CoinHours(headTime)
	if err != nil {
		return UxBalance{}, err
	}

	return UxBalance{
		Hash:           ux.Hash(),
		BkSeq:          ux.Head.BkSeq,
		Time:           ux.Head.Time,
		Address:        ux.Body.Address,
		Coins:          ux.Body.Coins,
		InitialHours:   ux.Body.Hours,
		Hours:          hours,
		SrcTransaction: ux.Body.SrcTransaction,
	}, nil
}

func uxBalancesSub(a, b []UxBalance) []UxBalance {
	var x []UxBalance

	bMap := make(map[cipher.SHA256]struct{}, len(b))
	for _, i := range b {
		bMap[i.Hash] = struct{}{}
	}

	for _, i := range a {
		if _, ok := bMap[i.Hash]; !ok {
			x = append(x, i)
		}
	}

	return x
}

// ChooseSpendsMinimizeUxOuts chooses uxout spends to satisfy an amount, using the least number of uxouts
//     -- PRO: Allows more frequent spending, less waiting for confirmations, useful for exchanges.
//     -- PRO: When transaction is volume is higher, transactions are prioritized by fee/size. Minimizing uxouts minimizes size.
//     -- CON: Would make the unconfirmed pool grow larger.
// Users with high transaction frequency will want to use this so that they will not need to wait as frequently
// for unconfirmed spends to complete before sending more.
// Alternatively, or in addition to this, they should batch sends into single transactions.
func ChooseSpendsMinimizeUxOuts(uxa []UxBalance, coins, hours uint64) ([]UxBalance, error) {
	return ChooseSpends(uxa, coins, hours, sortSpendsCoinsHighToLow)
}

// sortSpendsCoinsHighToLow sorts uxout spends with highest balance to lowest
func sortSpendsCoinsHighToLow(uxa []UxBalance) {
	sort.Slice(uxa, makeCmpUxOutByCoins(uxa, func(a, b uint64) bool {
		return a > b
	}))
}

// ChooseSpendsMaximizeUxOuts chooses uxout spends to satisfy an amount, using the most number of uxouts
// See the pros and cons of ChooseSpendsMinimizeUxOuts.
// This should be the default mode, because this keeps the unconfirmed pool smaller which will allow
// the network to scale better.
func ChooseSpendsMaximizeUxOuts(uxa []UxBalance, coins, hours uint64) ([]UxBalance, error) {
	return ChooseSpends(uxa, coins, hours, sortSpendsCoinsLowToHigh)
}

// sortSpendsCoinsLowToHigh sorts uxout spends with lowest balance to highest
func sortSpendsCoinsLowToHigh(uxa []UxBalance) {
	sort.Slice(uxa, makeCmpUxOutByCoins(uxa, func(a, b uint64) bool {
		return a < b
	}))
}

// sortSpendsHoursLowToHigh sorts uxout spends with lowest hours to highest
func sortSpendsHoursLowToHigh(uxa []UxBalance) {
	sort.Slice(uxa, makeCmpUxOutByHours(uxa, func(a, b uint64) bool {
		return a < b
	}))
}

func makeCmpUxOutByCoins(uxa []UxBalance, coinsCmp func(a, b uint64) bool) func(i, j int) bool {
	// Sort by:
	// coins highest or lowest depending on coinsCmp
	//  hours lowest
	//   oldest first
	//    tie break with hash comparison
	return func(i, j int) bool {
		a := uxa[i]
		b := uxa[j]

		if a.Coins == b.Coins {
			if a.Hours == b.Hours {
				if a.BkSeq == b.BkSeq {
					return cmpUxBalanceByUxID(a, b)
				}
				return a.BkSeq < b.BkSeq
			}
			return a.Hours < b.Hours
		}
		return coinsCmp(a.Coins, b.Coins)
	}
}

func makeCmpUxOutByHours(uxa []UxBalance, hoursCmp func(a, b uint64) bool) func(i, j int) bool {
	// Sort by:
	// hours highest or lowest depending on hoursCmp
	//  coins lowest
	//   oldest first
	//    tie break with hash comparison
	return func(i, j int) bool {
		a := uxa[i]
		b := uxa[j]

		if a.Hours == b.Hours {
			if a.Coins == b.Coins {
				if a.BkSeq == b.BkSeq {
					return cmpUxBalanceByUxID(a, b)
				}
				return a.BkSeq < b.BkSeq
			}
			return a.Coins < b.Coins
		}
		return hoursCmp(a.Hours, b.Hours)
	}
}

func cmpUxBalanceByUxID(a, b UxBalance) bool {
	cmp := bytes.Compare(a.Hash[:], b.Hash[:])
	if cmp == 0 {
		logger.Panic("Duplicate UxOut when sorting")
	}
	return cmp < 0
}

// ChooseSpends chooses uxouts from a list of uxouts.
// It first chooses the uxout with the most number of coins that has nonzero coinhours.
// It then chooses uxouts with zero coinhours, ordered by sortStrategy
// It then chooses remaining uxouts with nonzero coinhours, ordered by sortStrategy
func ChooseSpends(uxa []UxBalance, coins, hours uint64, sortStrategy func([]UxBalance)) ([]UxBalance, error) {
	if coins == 0 {
		return nil, ErrZeroSpend
	}

	if len(uxa) == 0 {
		return nil, ErrNoUnspents
	}

	for _, ux := range uxa {
		if ux.Coins == 0 {
			logger.Panic("UxOut coins are 0, can't spend")
			return nil, errors.New("UxOut coins are 0, can't spend")
		}
	}

	// Split UxBalances into those with and without hours
	var nonzero, zero []UxBalance
	for _, ux := range uxa {
		if ux.Hours == 0 {
			zero = append(zero, ux)
		} else {
			nonzero = append(nonzero, ux)
		}
	}

	// Abort if there are no uxouts with non-zero coinhours, they can't be spent yet
	if len(nonzero) == 0 {
		return nil, fee.ErrTxnNoFee
	}

	// Sort uxouts with hours lowest to highest and coins highest to lowest
	sortSpendsCoinsHighToLow(nonzero)

	var have Balance
	var spending []UxBalance

	// Use the first nonzero output. This output will have the least hours possible
	firstNonzero := nonzero[0]
	if firstNonzero.Hours == 0 {
		logger.Panic("balance has zero hours unexpectedly")
		return nil, errors.New("balance has zero hours unexpectedly")
	}

	nonzero = nonzero[1:]

	spending = append(spending, firstNonzero)

	have.Coins += firstNonzero.Coins
	have.Hours += firstNonzero.Hours

	if have.Coins >= coins && fee.RemainingHours(have.Hours) >= hours {
		return spending, nil
	}

	// Sort uxouts without hours according to the sorting strategy
	sortStrategy(zero)

	for _, ux := range zero {
		spending = append(spending, ux)

		have.Coins += ux.Coins
		have.Hours += ux.Hours

		if have.Coins >= coins {
			break
		}
	}

	if have.Coins >= coins && fee.RemainingHours(have.Hours) >= hours {
		return spending, nil
	}

	// Sort remaining uxouts with hours according to the sorting strategy
	sortStrategy(nonzero)

	for _, ux := range nonzero {
		spending = append(spending, ux)

		have.Coins += ux.Coins
		have.Hours += ux.Hours

		if have.Coins >= coins && fee.RemainingHours(have.Hours) >= hours {
			return spending, nil
		}
	}

	if have.Coins < coins {
		return nil, ErrInsufficientBalance
	}

	return nil, ErrInsufficientHours
}
