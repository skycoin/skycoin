package deterministic

import (
	"github.com/SkycoinProject/skycoin/src/wallet"
)

func walletOptionFunc(f func(*Wallet)) wallet.Option {
	return func(v interface{}) {
		w, ok := v.(*Wallet)
		if !ok {
			return
		}
		f(w)
	}
}
