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
	logger = logging.MustGetLogger("wallet")

	// ErrInsufficientBalance is returned if a wallet does not have enough balance for a spend
	ErrInsufficientBalance = errors.New("balance is not sufficient")
)

// CoinType represents the wallet coin type
type CoinType string

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

// NewWalletFilename check for collisions and retry if failure
func NewWalletFilename() string {
	timestamp := time.Now().Format(WalletTimestampFormat)
	// should read in wallet files and make sure does not exist
	padding := hex.EncodeToString((cipher.RandByte(2)))
	return fmt.Sprintf("%s_%s.%s", timestamp, padding, WalletExt)
}

// Wallet contains meta data and address entries.
// Meta:
// 		Filename
// 		Seed
//		Type - wallet type
//		Coin - coin type
type Wallet struct {
	Meta    map[string]string
	Entries []Entry
}

var version = "0.1"

// Options are wallet constructor options
type Options struct {
	Coin  CoinType
	Label string
	Seed  string
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
			"version":  version,
			"label":    opts.Label,
			"seed":     seed,
			"lastSeed": seed,
			"tm":       fmt.Sprintf("%v", time.Now().Unix()),
			"type":     "deterministic",
			"coin":     string(coin),
		},
	}

	return w, nil
}

// Load loads wallet from given file
func Load(wltFile string) (*Wallet, error) {
	w := Wallet{}
	if err := w.Load(wltFile); err != nil {
		return nil, err
	}

	return &w, nil
}

// newWalletFromReadable creates wallet from readable wallet
func newWalletFromReadable(r *ReadableWallet) (*Wallet, error) {
	ets, err := r.Entries.ToWalletEntries()
	if err != nil {
		return nil, err
	}

	w := Wallet{
		Meta:    r.Meta,
		Entries: ets,
	}

	if err := w.Validate(); err != nil {
		return nil, fmt.Errorf("invalid wallet %s: %v", w.GetFilename(), err)
	}

	return &w, nil
}

// Validate validates the wallet
func (w *Wallet) Validate() error {
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

	return nil
}

// GetType gets the wallet type
func (w *Wallet) GetType() string {
	return w.Meta["type"]
}

// GetFilename gets the wallet filename
func (w *Wallet) GetFilename() string {
	return w.Meta["filename"]
}

// SetFilename sets the wallet filename
func (w *Wallet) SetFilename(fn string) {
	w.Meta["filename"] = fn
}

// GetID gets the wallet id
func (w *Wallet) GetID() string {
	return w.Meta["filename"]
}

// GetLabel gets the wallet label
func (w *Wallet) GetLabel() string {
	return w.Meta["label"]
}

// SetLabel sets the wallet label
func (w *Wallet) SetLabel(label string) {
	w.Meta["label"] = label
}

func (w *Wallet) getLastSeed() string {
	return w.Meta["lastSeed"]
}

func (w *Wallet) setLastSeed(lseed string) {
	w.Meta["lastSeed"] = lseed
}

// GetVersion gets the wallet version
func (w *Wallet) GetVersion() string {
	return w.Meta["version"]
}

// NumEntries returns the number of entries
func (w *Wallet) NumEntries() int {
	return len(w.Entries)
}

// GenerateAddresses generate addresses of given number and adds them to the wallet
func (w *Wallet) GenerateAddresses(num uint64) []cipher.Address {
	if num == 0 {
		return []cipher.Address{}
	}

	var seckeys []cipher.SecKey
	var seed []byte
	if len(w.Entries) == 0 {
		seed, seckeys = cipher.GenerateDeterministicKeyPairsSeed([]byte(w.getLastSeed()), int(num))
	} else {
		var err error
		seed, err = hex.DecodeString(w.getLastSeed())
		if err != nil {
			logger.Panicf("decode hex seed failed: %v", err)
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
	return addrs
}

// ScanAddresses scans ahead N addresses to find one with non-zero coins
func (w *Wallet) ScanAddresses(scanN uint64, bg BalanceGetter) error {
	if scanN <= 0 {
		return nil
	}

	nExistingAddrs := uint64(w.NumEntries())

	// Generate the addresses to scan
	addrs := w.GenerateAddresses(scanN)

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
		w.Reset()
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

// Reset resets the wallet entries and move the lastSeed to origin
func (w *Wallet) Reset() {
	w.Entries = []Entry{}
	w.Meta["lastSeed"] = w.Meta["seed"]
}

// Save persists wallet to disk
func (w *Wallet) Save(dir string) error {
	r := NewReadableWallet(*w)
	return r.Save(filepath.Join(dir, w.GetFilename()))
}

// Load loads wallets from given wallet file
func (w *Wallet) Load(wltFile string) error {
	if _, err := os.Stat(wltFile); os.IsNotExist(err) {
		return fmt.Errorf("load wallet file failed, wallet %s doesn't exist", wltFile)
	}

	r := &ReadableWallet{}
	if err := r.Load(wltFile); err != nil {
		return err
	}

	// update filename meta info with the real filename
	r.Meta["filename"] = filepath.Base(wltFile)
	wlt, err := newWalletFromReadable(r)
	if err != nil {
		return err
	}

	*w = *wlt
	return nil
}

// Copy returns the copy of wallet
func (w *Wallet) Copy() Wallet {
	wlt := Wallet{Meta: make(map[string]string)}
	for k, v := range w.Meta {
		wlt.Meta[k] = v
	}

	for _, e := range w.Entries {
		wlt.Entries = append(wlt.Entries, e)
	}

	return wlt
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
			return nil, fmt.Errorf("address:%v does not exist in wallet:%v", au.Address, w.GetID())
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
