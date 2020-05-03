package wallet

import (
	"fmt"

	"github.com/SkycoinProject/skycoin/src/cipher"
)

var registeredAddressSecKeyDecoders = initAddressSecKeyDecoders()

type addressSecKeyDecoders struct {
	adapters map[CoinType]AddressSecKeyDecoder
}

func initAddressSecKeyDecoders() addressSecKeyDecoders {
	return addressSecKeyDecoders{
		adapters: map[CoinType]AddressSecKeyDecoder{
			CoinTypeSkycoin: skycoinDecoder{},
			CoinTypeBitcoin: bitcoinDecoder{},
		},
	}
}

func (a addressSecKeyDecoders) get(coinType CoinType) AddressSecKeyDecoder {
	adpt, ok := a.adapters[coinType]
	if !ok {
		return skycoinDecoder{}
	}
	return adpt
}

func (a addressSecKeyDecoders) add(coinType CoinType, ca AddressSecKeyDecoder) error {
	if _, ok := a.adapters[coinType]; ok {
		return fmt.Errorf("coin adapter for %s already registered", coinType)
	}
	a.adapters[coinType] = ca
	return nil
}

// RegisterCoinAdapter registers a new adapter
func RegisterCoinAdapter(coinType CoinType, ca AddressSecKeyDecoder) error {
	return registeredAddressSecKeyDecoders.add(coinType, ca)
}

// ResolveCoinAdapter returns a coin adapter by coin type,
// if the corresponding adapter of given coin type is not found, returns
// the skycoin adapter.
func ResolveCoinAdapter(coinType CoinType) AddressSecKeyDecoder {
	return registeredAddressSecKeyDecoders.get(coinType)
}

type AddressDecoder interface {
	DecodeBase58Address(addr string) (cipher.Addresser, error)
	AddressFromPubKey(key cipher.PubKey) cipher.Addresser
}

type SecKeyDecoder interface {
	SecKeyToHex(key cipher.SecKey) string
	SecKeyFromHex(key string) (cipher.SecKey, error)
}

type AddressSecKeyDecoder interface {
	AddressDecoder
	SecKeyDecoder
}

type skycoinDecoder struct{}

func (s skycoinDecoder) AddressFromPubKey(key cipher.PubKey) cipher.Addresser {
	return cipher.AddressFromPubKey(key)
}

func (s skycoinDecoder) DecodeBase58Address(addr string) (cipher.Addresser, error) {
	return cipher.DecodeBase58Address(addr)
}

func (s skycoinDecoder) SecKeyToHex(secKey cipher.SecKey) string {
	return secKey.Hex()
}

func (s skycoinDecoder) SecKeyFromHex(secKey string) (cipher.SecKey, error) {
	return cipher.SecKeyFromHex(secKey)
}

type bitcoinDecoder struct{}

func (b bitcoinDecoder) AddressFromPubKey(key cipher.PubKey) cipher.Addresser {
	return cipher.BitcoinAddressFromPubKey(key)
}

func (s bitcoinDecoder) DecodeBase58Address(addr string) (cipher.Addresser, error) {
	return cipher.DecodeBase58BitcoinAddress(addr)
}

func (s bitcoinDecoder) SecKeyToHex(secKey cipher.SecKey) string {
	return cipher.BitcoinWalletImportFormatFromSeckey(secKey)
}

func (s bitcoinDecoder) SecKeyFromHex(secKey string) (cipher.SecKey, error) {
	return cipher.SecKeyFromBitcoinWalletImportFormat(secKey)
}
