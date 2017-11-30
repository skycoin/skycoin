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

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
)

// Wallets wallets map
type Wallets map[string]*Wallet

var (
	// ErrWalletNameConflict represents the wallet name conflict error
	ErrWalletNameConflict = errors.New("wallet name would conflict with existing wallet, renaming")
)

// LoadWallets Loads all wallets contained in wallet dir.  If any regular file in wallet
// dir fails to load, loading is aborted and error returned.  Only files with
// extension WalletExt are considered. If encounter old wallet file, then backup
// the wallet file into dir/backup/
func LoadWallets(dir string) (Wallets, error) {
	// TODO -- don't load duplicate wallets.
	// TODO -- save a last_modified value in wallets to decide which to load
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// create backup dir if not exist
	bkpath := dir + "/backup/"
	if _, err := os.Stat(bkpath); os.IsNotExist(err) {
		// create the backup dir
		logger.Critical("create wallet backup dir, %v", bkpath)
		if err := os.Mkdir(bkpath, 0777); err != nil {
			return nil, err
		}
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
			w, err := rw.toWallet()
			if err != nil {
				return nil, err
			}
			logger.Info("Loaded wallet from %s", fullpath)
			w.setFilename(name)
			wallets[name] = w
		}
	}
	return wallets, nil
}

func backupWltFile(src, dst string) error {
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

// func mustUpdateWallet(wlt *Wallet, dir string, tm int64) {
// 	// update version meta data.
// 	wlt.Meta["version"] = version

// 	// update lastSeed meta data.
// 	lsd, seckeys := cipher.GenerateDeterministicKeyPairsSeed([]byte(wlt.Meta["seed"]), 1)
// 	if seckeys[0] != wlt.Entries[0].Secret {
// 		logger.Panic("update wallet failed, seckey not match")
// 	}

// 	wlt.Meta["lastSeed"] = hex.EncodeToString(lsd)

// 	// update tm meta data.
// 	wlt.Meta["tm"] = fmt.Sprintf("%v", tm)
// 	if err := wlt.Save(dir); err != nil {
// 		logger.Panic(err)
// 	}
// }

// Add adds wallet to current wallet
func (wlts Wallets) Add(w Wallet) error {
	if _, dup := wlts[w.Filename()]; dup {
		return errors.New("wallet name would conflict with existing wallet, renaming")
	}

	wlts[w.Filename()] = &w
	return nil
}

// Remove wallet of specific id
func (wlts Wallets) Remove(id string) {
	delete(wlts, id)
}

// Get returns wallet by wallet id
func (wlts Wallets) Get(id string) (*Wallet, bool) {
	if w, ok := wlts[id]; ok {
		return w, true
	}
	return &Wallet{}, false
}

// set sets a wallet into the map
func (wlts Wallets) set(w Wallet) {
	wlts[w.GetFilename()] = &w
}

// Update updates the given wallet, return error if not exist
func (wlts Wallets) Update(wltID string, updateFunc func(Wallet) Wallet) error {
	w, ok := wlts[wltID]
	if !ok {
		return errWalletNotExist(wltID)
	}

	newWlt := updateFunc(*w)
	wlts[wltID] = &newWlt
	return nil
}

// NewAddresses creates num addresses in given wallet
func (wlts *Wallets) NewAddresses(id string, num int, password string) ([]cipher.Address, error) {
	if w, ok := (*wlts)[id]; ok {
		return w.GenerateAddresses(password, num)
	}
	return nil, fmt.Errorf("wallet: %v does not exist", id)
}

// Save check for name conflicts!
// resolve conflicts for saving wallets who have different names
func (wlts Wallets) Save(dir string) map[string]error {
	errs := make(map[string]error)
	for id, w := range wlts {
		if err := Save(dir, w); err != nil {
			errs[id] = err
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// ToReadable converts Wallets to *ReadableWallet array
func (wlts Wallets) ToReadable() []*ReadableWallet {
	var rw []*ReadableWallet
	for _, w := range wlts {
		rw = append(rw, NewReadableWallet(*w))
	}
	sort.Sort(ByTm(rw))
	return rw
}

// Update updates the given wallet, return error if not exist
func (wlts Wallets) update(id string, updateFunc func(Wallet) Wallet) error {
	w, ok := wlts[id]
	if !ok {
		return errWalletNotExist(id)
	}

	newWlt := updateFunc(*w)
	wlts[id] = &newWlt
	return nil
}
