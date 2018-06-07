package main

import file "github.com/skycoin/skycoin/src/util/file"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_file_InitDataDir
func SKY_file_InitDataDir(_dir string, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	dir := _dir
	__arg1, ____return_err := file.InitDataDir(dir)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}

//export SKY_file_UserHome
func SKY_file_UserHome(_arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := file.UserHome()
	copyString(__arg0, _arg0)
	return
}

//export SKY_file_ResolveResourceDirectory
func SKY_file_ResolveResourceDirectory(_path string, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	path := _path
	__arg1 := file.ResolveResourceDirectory(path)
	copyString(__arg1, _arg1)
	return
}

//export SKY_file_DetermineResourcePath
func SKY_file_DetermineResourcePath(_staticDir string, _resourceDir string, _devDir string, _arg3 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	staticDir := _staticDir
	resourceDir := _resourceDir
	devDir := _devDir
	__arg3, ____return_err := file.DetermineResourcePath(staticDir, resourceDir, devDir)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg3, _arg3)
	}
	return
}
