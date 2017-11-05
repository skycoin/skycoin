package wallet

import (
	//"fmt"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
)

// Wallets wallets map
type Wallets map[string]*Wallet

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
	for i, e := range entries {
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
			// check the wallet version
			if w.GetVersion() != version {
				logger.Info("Update wallet %v", fullpath)
				bkFile := filepath.Join(bkpath, w.GetFilename())
				if err := backupWltFile(fullpath, bkFile); err != nil {
					return nil, err
				}

				// update wallet to new version.
				tm := time.Now().Unix() + int64(i)
				mustUpdateWallet(&w, dir, tm)
			}

			wallets[name] = &w
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

func mustUpdateWallet(wlt *Wallet, dir string, tm int64) {
	// update version meta data.
	wlt.Meta["version"] = version

	// update lastSeed meta data.
	lsd, seckeys := cipher.GenerateDeterministicKeyPairsSeed([]byte(wlt.Meta["seed"]), 1)
	if seckeys[0] != wlt.Entries[0].Secret {
		logger.Panic("update wallet failed, seckey not match")
	}

	wlt.Meta["lastSeed"] = hex.EncodeToString(lsd)

	// update tm meta data.
	wlt.Meta["tm"] = fmt.Sprintf("%v", tm)
	if err := wlt.Save(dir); err != nil {
		logger.Panic(err)
	}
}

// Add add walet to current wallet
func (wlts Wallets) Add(w Wallet) error {
	if _, dup := wlts[w.GetFilename()]; dup {
		return errors.New("wallet name would conflict with existing wallet, renaming")
	}

	wlts[w.GetFilename()] = &w
	return nil
}

// Remove wallet of specific id
func (wlts Wallets) Remove(id string) {
	delete(wlts, id)
}

// Get returns wallet by wallet id
func (wlts Wallets) Get(wltID string) (*Wallet, bool) {
	if w, ok := wlts[wltID]; ok {
		return w, true
	}
	return &Wallet{}, false
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
func (wlts *Wallets) NewAddresses(wltID string, num int) ([]cipher.Address, error) {
	if w, ok := (*wlts)[wltID]; ok {
		return w.GenerateAddresses(num), nil
	}
	return nil, fmt.Errorf("wallet: %v does not exist", wltID)
}

// Save check for name conflicts!
// resolve conflicts for saving wallets who have different names
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

func (wlts Wallets) toReadable(f ReadableWalletCtor) []*ReadableWallet {
	var rw []*ReadableWallet
	for _, w := range wlts {
		rw = append(rw, f(*w))
	}
	sort.Sort(ByTm(rw))
	return rw
}

// ToReadable converts Wallets to *ReadableWallet array
func (wlts Wallets) ToReadable() []*ReadableWallet {
	return wlts.toReadable(NewReadableWallet)
}
