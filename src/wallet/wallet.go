package wallet

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	// ErrInvalidEncryptedFieldValue is returned if a wallet's Meta.encrypted value is invalid.
	ErrInvalidEncryptedFieldValue = errors.New(`encrypted field value is not valid, must be "true", "false" or ""`)
	// ErrWalletEncrypted is returned when trying to generate addresses or sign tx in encrypted wallet
	ErrWalletEncrypted = errors.New("wallet is encrypted")
	// ErrWalletNotEncrypted is returned when trying to decrypt unencrypted wallet
	ErrWalletNotEncrypted = errors.New("wallet is not encrypted")
	// ErrRequirePassword find no password when creating wallet
	ErrRequirePassword = errors.New("password is required")
	// ErrInvalidWalletVersion represents invalid wallet version erro
	ErrInvalidWalletVersion = errors.New("invalid wallet version")
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

// CoinType represents the wallet coin type
type CoinType string

// Options are wallet constructor options
type Options struct {
	Coin       CoinType
	Label      string
	Seed       string
	Encrypt    bool
	Password   []byte
	AddressNum uint64 // Generate N addresses when create wallet
}

// Option NewWallet optional arguments type
type Option func(w *Wallet)

// newWalletFilename check for collisions and retry if failure
func newWalletFilename() string {
	timestamp := time.Now().Format(WalletTimestampFormat)
	// should read in wallet files and make sure does not exist
	padding := hex.EncodeToString((cipher.RandByte(2)))
	return fmt.Sprintf("%s_%s.%s", timestamp, padding, WalletExt)
}

// Wallet contains meta data and address entries.
//
// Meta:
//      filename
//      version
//      label
// 		encrypted - whether this wallet is encrypted
//      seed
//      lastSeed - seed for generating next address
//      tm - timestamp when creating the wallet
//      type - wallet type
//      coin - coin type
type Wallet struct {
	Meta    map[string]string
	Entries []Entry
}

// NewWallet generates Deterministic Wallet
// generates a random seed if seed is ""
func NewWallet(wltName string, opts Options) (*Wallet, error) {
	seed := opts.Seed
	if seed == "" {
		return nil, errors.New("seed required")
	}

	coin := opts.Coin
	if coin == "" {
		coin = CoinTypeSkycoin
	}

	w := &Wallet{
		Meta: map[string]string{
			"filename": wltName,
			"version":  Version,
			"label":    opts.Label,
			"seed":     seed,
			"lastSeed": seed,
			"tm":       fmt.Sprintf("%v", time.Now().Unix()),
			"type":     "deterministic",
			"coin":     string(coin),
		},
	}

	// Generate addresses
	if _, err := w.GenerateAddresses(opts.AddressNum); err != nil {
		return nil, fmt.Errorf("generate addresses failed when creating wallets: %v", err)
	}

	if !opts.Encrypt {
		return w, nil
	}

	if opts.Password == nil {
		return nil, errors.New("password is required for creating wallet with encryption")
	}

	if err := w.lock(opts.Password); err != nil {
		return nil, fmt.Errorf("lock wallet failed: %v", err)
	}

	// Update the encrypted meta field
	w.setEncrypted(true)

	return w, nil
}

// lock encrypts the wallet with password
func (w *Wallet) lock(password []byte) error {
	if password == nil {
		return ErrRequirePassword
	}

	if w.IsEncrypted() {
		return ErrWalletEncrypted
	}

	// Encrypt the seed
	ss, err := Encrypt([]byte(w.seed()), password)
	if err != nil {
		return err
	}

	w.setSeed(ss)

	// Encrypt the last seed
	sls, err := Encrypt([]byte(w.lastSeed()), password)
	if err != nil {
		return err
	}

	w.setLastSeed(sls)

	// encrypt private keys in entries
	for i, e := range w.Entries {
		se, err := Encrypt(e.Secret[:], password)
		if err != nil {
			return err
		}

		// Set the encrypted seckey value
		w.Entries[i].EncryptedSeckey = se
		// Clear the entry.Secret
		w.Entries[i].Secret = cipher.SecKey{}
	}

	w.setEncrypted(true)

	return nil
}

