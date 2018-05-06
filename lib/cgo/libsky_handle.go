package main

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

import (
	cli "github.com/skycoin/skycoin/src/api/cli"
	"unsafe"
)

type Handle uint64

var (
	handleMap = make(map[Handle]interface{})
)

func registerHandle(obj interface{}) Handle {
	ptr := &obj
	handle := *(*Handle)(unsafe.Pointer(&ptr))
	handleMap[handle] = obj
	return handle
}

func lookupHandleObj(handle Handle) (interface{}, bool) {
	obj, ok := handleMap[handle]
	return obj, ok
}

func registerPasswordReaderHandle(obj cli.PasswordReader) C.PasswordReader__Handle {
	return (C.PasswordReader__Handle)(registerHandle(obj))
}

func lookupPasswordReaderHandle(handle C.PasswordReader__Handle) (cli.PasswordReader, bool) {
	obj, ok := lookupHandleObj(Handle(handle))
	if ok {
		if obj, isOK := (obj).(cli.PasswordReader); isOK {
			return obj, true
		}
	}
	return nil, false
}

func closeHandle(handle Handle) {
	delete(handleMap, handle)
}

//export SKY_handle_close
func SKY_handle_close(handle C.Handle) {
	closeHandle(Handle(handle))
}
