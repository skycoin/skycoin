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
	bip39 "github.com/skycoin/skycoin/src/cipher/go-bip39"
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
	//should read in wallet files and make sure does not exist
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

// Option NewWallet optional arguments type
type Option func(w *Wallet)

// NewWallet generates Deterministic Wallet
// generates a random seed if seed is ""
func NewWallet(wltName string, opts ...Option) (*Wallet, error) {
	// generaten bip39 as default seed
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil, fmt.Errorf("generate bip39 entropy failed, err:%v", err)
	}

	seed, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("generate bip39 seed failed, err:%v", err)
	}

	w := &Wallet{
		Meta: map[string]string{
			"filename": wltName,
			"version":  version,
			"label":    "",
			"seed":     seed,
			"lastSeed": seed,
			"tm":       fmt.Sprintf("%v", time.Now().Unix()),
			"type":     "deterministic",
			"coin":     string(CoinTypeSkycoin),
		},
	}

	for _, opt := range opts {
		opt(w)
	}

	return w, nil
}

// OptCoin NewWallet function's optional argument
func OptCoin(coin string) Option {
	return func(w *Wallet) {
		w.Meta["coin"] = coin
	}
}

// OptLabel NewWallet function's optional argument
func OptLabel(label string) Option {
	return func(w *Wallet) {
		w.Meta["label"] = label
	}
}