// unlock decrypts the wallet into a temporary decrypted copy of the wallet
// It returns an error if decryption fails
// The temporary decrypted wallet should be erased from memory when done.
func (w *Wallet) unlock(password []byte) (*Wallet, error) {
	if !w.IsEncrypted() {
		return nil, ErrWalletNotEncrypted
	}

	if password == nil {
		return nil, errors.New("password is required to decrypt wallet")
	}

	wlt := w.clone()

	// decrypt the seed
	s, err := Decrypt(wlt.seed(), password)
	if err != nil {
		return nil, err
	}
	wlt.setSeed(string(s))

	// decrypt lastSeed
	ls, err := Decrypt(wlt.lastSeed(), password)
	if err != nil {
		return nil, err
	}
	wlt.setLastSeed(string(ls))

	// decrypt the entries
	for i := range wlt.Entries {
		sk, err := Decrypt(wlt.Entries[i].EncryptedSeckey, password)
		if err != nil {
			return nil, err
		}
		copy(wlt.Entries[i].Secret[:], sk[:])
		wlt.Entries[i].EncryptedSeckey = ""
	}
	wlt.setEncrypted(false)

	return wlt, nil
}

// guard will do:
// 1. unlock the encrypted wallet
// 2. process with the decrypted wallet by calling the callback function
// 3. lock the wallet at the end again
// If the wallet is not encrypted, it would return ErrWalletNotEncrypted error
func (w *Wallet) guard(password []byte, f func(w *Wallet) error) (err error) {
	if !w.IsEncrypted() {
		return ErrWalletNotEncrypted
	}

	if password == nil {
		return ErrRequirePassword
	}

	var wlt *Wallet
	wlt, err = w.unlock(password)
	if err != nil {
		return fmt.Errorf("unlock wallet failed: %v", err)
	}

	defer func() {
		if lockErr := wlt.lock(password); lockErr != nil {
			err = fmt.Errorf("lock wallet failed: %v", lockErr)
			return
		}

		*w = *wlt
	}()

	if err := f(wlt); err != nil {
		return fmt.Errorf("process wallet failed: %v", err)
	}

	return
}

// Load loads wallet from given file
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
	return r.toWallet()
}

// Save saves the wallet to given dir
func Save(dir string, w *Wallet) error {
	r := NewReadableWallet(w)
	return r.Save(filepath.Join(dir, w.Filename()))
}

// reset resets the wallet entries and move the lastSeed to origin
func (w *Wallet) reset() {
	w.Entries = []Entry{}
	w.setLastSeed(w.seed())
}

// Validate validates the wallet
func (w *Wallet) validate() error {
	if _, ok := w.Meta["filename"]; !ok {
		return errors.New("filename not set")
	}
	if _, ok := w.Meta["seed"]; !ok {
		return errors.New("seed field not set")
	}

	walletType, ok := w.Meta["type"]
	if !ok {
		return errors.New("type field not set")
	}
	if walletType != "deterministic" {
		return errors.New("wallet type invalid")
	}

	if _, ok := w.Meta["coin"]; !ok {
		return errors.New("coin field not set")
	}

	switch w.Meta["encrypted"] {
	case "true", "false", "":
	default:
		return ErrInvalidEncryptedFieldValue
	}

	return nil
}

// Type gets the wallet type
func (w *Wallet) Type() string {
	return w.Meta["type"]
}

// Version gets the wallet version
func (w *Wallet) Version() string {
	return w.Meta["version"]
}

func (w *Wallet) setVersion(v string) {
	w.Meta["version"] = v
}

// Filename gets the wallet filename
func (w *Wallet) Filename() string {
	return w.Meta["filename"]
}

// setFilename sets the wallet filename
func (w *Wallet) setFilename(fn string) {
	w.Meta["filename"] = fn
}

// Label gets the wallet label
func (w *Wallet) Label() string {
	return w.Meta["label"]
}

// setLabel sets the wallet label
func (w *Wallet) setLabel(label string) {
	w.Meta["label"] = label
}

// lastSeed returns the last seed
func (w *Wallet) lastSeed() string {
	return w.Meta["lastSeed"]
}

func (w *Wallet) setLastSeed(lseed string) {
	w.Meta["lastSeed"] = lseed
}

