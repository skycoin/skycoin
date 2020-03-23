package wallet

import (
	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
)

var registeredCoinAdapters = initCoinAdapters()

type coinAdapters struct {
	adapters map[CoinType]coinAdapter
}

func initCoinAdapters() coinAdapters {
	return coinAdapters{
		adapters: map[CoinType]coinAdapter{
			CoinTypeSkycoin: skycoinAdapter{},
			CoinTypeBitcoin: bitcoinAdapter{},
		},
	}
}

func (a coinAdapters) get(coinType CoinType) coinAdapter {
	adpt, ok := a.adapters[coinType]
	if !ok {
		// if no adapter found, returns the default skycoin adapter
		return skycoinAdapter{}
	}
	return adpt
}

// RegisterCoinAdapter registers a new adapter
func RegisterCoinAdapter(coinType CoinType, ca coinAdapter) {
	//
}

func resolveCoinAdapter(coinType CoinType) coinAdapter {
	return registeredCoinAdapters.get(coinType)
}

type coinAdapter interface {
	Bip44CoinType() bip44.CoinType
	AddressFromPubKey(key cipher.PubKey) cipher.Addresser
}

type skycoinAdapter struct{}

func (s skycoinAdapter) Bip44CoinType() bip44.CoinType {
	return bip44.CoinTypeSkycoin
}

func (s skycoinAdapter) AddressFromPubKey(key cipher.PubKey) cipher.Addresser {
	return cipher.AddressFromPubKey(key)
}

type bitcoinAdapter struct{}

func (b bitcoinAdapter) Bip44CoinType() bip44.CoinType {
	return bip44.CoinTypeBitcoin
}

func (b bitcoinAdapter) AddressFromPubKey(key cipher.PubKey) cipher.Addresser {
	return cipher.BitcoinAddressFromPubKey(key)
}
