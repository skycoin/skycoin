package wallet

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
)

// Wallets wallets map
type Wallets map[string]Walleter

// loadWallets Loads all wallets contained in wallet dir.  If any regular file in wallet
// dir fails to load, loading is aborted and error returned.  Only files with
// extension WalletExt are considered.
func loadWallets(dir string) (Wallets, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		logger.WithError(err).WithField("dir", dir).Error("loadWallets: ioutil.ReadDir failed")
		return nil, err
	}

	wallets := Wallets{}
	for _, e := range entries {
		if e.Mode().IsRegular() {
			name := e.Name()
			if !strings.HasSuffix(name, WalletExt) {
				logger.WithField("filename", name).Info("loadWallets: skipping file")
				continue
			}

			fullpath := filepath.Join(dir, name)
			w, err := Load(fullpath)
			if err != nil {
				logger.WithError(err).WithField("filename", fullpath).Error("loadWallets: loadWallet failed")
				return nil, err
			}

			logger.WithField("filename", fullpath).Info("loadWallets: loaded wallet")

			wallets[name] = w
		}
	}

	for name, w := range wallets {
		if err := w.Validate(); err != nil {
			logger.WithError(err).WithField("name", name).Error("loadWallets: wallet.Validate failed")
			return nil, err
		}

		if w.Coin() != CoinTypeSkycoin {
			err := fmt.Errorf("LoadWallets only support skycoin wallets, %s is a %s wallet", name, w.Coin())
			logger.WithError(err).WithField("name", name).Error()
			return nil, err
		}
	}

	return wallets, nil
}

// add add walet to current wallet
func (wlts Wallets) add(w Walleter) error {
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
func (wlts Wallets) get(id string) Walleter {
	return wlts[id]
}

// set sets a wallet into the map
func (wlts Wallets) set(w Walleter) {
	wlts[w.Filename()] = w.Clone()
}

// containsDuplicate returns true if there is a duplicate wallet
// (identified by the first address in the wallet) and return the ID of that wallet
// and the first address if true
func (wlts Wallets) containsDuplicate() (string, cipher.Address, bool) {
	m := make(map[cipher.Address]struct{}, len(wlts))
	for wltID, wlt := range wlts {
		if wlt.EntriesLen() == 0 {
			continue
		}
		e := wlt.GetEntryAt(0)
		addr := e.SkycoinAddress()
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
		if wlt.EntriesLen() == 0 {
			return wltID, true
		}
	}
	return "", false
}
