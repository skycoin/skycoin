package main

import (
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

func copyToBuffer(sourceType reflect.Value, p unsafe.Pointer, c int){
}

func copyString(source string, dest *C.GoString_){
}