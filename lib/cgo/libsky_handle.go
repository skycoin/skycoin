package main

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_handle_close
func SKY_handle_close(handle *C.Handle){
	closeHandle(Handle(*handle))
}