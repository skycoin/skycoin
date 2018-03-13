package main

import (
	"unsafe"
)

type Handle uint64

var (
	handleMap = make(map[Handle]interface{})
)

func openHandle(obj interface{}) Handle {
	ptr := &obj
	handle := *(*Handle)(unsafe.Pointer(&ptr))
	handleMap[handle] = obj
	return handle
}

func lookupHandleObj(handle Handle) (interface{}, bool) {
	obj, ok := handleMap[handle]
	return obj, ok
}

func closeHandle(handle Handle) {
	delete(handleMap, handle)
}

func inplaceArrayObj(p unsafe.Pointer, length int) interface{} {
	// Create slice without copying data
	return (*[1 << 30]byte)(p)[:length:length]

}
