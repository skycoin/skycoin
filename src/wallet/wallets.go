package wallet

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/amherag/skycoin/src/cipher"
)

// Wallets wallets map
type Wallets map[string]*Wallet

// LoadWallets Loads all wallets contained in wallet dir.  If any regular file in wallet
// dir fails to load, loading is aborted and error returned.  Only files with
// extension WalletExt are considered.
func LoadWallets(dir string) (Wallets, error) {
	// TODO -- return error if duplicate wallet (by first address) is found
	// TODO -- but make sure that the client has a good warning to the user if wallet loading fails
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	wallets := Wallets{}
	for _, e := range entries {
		if e.Mode().IsRegular() {
			name := e.Name()
			if !strings.HasSuffix(name, WalletExt) {
				continue
			}

			fullpath := filepath.Join(dir, name)
			w, err := loadWallet(fullpath)
			if err != nil {
				return nil, err
			}

			wallets[name] = w
		}
	}
	return wallets, nil
}

func loadWallet(fn string) (*Wallet, error) {
	rw, err := LoadReadableWallet(fn)
	if err != nil {
		return nil, err
	}

	// Normalize coin types (older wallets used different names for the coin type)
	switch strings.ToLower(rw.Meta[metaCoin]) {
	case "sky", "skycoin":
		rw.Meta[metaCoin] = string(CoinTypeSkycoin)
	case "btc", "bitcoin":
		rw.Meta[metaCoin] = string(CoinTypeBitcoin)
	}

	w, err := rw.ToWallet()
	if err != nil {
		return nil, err
	}

	coinType := w.coin()
	if coinType != CoinTypeSkycoin {
		return nil, fmt.Errorf("LoadWallets only support skycoin wallets, %s is a %s wallet", fn, coinType)
	}

	logger.Infof("Loaded wallet from %s", fn)
	w.setFilename(filepath.Base(fn))

	return w, nil
}

// add add walet to current wallet
func (wlts Wallets) add(w *Wallet) error {
	if _, dup := wlts[w.Filename()]; dup {
		return ErrWalletNameConflict
	}

	wlts[w.Filename()] = w
	return nil
}

// remove wallet of specific id
func (wlts Wallets) remove(id string) {
	delete(wlts, id)
}

// get returns wallet by wallet id
func (wlts Wallets) get(id string) *Wallet {
	return wlts[id]
}

// set sets a wallet into the map
func (wlts Wallets) set(w *Wallet) {
	wlts[w.Filename()] = w.clone()
}

// ToReadable converts Wallets to *ReadableWallet array, sorting them by timestamp
func (wlts Wallets) ToReadable() []*ReadableWallet {
	var rw []*ReadableWallet
	for _, w := range wlts {
		rw = append(rw, NewReadableWallet(w))
	}

	sort.Slice(rw, func(i int, j int) bool {
		a := rw[i].timestamp()
		b := rw[j].timestamp()

		if a == b {
			return rw[i].filename() < rw[j].filename()
		}

		return a < b
	})

	return rw
}

// containsDuplicate returns true if there is a duplicate wallet
// (identified by the first address in the wallet) and return the ID of that wallet
// and the first address if true
func (wlts Wallets) containsDuplicate() (string, cipher.Address, bool) {
	m := make(map[cipher.Address]struct{}, len(wlts))
	for wltID, wlt := range wlts {
		if len(wlt.Entries) == 0 {
			continue
		}
		addr := wlt.Entries[0].SkycoinAddress()
		if _, ok := m[addr]; ok {
			return wltID, addr, true
		}

		m[addr] = struct{}{}
	}

	return "", cipher.Address{}, false
}

// containsEmpty returns true there is an empty wallet and the ID of that wallet if true
func (wlts Wallets) containsEmpty() (string, bool) {
	for wltID, wlt := range wlts {
		if len(wlt.Entries) == 0 {
			return wltID, true
		}
	}
	return "", false
}
