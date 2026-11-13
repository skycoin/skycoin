// Package wallet src/wallet/decoders.go c1-coin-wallet
package wallet

import (
	"fmt"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/eth"
)

var registeredAddressSecKeyDecoders = initAddressSecKeyDecoders()

type addressSecKeyDecoders struct {
	adapters map[CoinType]AddressSecKeyDecoder
}

func initAddressSecKeyDecoders() addressSecKeyDecoders {
	return addressSecKeyDecoders{
		adapters: map[CoinType]AddressSecKeyDecoder{
			CoinTypeSkycoin:       skycoinDecoder{},
			CoinTypeBitcoin:       bitcoinDecoder{},
			CoinTypeBitcoinSegwit: bitcoinSegwitDecoder{},
			CoinTypeEthereum:      ethereumDecoder{},
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
	AddressFromPubKey(key cipher.PubKey) (cipher.Addresser, error)
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

func (s skycoinDecoder) AddressFromPubKey(key cipher.PubKey) (cipher.Addresser, error) {
	return cipher.AddressFromPubKey(key), nil
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

func (b bitcoinDecoder) AddressFromPubKey(key cipher.PubKey) (cipher.Addresser, error) {
	return cipher.BitcoinAddressFromPubKey(key), nil
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

type bitcoinSegwitDecoder struct{}

func (b bitcoinSegwitDecoder) AddressFromPubKey(key cipher.PubKey) (cipher.Addresser, error) {
	return cipher.BitcoinSegwitAddressFromPubKey(key), nil
}

func (b bitcoinSegwitDecoder) DecodeBase58Address(addr string) (cipher.Addresser, error) {
	return cipher.DecodeBech32BitcoinAddress(addr)
}

func (b bitcoinSegwitDecoder) SecKeyToHex(secKey cipher.SecKey) string {
	return cipher.BitcoinWalletImportFormatFromSeckey(secKey)
}

func (b bitcoinSegwitDecoder) SecKeyFromHex(secKey string) (cipher.SecKey, error) {
	return cipher.SecKeyFromBitcoinWalletImportFormat(secKey)
}

type ethereumDecoder struct{}

func (e ethereumDecoder) AddressFromPubKey(key cipher.PubKey) (cipher.Addresser, error) {
	a, err := eth.AddressFromPubKey(key)
	if err != nil {
		return nil, fmt.Errorf("ethereum address from pubkey: %w", err)
	}
	return a, nil
}

func (e ethereumDecoder) DecodeBase58Address(addr string) (cipher.Addresser, error) {
	return eth.ParseAddress(addr)
}

func (e ethereumDecoder) SecKeyToHex(secKey cipher.SecKey) string {
	return secKey.Hex()
}

func (e ethereumDecoder) SecKeyFromHex(secKey string) (cipher.SecKey, error) {
	return cipher.SecKeyFromHex(strings.TrimPrefix(secKey, "0x"))
}
