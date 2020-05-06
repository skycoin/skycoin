package bip44wallet

import (
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
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

func moreOptionFunc(f func(*moreOptions)) wallet.Option {
	return func(v interface{}) {
		o, ok := v.(*moreOptions)
		if !ok {
			return
		}
		f(o)
	}
}

func Version(version string) wallet.Option {
	return walletOptionFunc(func(w *Wallet) {
		w.Meta[wallet.MetaVersion] = version
	})
}

// CryptoType is an option to set the wallet crypto type
func CryptoType(cryptoType crypto.CryptoType) wallet.Option {
	return walletOptionFunc(func(w *Wallet) {
		w.Meta[wallet.MetaCryptoType] = string(cryptoType)
	})
}

func CoinType(coinType wallet.CoinType) wallet.Option {
	return walletOptionFunc(func(w *Wallet) {
		w.Meta[wallet.MetaCoin] = string(coinType)
	})
}

func Bip44CoinType(bip44Coin bip44.CoinType) wallet.Option {
	return walletOptionFunc(func(w *Wallet) {
		w.Meta.SetBip44Coin(bip44Coin)
	})
}

func Decoder(d wallet.Decoder) wallet.Option {
	return walletOptionFunc(func(w *Wallet) {
		w.decoder = d
	})
}

type moreOptions struct {
	Encrypt   bool
	Password  []byte
	GenerateN uint64
	ScanN     uint64
	TF        wallet.TransactionsFinder
}

func Encrypt(encrypt bool) wallet.Option {
	return moreOptionFunc(func(opts *moreOptions) {
		opts.Encrypt = encrypt
	})
}

func Password(password []byte) wallet.Option {
	return moreOptionFunc(func(opts *moreOptions) {
		opts.Password = password
	})
}

func ScanN(n uint64) wallet.Option {
	return moreOptionFunc(func(opts *moreOptions) {
		opts.ScanN = n
	})
}

func TransactionsFinder(tf wallet.TransactionsFinder) wallet.Option {
	return moreOptionFunc(func(opts *moreOptions) {
		opts.TF = tf
	})
}

func GenerateN(n uint64) wallet.Option {
	return moreOptionFunc(func(opts *moreOptions) {
		opts.GenerateN = n
	})
}
