package wallet

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Wallets wallets map
type Wallets map[string]Wallet

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
func (wlts Wallets) add(w Wallet) error {
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
func (wlts Wallets) get(id string) Wallet {
	return wlts[id]
}

// set sets a wallet into the map
func (wlts Wallets) set(w Wallet) {
	wlts[w.Filename()] = w.Clone()
}

// containsDuplicate returns true if there is a duplicate wallet identified by
// the wallet's fingerprint. This is to detect duplicate generative wallets;
// wallets with no defined generation method do not have a concept of being
// a duplicate of another wallet
func (wlts Wallets) containsDuplicate() (string, string, bool) {
	m := make(map[string]struct{}, len(wlts))
	for wltID, wlt := range wlts {
		fp := wlt.Fingerprint()
		if fp == "" {
			continue
		}

		if _, ok := m[fp]; ok {
			return wltID, fp, true
		}

		m[fp] = struct{}{}
	}

	return "", "", false
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
