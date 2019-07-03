/*
Package wallet implements wallets and the wallet database service
*/
package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
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
	// ErrInvalidWalletType is returned for invalid wallet types
	ErrInvalidWalletType = NewError(errors.New("invalid wallet type"))
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

	// WalletTypeDeterministic deterministic wallet type.
	// Uses the original Skycoin deterministic key generator.
	WalletTypeDeterministic = "deterministic"
	// WalletTypeCollection collection wallet type.
	// Does not use any key generator; keys must be added explicitly
	WalletTypeCollection = "collection"
	// WalletTypeBip44 bip44 HD wallet type.
	// Follow the bip44 spec.
	WalletTypeBip44 = "bip44"
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

// IsValidWalletType returns true if a wallet type is recognized
func IsValidWalletType(t string) bool {
	switch t {
	case WalletTypeDeterministic, WalletTypeCollection:
		return true
	default:
		return false
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

// CoinType represents the wallet coin type, which refers to the pubkey2addr method used
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
	Type       string     // wallet type: deterministic, collection. Refers to which key generation mechanism is used.
	Coin       CoinType   // coin type: skycoin, bitcoin, etc. Refers to which pubkey2addr method is used.
	Label      string     // wallet label.
	Seed       string     // wallet seed.
	Encrypt    bool       // whether the wallet need to be encrypted.
	Password   []byte     // password that would be used for encryption, and would only be used when 'Encrypt' is true.
	CryptoType CryptoType // wallet encryption type, scrypt-chacha20poly1305 or sha256-xor.
	ScanN      uint64     // number of addresses that're going to be scanned for a balance. The highest address with a balance will be used.
	GenerateN  uint64     // number of addresses to generate, regardless of balance
}

// newWallet creates a wallet instance with given name and options.
func newWallet(wltName string, opts Options, tf TransactionsFinder) (Wallet, error) {
	wltType := opts.Type
	if wltType == "" {
		wltType = WalletTypeDeterministic
	}
	if !IsValidWalletType(wltType) {
		return nil, ErrInvalidWalletType
	}

	switch wltType {
	case WalletTypeDeterministic:
		if opts.Seed == "" {
			return nil, ErrMissingSeed
		}

		if opts.ScanN > 0 && tf == nil {
			return nil, ErrNilTransactionsFinder
		}

	case WalletTypeCollection:
		if opts.Seed != "" {
			return nil, NewError(errors.New("seed should not be provided for \"collection\" wallets"))
		}

	default:
		return nil, ErrInvalidWalletType
	}

	coin := opts.Coin
	if coin == "" {
		coin = CoinTypeSkycoin
	}
	coin, err := ResolveCoinType(string(coin))
	if err != nil {
		return nil, err
	}

	meta := Meta{
		metaFilename:   wltName,
		metaVersion:    Version,
		metaLabel:      opts.Label,
		metaSeed:       opts.Seed,
		metaLastSeed:   opts.Seed,
		metaTimestamp:  strconv.FormatInt(time.Now().Unix(), 10),
		metaType:       wltType,
		metaCoin:       string(coin),
		metaEncrypted:  "false",
		metaCryptoType: "",
		metaSecrets:    "",
	}

	// Create the wallet
	var w Wallet
	switch wltType {
	case WalletTypeDeterministic:
		w = newDeterministicWallet(meta)
	case WalletTypeCollection:
		w = newCollectionWallet(meta)
	default:
		logger.Panic("unhandled wltType")
	}

	// Generate wallet addresses
	switch wltType {
	case WalletTypeDeterministic:
		generateN := opts.GenerateN
		if generateN == 0 {
			generateN = 1
		}

		logger.WithField("generateN", generateN).Info("Generating addresses for deterministic wallet")

		if _, err := w.GenerateAddresses(generateN); err != nil {
			return nil, err
		}

		if opts.ScanN != 0 && coin != CoinTypeSkycoin {
			return nil, errors.New("Wallet scanning is not supported for Bitcoin wallets")
		}

		if opts.ScanN > generateN {
			// Scan for addresses with balances
			logger.WithField("scanN", opts.ScanN).Info("Scanning addresses for deterministic wallet")
			if err := w.ScanAddresses(opts.ScanN, tf); err != nil {
				return nil, err
			}
		}

	case WalletTypeCollection:
		if opts.GenerateN != 0 || opts.ScanN != 0 {
			return nil, NewError(errors.New("wallet scanning is not defined for \"collection\" wallets"))
		}

	default:
		logger.Panic("unhandled wltType")
	}

	// Check if the wallet should be encrypted
	if !opts.Encrypt {
		if len(opts.Password) != 0 {
			return nil, ErrMissingEncrypt
		}
		return w, nil
	}

	// Check if the password is provided
	if len(opts.Password) == 0 {
		return nil, ErrMissingPassword
	}

	// Check crypto type
	if opts.CryptoType == "" {
		opts.CryptoType = DefaultCryptoType
	}

	if _, err := getCrypto(opts.CryptoType); err != nil {
		return nil, err
	}

	// Encrypt the wallet
	if err := Lock(w, opts.Password, opts.CryptoType); err != nil {
		return nil, err
	}

	// Validate the wallet
	if err := w.Validate(); err != nil {
		return nil, err
	}

	return w, nil
}

// NewWallet creates wallet without scanning addresses
func NewWallet(wltName string, opts Options) (Wallet, error) {
	return newWallet(wltName, opts, nil)
}

// NewWalletScanAhead creates wallet and scan ahead N addresses
func NewWalletScanAhead(wltName string, opts Options, tf TransactionsFinder) (Wallet, error) {
	return newWallet(wltName, opts, tf)
}

// Lock encrypts the wallet with the given password and specific crypto type
func Lock(w Wallet, password []byte, cryptoType CryptoType) error {
	if len(password) == 0 {
		return ErrMissingPassword
	}

	if w.IsEncrypted() {
		return ErrWalletEncrypted
	}

	wlt := w.Clone()

	// Records seeds in secrets
	ss := make(Secrets)
	defer func() {
		// Wipes all unencrypted sensitive data
		ss.erase()
		wlt.Erase()
	}()

	wlt.PackSecrets(ss)

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

	// Sets wallet as encrypted
	wlt.SetEncrypted(cryptoType, string(encSecret))

	// Update the wallet to the latest version, which indicates encryption support
	wlt.SetVersion(Version)

	// Wipes unencrypted sensitive data
	wlt.Erase()

	// Wipes the secret fields in w
	w.Erase()

	// Replace the original wallet with new encrypted wallet
	w.CopyFrom(wlt)
	return nil
}

// Unlock decrypts the wallet into a temporary decrypted copy of the wallet
// Returns error if the decryption fails
// The temporary decrypted wallet should be erased from memory when done.
func Unlock(w Wallet, password []byte) (Wallet, error) {
	if !w.IsEncrypted() {
		return nil, ErrWalletNotEncrypted
	}

	if len(password) == 0 {
		return nil, ErrMissingPassword
	}

	wlt := w.Clone()

	// Gets the secrets string
	sstr := w.Secrets()
	if sstr == "" {
		return nil, errors.New("secrets missing from wallet")
	}

	ct := w.CryptoType()
	if ct == "" {
		return nil, errors.New("missing crypto type")
	}

	// Gets the crypto module
	crypto, err := getCrypto(ct)
	if err != nil {
		return nil, err
	}

	// Decrypts the secrets
	sb, err := crypto.Decrypt([]byte(sstr), password)
	if err != nil {
		return nil, ErrInvalidPassword
	}

	defer func() {
		// Wipe the data from the secrets bytes buffer
		for i := range sb {
			sb[i] = 0
		}
	}()

	// Deserialize into secrets
	ss := make(Secrets)
	defer ss.erase()
	if err := ss.deserialize(sb); err != nil {
		return nil, err
	}

	if err := wlt.UnpackSecrets(ss); err != nil {
		return nil, err
	}

	wlt.SetDecrypted()

	return wlt, nil
}

// Wallet defines the wallet API
type Wallet interface {
	Find(string) string
	Seed() string
	LastSeed() string
	Timestamp() int64
	SetTimestamp(int64)
	Coin() CoinType
	Type() string
	Label() string
	SetLabel(string)
	Filename() string
	IsEncrypted() bool
	SetEncrypted(cryptoType CryptoType, encryptedSecrets string)
	SetDecrypted()
	CryptoType() CryptoType
	Version() string
	SetVersion(string)
	AddressConstructor() func(cipher.PubKey) cipher.Addresser
	Secrets() string

	UnpackSecrets(ss Secrets) error
	PackSecrets(ss Secrets)

	Erase()
	Clone() Wallet
	CopyFrom(src Wallet)
	CopyFromRef(src Wallet)

	ToReadable() Readable

	Validate() error

	Fingerprint() string
	GetAddresses() []cipher.Addresser
	GetSkycoinAddresses() ([]cipher.Address, error)
	GetEntryAt(i int) Entry
	GetEntry(cipher.Address) (Entry, bool)
	HasEntry(cipher.Address) bool
	EntriesLen() int
	GetEntries() Entries

	GenerateAddresses(num uint64) ([]cipher.Addresser, error)
	GenerateSkycoinAddresses(num uint64) ([]cipher.Address, error)
	ScanAddresses(scanN uint64, tf TransactionsFinder) error
}

// GuardUpdate executes a function within the context of a read-write managed decrypted wallet.
// Returns ErrWalletNotEncrypted if wallet is not encrypted.
func GuardUpdate(w Wallet, password []byte, fn func(w Wallet) error) error {
	if !w.IsEncrypted() {
		return ErrWalletNotEncrypted
	}

	if len(password) == 0 {
		return ErrMissingPassword
	}

	cryptoType := w.CryptoType()
	wlt, err := Unlock(w, password)
	if err != nil {
		return err
	}

	defer wlt.Erase()

	if err := fn(wlt); err != nil {
		return err
	}

	if err := Lock(wlt, password, cryptoType); err != nil {
		return err
	}

	w.CopyFromRef(wlt)

	// Wipes all sensitive data
	w.Erase()
	return nil
}

// GuardView executes a function within the context of a read-only managed decrypted wallet.
// Returns ErrWalletNotEncrypted if wallet is not encrypted.
func GuardView(w Wallet, password []byte, f func(w Wallet) error) error {
	if !w.IsEncrypted() {
		return ErrWalletNotEncrypted
	}

	if len(password) == 0 {
		return ErrMissingPassword
	}

	wlt, err := Unlock(w, password)
	if err != nil {
		return err
	}

	defer wlt.Erase()

	return f(wlt)
}

type walletLoadMeta struct {
	Meta struct {
		Type string `json:"type"`
	} `json:"meta"`
}

type walletLoader interface {
	SetFilename(string)
	SetCoin(CoinType)
	Coin() CoinType
	ToWallet() (Wallet, error)
}

// Load loads wallet from a given file
func Load(filename string) (Wallet, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("wallet %q doesn't exist", filename)
	}

	// Load the wallet meta type field from JSON
	var m walletLoadMeta
	if err := file.LoadJSON(filename, &m); err != nil {
		logger.WithError(err).WithField("filename", filename).Error("Load: file.LoadJSON failed")
		return nil, err
	}

	if !IsValidWalletType(m.Meta.Type) {
		logger.WithError(ErrInvalidWalletType).WithFields(logrus.Fields{
			"filename":   filename,
			"walletType": m.Meta.Type,
		}).Error("wallet meta loaded from disk has invalid wallet type")
		return nil, fmt.Errorf("invalid wallet %q: %v", filename, ErrInvalidWalletType)
	}

	// Depending on the wallet type in the wallet metadata header, load the full wallet data
	var rw walletLoader
	var err error
	switch m.Meta.Type {
	case WalletTypeDeterministic:
		logger.WithField("filename", filename).Info("LoadReadableDeterministicWallet")
		rw, err = LoadReadableDeterministicWallet(filename)
	case WalletTypeCollection:
		logger.WithField("filename", filename).Info("LoadReadableCollectionWallet")
		rw, err = LoadReadableCollectionWallet(filename)
	default:
		logger.WithField("walletType", m.Meta.Type).Error("Load: unhandled wallet type")
		return nil, ErrInvalidWalletType
	}

	if err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"filename":   filename,
			"walletType": m.Meta.Type,
		}).Error("Load readable wallet failed")
		return nil, err
	}

	// Make sure "sky", "btc" normalize to "skycoin", "bitcoin"
	ct, err := ResolveCoinType(string(rw.Coin()))
	if err != nil {
		logger.WithError(err).WithField("coinType", rw.Coin()).Error("Load: invalid coin type")
		return nil, fmt.Errorf("invalid wallet %q: %v", filename, err)
	}
	rw.SetCoin(ct)

	rw.SetFilename(filepath.Base(filename))

	return rw.ToWallet()
}

// Save saves the wallet to a directory. The wallet's filename is read from its metadata.
func Save(w Wallet, dir string) error {
	rw := w.ToReadable()
	return file.SaveJSON(filepath.Join(dir, rw.Filename()), rw, 0600)
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
