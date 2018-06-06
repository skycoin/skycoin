package main

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_map_get
func SKY_map_get(gomap *C.GoStringMap_, key string, value *C.GoString_) (____error_code uint32) {
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

//export SKY_map_has_key
func SKY_map_has_key(gomap *C.GoStringMap_, key string) (found bool) {
	obj, ok := lookupHandle(C.Handle(*gomap))
	found = false
	if ok {
		if m, isMap := (obj).(map[string]string); isMap {
			_, found = m[key]
		}
	}
	return
}

//export SKY_map_close
func SKY_map_close(gomap *C.GoStringMap_) (____error_code uint32) {
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
