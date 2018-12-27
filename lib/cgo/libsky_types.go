package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

func toAddresserArray(addrs []cipher.Address) []cipher.Addresser {
	// TODO : Support for arrays of interface objects in cgogen
	var __addrs = make([]cipher.Addresser, len(addrs))
	for _, addr := range addrs {
		__addrs = append(__addrs, addr)
	}
	return __addrs
}