// OptSeed NewWallet function's optional argument
func OptSeed(sd string) Option {
	return func(w *Wallet) {
		if sd != "" {
			w.Meta["seed"] = sd
			w.Meta["lastSeed"] = sd
		}
	}
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
func (wlt Wallet) Validate() error {
	if _, ok := wlt.Meta["filename"]; !ok {
		return errors.New("filename not set")
	}
	if _, ok := wlt.Meta["seed"]; !ok {
		return errors.New("seed field not set")
	}

	walletType, ok := wlt.Meta["type"]
	if !ok {
		return errors.New("type field not set")
	}
	if walletType != "deterministic" {
		return errors.New("wallet type invalid")
	}

	if _, ok := wlt.Meta["coin"]; !ok {
		return errors.New("coin field not set")
	}

	return nil
}

// GetType gets the wallet type
func (wlt Wallet) GetType() string {
	return wlt.Meta["type"]
}

// GetFilename gets the wallet filename
func (wlt Wallet) GetFilename() string {
	return wlt.Meta["filename"]
}

// SetFilename sets the wallet filename
func (wlt *Wallet) SetFilename(fn string) {
	wlt.Meta["filename"] = fn
}

// GetID gets the wallet id
func (wlt Wallet) GetID() string {
	return wlt.Meta["filename"]
}

// GetLabel gets the wallet label
func (wlt Wallet) GetLabel() string {
	return wlt.Meta["label"]
}

// SetLabel sets the wallet label
func (wlt *Wallet) SetLabel(label string) {
	wlt.Meta["label"] = label
}

func (wlt Wallet) getLastSeed() string {
	return wlt.Meta["lastSeed"]
}

func (wlt *Wallet) setLastSeed(lseed string) {
	wlt.Meta["lastSeed"] = lseed
}

// GetVersion gets the wallet version
func (wlt *Wallet) GetVersion() string {
	return wlt.Meta["version"]
}

// NumEntries returns the number of entries
func (wlt Wallet) NumEntries() int {
	return len(wlt.Entries)
}

// GenerateAddresses generate addresses of given number
func (wlt *Wallet) GenerateAddresses(num uint64) []cipher.Address {
	if num == 0 {
		return []cipher.Address{}
	}

	var seckeys []cipher.SecKey
	var sd []byte
	var err error
	if len(wlt.Entries) == 0 {
		sd, seckeys = cipher.GenerateDeterministicKeyPairsSeed([]byte(wlt.getLastSeed()), int(num))
	} else {
		sd, err = hex.DecodeString(wlt.getLastSeed())
		if err != nil {
			logger.Panicf("decode hex seed failed,%v", err)
		}
		sd, seckeys = cipher.GenerateDeterministicKeyPairsSeed(sd, int(num))
	}
	wlt.setLastSeed(hex.EncodeToString(sd))
	addrs := make([]cipher.Address, len(seckeys))
	for i, s := range seckeys {
		p := cipher.PubKeyFromSecKey(s)
		a := cipher.AddressFromPubKey(p)
		addrs[i] = a
		wlt.Entries = append(wlt.Entries, Entry{
			Address: a,
			Secret:  s,
			Public:  p,
		})
	}
	return addrs
}

// GetAddresses returns all addresses in wallet
func (wlt *Wallet) GetAddresses() []cipher.Address {
	addrs := make([]cipher.Address, len(wlt.Entries))
	for i, e := range wlt.Entries {
		addrs[i] = e.Address
	}
	return addrs
}

// GetEntry returns entry of given address
func (wlt *Wallet) GetEntry(a cipher.Address) (Entry, bool) {
	for _, e := range wlt.Entries {
		if e.Address == a {
			return e, true
		}
	}
	return Entry{}, false
}

// AddEntry adds new entry
func (wlt *Wallet) AddEntry(entry Entry) error {
	// dup check
	for _, e := range wlt.Entries {
		if e.Address == entry.Address {
			return errors.New("duplicate address entry")
		}
	}

	wlt.Entries = append(wlt.Entries, entry)
	return nil
}

// Reset resets the wallet entries and move the lastSeed to origin
func (wlt *Wallet) Reset() {
	wlt.Entries = []Entry{}
	wlt.Meta["lastSeed"] = wlt.Meta["seed"]
}

// Save persists wallet to disk
func (wlt *Wallet) Save(dir string) error {
	r := NewReadableWallet(*wlt)
	return r.Save(filepath.Join(dir, wlt.GetFilename()))
}

// Load loads wallets from given wallet file
func (wlt *Wallet) Load(wltFile string) error {
	if _, err := os.Stat(wltFile); os.IsNotExist(err) {
		return fmt.Errorf("load wallet file failed, wallet %s doesn't exist", wltFile)
	}

	r := &ReadableWallet{}
	if err := r.Load(wltFile); err != nil {
		return err
	}

	// update filename meta info with the real filename
	r.Meta["filename"] = filepath.Base(wltFile)
	w, err := newWalletFromReadable(r)
	if err != nil {
		return err
	}

	*wlt = *w
	return nil
}

// Copy returns the copy of wallet
func (wlt *Wallet) Copy() Wallet {
	w := Wallet{Meta: make(map[string]string)}
	for k, v := range wlt.Meta {
		w.Meta[k] = v
	}

	for _, e := range wlt.Entries {
		w.Entries = append(w.Entries, e)
	}

	return w
}

// Validator validate if the wallet be able to create spending transaction
type Validator interface {
	// checks if any of the given addresses has unconfirmed spending transactions
	HasUnconfirmedSpendTx(addr []cipher.Address) (bool, error)
}

// CreateAndSignTransaction Creates a Transaction
// spending coins and hours from wallet
func (wlt *Wallet) CreateAndSignTransaction(vld Validator, unspent blockdb.UnspentGetter,
	headTime uint64, amt Balance, dest cipher.Address) (*coin.Transaction, error) {

	addrs := wlt.GetAddresses()
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
	spends, err := CreateSpendsMaximizeUxOuts(uxb, amt)
	if err != nil {
		return nil, err
	}

	// Add these unspents as tx inputs
	toSign := make([]cipher.SecKey, len(spends))
	spending := Balance{Coins: 0, Hours: 0}
	for i, au := range spends {
		entry, exists := wlt.GetEntry(au.Address)
		if !exists {
			return nil, fmt.Errorf("address:%v does not exist in wallet:%v", au.Address, wlt.GetID())
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
	changeCoins := spending.Coins - amt.Coins
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

	txn.PushOutput(dest, amt.Coins, addrHours[0])

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
	// TODO: Allow the caller to control coinhour distribution
	remainingHours := inputHours / fee.BurnFactor

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

// createSpendsOldestFirst chooses uxout spends to satisfy an amount, prioritizing older oxouts
func createSpendsOldestFirst(uxa []UxBalance, amt Balance) ([]UxBalance, error) {
	sort.Slice(uxa, func(i, j int) bool {
		a := uxa[i].BkSeq
		b := uxa[j].BkSeq
		// Use hash to break ties
		if a == b {
			return cmpUxOutByHash(uxa[i], uxa[j])
		}
		return a < b
	})

	return ChooseSpends(uxa, amt)
}

// CreateSpendsMinimizeUxOuts chooses uxout spends to satisfy an amount, using the least number of uxouts
//     -- PRO: Allows more frequent spending, less waiting for confirmations, useful for exchanges.
//     -- PRO: When transaction is volume is higher, transactions are prioritized by fee/size. Minimizing uxouts minimizes size.
//     -- CON: Would make the unconfirmed pool grow larger.
// Users with high transaction frequency will want to use this so that they will not need to wait as frequently
// for unconfirmed spends to complete before sending more.
// Alternatively, or in addition to this, they should batch sends into single transactions.
func CreateSpendsMinimizeUxOuts(uxa []UxBalance, amt Balance) ([]UxBalance, error) {
	sort.Slice(uxa, makeCmpUxOutByCoins(uxa, func(a, b uint64) bool {
		return a < b
	}))

	return ChooseSpends(uxa, amt)
}

// CreateSpendsMaximizeUxOuts chooses uxout spends to satisfy an amount, using the most number of uxouts
// See the pros and cons of CreateSpendsMinimizeUxOuts.
// This should be the default mode, because this keeps the unconfirmed pool smaller which will allow
// the network to scale better.
func CreateSpendsMaximizeUxOuts(uxa []UxBalance, amt Balance) ([]UxBalance, error) {
	sort.Slice(uxa, makeCmpUxOutByCoins(uxa, func(a, b uint64) bool {
		return a > b
	}))
	return ChooseSpends(uxa, amt)
}

// ChooseSpends chooses uxouts from a prioritized list of uxouts.
// uxOuts with zero coinhours should be sorted last, to avoid choosing unspents
// that all have no coin hours, since a valid transaction requires at least 1 coinhour.
// Make sure that the UxBalances have the coinhours updated based on
// blockchain head time before calling this.
func ChooseSpends(uxa []UxBalance, amt Balance) ([]UxBalance, error) {
	if amt.Coins == 0 {
		return nil, errors.New("zero spend amount")
	}

	var have Balance
	var spending []UxBalance
	for i := range uxa {
		b := Balance{
			Coins: uxa[i].Coins,
			Hours: uxa[i].Hours,
		}

		if b.Coins == 0 {
			logger.Error("UxOut coins are 0, can't spend")
			continue
		}

		have = have.Add(b)
		spending = append(spending, uxa[i])

		if have.Coins >= amt.Coins && have.Hours != 0 {
			break
		}
	}

	if amt.Coins > have.Coins {
		return nil, ErrInsufficientBalance
	}

	if have.Hours == 0 {
		return nil, fee.ErrTxnNoFee
	}

	return spending, nil
}

// Sorts UxOuts by those with zero coinhours last.
// Within uxouts that have coinhours and don't have coinhours, respecitvely, they
// they are sorted by ascending or descending coins (depending on coinsCmp).
// If coins are equal, then they are sorted by least hours first
// If hours are equal, then they are sorted by oldest first
// If they are equally old, the UxOut's hash is used to break the tie.
func makeCmpUxOutByCoins(uxa []UxBalance, coinsCmp func(a, b uint64) bool) func(i, j int) bool {
	cmpUxOutByCoins := func(a, b UxBalance) bool {
		// Sort by:
		// coins highest
		//  hours lowest, unless zero, then last
		//   oldest first
		//    tie break with hash comparison
		if a.Coins == b.Coins {
			if a.Hours == b.Hours {
				if a.BkSeq == b.BkSeq {
					return cmpUxOutByHash(a, b)
				}

				return a.BkSeq < b.BkSeq
			}

			// Sort by least hours, unless hours are zero, then sort them last
			if a.Hours == 0 {
				return false
			} else if b.Hours == 0 {
				return true
			}

			return a.Hours < b.Hours
		}

		return coinsCmp(a.Coins, b.Coins)
	}

	return func(i, j int) bool {
		a := uxa[i]
		b := uxa[j]

		if a.Hours == 0 {
			if a.Hours == b.Hours {
				// If they are both zero, sort as normal
				return cmpUxOutByCoins(a, b)
			}

			// If a's Hours are 0, sort it last
			return false
		} else if b.Hours == 0 {
			// If b's Hours are 0, sort it last
			return false
		}

		// If both a's and b's Hours are non-zero, sort as normal
		return cmpUxOutByCoins(a, b)
	}
}

func cmpUxOutByHash(a, b UxBalance) bool {
	cmp := bytes.Compare(a.Hash[:], b.Hash[:])
	if cmp == 0 {
		logger.Panic("Duplicate UxOut when sorting")
	}
	return cmp < 0

}
