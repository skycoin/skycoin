package visor

import (
	"github.com/amherag/skycoin/src/coin"
	"github.com/amherag/skycoin/src/params"
)

// TransactionIsLocked returns true if the transaction spends locked outputs
func TransactionIsLocked(inUxs coin.UxArray) bool {
	lockedAddrs := params.GetLockedDistributionAddresses()
	lockedAddrsMap := make(map[string]struct{})
	for _, a := range lockedAddrs {
		lockedAddrsMap[a] = struct{}{}
	}

	for _, o := range inUxs {
		uxAddr := o.Body.Address.String()
		if _, ok := lockedAddrsMap[uxAddr]; ok {
			return true
		}
	}

	return false
}
