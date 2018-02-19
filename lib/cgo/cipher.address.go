package main

/*

#include <string.h>
#include <stdlib.h>

typedef unsigned char Ripemd160[20];

typedef struct {
	unsigned char Version;
	Ripemd160 Key;
} Address;

typedef unsigned char Checksum[4];

typedef unsigned char PubKey[33];
typedef unsigned char SecKey[32];
typedef unsigned char Checksum[4];
*/
import "C"

import (
	"unsafe"

	"github.com/skycoin/skycoin/src/cipher"
)

/**
 * Functions in github.com/skycoin/skycoin/src/cipher/address.go
 */

//export SKY_Cipher_DecodeBase58Address
func SKY_Cipher_DecodeBase58Address(_strAddr string, _retAddr *C.Address) C.uint {
	addr, err := cipher.DecodeBase58Address(_strAddr)
	*_retAddr = *(*C.Address)(unsafe.Pointer(&addr))

	errCode := 1
	if err != nil {
		errCode = 0
	}
	return C.uint(errCode)
}

// export SKY_Cipher_AddressFromPubKey
func SKY_Cipher_AddressFromPubKey(_pubKey *C.PubKey, _retAddr *C.Address) {
	var pubKey cipher.PubKey
	pubKey = *(*cipher.PubKey)(unsafe.Pointer(_pubKey))
	addr := cipher.AddressFromPubKey(pubKey)
	*_retAddr = *(*C.Address)(unsafe.Pointer(&addr))
}

// export SKY_Cipher_AddressFromSecKey
func SKY_Cipher_AddressFromSecKey(_secKey *C.SecKey, _retAddr *C.Address) {
	var secKey cipher.SecKey
	secKey = *(*cipher.SecKey)(unsafe.Pointer(_secKey))
	addr := cipher.AddressFromSecKey(secKey)
	*_retAddr = *(*C.Address)(unsafe.Pointer(&addr))
}

// export SKY_Cipher_BitcoinDecodeBase58Address
func SKY_Cipher_BitcoinDecodeBase58Address(_strAddr string, _retAddr *C.Address) C.uint {
	addr, err := cipher.BitcoinDecodeBase58Address(_strAddr)
	*_retAddr = *(*C.Address)(unsafe.Pointer(&addr))

	errCode := 1
	if err != nil {
		errCode = 0
	}
	return C.uint(errCode)
}

// export SKY_Cipher_Address_Bytes
func SKY_Cipher_Address_Bytes(_addr *C.Address, _ret *C.uchar) {
	addr := (*cipher.Address)(unsafe.Pointer(_addr))
	bytes := addr.Bytes()
	C.memcpy(unsafe.Pointer(_ret), unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)))
}

// export SKY_Cipher_Address_BitcoinBytes
func SKY_Cipher_Address_BitcoinBytes(_addr *C.Address, _ret *C.uchar) {
	addr := (*cipher.Address)(unsafe.Pointer(_addr))
	bytes := addr.BitcoinBytes()
	C.memcpy(unsafe.Pointer(_ret), unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)))
}

// export SKY_Cipher_Address_Verify
func SKY_Cipher_Address_Verify(_addr *C.Address, _key *C.PubKey) C.uint {
	addr := (*cipher.Address)(unsafe.Pointer(_addr))
	key := (*cipher.PubKey)(unsafe.Pointer(&_key))
	err := addr.Verify(*key)
	if err != nil {
		return 1
	}
	return 0
}

// export SKY_Cipher_Address_String
func SKY_Cipher_Address_String(_addr *C.Address) string {
	addr := (*cipher.Address)(unsafe.Pointer(_addr))
	return addr.String()
}

// export SKY_Cipher_Address_BitcoinString
func SKY_Cipher_Address_BitcoinString(_addr *C.Address) string {
	addr := (*cipher.Address)(unsafe.Pointer(_addr))
	return addr.BitcoinString()
}

// export SKY_Cipher_Address_Checksum
func SKY_Cipher_Address_Checksum(_addr *C.Address, _ret *C.Checksum) {
	addr := (*cipher.Address)(unsafe.Pointer(_addr))
	cs := addr.Checksum()
	C.memcpy(unsafe.Pointer(_ret), unsafe.Pointer(&cs[0]), C.size_t(len(cs)))
}

// export SKY_Cipher_Address_BitcoinChecksum
func SKY_Cipher_Address_BitcoinChecksum(_addr *C.Address, _ret *C.Checksum) {
	addr := (*cipher.Address)(unsafe.Pointer(_addr))
	cs := addr.BitcoinChecksum()
	C.memcpy(unsafe.Pointer(_ret), unsafe.Pointer(&cs[0]), C.size_t(len(cs)))
}

/*
Bitcoin Functions
*/

// export SKY_Cipher_BitcoinAddressFromPubkey
func SKY_Cipher_BitcoinAddressFromPubkey(_pubkey *C.PubKey) string {
	pubkey := (*cipher.PubKey)(unsafe.Pointer(_pubkey))
	return cipher.BitcoinAddressFromPubkey(*pubkey)
}

// export SKY_Cipher_BitcoinWalletImportFormatFromSeckey
func SKY_Cipher_BitcoinWalletImportFormatFromSeckey(_seckey *C.SecKey) string {
	seckey := (*cipher.SecKey)(unsafe.Pointer(_seckey))
	return cipher.BitcoinWalletImportFormatFromSeckey(*seckey)
}

// export SKY_Cipher_BitcoinAddressFromBytes
func SKY_Cipher_BitcoinAddressFromBytes(_b *C.uchar, _sz C.size_t, _retAddr *C.Address) C.uint {
	b := C.GoBytes(unsafe.Pointer(_b), C.int(_sz))
	addr, err := cipher.BitcoinAddressFromBytes(b)
	if err != nil {
		*_retAddr = *(*C.Address)(unsafe.Pointer(&addr))
		return 1
	} else {
		return 0
	}
}

// export SKY_Cipher_SecKeyFromWalletImportFormat
func SKY_Cipher_SecKeyFromWalletImportFormat(_input string, _seckey *C.SecKey) C.uint {
	seckey, err := cipher.SecKeyFromWalletImportFormat(_input)
	if err != nil {
		*_seckey = *(*C.SecKey)(unsafe.Pointer(&seckey))
		return 1
	} else {
		return 0
	}
}
