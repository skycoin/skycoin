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
	"github.com/skycoin/skycoin/src/cipher"
)

/**
 * Cipher API
 */

//export DecodeBase58Address
func DecodeBase58Address(strAddr string, cAddr *C.Address) C.int {
	_, err := cipher.DecodeBase58Address(strAddr)
	if err != nil {
		return 0
	}
	// TODO: Copy memory
	return 1
}

func main() {}
