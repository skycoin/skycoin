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

// RegisterAddressSecKeyDecoder registers an address and seckey decoder
func RegisterAddressSecKeyDecoder(coinType CoinType, d AddressSecKeyDecoder) error {
	return registeredAddressSecKeyDecoders.add(coinType, d)
}

// ResolveAddressSecKeyDecoder returns an address and seckey decoder by coin type,
// if the corresponding decoder of is not found, returns
// the skycoin decoder.
func ResolveAddressSecKeyDecoder(coinType CoinType) AddressSecKeyDecoder {
	return registeredAddressSecKeyDecoders.get(coinType)
}

// ResolveAddressDecoder returns an address decoder by coin type.
func ResolveAddressDecoder(coinType CoinType) AddressDecoder {
	return registeredAddressSecKeyDecoders.get(coinType)
}

// ResolveSecKeyDecoder returns a SecKey decoder
func ResolveSecKeyDecoder(coinType CoinType) SecKeyDecoder {
	return registeredAddressSecKeyDecoders.get(coinType)
}

// AddressDecoder interface that wraps methods for encoding/decoding cipher.Addresser
type AddressDecoder interface {
	DecodeBase58Address(addr string) (cipher.Addresser, error)
	AddressFromPubKey(key cipher.PubKey) cipher.Addresser
}

// SecKeyDecoder interface that wraps the methods for encoding/decoding cipher.SecKey
type SecKeyDecoder interface {
	SecKeyToHex(key cipher.SecKey) string
	SecKeyFromHex(key string) (cipher.SecKey, error)
}

// AddressSecKeyDecoder interface that embedes the AddressDecoder and SecKeyDecoder
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

func (b bitcoinDecoder) DecodeBase58Address(addr string) (cipher.Addresser, error) {
	return cipher.DecodeBase58BitcoinAddress(addr)
}

func (b bitcoinDecoder) SecKeyToHex(secKey cipher.SecKey) string {
	return cipher.BitcoinWalletImportFormatFromSeckey(secKey)
}

func (b bitcoinDecoder) SecKeyFromHex(secKey string) (cipher.SecKey, error) {
	return cipher.SecKeyFromBitcoinWalletImportFormat(secKey)
}
