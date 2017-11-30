package wallet

import (
	"fmt"
	"os"
	"sync"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/blockdb"
)

func errWalletNotExist(wltName string) error {
	return fmt.Errorf("wallet %s doesn't exist", wltName)
}

// BalanceGetter interface for getting the balance of given addresses
type BalanceGetter interface {
	GetBalanceOfAddrs(addrs []cipher.Address) ([]BalancePair, error)
}

// Service wallet service struct
type Service struct {
	sync.RWMutex
	wallets        Wallets
	options        []Option
	firstAddrIDMap map[string]string // key: first address in wallet, value: wallet id

	WalletDirectory string
}

// NewService new wallet service
func NewService(walletDir string, options ...Option) (*Service, error) {
	serv := &Service{
		firstAddrIDMap: make(map[string]string),
	}
	if err := os.MkdirAll(walletDir, os.FileMode(0700)); err != nil {
		return nil, fmt.Errorf("failed to create wallet directory %s: %v", walletDir, err)
	}

	serv.WalletDirectory = walletDir
	for i := range options {
		serv.options = append(serv.options, options[i])
	}

	w, err := LoadWallets(serv.WalletDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to load all wallets: %v", err)
	}

	serv.wallets = serv.removeDup(w)

	if len(serv.wallets) == 0 {
		wltName := NewWalletFilename()
		// create default wallet
		w, err := serv.CreateWallet(wltName, OptLabel("Your Wallet"))
		if err != nil {
			return nil, err
		}

		if err := w.Save(serv.WalletDirectory); err != nil {
			return nil, fmt.Errorf("failed to save wallets to %s: %v", serv.WalletDirectory, err)
		}
	}

	return serv, nil
}

// CreateWallet creates wallet
func (serv *Service) CreateWallet(wltName string, options ...Option) (Wallet, error) {
	ops := make([]Option, 0, len(serv.options)+len(options))
	ops = append(ops, serv.options...)
	ops = append(ops, options...)
	w, err := NewWallet(wltName, ops...)
	if err != nil {
		return Wallet{}, err
	}

	// generate a default address
	w.GenerateAddresses(1)

	serv.Lock()
	defer serv.Unlock()
	// check dup
	if id, ok := serv.firstAddrIDMap[w.Entries[0].Address.String()]; ok {
		return Wallet{}, fmt.Errorf("duplicate wallet with %v", id)
	}

	if err := serv.wallets.Add(*w); err != nil {
		return Wallet{}, err
	}

	if err := w.Save(serv.WalletDirectory); err != nil {
		// remove the added wallet from serv.wallets.
		serv.wallets.Remove(w.GetID())
		return Wallet{}, err
	}

	serv.firstAddrIDMap[w.Entries[0].Address.String()] = w.GetID()

	return *w, nil
}

// LoadAndScanWallet loads wallet from seed and scan the first N address
func (serv *Service) LoadAndScanWallet(wltName string, seed string, scanN uint64, bg BalanceGetter, options ...Option) (Wallet, error) {
	ops := make([]Option, 0, len(serv.options)+len(options))
	ops = append(ops, serv.options...)
	ops = append(ops, options...)
	ops = append(ops, OptSeed(seed))
	w, err := NewWallet(wltName, ops...)
	if err != nil {
		return Wallet{}, err
	}

	// generate a default address
	w.GenerateAddresses(1)

	serv.Lock()
	defer serv.Unlock()
	// check dup
	if id, ok := serv.firstAddrIDMap[w.Entries[0].Address.String()]; ok {
		return Wallet{}, fmt.Errorf("duplicate wallet with %v", id)
	}

	// generate the remaining addresses that are need to scan
	w.GenerateAddresses(scanN - 1)

	// check balance from the last one till we find the
	// address that has coins
	addrs := w.GetAddresses()
	bals, err := bg.GetBalanceOfAddrs(addrs)
	if err != nil {
		return Wallet{}, err
	}

	var keepNum uint64 = 1
	for i := len(bals) - 1; i >= 0; i-- {
		if bals[i].Confirmed.Coins > 0 || bals[i].Predicted.Coins > 0 {
			keepNum = uint64(i + 1)
			break
		}
	}

	// reset the wallet if scan number > 1 and not equal to the keep number
	if scanN > 1 && keepNum != scanN {
		w.Reset()
		w.GenerateAddresses(uint64(keepNum))
	}

	if err := serv.wallets.Add(*w); err != nil {
		return Wallet{}, err
	}

	if err := w.Save(serv.WalletDirectory); err != nil {
		// remove the added wallet from serv.wallets.
		serv.wallets.Remove(w.GetID())
		return Wallet{}, err
	}

	serv.firstAddrIDMap[w.Entries[0].Address.String()] = w.GetID()

	return *w, nil
}

