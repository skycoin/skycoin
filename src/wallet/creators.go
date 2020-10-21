package wallet

import (
	"fmt"
	"sync"
)

var walletCreators creators

// RegisterCreator register wallet creator
func RegisterCreator(walletType string, c Creator) error {
	return walletCreators.add(walletType, c)
}

func getCreator(walletType string) (Creator, bool) {
	return walletCreators.get(walletType)
}

// Creator is the interface that wraps the Create and Type methods.
//
// Create creates a wallet base on the parameters
type Creator interface {
	Create(filename, label, seed string, options Options) (Wallet, error)
}

type creators struct {
	l   sync.Mutex
	wcs map[string]Creator
}

// Add add a new creator to the wallet creator list
func (cs *creators) add(walletType string, c Creator) error {
	cs.l.Lock()
	defer cs.l.Unlock()
	if cs.wcs == nil {
		cs.wcs = map[string]Creator{}
	}

	if _, ok := cs.wcs[walletType]; ok {
		return fmt.Errorf("wallet creator for %s already exists", walletType)
	}

	cs.wcs[walletType] = c
	return nil
}

// Get returns the wallet creator base on the type
func (cs *creators) get(walletType string) (Creator, bool) {
	cs.l.Lock()
	defer cs.l.Unlock()
	c, ok := cs.wcs[walletType]
	return c, ok
}
