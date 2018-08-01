package main

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_map_Get
func SKY_map_Get(gomap *C.GoStringMap_, key string, value *C.GoString_) (____error_code uint32) {
	obj, ok := lookupHandle(C.Handle(*gomap))
	____error_code = SKY_ERROR
	if ok {
		if m, isMap := (obj).(map[string]string); isMap {
			result, ok := m[key]
			if ok {
				copyString(result, value)
				____error_code = SKY_OK
			}
		}
	}
	return
}

//export SKY_map_HasKey
func SKY_map_HasKey(gomap *C.GoStringMap_, key string) (found bool) {
	obj, ok := lookupHandle(C.Handle(*gomap))
	found = false
	if ok {
		if m, isMap := (obj).(map[string]string); isMap {
			_, found = m[key]
		}
	}
	return
}

//export SKY_map_Close
func SKY_map_Close(gomap *C.GoStringMap_) (____error_code uint32) {
	obj, ok := lookupHandle(C.Handle(*gomap))
	____error_code = SKY_ERROR
	if ok {
		if _, isMap := (obj).(map[string]string); isMap {
			closeHandle(Handle(*gomap))
			____error_code = SKY_OK
		}
	}
	return

}