// NewAddresses generate address entries in given wallet,
// return nil if wallet does not exist.
func (serv *Service) NewAddresses(wltID string, num uint64) ([]cipher.Address, error) {
	serv.Lock()
	defer serv.Unlock()
	w, ok := serv.wallets.Get(wltID)
	if !ok {
		return []cipher.Address{}, errWalletNotExist(wltID)
	}

	addrs := w.GenerateAddresses(num)
	if err := w.Save(serv.WalletDirectory); err != nil {
		return []cipher.Address{}, err
	}

	return addrs, nil
}

// GetAddresses returns all addresses in given wallet
func (serv *Service) GetAddresses(wltID string) ([]cipher.Address, error) {
	serv.RLock()
	defer serv.RUnlock()
	w, ok := serv.wallets.Get(wltID)
	if !ok {
		return []cipher.Address{}, errWalletNotExist(wltID)
	}

	return w.GetAddresses(), nil
}

// GetWallet returns wallet by id
func (serv *Service) GetWallet(wltID string) (Wallet, bool) {
	serv.RLock()
	defer serv.RUnlock()
	w, ok := serv.wallets.Get(wltID)
	if !ok {
		return Wallet{}, false
	}
	return w.Copy(), true
}

// GetWallets returns all wallet
func (serv *Service) GetWallets() Wallets {
	wlts := make(Wallets, len(serv.wallets))
	for k, w := range serv.wallets {
		nw := w.Copy()
		wlts[k] = &nw
	}
	return wlts
}

// ReloadWallets reload wallets
func (serv *Service) ReloadWallets() error {
	serv.Lock()
	defer serv.Unlock()
	wallets, err := LoadWallets(serv.WalletDirectory)
	if err != nil {
		return err
	}

	serv.firstAddrIDMap = make(map[string]string)
	serv.wallets = serv.removeDup(wallets)
	return nil
}

// GetWalletsReadable returns readable wallets
func (serv *Service) GetWalletsReadable() []*ReadableWallet {
	serv.RLock()
	defer serv.RUnlock()
	return serv.wallets.ToReadable()
}

// CreateAndSignTransaction creates and sign transaction from wallet
func (serv *Service) CreateAndSignTransaction(wltID string,
	vld Validator,
	unspent blockdb.UnspentGetter,
	headTime uint64,
	amt Balance,
	dest cipher.Address) (*coin.Transaction, error) {
	serv.RLock()
	defer serv.RUnlock()
	w, ok := serv.wallets.Get(wltID)
	if !ok {
		return nil, errWalletNotExist(wltID)
	}

	return w.CreateAndSignTransaction(vld, unspent, headTime, amt, dest)
}

// UpdateWalletLabel updates the wallet label
func (serv *Service) UpdateWalletLabel(wltID, label string) error {
	serv.Lock()
	defer serv.Unlock()
	var wlt Wallet
	if err := serv.wallets.Update(wltID, func(w Wallet) Wallet {
		w.SetLabel(label)
		wlt = w
		return w
	}); err != nil {
		return err
	}

	return wlt.Save(serv.WalletDirectory)
}

func (serv *Service) removeDup(wlts Wallets) Wallets {
	var rmWltIDS []string
	// remove dup wallets
	for wltID, wlt := range wlts {
		if len(wlt.Entries) == 0 {
			// empty wallet
			rmWltIDS = append(rmWltIDS, wltID)
			continue
		}

		addr := wlt.Entries[0].Address.String()
		id, ok := serv.firstAddrIDMap[addr]
		if ok {
			// check whose entries number is bigger
			pw, _ := wlts.Get(id)
			if len(pw.Entries) >= len(wlt.Entries) {
				rmWltIDS = append(rmWltIDS, wltID)
				continue
			}

			// replace the old wallet with the new one
			// records the wallet id that need to remove
			rmWltIDS = append(rmWltIDS, id)
			// update wallet id
			serv.firstAddrIDMap[addr] = wltID
			continue
		}

		serv.firstAddrIDMap[addr] = wltID
	}

	// remove the duplicate and empty wallet
	for _, id := range rmWltIDS {
		wlts.Remove(id)
	}

	return wlts
}
