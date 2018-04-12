package main

import (
	"reflect"
	"unsafe"
	cipher "github.com/skycoin/skycoin/src/cipher"
	"hash"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

func libErrorCode(error) uint32{
	return 0
}

func catchApiPanic(uint32, interface{})  uint32 {
	return 0
}

func copyString(source string, dest *C.GoString_){
}

func copyToGoSlice(v reflect.Value, s *C.GoSlice_){
}

func copyToBuffer(sourceType reflect.Value, p unsafe.Pointer, c uint){
}

func copyToInterface(a *C.GoInterface_) interface{}{
	return nil
}

func inplacePubKeySlice(*C.cipher__PubKeySlice) (ret cipher.PubKeySlice) {
	return 
}

func inplaceAddress(*C.cipher__Address) (ret *cipher.Address) {
	return 
}

func copyToFunc(f C.Handle) func() hash.Hash {
	return nil
}

func copyStringMap( source map[string]string, dest *C.GoMap_ ){
}

func main(){
}

//TODO: stdevEclipse Get Sizes
var (
	SKY_OK				uint32 	= 0
	SizeofPubKey 		uint 	= 32
	SizeofRipemd160 	uint 	= 32
	SizeofSecKey 		uint 	= 32
	SizeofSig			uint 	= 32
	SizeofSHA256		uint 	= 32
)

