package main

/*

typedef unsigned char Ripemd160[20];

typedef struct {
	unsigned char Version;
	Ripemd160 Key;
} Address;

*/
import "C"

import (
	"unsafe"

	"github.com/skycoin/skycoin/src/cipher"
)

/**
 * Cipher API
 */

//export DecodeBase58Address
func DecodeBase58Address(strAddr string, retAddr *C.Address) C.int {
	addr, err := cipher.DecodeBase58Address(strAddr)
	*retAddr = *(*C.Address)(unsafe.Pointer(&addr))

	errCode := 1
	if err != nil {
		errCode = 0
	}
	return C.int(errCode)
}

func main() {}
