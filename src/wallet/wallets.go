package wallet

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/skycoin/skycoin/src/util/file"
)

// Wallets wallets map
type Wallets map[string]*Wallet

// LoadWallets Loads all wallets contained in wallet dir.  If any regular file in wallet
// dir fails to load, loading is aborted and error returned.  Only files with
// extension WalletExt are considered.
func LoadWallets(dir string) (Wallets, error) {
	// TODO -- don't load duplicate wallets.
	// TODO -- save a last_modified value in wallets to decide which to load
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
			rw, err := LoadReadableWallet(fullpath)
			if err != nil {
				return nil, err
			}
			w, err := rw.ToWallet()
			if err != nil {
				return nil, err
			}
			logger.Infof("Loaded wallet from %s", fullpath)
			w.setFilename(name)

			if isLoaded, fileName := isWalletLoaded(wallets, w); isLoaded {
				err = fmt.Errorf("Duplicate Walletfiles: '" + fileName + "' && '" + name + "'")
				return nil, err
			} else {
				wallets[name] = w
			}
		}
	}
	return wallets, nil
}

func isWalletLoaded(wlts Wallets, wlt *Wallet) (bool, string) {
	logger.Info("Checking if Wallet is already loaded...")
	for _, wltItem := range wlts {

		if wltItem.IsEncrypted() && wlt.IsEncrypted() {
			if wltItem.secrets() == wlt.secrets() {
				return true, wltItem.Filename()
			}
		} else if !wltItem.IsEncrypted() && !wlt.IsEncrypted() {
			if wltItem.seed() == wlt.seed() {
				return true, wltItem.Filename()
			}
		}
	}

	return false, ""
}

func backupWltFile(src, dst string) error { // nolint: deadcode,unused,megacheck
	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("%v file already exist", dst)
	}

	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	n, err := file.CopyFile(dst, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	// check if the content bytes are equal.
	if n != int64(len(b)) {
		return errors.New("copy file failed")
	}
	return nil
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
func (wlts Wallets) get(id string) (*Wallet, bool) {
	if w, ok := wlts[id]; ok {
		return w, true
	}
	return nil, false
}

// set sets a wallet into the map
func (wlts Wallets) set(w *Wallet) {
	wlts[w.Filename()] = w.clone()
}

// ToReadable converts Wallets to *ReadableWallet array
func (wlts Wallets) ToReadable() []*ReadableWallet {
	var rw []*ReadableWallet
	for _, w := range wlts {
		rw = append(rw, NewReadableWallet(w))
	}

	sort.Slice(rw, func(i int, j int) bool {
		return rw[i].time() < rw[j].time()
	})
	return rw
}
