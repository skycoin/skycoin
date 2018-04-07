package main

/*
#include <string.h>
#include <stdlib.h>

#include "../../include/skytypes.h"

*/
import "C"

import (
	"unsafe"

	"github.com/skycoin/skycoin/src/cipher"
)

/**
 * Functions in github.com/skycoin/skycoin/src/cipher/address.go
 */

//export SKY_cipher_DecodeBase58Address
func SKY_cipher_DecodeBase58Address(_addr string, _arg1 *C.Address) uint32 {
	addr, err := cipher.DecodeBase58Address(_addr)
	errcode := libErrorCode(err)
	if err == nil {
		*_arg1 = *(*C.Address)(unsafe.Pointer(&addr))
	}
	return errcode
}

//export SKY_cipher_AddressFromPubKey
func SKY_cipher_AddressFromPubKey(_pubKey *C.PubKey, _arg1 *C.Address) {
	pubKey := (*cipher.PubKey)(unsafe.Pointer(_pubKey))

	addr := cipher.AddressFromPubKey(*pubKey)
	*_arg1 = *(*C.Address)(unsafe.Pointer(&addr))
}

//export SKY_cipher_AddressFromSecKey
func SKY_cipher_AddressFromSecKey(_secKey *C.SecKey, _arg1 *C.Address) {
	var secKey cipher.SecKey
	secKey = *(*cipher.SecKey)(unsafe.Pointer(_secKey))
	addr := cipher.AddressFromSecKey(secKey)
	*_arg1 = *(*C.Address)(unsafe.Pointer(&addr))
}

//export SKY_cipher_BitcoinDecodeBase58Address
func SKY_cipher_BitcoinDecodeBase58Address(_addr string, _arg1 *C.Address) uint32 {
	addr, err := cipher.BitcoinDecodeBase58Address(_addr)
	errcode := libErrorCode(err)
	if err == nil {
		*_arg1 = *(*C.Address)(unsafe.Pointer(&addr))
	}
	return errcode
}

//export SKY_cipher_Address_Bytes
func SKY_cipher_Address_Bytes(_addr *C.Address, _arg0 *C.GoSlice_) {
	addr := (*cipher.Address)(unsafe.Pointer(_addr))
	bytes := addr.Bytes()
	bytes_len := len(bytes)
	if bytes_len > int(_arg0.cap) {
		// Negative len on cap overflow
		_arg0.len = _arg0.cap - C.GoInt_(bytes_len)
	} else {
		C.memcpy(unsafe.Pointer(_arg0.data), unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)))
		_arg0.len = C.GoInt_(bytes_len)
	}
}

//export SKY_cipher_Address_BitcoinBytes
func SKY_cipher_Address_BitcoinBytes(_addr *C.Address, _arg0 *C.GoSlice_) {
	addr := (*cipher.Address)(unsafe.Pointer(_addr))
	bytes := addr.BitcoinBytes()
	bytes_len := len(bytes)
	if bytes_len > int(_arg0.cap) {
		// Negative len on cap overflow
		_arg0.len = _arg0.cap - C.GoInt_(bytes_len)
	} else {
		C.memcpy(unsafe.Pointer(_arg0.data), unsafe.Pointer(&bytes[0]), C.size_t(bytes_len))
		_arg0.len = C.GoInt_(bytes_len)
	}
}

//export SKY_cipher_Address_Verify
func SKY_cipher_Address_Verify(_addr *C.Address, _key *C.PubKey) uint32 {
	addr := (*cipher.Address)(unsafe.Pointer(_addr))
	key := (*cipher.PubKey)(unsafe.Pointer(_key))
	err := addr.Verify(*key)
	return libErrorCode(err)
}

//export SKY_cipher_Address_String
func SKY_cipher_Address_String(_addr *C.Address, _arg1 *C.GoString_) {
	addr := (*cipher.Address)(unsafe.Pointer(_addr))
	s := addr.String()
	copyString(s, _arg1)
}

//export SKY_cipher_Address_BitcoinString
func SKY_cipher_Address_BitcoinString(_addr *C.Address, _arg1 *C.GoString_) {
	addr := (*cipher.Address)(unsafe.Pointer(_addr))
	s := addr.BitcoinString()
	copyString(s, _arg1)
}

//export SKY_cipher_Address_Checksum
func SKY_cipher_Address_Checksum(_addr *C.Address, _arg0 *C.Checksum) {
	addr := (*cipher.Address)(unsafe.Pointer(_addr))
	cs := addr.Checksum()
	C.memcpy(unsafe.Pointer(_arg0), unsafe.Pointer(&cs[0]), C.size_t(len(cs)))
}

//export SKY_cipher_Address_BitcoinChecksum
func SKY_cipher_Address_BitcoinChecksum(_addr *C.Address, _arg0 *C.Checksum) {
	addr := (*cipher.Address)(unsafe.Pointer(_addr))
	cs := addr.BitcoinChecksum()
	C.memcpy(unsafe.Pointer(_arg0), unsafe.Pointer(&cs[0]), C.size_t(len(cs)))
}

/*
Bitcoin Functions
*/

//export SKY_cipher_BitcoinAddressFromPubkey
func SKY_cipher_BitcoinAddressFromPubkey(_pubkey *C.PubKey, _arg1 *C.GoString_) {
	pubkey := (*cipher.PubKey)(unsafe.Pointer(_pubkey))
	s := cipher.BitcoinAddressFromPubkey(*pubkey)
	copyString(s, _arg1)
}

//export SKY_cipher_BitcoinWalletImportFormatFromSeckey
func SKY_cipher_BitcoinWalletImportFormatFromSeckey(_seckey *C.SecKey, _arg1 *C.GoString_) {
	seckey := (*cipher.SecKey)(unsafe.Pointer(_seckey))
	s := cipher.BitcoinWalletImportFormatFromSeckey(*seckey)
	copyString(s, _arg1)
}

//export SKY_cipher_BitcoinAddressFromBytes
func SKY_cipher_BitcoinAddressFromBytes(_b []byte, _arg1 *C.Address) uint32 {
	addr, err := cipher.BitcoinAddressFromBytes(_b)
	errcode := libErrorCode(err)
	if err == nil {
		*_arg1 = *(*C.Address)(unsafe.Pointer(&addr))
	}
	return errcode
}

//export SKY_cipher_SecKeyFromWalletImportFormat
func SKY_cipher_SecKeyFromWalletImportFormat(_input string, _arg1 *C.SecKey) uint32 {
	seckey, err := cipher.SecKeyFromWalletImportFormat(_input)
	errcode := libErrorCode(err)
	if err == nil {
		*_arg1 = *(*C.SecKey)(unsafe.Pointer(&seckey))
	}
	return errcode
}
