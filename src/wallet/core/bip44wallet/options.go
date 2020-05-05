package bip44wallet

import (
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
)

func optionFunc(f func(*Wallet)) wallet.Option {
	return func(v interface{}) {
		f(v.(*Wallet))
	}
}

func Version(version string) wallet.Option {
	return optionFunc(func(w *Wallet) {
		w.Meta[wallet.MetaVersion] = version
	})
}

// CryptoType is an option to set the wallet crypto type
func CryptoType(cryptoType crypto.CryptoType) wallet.Option {
	return optionFunc(func(w *Wallet) {
		w.Meta[wallet.MetaCryptoType] = string(cryptoType)
	})
}

func CoinType(coinType wallet.CoinType) wallet.Option {
	return optionFunc(func(w *Wallet) {
		w.Meta[wallet.MetaCoin] = string(coinType)
	})
}

func Bip44CoinType(bip44Coin bip44.CoinType) wallet.Option {
	return optionFunc(func(w *Wallet) {
		w.Meta.SetBip44Coin(bip44Coin)
	})
}

func Decoder(d wallet.Decoder) wallet.Option {
	return optionFunc(func(w *Wallet) {
		w.decoder = d
	})
}
