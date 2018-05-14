package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_cipher_AddressFromPubKey
func SKY_cipher_AddressFromPubKey(_pubKey *C.cipher__PubKey, _arg1 *C.cipher__Address) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pubKey := *(*cipher.PubKey)(unsafe.Pointer(_pubKey))
	__arg1 := cipher.AddressFromPubKey(pubKey)
	*_arg1 = *(*C.cipher__Address)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_cipher_AddressFromSecKey
func SKY_cipher_AddressFromSecKey(_secKey *C.cipher__SecKey, _arg1 *C.cipher__Address) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	secKey := *(*cipher.SecKey)(unsafe.Pointer(_secKey))
	__arg1 := cipher.AddressFromSecKey(secKey)
	*_arg1 = *(*C.cipher__Address)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_cipher_DecodeBase58Address
func SKY_cipher_DecodeBase58Address(_addr string, _arg1 *C.cipher__Address) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addr := _addr
	__arg1, ____return_err := cipher.DecodeBase58Address(addr)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.cipher__Address)(unsafe.Pointer(&__arg1))
	}
	return
}

//export SKY_cipher_MustDecodeBase58Address
func SKY_cipher_MustDecodeBase58Address(_addr string, _arg1 *C.cipher__Address) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addr := _addr
	__arg1 := cipher.MustDecodeBase58Address(addr)
	*_arg1 = *(*C.cipher__Address)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_cipher_BitcoinDecodeBase58Address
func SKY_cipher_BitcoinDecodeBase58Address(_addr string, _arg1 *C.cipher__Address) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addr := _addr
	__arg1, ____return_err := cipher.BitcoinDecodeBase58Address(addr)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.cipher__Address)(unsafe.Pointer(&__arg1))
	}
	return
}

//export SKY_cipher_BitcoinMustDecodeBase58Address
func SKY_cipher_BitcoinMustDecodeBase58Address(_addr string, _arg1 *C.cipher__Address) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addr := _addr
	__arg1 := cipher.BitcoinMustDecodeBase58Address(addr)
	*_arg1 = *(*C.cipher__Address)(unsafe.Pointer(&__arg1))
	return
}

//export SKY_cipher_AddressFromBytes
func SKY_cipher_AddressFromBytes(_b []byte, _arg1 *C.cipher__Address) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(&_b))
	__arg1, ____return_err := cipher.AddressFromBytes(b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.cipher__Address)(unsafe.Pointer(&__arg1))
	}
	return
}

//export SKY_cipher_MustAddressFromBytes
func SKY_cipher_MustAddressFromBytes(_b []byte, _arg1 *C.cipher__Address) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(&_b))
	__arg1, ____return_err := cipher.MustAddressFromBytes(b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.cipher__Address)(unsafe.Pointer(&__arg1))
	}
	return
}

//export SKY_cipher_Address_Null
func SKY_cipher_Address_Null(_addr *C.cipher__Address, _arg0 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addr := *inplaceAddress(_addr)
	__arg0 := addr.Null()
	*_arg0 = __arg0
	return
}

//export SKY_cipher_Address_Bytes
func SKY_cipher_Address_Bytes(_addr *C.cipher__Address, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addr := inplaceAddress(_addr)
	__arg0 := addr.Bytes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_cipher_Address_BitcoinBytes
func SKY_cipher_Address_BitcoinBytes(_addr *C.cipher__Address, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addr := inplaceAddress(_addr)
	__arg0 := addr.BitcoinBytes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}

//export SKY_cipher_Address_Verify
func SKY_cipher_Address_Verify(_addr *C.cipher__Address, _key *C.cipher__PubKey) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addr := *inplaceAddress(_addr)
	key := *(*cipher.PubKey)(unsafe.Pointer(_key))
	____return_err := addr.Verify(key)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_cipher_Address_String
func SKY_cipher_Address_String(_addr *C.cipher__Address, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addr := *inplaceAddress(_addr)
	__arg0 := addr.String()
	copyString(__arg0, _arg0)
	return
}

//export SKY_cipher_Address_BitcoinString
func SKY_cipher_Address_BitcoinString(_addr *C.cipher__Address, _arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addr := *inplaceAddress(_addr)
	__arg0 := addr.BitcoinString()
	copyString(__arg0, _arg0)
	return
}

//export SKY_cipher_Address_Checksum
func SKY_cipher_Address_Checksum(_addr *C.cipher__Address, _arg0 *C.cipher__Checksum) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addr := inplaceAddress(_addr)
	__arg0 := addr.Checksum()
	*_arg0 = *(*C.cipher__Checksum)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_cipher_Address_BitcoinChecksum
func SKY_cipher_Address_BitcoinChecksum(_addr *C.cipher__Address, _arg0 *C.cipher__Checksum) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	addr := inplaceAddress(_addr)
	__arg0 := addr.BitcoinChecksum()
	*_arg0 = *(*C.cipher__Checksum)(unsafe.Pointer(&__arg0))
	return
}

//export SKY_cipher_BitcoinAddressFromPubkey
func SKY_cipher_BitcoinAddressFromPubkey(_pubkey *C.cipher__PubKey, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pubkey := *(*cipher.PubKey)(unsafe.Pointer(_pubkey))
	__arg1 := cipher.BitcoinAddressFromPubkey(pubkey)
	copyString(__arg1, _arg1)
	return
}

//export SKY_cipher_BitcoinWalletImportFormatFromSeckey
func SKY_cipher_BitcoinWalletImportFormatFromSeckey(_seckey *C.cipher__SecKey, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	seckey := *(*cipher.SecKey)(unsafe.Pointer(_seckey))
	__arg1 := cipher.BitcoinWalletImportFormatFromSeckey(seckey)
	copyString(__arg1, _arg1)
	return
}

//export SKY_cipher_BitcoinAddressFromBytes
func SKY_cipher_BitcoinAddressFromBytes(_b []byte, _arg1 *C.cipher__Address) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	b := *(*[]byte)(unsafe.Pointer(&_b))
	__arg1, ____return_err := cipher.BitcoinAddressFromBytes(b)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.cipher__Address)(unsafe.Pointer(&__arg1))
	}
	return
}

//export SKY_cipher_SecKeyFromWalletImportFormat
func SKY_cipher_SecKeyFromWalletImportFormat(_input string, _arg1 *C.cipher__SecKey) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	input := _input
	__arg1, ____return_err := cipher.SecKeyFromWalletImportFormat(input)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofSecKey))
	}
	return
}

//export SKY_cipher_MustSecKeyFromWalletImportFormat
func SKY_cipher_MustSecKeyFromWalletImportFormat(_input string, _arg1 *C.cipher__SecKey) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	input := _input
	__arg1 := cipher.MustSecKeyFromWalletImportFormat(input)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofSecKey))
	return
}
