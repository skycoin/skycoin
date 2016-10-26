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
	"github.com/skycoin/skycoin/src/cipher/secp256k1-go"
)

// Wallet contains meta data and address entries.
// Meta:
// 		Filename
// 		Seed
//		Type - wallet type
//		Coin - coin type
type Wallet struct {
	Meta    map[string]string
	Entries []WalletEntry
}

var version = "0.1"

// NewWallet generates Deterministic Wallet
// generates a random seed if seed is ""
func NewWallet(seed, wltName, label string) Wallet {
	//if seed is blank, generate a new seed
	if seed == "" {
		seedRaw := cipher.SumSHA256(secp256k1.RandByte(64))
		seed = hex.EncodeToString(seedRaw[:])
	}

	// generate the first address.
	// pub, sec := cipher.GenerateDeterministicKeyPair([]byte(seed[:]))
	return Wallet{
		Meta: map[string]string{
			"filename": wltName,
			"version":  version,
			"label":    label,
			"seed":     seed,
			"lastSeed": seed,
			"tm":       fmt.Sprintf("%v", time.Now().Unix()),
			"type":     "deterministic",
			"coin":     "sky"},
	}
}

func Load(wltFile string) (*Wallet, error) {
	// check file's existence
	if _, err := os.Stat(wltFile); os.IsNotExist(err) {
		return nil, err
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
		return nil, err
	}
	return &wlt, nil
}

func NewWalletFromReadable(r *ReadableWallet) Wallet {
	w := Wallet{
		Meta:    r.Meta,
		Entries: r.Entries.ToWalletEntries(),
	}

	err := w.Validate()
	if err != nil {
		log.Panic("Wallet %s invalid: %v", w.GetFilename, err)
	}
	return w
}

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

	coinType, ok := wlt.Meta["coin"]
	if !ok {
		return errors.New("coin field not set")
	}
	if coinType != "sky" {
		return errors.New("coin type invalid")
	}

	return nil

}

func (wlt Wallet) GetType() string {
	return wlt.Meta["type"]
}

func (wlt Wallet) GetFilename() string {
	return wlt.Meta["filename"]
}

func (wlt *Wallet) SetFilename(fn string) {
	wlt.Meta["filename"] = fn
}

func (wlt Wallet) GetID() string {
	return wlt.Meta["filename"]
}

func (wlt Wallet) GetLabel() string {
	return wlt.Meta["label"]
}

func (wlt Wallet) getLastSeed() string {
	return wlt.Meta["lastSeed"]
}

func (wlt *Wallet) setLastSeed(lseed string) {
	wlt.Meta["lastSeed"] = lseed
}

func (wlt *Wallet) GetVersion() string {
	return wlt.Meta["version"]
}

func (wlt Wallet) NumEntries() int {
	return len(wlt.Entries)
}

func (wlt *Wallet) GenerateAddresses(num int) []cipher.Address {
	var seckeys []cipher.SecKey
	var sd []byte
	var err error
	if len(wlt.Entries) == 0 {
		sd, seckeys = cipher.GenerateDeterministicKeyPairsSeed([]byte(wlt.getLastSeed()), num)
	} else {
		sd, err = hex.DecodeString(wlt.getLastSeed())
		if err != nil {
			log.Panicf("decode hex seed faild,%v", err)
		}
		sd, seckeys = cipher.GenerateDeterministicKeyPairsSeed(sd, num)
	}
	wlt.setLastSeed(hex.EncodeToString(sd))
	addrs := make([]cipher.Address, len(seckeys))
	for i, s := range seckeys {
		p := cipher.PubKeyFromSecKey(s)
		a := cipher.AddressFromPubKey(p)
		addrs[i] = a
		wlt.Entries = append(wlt.Entries, WalletEntry{
			Address: a,
			Secret:  s,
			Public:  p,
		})
	}
	return addrs
}

func (wlt *Wallet) GetAddresses() []cipher.Address {
	addrs := make([]cipher.Address, len(wlt.Entries))
	for i, e := range wlt.Entries {
		addrs[i] = e.Address
	}
	return addrs
}

func (wlt *Wallet) GetAddressSet() map[cipher.Address]byte {
	set := make(map[cipher.Address]byte)
	for _, e := range wlt.Entries {
		set[e.Address] = byte(1)
	}
	return set
}

func (wlt *Wallet) GetEntry(a cipher.Address) (WalletEntry, bool) {
	for _, e := range wlt.Entries {
		if e.Address == a {
			return e, true
		}
	}
	return WalletEntry{}, false
}

func (wlt *Wallet) Save(dir string) error {
	r := NewReadableWallet(*wlt)
	return r.Save(filepath.Join(dir, wlt.GetFilename()))
}

func (wlt *Wallet) Load(dir string) error {
	r := &ReadableWallet{}
	if err := r.Load(filepath.Join(dir, wlt.GetFilename())); err != nil {
		return err
	}
	*wlt = NewWalletFromReadable(r)
	return nil
}
