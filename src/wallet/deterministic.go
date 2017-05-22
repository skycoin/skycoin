package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	bip39 "github.com/skycoin/skycoin/src/cipher/go-bip39"
	"github.com/skycoin/skycoin/src/util"
)

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
func NewWallet(wltName string, opts ...Option) Wallet {
	//old seed generation
	//seedRaw := cipher.SumSHA256(secp256k1.RandByte(64))
	//seed := hex.EncodeToString(seedRaw[:])

	// generaten bip39 as default seed
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		log.Panicf("generate bip39 entropy failed, err:%v", err)
	}

	seed, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Panicf("generate bip39 seed failed, err:%v", err)
	}

	w := Wallet{
		Meta: map[string]string{
			"filename":     wltName,
			"version":      version,
			"label":        "",
			"walletFolder": util.UserHome() + "/.skycoin/wallets",
			"seed":         seed,
			"lastSeed":     seed,
			"tm":           fmt.Sprintf("%v", time.Now().Unix()),
			"type":         "deterministic",
			"coin":         "sky"},
	}

	for _, opt := range opts {
		opt(&w)
	}

	return w
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
		log.Panicf("Wallet %s invalid: %v", w.GetFilename(), err)
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

	// if _, ok := wlt.Meta["lastSeed"]; !ok {
	// 	return errors.New("lastSeed not set")
	// }

	walletType, ok := wlt.Meta["type"]
	if !ok {
		return errors.New("type not set")
	}
	if walletType != "deterministic" {
		return errors.New("wallet type invalid")
	}

	// coinType, ok := wlt.Meta["coin"]
	if _, ok := wlt.Meta["coin"]; !ok {
		return errors.New("coin field not set")
	}
	// if coinType != "sky" {
	// 	return errors.New("coin type invalid")
	// }

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
			log.Panicf("decode hex seed failed,%v", err)
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

// GetAddressSet returns address in map
func (wlt *Wallet) GetAddressSet() map[cipher.Address]byte {
	set := make(map[cipher.Address]byte)
	for _, e := range wlt.Entries {
		set[e.Address] = byte(1)
	}
	return set
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
