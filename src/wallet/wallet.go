/*
Package wallet implements wallets and the wallet database service
*/
package wallet

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"encoding/hex"

	"github.com/skycoin/skycoin/src/cipher"

	"github.com/skycoin/skycoin/src/util/logging"
)

// Error wraps wallet-related errors.
// It wraps errors caused by user input, but not errors caused by programmer input or internal issues.
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
	// ErrWalletRecoverSeedWrong is returned if the seed does not match the specified wallet when recovering
	ErrWalletRecoverSeedWrong = NewError(errors.New("wallet recovery seed is wrong"))
	// ErrNilTransactionsFinder is returned if Options.ScanN > 0 but a nil TransactionsFinder was provided
	ErrNilTransactionsFinder = NewError(errors.New("scan ahead requested but balance getter is nil"))
	// ErrWalletNotDeterministic is returned if a wallet's type is not deterministic but it is necessary for the requested operation
	ErrWalletNotDeterministic = NewError(errors.New("wallet type is not deterministic"))
	// ErrInvalidCoinType is returned for invalid coin types
	ErrInvalidCoinType = NewError(errors.New("invalid coin type"))
)

const (
	// WalletExt wallet file extension
	WalletExt = "wlt"

	// WalletTimestampFormat wallet timestamp layout
	WalletTimestampFormat = "2006_01_02"

	// CoinTypeSkycoin skycoin type
	CoinTypeSkycoin CoinType = "skycoin"
	// CoinTypeBitcoin bitcoin type
	CoinTypeBitcoin CoinType = "bitcoin"

	// WalletTypeDeterministic deterministic wallet type
	WalletTypeDeterministic = "deterministic"
)

// ResolveCoinType normalizes a coin type string to a CoinType constant
func ResolveCoinType(s string) (CoinType, error) {
	switch strings.ToLower(s) {
	case "sky", "skycoin":
		return CoinTypeSkycoin, nil
	case "btc", "bitcoin":
		return CoinTypeBitcoin, nil
	default:
		return CoinType(""), ErrInvalidCoinType
	}
}

