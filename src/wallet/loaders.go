package wallet

import (
	"fmt"
	"sync"
)

var walletLoaders loaders

// RegisterLoader registers a wallet Loader
func RegisterLoader(walletType string, l Loader) error {
	return walletLoaders.add(walletType, l)
}

func getLoader(walletType string) (Loader, bool) {
	return walletLoaders.get(walletType)
}

// Loader is the interface that wraps the Load method.
//
// Load loads wallet from data bytes
type Loader interface {
	Load(data []byte) (Wallet, error)
}

type loaders struct {
	sync.Mutex
	wls map[string]Loader
}

// Add add a new wallet type Loader
func (ls *loaders) add(walletType string, l Loader) error {
	ls.Lock()
	defer ls.Unlock()
	if ls.wls == nil {
		ls.wls = map[string]Loader{}
	}

	if _, ok := ls.wls[walletType]; ok {
		return fmt.Errorf("wallet loader for %s already exists", walletType)
	}

	ls.wls[walletType] = l
	return nil
}

// Get returns the wallet Loader base on wallet type
func (ls *loaders) get(walletType string) (Loader, bool) {
	ls.Lock()
	defer ls.Unlock()
	c, ok := ls.wls[walletType]
	return c, ok
}
