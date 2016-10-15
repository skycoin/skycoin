package wallet

import (
	//"fmt"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
)

type Wallets map[string]*Wallet

// LoadWallets Loads all wallets contained in wallet dir.  If any regular file in wallet
// dir fails to load, loading is aborted and error returned.  Only files with
// extension WalletExt are considered
func LoadWallets(dir string) (Wallets, error) {
	// TODO -- don't load duplicate wallets.
	// TODO -- save a last_modified value in wallets to decide which to load
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	//have := make(map[WalletID]Wallet, len(entries))
	wallets := make(Wallets, 0)
	for _, e := range entries {
		if e.Mode().IsRegular() {
			name := e.Name()
			if !strings.HasSuffix(name, WalletExt) {
				continue
			}
			fullpath := filepath.Join(dir, name)
			rw, err := LoadReadableWallet(fullpath)
			if err != nil {
				return nil, err
			}
			w, err := rw.ToWallet()
			if err != nil {
				return nil, err
			}
			logger.Info("Loaded wallet from %s", fullpath)
			w.SetFilename(name)
			wallets[name] = &w
		}
	}
	return wallets, nil
}

// Add add walet to current wallet
func (wlts *Wallets) Add(w Wallet) error {
	if _, dup := (*wlts)[w.GetFilename()]; dup {
		return errors.New("Wallets.Add, Wallet name would conflict with existing wallet, renaming")
	}

	(*wlts)[w.GetFilename()] = &w
	return nil
}

func (wlts *Wallets) Get(wltID string) (Wallet, bool) {
	if w, ok := (*wlts)[wltID]; ok {
		return *w, true
	}
	return Wallet{}, false
}

func (wlts *Wallets) NewAddresses(wltID string, num int) ([]cipher.Address, error) {
	if w, ok := (*wlts)[wltID]; ok {
		return w.GenerateAddresses(num), nil
	}
	return nil, fmt.Errorf("wallet: %v does not exist", wltID)
}

//check for name conflicts!
//resolve conflicts for saving wallets who have different names
func (wlts Wallets) Save(dir string) map[string]error {
	errs := make(map[string]error)
	for id, w := range wlts {
		if err := w.Save(dir); err != nil {
			errs[id] = err
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// GetAddressSet get all addresses.
func (wlts Wallets) GetAddressSet() map[cipher.Address]byte {
	set := make(map[cipher.Address]byte)
	for _, w := range wlts {
		for _, a := range w.GetAddresses() {
			set[a] = byte(1)
		}
	}
	return set
}

func (wlts Wallets) toReadable(f ReadableWalletCtor) []*ReadableWallet {
	var rw []*ReadableWallet
	for _, w := range wlts {
		rw = append(rw, f(*w))
	}
	sort.Sort(ByTm(rw))
	return rw
}

func (wlts Wallets) ToReadable() []*ReadableWallet {
	return wlts.toReadable(NewReadableWallet)
}
