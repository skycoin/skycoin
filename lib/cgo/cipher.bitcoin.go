package main

/*
#include <string.h>
#include <stdlib.h>

#include "skytypes.h"

*/
import "C"

import (
	"reflect"
	"unsafe"

	cipher "github.com/skycoin/skycoin/src/cipher"
)

/**
 * Functions in github.com/skycoin/skycoin/src/cipher/bitcoin.go
 */

//export SKY_cipher_DecodeBase58BitcoinAddress
func SKY_cipher_DecodeBase58BitcoinAddress(_addr string, _arg1 *C.cipher__BitcoinAddress) uint32 {
	addr, err := cipher.DecodeBase58BitcoinAddress(_addr)
	errcode := libErrorCode(err)
	if err == nil {
		*_arg1 = *(*C.cipher__BitcoinAddress)(unsafe.Pointer(&addr))
	}
	return errcode
}

//export SKY_cipher_BitcoinAddressFromPubKey
func SKY_cipher_BitcoinAddressFromPubKey(_pubKey *C.cipher__PubKey, _arg1 *C.cipher__BitcoinAddress) {
	pubKey := (*cipher.PubKey)(unsafe.Pointer(_pubKey))

	addr := cipher.BitcoinAddressFromPubKey(*pubKey)
	*_arg1 = *(*C.cipher__BitcoinAddress)(unsafe.Pointer(&addr))
}

//export SKY_cipher_BitcoinAddressFromSecKey
func SKY_cipher_BitcoinAddressFromSecKey(_secKey *C.cipher__SecKey, _arg1 *C.cipher__BitcoinAddress) uint32 {
	secKey := (*cipher.SecKey)(unsafe.Pointer(_secKey))

	addr, err := cipher.BitcoinAddressFromSecKey(*secKey)
	if err == nil {
		*_arg1 = *(*C.cipher__BitcoinAddress)(unsafe.Pointer(&addr))
	}
	return libErrorCode(err)
}

//export SKY_cipher_BitcoinWalletImportFormatFromSeckey
func SKY_cipher_BitcoinWalletImportFormatFromSeckey(_seckey *C.cipher__SecKey, _arg1 *C.GoString_) {
	seckey := (*cipher.SecKey)(unsafe.Pointer(_seckey))
	s := cipher.BitcoinWalletImportFormatFromSeckey(*seckey)
	copyString(s, _arg1)
}

//export SKY_cipher_BitcoinAddressFromBytes
func SKY_cipher_BitcoinAddressFromBytes(_b []byte, _arg1 *C.cipher__BitcoinAddress) uint32 {
	addr, err := cipher.BitcoinAddressFromBytes(_b)
	errcode := libErrorCode(err)
	if err == nil {
		*_arg1 = *(*C.cipher__BitcoinAddress)(unsafe.Pointer(&addr))
	}
	return errcode
}

//export SKY_cipher_SecKeyFromBitcoinWalletImportFormat
func SKY_cipher_SecKeyFromBitcoinWalletImportFormat(_input string, _arg1 *C.cipher__SecKey) uint32 {
	seckey, err := cipher.SecKeyFromBitcoinWalletImportFormat(_input)
	errcode := libErrorCode(err)
	if err == nil {
		*_arg1 = *(*C.cipher__SecKey)(unsafe.Pointer(&seckey))
	}
	return errcode
}

//export SKY_cipher_BitcoinAddress_Null
func SKY_cipher_BitcoinAddress_Null(_addr *C.cipher__BitcoinAddress) bool {
	addr := (*cipher.BitcoinAddress)(unsafe.Pointer(_addr))
	return addr.Null()
}

//export SKY_cipher_BitcoinAddress_Bytes
func SKY_cipher_BitcoinAddress_Bytes(_addr *C.cipher__BitcoinAddress, _arg0 *C.GoSlice_) {
	addr := (*cipher.BitcoinAddress)(unsafe.Pointer(_addr))
	bytes := addr.Bytes()
	copyToGoSlice(reflect.ValueOf(bytes), _arg0)
}

//export SKY_cipher_BitcoinAddress_Verify
func SKY_cipher_BitcoinAddress_Verify(_addr *C.cipher__BitcoinAddress, _key *C.cipher__PubKey) uint32 {
	addr := (*cipher.BitcoinAddress)(unsafe.Pointer(_addr))
	key := (*cipher.PubKey)(unsafe.Pointer(_key))
	err := addr.Verify(*key)
	return libErrorCode(err)
}

//export SKY_cipher_BitcoinAddress_String
func SKY_cipher_BitcoinAddress_String(_addr *C.cipher__BitcoinAddress, _arg1 *C.GoString_) {
	addr := (*cipher.BitcoinAddress)(unsafe.Pointer(_addr))
	s := addr.String()
	copyString(s, _arg1)
}

//export SKY_cipher_BitcoinAddress_Checksum
func SKY_cipher_BitcoinAddress_Checksum(_addr *C.cipher__BitcoinAddress, _arg0 *C.cipher__Checksum) {
	addr := (*cipher.BitcoinAddress)(unsafe.Pointer(_addr))
	cs := addr.Checksum()
	C.memcpy(unsafe.Pointer(_arg0), unsafe.Pointer(&cs[0]), C.size_t(len(cs)))
}