// wallet meta fields
const (
	metaVersion    = "version"    // wallet version
	metaFilename   = "filename"   // wallet file name
	metaLabel      = "label"      // wallet label
	metaTimestamp  = "tm"         // the timestamp when creating the wallet
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

// NewWalletFilename generates a filename from the current time and random bytes
func NewWalletFilename() string {
	timestamp := time.Now().Format(WalletTimestampFormat)
	// should read in wallet files and make sure does not exist
	padding := hex.EncodeToString((cipher.RandByte(2)))
	return fmt.Sprintf("%s_%s.%s", timestamp, padding, WalletExt)
}

// Options options that could be used when creating a wallet
type Options struct {
	Coin       CoinType   // coin type, skycoin, bitcoin, etc.
	Label      string     // wallet label.
	Seed       string     // wallet seed.
	Encrypt    bool       // whether the wallet need to be encrypted.
	Password   []byte     // password that would be used for encryption, and would only be used when 'Encrypt' is true.
	CryptoType CryptoType // wallet encryption type, scrypt-chacha20poly1305 or sha256-xor.
	ScanN      uint64     // number of addresses that're going to be scanned for a balance. The highest address with a balance will be used.
	GenerateN  uint64     // number of addresses to generate, regardless of balance
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
func newWallet(wltName string, opts Options, tf TransactionsFinder) (*Wallet, error) {
	if opts.Seed == "" {
		return nil, ErrMissingSeed
	}

	if opts.ScanN > 0 && tf == nil {
		return nil, ErrNilTransactionsFinder
	}

	coin := opts.Coin
	if coin == "" {
		coin = CoinTypeSkycoin
	}

	switch coin {
	case CoinTypeSkycoin, CoinTypeBitcoin:
	default:
		return nil, fmt.Errorf("Invalid coin type %q", coin)
	}

	w := &Wallet{
		Meta: map[string]string{
			metaFilename:   wltName,
			metaVersion:    Version,
			metaLabel:      opts.Label,
			metaSeed:       opts.Seed,
			metaLastSeed:   opts.Seed,
			metaTimestamp:  strconv.FormatInt(time.Now().Unix(), 10),
			metaType:       WalletTypeDeterministic,
			metaCoin:       string(coin),
			metaEncrypted:  "false",
			metaCryptoType: "",
			metaSecrets:    "",
		},
	}

	// Create a default wallet
	generateN := opts.GenerateN
	if generateN == 0 {
		generateN = 1
	}
	if _, err := w.GenerateAddresses(generateN); err != nil {
		return nil, err
	}

	if opts.ScanN != 0 && coin != CoinTypeSkycoin {
		return nil, errors.New("Wallet address scanning is not supported for Bitcoin wallets")
	}

	if opts.ScanN > generateN {
		// Scan for addresses with balances
		if _, err := w.ScanAddresses(opts.ScanN, tf); err != nil {
			return nil, err
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
	return newWallet(wltName, opts, nil)
}

// NewWalletScanAhead creates wallet and scan ahead N addresses
func NewWalletScanAhead(wltName string, opts Options, tf TransactionsFinder) (*Wallet, error) {
	return newWallet(wltName, opts, tf)
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
		wlt.Erase()
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
	wlt.Erase()

	// Wipes the secret fields in w
	w.Erase()

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
	w.Entries = append(w.Entries, src.Entries...)
}

// Erase wipes secret fields in wallet
func (w *Wallet) Erase() {
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

// GuardUpdate executes a function within the context of a read-write managed decrypted wallet.
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

	defer wlt.Erase()

	if err := fn(wlt); err != nil {
		return err
	}

	if err := wlt.Lock(password, cryptoType); err != nil {
		return err
	}

	*w = *wlt
	// Wipes all sensitive data
	w.Erase()
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

	defer wlt.Erase()

	return f(wlt)
}

// Load loads wallet from a given file
func Load(wltFile string) (*Wallet, error) {
	if _, err := os.Stat(wltFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("wallet %s doesn't exist", wltFile)
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
	if fn := w.Meta[metaFilename]; fn == "" {
		return errors.New("filename not set")
	}

	if tm := w.Meta[metaTimestamp]; tm != "" {
		_, err := strconv.ParseInt(tm, 10, 64)
		if err != nil {
			return errors.New("invalid timestamp")
		}
	}

	walletType, ok := w.Meta[metaType]
	if !ok {
		return errors.New("type field not set")
	}
	if walletType != WalletTypeDeterministic {
		return errors.New("wallet type invalid")
	}

	if coinType := w.Meta[metaCoin]; coinType == "" {
		return errors.New("coin field not set")
	}

	var isEncrypted bool
	if encStr, ok := w.Meta[metaEncrypted]; ok {
		// validate the encrypted value
		var err error
		isEncrypted, err = strconv.ParseBool(encStr)
		if err != nil {
			return errors.New("encrypted field is not a valid bool")
		}
	}

	// checks if the secrets field is empty
	if isEncrypted {
		cryptoType, ok := w.Meta[metaCryptoType]
		if !ok {
			return errors.New("crypto type field not set")
		}

		if _, err := getCrypto(CryptoType(cryptoType)); err != nil {
			return errors.New("unknown crypto type")
		}

		if s := w.Meta[metaSecrets]; s == "" {
			return errors.New("wallet is encrypted, but secrets field not set")
		}
	} else {
		if s := w.Meta[metaSeed]; s == "" {
			return errors.New("seed missing in unencrypted wallet")
		}

		if s := w.Meta[metaLastSeed]; s == "" {
			return errors.New("lastSeed missing in unencrypted wallet")
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

func (w *Wallet) coin() CoinType {
	return CoinType(w.Meta[metaCoin])
}

func (w *Wallet) addressConstructor() func(cipher.PubKey) cipher.Addresser {
	switch w.coin() {
	case CoinTypeSkycoin:
		return func(pk cipher.PubKey) cipher.Addresser {
			return cipher.AddressFromPubKey(pk)
		}
	case CoinTypeBitcoin:
		return func(pk cipher.PubKey) cipher.Addresser {
			return cipher.BitcoinAddressFromPubKey(pk)
		}
	default:
		logger.Panicf("Invalid wallet coin type %q", w.coin())
		return nil
	}
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

func (w *Wallet) timestamp() int64 {
	// Intentionally ignore the error when parsing the timestamp,
	// if it isn't valid or is missing it will be set to 0.
	// Also, this value is validated by wallet.Validate()
	x, _ := strconv.ParseInt(w.Meta[metaTimestamp], 10, 64) //nolint:errcheck
	return x
}

func (w *Wallet) setTimestamp(t int64) {
	w.Meta[metaTimestamp] = strconv.FormatInt(t, 10)
}

// GenerateAddresses generates addresses
func (w *Wallet) GenerateAddresses(num uint64) ([]cipher.Addresser, error) {
	if num == 0 {
		return nil, nil
	}

	if w.IsEncrypted() {
		return nil, ErrWalletEncrypted
	}

	var seckeys []cipher.SecKey
	var seed []byte
	if len(w.Entries) == 0 {
		seed, seckeys = cipher.MustGenerateDeterministicKeyPairsSeed([]byte(w.seed()), int(num))
	} else {
		sd, err := hex.DecodeString(w.lastSeed())
		if err != nil {
			return nil, fmt.Errorf("decode hex seed failed: %v", err)
		}
		seed, seckeys = cipher.MustGenerateDeterministicKeyPairsSeed(sd, int(num))
	}

	w.setLastSeed(hex.EncodeToString(seed))

	addrs := make([]cipher.Addresser, len(seckeys))
	makeAddress := w.addressConstructor()
	for i, s := range seckeys {
		p := cipher.MustPubKeyFromSecKey(s)
		a := makeAddress(p)
		addrs[i] = a
		w.Entries = append(w.Entries, Entry{
			Address: a,
			Secret:  s,
			Public:  p,
		})
	}
	return addrs, nil
}

// GenerateSkycoinAddresses generates Skycoin addresses. If the wallet's coin type is not Skycoin, returns an error
func (w *Wallet) GenerateSkycoinAddresses(num uint64) ([]cipher.Address, error) {
	if w.coin() != CoinTypeSkycoin {
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

// ScanAddresses scans ahead N addresses, truncating up to the highest address with a non-zero balance.
// If any address has a nonzero balance, it rescans N more addresses from that point, until a entire
// sequence of N addresses has no balance.
func (w *Wallet) ScanAddresses(scanN uint64, tf TransactionsFinder) (uint64, error) {
	if w.IsEncrypted() {
		return 0, ErrWalletEncrypted
	}

	if scanN == 0 {
		return 0, nil
	}

	w2 := w.clone()

	nExistingAddrs := uint64(len(w2.Entries))
	nAddAddrs := uint64(0)
	n := scanN
	extraScan := uint64(0)

	for {
		// Generate the addresses to scan
		addrs, err := w2.GenerateSkycoinAddresses(n)
		if err != nil {
			return 0, err
		}

		// Find if these addresses had any activity
		active, err := tf.AddressesActivity(addrs)
		if err != nil {
			return 0, err
		}

		// Check balance from the last one until we find the address that has activity
		var keepNum uint64
		for i := len(active) - 1; i >= 0; i-- {
			if active[i] {
				keepNum = uint64(i + 1)
				break
			}
		}

		if keepNum == 0 {
			break
		}

		nAddAddrs += keepNum + extraScan

		// extraScan is the number of addresses with a zero balance beyond the
		// last address with a nonzero balance
		extraScan = n - keepNum

		// n is the number of addresses to scan the next iteration
		n = scanN - extraScan
	}

	// Regenerate addresses up to nExistingAddrs + nAddAddrss.
	// This is necessary to keep the lastSeed updated.
	w2.reset()
	if _, err := w2.GenerateSkycoinAddresses(nExistingAddrs + nAddAddrs); err != nil {
		return 0, err
	}

	*w = *w2

	return nAddAddrs, nil
}

// GetAddresses returns all addresses in wallet
func (w *Wallet) GetAddresses() []cipher.Addresser {
	addrs := make([]cipher.Addresser, len(w.Entries))
	for i, e := range w.Entries {
		addrs[i] = e.Address
	}
	return addrs
}

// GetSkycoinAddresses returns all Skycoin addresses in wallet. The wallet's coin type must be Skycoin.
func (w *Wallet) GetSkycoinAddresses() ([]cipher.Address, error) {
	if w.coin() != CoinTypeSkycoin {
		return nil, errors.New("Wallet coin type is not Skycoin")
	}

	addrs := make([]cipher.Address, len(w.Entries))
	for i, e := range w.Entries {
		addrs[i] = e.SkycoinAddress()
	}
	return addrs, nil
}

// GetEntry returns entry of given address
func (w *Wallet) GetEntry(a cipher.Address) (Entry, bool) {
	for _, e := range w.Entries {
		if e.SkycoinAddress() == a {
			return e, true
		}
	}
	return Entry{}, false
}

// HasEntry returns true if the wallet has an Entry with a given cipher.Address.
func (w *Wallet) HasEntry(a cipher.Address) bool {
	// This doesn't use GetEntry() to avoid copying an Entry in the return value,
	// which may contain a secret key
	for _, e := range w.Entries {
		if e.SkycoinAddress() == a {
			return true
		}
	}
	return false
}

// AddEntry adds new entry
func (w *Wallet) AddEntry(entry Entry) error {
	// dup check
	for _, e := range w.Entries {
		if e.SkycoinAddress() == entry.SkycoinAddress() {
			return errors.New("duplicate address entry")
		}
	}

	w.Entries = append(w.Entries, entry)
	return nil
}

// clone returns the clone of self
func (w *Wallet) clone() *Wallet {
	wlt := Wallet{
		Meta: make(map[string]string),
	}
	for k, v := range w.Meta {
		wlt.Meta[k] = v
	}

	wlt.Entries = append(wlt.Entries, w.Entries...)

	return &wlt
}
