// Package transaction src/transaction/distribution.go c1-coin-chain
package transaction

import (
	"fmt"

	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
)

// TransactionIsLocked returns true if the transaction spends locked outputs
func TransactionIsLocked(d params.Distribution, inUxs coin.UxArray) (bool, error) { //nolint:revive
	lockedAddrs, err := d.LockedAddresses()
	if err != nil {
		return false, fmt.Errorf("locked distribution addresses: %w", err)
	}
	lockedAddrsMap := make(map[string]struct{})
	for _, a := range lockedAddrs {
		lockedAddrsMap[a] = struct{}{}
	}

	for _, o := range inUxs {
		uxAddr := o.Body.Address.String()
		if _, ok := lockedAddrsMap[uxAddr]; ok {
			return true, nil
		}
	}

	return false, nil
}