func (w *Wallet) seed() string {
	return w.Meta["seed"]
}

func (w *Wallet) setSeed(seed string) {
	w.Meta["seed"] = seed
}

// GenerateAddresses generates addresses
func (w *Wallet) GenerateAddresses(num uint64) ([]cipher.Address, error) {
	if w.IsEncrypted() {
		return nil, ErrWalletEncrypted
	}

	if num == 0 {
		return nil, nil
	}

	var seckeys []cipher.SecKey
	var seed []byte
	var err error
	if len(w.Entries) == 0 {
		seed, seckeys = cipher.GenerateDeterministicKeyPairsSeed([]byte(w.lastSeed()), int(num))
	} else {
		seed, err = hex.DecodeString(w.lastSeed())
		if err != nil {
			return nil, fmt.Errorf("decode hex seed failed: %v", err)
		}
		seed, seckeys = cipher.GenerateDeterministicKeyPairsSeed(seed, int(num))
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

// // CreateAndSignTransaction Creates a Transaction
// // spending coins and hours from wallet
// func (w *Wallet) CreateAndSignTransaction(vld Validator, unspent blockdb.UnspentGetter,
// 	headTime, coins uint64, dest cipher.Address) (*coin.Transaction, error) {
// 	if w.IsEncrypted() {
// 		return nil, ErrWalletEncrypted
// 	}

// 	return w.createAndSignTransaction(vld, unspent, headTime, coins, dest)
// }

// // CreateAndSignTransactionEncrypted creates and signs the transaction
// func (w *Wallet) CreateAndSignTransactionEncrypted(vld Validator, unspent blockdb.UnspentGetter,
// 	headTime, coins uint64, dest cipher.Address, password []byte) (*coin.Transaction, error) {
// 	var tx *coin.Transaction
// 	if err := w.guard(password, func(wlt *Wallet) error {
// 		var err error
// 		tx, err = wlt.createAndSignTransaction(vld, unspent, headTime, coins, dest)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}); err != nil {
// 		return nil, err
// 	}

// 	return tx, nil
// }

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
		return nil, errors.New("please spend after your pending transaction is confirmed")
	}

	txn := coin.Transaction{}
	auxs := unspent.GetUnspentsOfAddrs(addrs)

	// Determine which unspents to spend.
	// Use the MaximizeUxOuts strategy, this will keep the uxout pool smaller
	uxa := auxs.Flatten()
	uxb := NewUxBalances(headTime, uxa)
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

	logger.Info("wallet.CreateAndSignTransaction: spending.Hours=%d, fee.VerifyTransactionFeeForHours(%d, %d)", spending.Hours, outputHours, spending.Hours-outputHours)
	if err := fee.VerifyTransactionFeeForHours(outputHours, spending.Hours-outputHours); err != nil {
		logger.Warning("wallet.CreateAndSignTransaction: fee.VerifyTransactionFeeForHours failed: %v", err)
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

func (w *Wallet) setEncrypted(encrypt bool) {
	if encrypt {
		w.Meta["encrypted"] = "true"
	} else {
		w.Meta["encrypted"] = "false"
	}
}

// IsEncrypted checks whether the wallet is encrypted.
// Check the "encrypted" meta field:
//     - return true if "true".
//     - return false if "false" or "".
func (w *Wallet) IsEncrypted() bool {
	return checkEncrypted(w.Meta["encrypted"])
}

func checkEncrypted(v string) bool {
	switch v {
	// return false if it's value is "false" or empty string, cause old wallets do
	// not have this field.
	case "true":
		return true
	case "false", "":
		return false
	default:
		panic(ErrInvalidEncryptedFieldValue)
	}
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
func NewUxBalances(headTime uint64, uxa coin.UxArray) []UxBalance {
	uxb := make([]UxBalance, len(uxa))
	for i, ux := range uxa {
		b := UxBalance{
			Hash:    ux.Hash(),
			BkSeq:   ux.Head.BkSeq,
			Address: ux.Body.Address,
			Coins:   ux.Body.Coins,
			Hours:   ux.CoinHours(headTime),
		}

		uxb[i] = b
	}

	return uxb
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
