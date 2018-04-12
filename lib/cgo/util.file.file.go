package main

import file "github.com/skycoin/skycoin/src/util/file"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_file_InitDataDir
func SKY_file_InitDataDir(_dir string, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	dir := _dir
	__arg1, ____return_err := file.InitDataDir(dir)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}

// export SKY_file_UserHome
func SKY_file_UserHome(_arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	__arg0 := file.UserHome()
	copyString(__arg0, _arg0)
	return
}

// export SKY_file_LoadJSON
func SKY_file_LoadJSON(_filename string, _thing interface{}) (____error_code uint32) {
	____error_code = 0
	filename := _filename
	____return_err := file.LoadJSON(filename, thing)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_file_ResolveRetargetDirectory
func SKY_file_ResolveRetargetDirectory(_path string, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	path := _path
	__arg1 := file.ResolveRetargetDirectory(path)
	copyString(__arg1, _arg1)
	return
}

// export SKY_file_DetermineRetargetPath
func SKY_file_DetermineRetargetPath(_staticDir string, _retargetDir string, _devDir string, _arg3 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	staticDir := _staticDir
	retargetDir := _retargetDir
	devDir := _devDir
	__arg3, ____return_err := file.DetermineRetargetPath(staticDir, retargetDir, devDir)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg3, _arg3)
	}
	return
}
