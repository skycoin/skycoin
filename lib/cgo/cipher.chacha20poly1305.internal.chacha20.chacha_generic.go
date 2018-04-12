package main
/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"
/*
TODO: stdevEclipse Cant import internal
// export SKY_chacha20_XORKeyStream
func SKY_chacha20_XORKeyStream(_out, _in *C.GoSlice_, _counter *[]byte, _key *[]byte) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	//TODO: stdevEclipse Check Pointer typecast
	out := *(*[]byte)(unsafe.Pointer(_out))
	in := *(*[]byte)(unsafe.Pointer(_in))
	counter := *(*[]byte)(unsafe.Pointer(_counter))
	key := *(*[]byte)(unsafe.Pointer(_key))
	chacha20poly1305.XORKeyStream(out, in, counter, key)
	return
}
*/