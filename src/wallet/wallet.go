package wallet

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"encoding/hex"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/blockdb"

	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	// Version represents the current wallet version
	Version = "0.2"

	logger = logging.MustGetLogger("wallet")

	// ErrInsufficientBalance is returned if a wallet does not have enough balance for a spend
	ErrInsufficientBalance = errors.New("balance is not sufficient")
	// ErrSpendingUnconfirmed is returned if caller attempts to spend unconfirmed outputs
	ErrSpendingUnconfirmed = errors.New("please spend after your pending transaction is confirmed")
	// ErrInvalidEncryptedField is returned if a wallet's Meta.encrypted value is invalid.
	ErrInvalidEncryptedField = errors.New(`encrypted field value is not valid, must be "true", "false" or ""`)
	// ErrWalletEncrypted is returned when trying to generate addresses or sign tx in encrypted wallet
	ErrWalletEncrypted = errors.New("wallet is encrypted")
	// ErrWalletNotEncrypted is returned when trying to decrypt unencrypted wallet
	ErrWalletNotEncrypted = errors.New("wallet is not encrypted")
	// ErrMissingPassword is returned when trying to create wallet with encryption, but password is not provided.
	ErrMissingPassword = errors.New("missing password")
	// ErrMissingEncrypt is returned when trying to create wallet with password, but options.Encrypt is not set.
	ErrMissingEncrypt = errors.New("missing encrypt")
	// ErrInvalidPassword is returned if decrypts secrets failed
	ErrInvalidPassword = errors.New("invalid password")
	// ErrMissingSeed is returned when trying to create wallet without a seed
	ErrMissingSeed = errors.New("missing seed")
	// ErrMissingAuthenticated is returned if try to decrypt a scrypt chacha20poly1305 encrypted wallet, and find no authenticated metadata.
	ErrMissingAuthenticated = errors.New("missing authenticated metadata")
	// ErrWrongCryptoType is returned when decrypting wallet with wrong crypto method
	ErrWrongCryptoType = errors.New("wrong crypto type")
	// ErrWalletNotExist is returned if a wallet does not exist
	ErrWalletNotExist = errors.New("wallet doesn't exist")
	// ErrWalletAPIDisabled is returned when trying to do wallet actions while the EnableWalletAPI option is false
	ErrWalletAPIDisabled = errors.New("wallet api is disabled")
	// ErrSeedAPIDisabled is returned when trying to get seed of wallet while the EnableWalletAPI or EnableSeedAPI is false
	ErrSeedAPIDisabled = errors.New("wallet seed api is disabled")
	// ErrWalletNameConflict represents the wallet name conflict error
	ErrWalletNameConflict = errors.New("wallet name would conflict with existing wallet, renaming")
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
	if err := w.lock(opts.Password, opts.CryptoType); err != nil {
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

// lock encrypts the wallet with the given password and specific crypto type
func (w *Wallet) lock(password []byte, cryptoType CryptoType) error {
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

// unlock decrypts the wallet into a temporary decrypted copy of the wallet
// Returns error if the decryption fails
// The temporary decrypted wallet should be erased from memory when done.
func (w *Wallet) unlock(password []byte) (*Wallet, error) {
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
		logger.Errorf("Decrypt wallet failed: %v", err)
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

// guardUpdate executes a function within the context of a read-wirte managed decrypted wallet.
// Returns ErrWalletNotEncrypted if wallet is not encrypted.
func (w *Wallet) guardUpdate(password []byte, fn func(w *Wallet) error) error {
	if !w.IsEncrypted() {
		return ErrWalletNotEncrypted
	}

	if len(password) == 0 {
		return ErrMissingPassword
	}

	cryptoType := w.cryptoType()
	wlt, err := w.unlock(password)
	if err != nil {
		return err
	}

	defer wlt.erase()

	if err := fn(wlt); err != nil {
		return err
	}

	if err := wlt.lock(password, cryptoType); err != nil {
		return err
	}

	*w = *wlt
	// Wipes all sensitive data
	w.erase()
	return nil
}

// guardView executes a function within the context of a read-only managed decrypted wallet.
// Returns ErrWalletNotEncrypted if wallet is not encrypted.
func (w *Wallet) guardView(password []byte, f func(w *Wallet) error) error {
	if !w.IsEncrypted() {
		return ErrWalletNotEncrypted
	}

	if len(password) == 0 {
		return ErrMissingPassword
	}

	wlt, err := w.unlock(password)
	if err != nil {
		return err
	}

	defer wlt.erase()

	if err := f(wlt); err != nil {
		return err
	}
	return nil
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
func (w *Wallet) CreateAndSignTransaction(vld Validator, unspent blockdb.UnspentGetter,
	headTime, coins uint64, dest cipher.Address) (*coin.Transaction, error) {
	if w.IsEncrypted() {
		return nil, ErrWalletEncrypted
	}

	addrs := w.GetAddresses()
	ok, err := vld.HasUnconfirmedSpendTx(addrs)
	if err != nil {
		return nil, fmt.Errorf("checking unconfirmed spending failed: %v", err)
	}

	if ok {
		return nil, ErrSpendingUnconfirmed
	}

	txn := coin.Transaction{}
	auxs := unspent.GetUnspentsOfAddrs(addrs)

	// Determine which unspents to spend.
	// Use the MaximizeUxOuts strategy, this will keep the uxout pool smaller
	uxa := auxs.Flatten()
	uxb, err := NewUxBalances(headTime, uxa)
	if err != nil {
		return nil, err
	}

	spends, err := ChooseSpendsMaximizeUxOuts(uxb, coins)
	if err != nil {
		return nil, err
	}

	// Add these unspents as tx inputs
	toSign := make([]cipher.SecKey, len(spends))
	spending := Balance{Coins: 0, Hours: 0}
	for i, au := range spends {
		entry, exists := w.GetEntry(au.Address)
		if !exists {
			return nil, fmt.Errorf("address:%v does not exist in wallet:%v", au.Address, w.Filename())
		}

		txn.PushInput(au.Hash)

		if w.IsEncrypted() {
			return nil, ErrWalletEncrypted
		}

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

// UxBalance is an intermediate representation of a UxOut for sorting and spend choosing
type UxBalance struct {
	Hash    cipher.SHA256
	BkSeq   uint64
	Address cipher.Address
	Coins   uint64
	Hours   uint64
}

// NewUxBalances converts coin.UxArray to []UxBalance.
// headTime is required to calculate coin hours.
func NewUxBalances(headTime uint64, uxa coin.UxArray) ([]UxBalance, error) {
	uxb := make([]UxBalance, len(uxa))
	for i, ux := range uxa {
		hours, err := ux.CoinHours(headTime)
		if err != nil {
			return nil, err
		}

		b := UxBalance{
			Hash:    ux.Hash(),
			BkSeq:   ux.Head.BkSeq,
			Address: ux.Body.Address,
			Coins:   ux.Body.Coins,
			Hours:   hours,
		}

		uxb[i] = b
	}

	return uxb, nil
}

// ChooseSpendsMinimizeUxOuts chooses uxout spends to satisfy an amount, using the least number of uxouts
//     -- PRO: Allows more frequent spending, less waiting for confirmations, useful for exchanges.
//     -- PRO: When transaction is volume is higher, transactions are prioritized by fee/size. Minimizing uxouts minimizes size.
//     -- CON: Would make the unconfirmed pool grow larger.
// Users with high transaction frequency will want to use this so that they will not need to wait as frequently
// for unconfirmed spends to complete before sending more.
// Alternatively, or in addition to this, they should batch sends into single transactions.
func ChooseSpendsMinimizeUxOuts(uxa []UxBalance, coins uint64) ([]UxBalance, error) {
	return ChooseSpends(uxa, coins, sortSpendsCoinsHighToLow)
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
func ChooseSpendsMaximizeUxOuts(uxa []UxBalance, coins uint64) ([]UxBalance, error) {
	return ChooseSpends(uxa, coins, sortSpendsCoinsLowToHigh)
}

// sortSpendsCoinsLowToHigh sorts uxout spends with lowest balance to highest
func sortSpendsCoinsLowToHigh(uxa []UxBalance) {
	sort.Slice(uxa, makeCmpUxOutByCoins(uxa, func(a, b uint64) bool {
		return a < b
	}))
}

// Sorts UxOuts by those with zero coinhours last.
// Within uxouts that have coinhours and don't have coinhours, respecitvely, they
// they are sorted by ascending or descending coins (depending on coinsCmp).
// If coins are equal, then they are sorted by least hours first
// If hours are equal, then they are sorted by oldest first
// If they are equally old, the UxOut's hash is used to break the tie.
func makeCmpUxOutByCoins(uxa []UxBalance, coinsCmp func(a, b uint64) bool) func(i, j int) bool {
	// Sort by:
	// coins highest or lowest depending on coinsCmp
	//  hours lowest, unless zero, then last
	//   oldest first
	//    tie break with hash comparison
	return func(i, j int) bool {
		a := uxa[i]
		b := uxa[j]

		if a.Coins == b.Coins {
			if a.Hours == b.Hours {
				if a.BkSeq == b.BkSeq {
					return cmpUxOutByHash(a, b)
				}
				return a.BkSeq < b.BkSeq
			}
			return a.Hours < b.Hours
		}
		return coinsCmp(a.Coins, b.Coins)
	}
}

func cmpUxOutByHash(a, b UxBalance) bool {
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
func ChooseSpends(uxa []UxBalance, coins uint64, sortStrategy func([]UxBalance)) ([]UxBalance, error) {
	if coins == 0 {
		return nil, errors.New("zero spend amount")
	}

	if len(uxa) == 0 {
		return nil, errors.New("no unspents to spend")
	}

	for _, ux := range uxa {
		if ux.Coins == 0 {
			logger.Panic("UxOut coins are 0, can't spend")
			return nil, errors.New("UxOut coins are 0, can't spend")
		}
	}

	// Split split UxBalances into those with and without hours
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

	// Sort uxouts with hours, highest coins to lowest
	sortSpendsCoinsHighToLow(nonzero)

	var have Balance
	var spending []UxBalance

	firstNonzero := nonzero[0]
	if firstNonzero.Hours == 0 {
		logger.Panic("balance has zero hours unexpectedly")
		return nil, errors.New("balance has zero hours unexpectedly")
	}

	nonzero = nonzero[1:]

	spending = append(spending, firstNonzero)

	have.Coins += firstNonzero.Coins
	have.Hours += firstNonzero.Hours

	if have.Coins >= coins {
		return spending, nil
	}

	// Sort uxouts without hours according to the sorting strategy
	sortStrategy(zero)

	for _, ux := range zero {
		spending = append(spending, ux)

		have.Coins += ux.Coins
		have.Hours += ux.Hours

		if have.Coins >= coins {
			return spending, nil
		}
	}

	// Sort remaining uxouts with hours according to the sorting strategy
	sortStrategy(nonzero)

	for _, ux := range nonzero {
		spending = append(spending, ux)

		have.Coins += ux.Coins
		have.Hours += ux.Hours

		if have.Coins >= coins {
			return spending, nil
		}
	}

	return nil, ErrInsufficientBalance
}
