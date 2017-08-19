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

	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	logger = logging.MustGetLogger("wallet")
)

// WalletExt  wallet file extension
const WalletExt = "wlt"

// WalletTimestampFormat  wallet timestamp layout
const WalletTimestampFormat = "2006_01_02"

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
			"coin":     "sky"},
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
	// check file's existence
	if _, err := os.Stat(wltFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("load wallet file failed, %v", err)
	}
	wlt := Wallet{
		Meta: make(map[string]string),
	}
	wlt.SetFilename(filepath.Base(wltFile))
	dir, err := filepath.Abs(filepath.Dir(wltFile))
	if err != nil {
		return nil, err
	}
	if err := wlt.Load(dir); err != nil {
		return nil, fmt.Errorf("load wallet file failed, %v", err)
	}
	return &wlt, nil
}

// NewWalletFromReadable creates wallet from readable wallet
func NewWalletFromReadable(r *ReadableWallet) Wallet {
	w := Wallet{
		Meta:    r.Meta,
		Entries: r.Entries.ToWalletEntries(),
	}

	err := w.Validate()
	if err != nil {
		logger.Panicf("Wallet %s invalid: %v", w.GetFilename(), err)
	}
	return w
}

// Validate validates the wallet
func (wlt Wallet) Validate() error {
	if _, ok := wlt.Meta["filename"]; !ok {
		return errors.New("filename not set")
	}
	if _, ok := wlt.Meta["seed"]; !ok {
		return errors.New("seed not set")
	}

	walletType, ok := wlt.Meta["type"]
	if !ok {
		return errors.New("type not set")
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
func (wlt *Wallet) GenerateAddresses(num int) []cipher.Address {
	var seckeys []cipher.SecKey
	var sd []byte
	var err error
	if len(wlt.Entries) == 0 {
		sd, seckeys = cipher.GenerateDeterministicKeyPairsSeed([]byte(wlt.getLastSeed()), num)
	} else {
		sd, err = hex.DecodeString(wlt.getLastSeed())
		if err != nil {
			logger.Panicf("decode hex seed failed,%v", err)
		}
		sd, seckeys = cipher.GenerateDeterministicKeyPairsSeed(sd, num)
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

// Save persists wallet to disk
func (wlt *Wallet) Save(dir string) error {
	r := NewReadableWallet(*wlt)
	return r.Save(filepath.Join(dir, wlt.GetFilename()))
}

// Load loads wallets from given dir
func (wlt *Wallet) Load(dir string) error {
	r := &ReadableWallet{}
	if err := r.Load(filepath.Join(dir, wlt.GetFilename())); err != nil {
		return err
	}
	r.Meta["filename"] = wlt.GetFilename()
	*wlt = NewWalletFromReadable(r)
	return nil
}

// Validator validate if the wallet be able to create spending transaction
type Validator interface {
	// checks if any of the given addresses has unconfirmed spending transactions
	HasUnconfirmedSpendTx(addr []cipher.Address) (bool, error)
}

// CreateAndSignTransaction Creates a Transaction
// spending coins and hours from wallet
func (wlt *Wallet) CreateAndSignTransaction(
	vld Validator,
	unspent blockdb.UnspentGetter,
	headTime uint64,
	amt Balance,
	dest cipher.Address) (*coin.Transaction, error) {

	addrs := wlt.GetAddresses()
	ok, err := vld.HasUnconfirmedSpendTx(addrs)
	if err != nil {
		return nil, fmt.Errorf("Checking unconfirmed spending failed: %v", err)
	}

	if ok {
		return nil, errors.New("Please spend after your pending transaction is confirmed")
	}

	txn := coin.Transaction{}
	auxs := unspent.GetUnspentsOfAddrs(wlt.GetAddresses())

	// Determine which unspents to spend
	spends, err := createSpends(headTime, auxs.Flatten(), amt)
	if err != nil {
		return nil, err
	}

	// Add these unspents as tx inputs
	toSign := make([]cipher.SecKey, len(spends))
	spending := Balance{Coins: 0, Hours: 0}
	for i, au := range spends {
		entry, exists := wlt.GetEntry(au.Body.Address)
		if !exists {
			return nil, fmt.Errorf("Address:%v does not exist in wallet:%v", au.Body.Address, wlt.GetID())
		}

		txn.PushInput(au.Hash())
		toSign[i] = entry.Secret
		spending.Coins += au.Body.Coins
		spending.Hours += au.CoinHours(headTime)
	}

	//keep 1/4th of hours as change
	//send half to each address
	var changeHours = uint64(spending.Hours / 4)

	if amt.Coins == spending.Coins {
		txn.PushOutput(dest, amt.Coins, changeHours/2)
		txn.SignInputs(toSign)
		txn.UpdateHeader()
		return &txn, nil
	}

	change := NewBalance(spending.Coins-amt.Coins, changeHours/2)
	// TODO -- send change to a new address
	changeAddr := spends[0].Body.Address

	//create transaction
	txn.PushOutput(changeAddr, change.Coins, change.Hours)
	txn.PushOutput(dest, amt.Coins, changeHours/2)
	txn.SignInputs(toSign)
	txn.UpdateHeader()
	return &txn, nil
}

func createSpends(headTime uint64, uxa coin.UxArray,
	amt Balance) (coin.UxArray, error) {
	if amt.Coins == 0 {
		return nil, errors.New("Zero spend amount")
	}
	if amt.Coins%1e6 != 0 {
		return nil, errors.New("Coins must be multiple of 1e6")
	}

	sort.Sort(uxOutByTimeDesc(uxa))

	have := Balance{Coins: 0, Hours: 0}
	spending := make(coin.UxArray, 0)
	for i := range uxa {
		b := Balance{
			Coins: uxa[i].Body.Coins,
			Hours: uxa[i].CoinHours(headTime),
		}

		if b.Coins == 0 || b.Coins%1e6 != 0 {
			logger.Error("UxOut coins are 0 or not multiple of 1e6, can't spend")
			continue
		}
		have = have.Add(b)
		spending = append(spending, uxa[i])

		if have.Coins >= amt.Coins {
			break
		}
	}

	if amt.Coins > have.Coins {
		return nil, errors.New("Not enough confirmed coins")
	}

	return spending, nil
}

// OldestUxOut sorts a UxArray oldest to newest.
type uxOutByTimeDesc coin.UxArray

func (ouo uxOutByTimeDesc) Len() int      { return len(ouo) }
func (ouo uxOutByTimeDesc) Swap(i, j int) { ouo[i], ouo[j] = ouo[j], ouo[i] }
func (ouo uxOutByTimeDesc) Less(i, j int) bool {
	a := ouo[i].Head.BkSeq
	b := ouo[j].Head.BkSeq
	// Use hash to break ties
	if a == b {
		ih := ouo[i].Hash()
		jh := ouo[j].Hash()
		cmp := bytes.Compare(ih[:], jh[:])
		if cmp == 0 {
			logger.Panic("Duplicate UxOut when sorting")
		}
		return cmp < 0
	}
	return a < b
}
