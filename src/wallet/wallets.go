package wallet

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
)

type Wallets []Wallet

// Loads all wallets contained in wallet dir.  If any regular file in wallet
// dir fails to load, loading is aborted and error returned.  Only files with
// extension WalletExt are considered
func LoadWallets(dir string) (Wallets, error) {
	// TODO -- don't load duplicate wallets.
	// TODO -- save a last_modified value in wallets to decide which to load
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	have := make(map[WalletID]Wallet, len(entries))
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
			id := w.GetID()
			if kw, ok := have[id]; ok {
				return nil, fmt.Errorf("Duplicate wallet file detected. "+
					"%s and %s are the same wallet.", kw.GetFilename(), name)
			}
			have[id] = w
			w.SetFilename(name)
			wallets = append(wallets, w)
		}
	}
	return wallets, nil
}

func (self *Wallets) Add(w Wallet) {
	*self = append(*self, w)
}

func (self Wallets) Get(walletID WalletID) *Wallet {
	for _, w := range self {
		if w.GetID() == walletID {
			return &w
		}
	}
	return nil
}

func (self Wallets) Save(dir string) map[WalletID]error {
	errs := make(map[WalletID]error)
	for _, w := range self {
		if err := w.Save(dir); err != nil {
			errs[w.GetID()] = err
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func (self Wallets) GetAddressSet() map[cipher.Address]byte {
	set := make(AddressSet)
	for _, w := range self {
		set.Update(w.GetAddressSet())
	}
	return set
}

func (self Wallets) toReadable(f ReadableWalletCtor) []*ReadableWallet {
	rw := make([]*ReadableWallet, len(self))
	for i, w := range self {
		rw[i] = f(w)
	}
	return rw
}

func (self Wallets) ToReadable() []*ReadableWallet {
	return self.toReadable(NewReadableWallet)
}

//func (self Wallets) ToPublicReadable() []*ReadableWallet {
//	return self.toReadable(NewPublicReadableWallet)
//}
