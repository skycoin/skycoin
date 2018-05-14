package main

import (
	bucket "github.com/skycoin/skycoin/src/visor/dbutil"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_bucket_Itob
func SKY_bucket_Itob(_v uint64, _arg1 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	v := _v
	__arg1 := bucket.Itob(v)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

//export SKY_bucket_Btoi
func SKY_bucket_Btoi(_v []byte, _arg1 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	v := *(*[]byte)(unsafe.Pointer(&_v))
	__arg1 := bucket.Btoi(v)
	*_arg1 = __arg1
	return
}
