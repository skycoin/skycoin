package main

import (
	file "github.com/skycoin/skycoin/src/file"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_file_InitDataDir
func SKY_file_InitDataDir(_dir string, _arg1 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	dir := _dir
	__arg1, ____return_err := file.InitDataDir(dir)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg1, _arg1)
	}
	return
}

// export SKY_file_UserHome
func SKY_file_UserHome(_arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := file.UserHome()
	copyString(__arg0, _arg0)
	return
}

// export SKY_file_LoadJSON
func SKY_file_LoadJSON(_filename string, _thing interface{}) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	filename := _filename
	____return_err := file.LoadJSON(filename, thing)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_file_SaveJSON
func SKY_file_SaveJSON(_filename string, _thing interface{}, _mode *C.FileMode) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	filename := _filename
	____return_err := file.SaveJSON(filename, thing, mode)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_file_SaveJSONSafe
func SKY_file_SaveJSONSafe(_filename string, _thing interface{}, _mode *C.FileMode) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	filename := _filename
	____return_err := file.SaveJSONSafe(filename, thing, mode)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_file_SaveBinary
func SKY_file_SaveBinary(_filename string, _data *C.GoSlice_, _mode *C.FileMode) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	filename := _filename
	data := *(*[]byte)(unsafe.Pointer(_data))
	____return_err := file.SaveBinary(filename, data, mode)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_file_ResolveResourceDirectory
func SKY_file_ResolveResourceDirectory(_path string, _arg1 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	path := _path
	__arg1 := file.ResolveResourceDirectory(path)
	copyString(__arg1, _arg1)
	return
}

// export SKY_file_DetermineResourcePath
func SKY_file_DetermineResourcePath(_staticDir string, _resourceDir string, _devDir string, _arg3 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	staticDir := _staticDir
	resourceDir := _resourceDir
	devDir := _devDir
	__arg3, ____return_err := file.DetermineResourcePath(staticDir, resourceDir, devDir)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyString(__arg3, _arg3)
	}
	return
}

// export SKY_file_CopyFile
func SKY_file_CopyFile(_dst string, _src *C.Reader, _arg2 *int64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	dst := _dst
	__arg2, ____return_err := file.CopyFile(dst, src)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = __arg2
	}
	return
}
