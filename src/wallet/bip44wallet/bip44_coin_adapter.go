package bip44wallet

import (
	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
	"github.com/SkycoinProject/skycoin/src/wallet/meta"
)

var registeredCoinAdapters = initCoinAdapters()

type coinAdapters struct {
	adapters map[meta.CoinType]coinAdapter
}

func initCoinAdapters() coinAdapters {
	return coinAdapters{
		adapters: map[meta.CoinType]coinAdapter{
			meta.CoinTypeSkycoin: skycoinAdapter{},
			meta.CoinTypeBitcoin: bitcoinAdapter{},
		},
	}
}

func (a coinAdapters) get(coinType meta.CoinType) coinAdapter {
	adpt, ok := a.adapters[coinType]
	if !ok {
		// if no adapter found, returns the default skycoin adapter
		return skycoinAdapter{}
	}
	return adpt
}

// RegisterCoinAdapter registers a new adapter
func RegisterCoinAdapter(coinType meta.CoinType, ca coinAdapter) {
	//
}

func resolveCoinAdapter(coinType meta.CoinType) coinAdapter {
	return registeredCoinAdapters.get(coinType)
}

type coinAdapter interface {
	Bip44CoinType() bip44.CoinType
	AddressFromPubKey(key cipher.PubKey) cipher.Addresser
	DecodeBase58Address(addr string) (cipher.Addresser, error)
	SecKeyToHex(secKey cipher.SecKey) string
	SecKeyFromHex(secKey string) (cipher.SecKey, error)
}

type skycoinAdapter struct{}

func (s skycoinAdapter) Bip44CoinType() bip44.CoinType {
	return bip44.CoinTypeSkycoin
}

func (s skycoinAdapter) AddressFromPubKey(key cipher.PubKey) cipher.Addresser {
	return cipher.AddressFromPubKey(key)
}

func (s skycoinAdapter) DecodeBase58Address(addr string) (cipher.Addresser, error) {
	return cipher.DecodeBase58Address(addr)
}

func (s skycoinAdapter) SecKeyToHex(secKey cipher.SecKey) string {
	return secKey.Hex()
}

func (s skycoinAdapter) SecKeyFromHex(secKey string) (cipher.SecKey, error) {
	return cipher.SecKeyFromHex(secKey)
}

type bitcoinAdapter struct{}

func (b bitcoinAdapter) Bip44CoinType() bip44.CoinType {
	return bip44.CoinTypeBitcoin
}

func (b bitcoinAdapter) AddressFromPubKey(key cipher.PubKey) cipher.Addresser {
	return cipher.BitcoinAddressFromPubKey(key)
}

func (s bitcoinAdapter) DecodeBase58Address(addr string) (cipher.Addresser, error) {
	return cipher.DecodeBase58BitcoinAddress(addr)
}

func (s bitcoinAdapter) SecKeyToHex(secKey cipher.SecKey) string {
	return cipher.BitcoinWalletImportFormatFromSeckey(secKey)
}

func (s bitcoinAdapter) SecKeyFromHex(secKey string) (cipher.SecKey, error) {
	return cipher.SecKeyFromBitcoinWalletImportFormat(secKey)
}
