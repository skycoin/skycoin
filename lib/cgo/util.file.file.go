package main

import (
	file "github.com/skycoin/skycoin/src/util/file"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_file_InitDataDir
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

// export SKY_file_UserHome
func SKY_file_UserHome(_arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := file.UserHome()
	copyString(__arg0, _arg0)
	return
}

// export SKY_file_LoadJSON
func SKY_file_LoadJSON(_filename string, _thing interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	filename := _filename
	____return_err := file.LoadJSON(filename, thing)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_file_SaveJSON
func SKY_file_SaveJSON(_filename string, _thing interface{}, _mode *C.FileMode) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	filename := _filename
	____return_err := file.SaveJSON(filename, thing, mode)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_file_SaveJSONSafe
func SKY_file_SaveJSONSafe(_filename string, _thing interface{}, _mode *C.FileMode) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	filename := _filename
	____return_err := file.SaveJSONSafe(filename, thing, mode)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_file_SaveBinary
func SKY_file_SaveBinary(_filename string, _data *C.GoSlice_, _mode *C.FileMode) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	filename := _filename
	data := *(*[]byte)(unsafe.Pointer(_data))
	____return_err := file.SaveBinary(filename, data, mode)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_file_ResolveRetargetDirectory
func SKY_file_ResolveRetargetDirectory(_path string, _arg1 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	path := _path
	__arg1 := file.ResolveRetargetDirectory(path)
	copyString(__arg1, _arg1)
	return
}

// export SKY_file_DetermineRetargetPath
func SKY_file_DetermineRetargetPath(_staticDir string, _retargetDir string, _devDir string, _arg3 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
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

// export SKY_file_CopyFile
func SKY_file_CopyFile(_dst string, _src *C.Reader, _arg2 *int64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	dst := _dst
	__arg2, ____return_err := file.CopyFile(dst, src)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = __arg2
	}
	return
}
