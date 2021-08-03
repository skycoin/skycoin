/*
Package wallet implements wallets and the wallet database service

Values of the Wallet interface can be created by calling function NewWallet,
or by loading from `[]byte` that containing wallet data of type such as
"deterministic", "collection", "bip44" or "xpubwallet". Loading any particular
type of wallet requires the prior registration of a loader. Registration is typically
automatic as a side effect of initializing that wallet's package so that, to load a
"deterministic" wallet, it suffices to have
	import _ "github.com/skycoin/skycoin/src/wallet/deterministic"
in a program's main package. The _ means to import a package purely for its
initialization side effects.
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

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/cipher/crypto"
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
	Version = "0.4"

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
	// ErrMissingLabel is returned when trying to create wallet without label
	ErrMissingLabel = NewError(errors.New("missing label"))
	// ErrMissingAuthenticated is returned if try to decrypt a scrypt chacha20poly1305 encrypted wallet, and find no authenticated metadata.
	ErrMissingAuthenticated = NewError(errors.New("missing authenticated metadata"))
	// ErrMissingXPub is returned if try to create a XPub wallet without providing xpub key
	ErrMissingXPub = NewError(errors.New("missing xpub"))
	// ErrWrongCryptoType is returned when decrypting wallet with wrong crypto method
	ErrWrongCryptoType = NewError(errors.New("wrong crypto type"))
	// ErrWalletNotExist is returned if a wallet does not exist
	ErrWalletNotExist = NewError(errors.New("wallet doesn't exist"))
	// ErrSeedUsed is returned if a wallet already exists with the same seed
	ErrSeedUsed = NewError(errors.New("a wallet already exists with this seed"))
	// ErrXPubKeyUsed is returned if a wallet already exists with the same xpub key
	ErrXPubKeyUsed = NewError(errors.New("a wallet already exists with this xpub key"))
	// ErrWalletAPIDisabled is returned when trying to do wallet actions while the EnableWalletAPI option is false
	ErrWalletAPIDisabled = NewError(errors.New("wallet api is disabled"))
	// ErrSeedAPIDisabled is returned when trying to get seed of wallet while the EnableWalletAPI or EnableSeedAPI is false
	ErrSeedAPIDisabled = NewError(errors.New("wallet seed api is disabled"))
	// ErrWalletNameConflict represents the wallet name conflict error
	ErrWalletNameConflict = NewError(errors.New("wallet name would conflict with existing wallet, renaming"))
	// ErrWalletRecoverSeedWrong is returned if the seed or seed passphrase does not match the specified wallet when recovering
	ErrWalletRecoverSeedWrong = NewError(errors.New("wallet recovery seed or seed passphrase is wrong"))
	// ErrWalletSeedPassphrase is returned when using seed passphrase for none bip44 wallet
	ErrWalletSeedPassphrase = NewError(errors.New("seedPassphrase is only used for \"bip44\" wallets"))
	// ErrNilTransactionsFinder is returned if Options.ScanN > 0 but a nil TransactionsFinder was provided
	ErrNilTransactionsFinder = NewError(errors.New("scan ahead requested but balance getter is nil"))
	// ErrInvalidCoinType is returned for invalid coin types
	ErrInvalidCoinType = NewError(errors.New("invalid coin type"))
	// ErrInvalidWalletType is returned for invalid wallet types
	ErrInvalidWalletType = NewError(errors.New("invalid wallet type"))
	// ErrWalletTypeNotRecoverable is returned by RecoverWallet is the wallet type does not support recovery
	ErrWalletTypeNotRecoverable = NewError(errors.New("wallet type is not recoverable"))
	// ErrWalletPermission is returned when updating a wallet without writing permission
	ErrWalletPermission = NewError(errors.New("saving wallet permission denied"))

	// ErrEntryNotFound is returned by GetEntry is the wallet does not contains the entry
	ErrEntryNotFound = errors.New("entry not found")
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
	// WalletTypeXPub xpub HD wallet type.
	// Allows generating addresses without a secret key
	WalletTypeXPub = "xpub"
)

// CoinType represents the wallet coin type, which refers to the pubkey2addr method used
type CoinType string

// NewWalletFilename generates a filename from the current time and random bytes
func NewWalletFilename() string {
	timestamp := time.Now().Format(WalletTimestampFormat)
	// should read in wallet files and make sure does not exist
	padding := hex.EncodeToString(cipher.RandByte(2))
	return fmt.Sprintf("%s_%s.%s", timestamp, padding, WalletExt)
}

// Options options that could be used when creating a wallet
type Options struct {
	Version               string
	Type                  string            // wallet type: deterministic, collection. Refers to which key generation mechanism is used.
	Coin                  CoinType          // coin type: skycoin, bitcoin, etc. Refers to which pubkey2addr method is used.
	Bip44Coin             *bip44.CoinType   // bip44 path coin type
	Label                 string            // wallet label
	Seed                  string            // wallet seed
	SeedPassphrase        string            // wallet seed passphrase (bip44 wallets only)
	Encrypt               bool              // whether the wallet need to be encrypted.
	Password              []byte            // password that would be used for encryption, and would only be used when 'Encrypt' is true.
	CryptoType            crypto.CryptoType // wallet encryption type, scrypt-chacha20poly1305 or sha256-xor.
	ScanN                 uint64            // number of addresses that're going to be scanned for a balance. The highest address with a balance will be used.
	GenerateN             uint64            // number of addresses to generate, regardless of balance
	XPub                  string            // xpub key (xpub wallets only)
	Decoder               Decoder
	TF                    TransactionsFinder
	Temp                  bool            // whether the wallet is created temporary in memory.
	CollectionPrivateKeys []cipher.SecKey // private keys for collection wallet
}

func (opts Options) Validate() error {
	if opts.Type == WalletTypeDeterministic && opts.SeedPassphrase != "" {
		return ErrWalletSeedPassphrase
	}
	return nil
}

//go:generate mockery -name Wallet -case underscore -inpkg -testonly

// Wallet defines the wallet API
type Wallet interface {
	Seed() string
	LastSeed() string
	SeedPassphrase() string
	Timestamp() int64
	SetTimestamp(int64)
	Coin() CoinType
	SetCoin(coinType CoinType)
	// Type returns the wallet type, e.g. bip44, deterministic, collection
	Type() string
	// Bip44Coin returns the coin_type part of bip44 path
	Bip44Coin() *bip44.CoinType
	SetBip44Coin(ct bip44.CoinType)
	Label() string
	SetLabel(string)
	Filename() string
	SetFilename(string)
	IsEncrypted() bool
	// CryptoType returns the crypto type for encrypting/decrypting the wallet
	CryptoType() crypto.CryptoType
	SetCryptoType(ct crypto.CryptoType)
	// SetDecoder sets the wallet decoder
	SetDecoder(d Decoder)
	// Version returns the wallet version
	Version() string
	// Secrets returns the wallet secrets data
	Secrets() string
	// XPub returns the xpub key of a xpub wallet
	XPub() string
	// Lock encrypts the wallet
	Lock(password []byte) error
	// Unlock decrypts the wallets, returns an copy of the decrypted wallet
	Unlock(password []byte) (Wallet, error)
	// Erase wipes sensitive data
	Erase()
	// Clone returns a copy of the wallet
	Clone() Wallet
	// CopyFrom copies the src wallet to w
	// CopyFrom(src Wallet)
	// CopyFromRef copies the src wallet with a pointer dereference
	CopyFromRef(src Wallet)
	Fingerprint() string
	// ScanAddresses scans ahead given number of addresses
	ScanAddresses(scanN uint64, tf TransactionsFinder) ([]cipher.Addresser, error)
	// GetAddresses returns all addresses.
	// for bip44 wallet, if no options are specified, addresses on external chain of account
	// with index 0 will be returned.
	GetAddresses(options ...Option) ([]cipher.Addresser, error)
	// GenerateAddresses generates N addresses,
	// for bip44 wallet, if no options are specified, addresses will be generated
	// on external chain of account with index 0.
	GenerateAddresses(num uint64, options ...Option) ([]cipher.Addresser, error)
	// Entries returns entries,
	// for bip44 wallet if no options are used, entries on external chain of account
	// with index 0 will be returned.
	GetEntries(options ...Option) (Entries, error)
	// GetEntryAt returns the entry of given index,
	// for bip44 wallet, if no options are specified, the entry on external chain of the account
	// with index 0 will be returned.
	GetEntryAt(i int, options ...Option) (Entry, error)
	// GetEntry return the entry by address
	// for bip44 wallet, if no options are specified, it will search the external chain of account
	// of index 0.
	GetEntry(addr cipher.Addresser, options ...Option) (Entry, error)
	// HasEntry returns whether the entry exists in the wallet
	// for bip44 wallet, if no options are specified, it will check the external chain of account
	// of index 0.
	HasEntry(addr cipher.Addresser, options ...Option) (bool, error)
	// EntriesLen returns the entries length
	// for bip44 wallet, if no options are specified, the length of the entries on external chain of account
	// with index 0 will be returned.
	EntriesLen(options ...Option) (int, error)
	// Accounts returns the list of account for bip44 wallet
	Accounts() []Bip44Account
	// Serialize serialize the wallet to bytes, and error if any
	Serialize() ([]byte, error)
	// Deserialize deserialize the data to a Wallet, and error if any
	Deserialize(data []byte) error
	// IsTemp returns whether the wallet is a temporary wallet
	IsTemp() bool
	// SetTemp sets wallet temporary flag
	SetTemp(temp bool)
}

// Decoder is the interface that wraps the Encode and Decode methods.
// Encode method encodes the wallet to bytes, Decode method decodes bytes to bip44 wallet.
type Decoder interface {
	Encode(w Wallet) ([]byte, error)
	Decode(b []byte) (Wallet, error)
}

// NewWallet creates a new wallet
func NewWallet(filename, label, seed string, options Options) (Wallet, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}

	c, ok := getCreator(options.Type)
	if !ok {
		return nil, fmt.Errorf("wallet.NewWallet failed, wallet type %q is not supported", options.Type)
	}

	return c.Create(filename, label, seed, options)
}

// Bip44Account represents the wallet account
type Bip44Account struct {
	Name  string
	Index uint32
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

	wlt, err := w.Unlock(password)
	if err != nil {
		return err
	}

	defer wlt.Erase()

	if err := fn(wlt); err != nil {
		return err
	}

	if err := wlt.Lock(password); err != nil {
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

	wlt, err := w.Unlock(password)
	if err != nil {
		return err
	}

	// TODO: consider to catch panic and erase sensitive data in recovery function
	// in case the f(wlt) function get panic

	defer wlt.Erase()

	return f(wlt)
}

type walletLoadMeta struct {
	Meta struct {
		Type    string `json:"type"`
		Version string `json:"version"`
	} `json:"meta"`
}

// Save saves the wallet to a directory. The wallet's filename is read from its metadata.
func Save(w Wallet, dir string) error {
	if w.IsTemp() {
		return nil
	}

	data, err := w.Serialize()
	if err != nil {
		return err
	}
	return file.SaveBinary(filepath.Join(dir, w.Filename()), data, 0600)
}

// Load loads wallet from a file
func Load(filename string) (Wallet, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("wallet %q doesn't exist", filename)
	}

	// Load the wallet meta type field from JSON
	m, err := loadWalletMeta(filename)
	if err != nil {
		return nil, err
	}

	if m.Meta.Type == "" {
		err := errors.New("missing meta.type field")
		logger.WithError(err).WithField("filename", filename)
		return nil, err
	}

	// Depending on the wallet type in the wallet metadata header, load the full wallet data
	l, ok := getLoader(m.Meta.Type)
	if !ok {
		logger.Errorf("wallet loader for type of %q not found", m.Meta.Type)
		return nil, nil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	w, err := l.Load(data)
	if err != nil {
		return nil, err
	}

	w.SetFilename(filepath.Base(filename))
	return w, nil
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
			m, err := loadWalletMeta(f)
			if err != nil {
				return err
			}

			if m.Meta.Version == "0.1" {
				if err := os.Remove(bf); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func loadWalletMeta(filename string) (*walletLoadMeta, error) {
	var m walletLoadMeta
	if err := file.LoadJSON(filename, &m); err != nil {
		logger.WithError(err).WithField("filename", filename).Error("Load: file.LoadJSON failed")
		return nil, err
	}

	return &m, nil
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

// IsValidWalletType returns true if a wallet type is recognized
func IsValidWalletType(t string) bool {
	switch t {
	case WalletTypeDeterministic,
		WalletTypeCollection,
		WalletTypeBip44,
		WalletTypeXPub:
		return true
	default:
		return false
	}
}

// AddressConstructor returns a function to create a cipher.Addresser from a cipher.PubKey
func AddressConstructor(m Meta) func(cipher.PubKey) cipher.Addresser {
	switch m.Coin() {
	case CoinTypeSkycoin:
		return func(pk cipher.PubKey) cipher.Addresser {
			return cipher.AddressFromPubKey(pk)
		}
	case CoinTypeBitcoin:
		return func(pk cipher.PubKey) cipher.Addresser {
			return cipher.BitcoinAddressFromPubKey(pk)
		}
	default:
		logger.Panicf("Invalid wallet coin type %q", m.Coin())
		return nil
	}
}

// ValidateMeta validates the common meta data when initializing a wallet
func ValidateMeta(m Meta) error {
	if fn := m[MetaFilename]; fn == "" {
		return errors.New("filename not set")
	}

	if tm := m[MetaTimestamp]; tm != "" {
		_, err := strconv.ParseInt(tm, 10, 64)
		if err != nil {
			return errors.New("invalid timestamp")
		}
	}

	if tp := m[MetaType]; tp == "" {
		return errors.New("type field not set")
	}

	if coinType := m[MetaCoin]; coinType == "" {
		return errors.New("coin field not set")
	}

	var isEncrypted bool
	if encStr, ok := m[MetaEncrypted]; ok {
		// validate the encrypted value
		var err error
		isEncrypted, err = strconv.ParseBool(encStr)
		if err != nil {
			return errors.New("encrypted field is not a valid bool")
		}
	}

	if isEncrypted {
		if s := m[MetaSecrets]; s == "" {
			return errors.New("wallet is encrypted, but secrets field not set")
		}
	} else {
		if s := m[MetaSecrets]; s != "" {
			return errors.New("secrets should not be in unencrypted wallets")
		}
	}

	return nil
}

// ValidateMetaCryptoType validates meta crypto type
func ValidateMetaCryptoType(m Meta) error {
	cryptoType, ok := m[MetaCryptoType]
	if !ok {
		return errors.New("crypto type field not set")
	}

	if _, err := crypto.GetCrypto(crypto.CryptoType(cryptoType)); err != nil {
		return errors.New("unknown crypto type")
	}

	return nil
}

// ValidateMetaSeed validate meta seed
func ValidateMetaSeed(m Meta) error {
	if m.IsEncrypted() {
		if s := m[MetaSeed]; s != "" {
			return errors.New("seed should not be visible in encrypted wallets")
		}
	} else {
		if s := m[MetaSeed]; s == "" {
			return ErrMissingSeed
		}
	}
	return nil
}

// SkycoinAddresses converts the addresses to skycoin addresses
func SkycoinAddresses(addrs []cipher.Addresser) []cipher.Address {
	skyAddrs := make([]cipher.Address, len(addrs))
	for i, a := range addrs {
		skyAddrs[i] = a.(cipher.Address)
	}
	return skyAddrs
}
