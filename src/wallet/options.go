package wallet

import (
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
)

// Option represents the general options, it can be used to set optional
// parameters while creating a new wallet. Also, can be used to get
// entries service of a wallet.
type Option func(interface{})

// Bip44EntriesOptions represents the options that will be used
// by bip44 to get entries service
type Bip44EntriesOptions struct {
	Account uint32
	Change  uint32
}

// OptionAccount is the option type for specifying a bip44 account
func OptionAccount(index uint32) Option {
	return func(opts interface{}) {
		bip44, ok := opts.(*Bip44EntriesOptions)
		if !ok {
			return
		}
		bip44.Account = index
	}
}

// OptionChange is the option type for specifying a bip44 chain
func OptionChange(change bool) Option {
	return func(opts interface{}) {
		bip44, ok := opts.(*Bip44EntriesOptions)
		if !ok {
			return
		}

		var chain uint32
		if change {
			chain = 1
		}
		bip44.Change = chain
	}
}

func walletOptionFunc(f func(Wallet)) Option {
	return func(v interface{}) {
		w, ok := v.(Wallet)
		if !ok {
			return
		}
		f(w)
	}
}

// OptionCryptoType is the option type for setting wallet crypto type
func OptionCryptoType(cryptoType crypto.CryptoType) Option {
	return walletOptionFunc(func(w Wallet) {
		w.SetCryptoType(cryptoType)
	})
}

// OptionCoinType is the option type for setting  wallet coin  type
func OptionCoinType(coinType CoinType) Option {
	return walletOptionFunc(func(w Wallet) {
		w.SetCoin(coinType)
	})
}

// OptionDecoder is the option type for setting wallet decoder
func OptionDecoder(d Decoder) Option {
	return walletOptionFunc(func(w Wallet) {
		w.SetDecoder(d)
	})
}

// OptionBip44Coin is the option type for setting bip44 coin type for bip44 wallet
func OptionBip44Coin(ct *bip44.CoinType) Option {
	return walletOptionFunc(func(w Wallet) {
		w.SetBip44Coin(*ct)
	})
}

// AdvancedOptions are advanced options that can be used when creating a new wallet
type AdvancedOptions struct {
	DefaultBip44AccountName string
	Encrypt                 bool
	Password                []byte
	GenerateN               uint64
	ScanN                   uint64
	TF                      TransactionsFinder
}

// advancedOptionFunc is a helper function that assert the
// interface in wallet.Option to AdvancedOptions, so that
// the caller can use AdvancedOptions directly.
func advancedOptionFunc(f func(*AdvancedOptions)) Option {
	return func(v interface{}) {
		o, ok := v.(*AdvancedOptions)
		if !ok {
			return
		}
		f(o)
	}
}

// OptionDefaultBip44AccountName can be used to set the bip44 default account name
func OptionDefaultBip44AccountName(name string) Option {
	return advancedOptionFunc(func(opts *AdvancedOptions) {
		opts.DefaultBip44AccountName = name
	})
}

// OptionEncrypt can be used to set whether the wallet should be encrypted when creating a new wallet
func OptionEncrypt(encrypt bool) Option {
	return advancedOptionFunc(func(opts *AdvancedOptions) {
		opts.Encrypt = encrypt
	})
}

// OptionPassword can be used to set the password for encrypting when creating a new wallet.
func OptionPassword(password []byte) Option {
	return advancedOptionFunc(func(opts *AdvancedOptions) {
		opts.Password = password
	})
}

// OptionScanN can be used to set the scanning number when creating a new wallet
func OptionScanN(n uint64) Option {
	return advancedOptionFunc(func(opts *AdvancedOptions) {
		opts.ScanN = n
	})
}

// OptionTransactionsFinder can be used to set the transactions finder when creating a new wallet
func OptionTransactionsFinder(tf TransactionsFinder) Option {
	return advancedOptionFunc(func(opts *AdvancedOptions) {
		opts.TF = tf
	})
}

// OptionGenerateN can be used to set the initial number of addresses to generate
// when creating a new wallet
func OptionGenerateN(n uint64) Option {
	return advancedOptionFunc(func(opts *AdvancedOptions) {
		opts.GenerateN = n
	})
}
