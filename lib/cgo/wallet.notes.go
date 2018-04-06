package main

import (
	cipher "github.com/skycoin/skycoin/src/cipher"
	wallet "github.com/skycoin/skycoin/src/wallet"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_wallet_NewNotesFilename
func SKY_wallet_NewNotesFilename(_arg0 *C.GoString_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	__arg0 := wallet.NewNotesFilename()
	copyString(__arg0, _arg0)
	return
}

// export SKY_wallet_LoadNotes
func SKY_wallet_LoadNotes(_dir string, _arg1 *C.Notes) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	dir := _dir
	__arg1, ____return_err := wallet.LoadNotes(dir)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofNotes))
	}
	return
}

// export SKY_wallet_LoadReadableNotes
func SKY_wallet_LoadReadableNotes(_filename string, _arg1 *C.ReadableNotes) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	filename := _filename
	__arg1, ____return_err := wallet.LoadReadableNotes(filename)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg1)[:]), unsafe.Pointer(_arg1), uint(SizeofReadableNotes))
	}
	return
}

// export SKY_wallet_ReadableNotes_Load
func SKY_wallet_ReadableNotes_Load(_rns *C.ReadableNotes, _filename string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rns := (*cipher.ReadableNotes)(unsafe.Pointer(_rns))
	filename := _filename
	____return_err := rns.Load(filename)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_wallet_ReadableNotes_ToNotes
func SKY_wallet_ReadableNotes_ToNotes(_rns *C.ReadableNotes, _arg0 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rns := *(*cipher.ReadableNotes)(unsafe.Pointer(_rns))
	__arg0, ____return_err := rns.ToNotes()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	}
	return
}

// export SKY_wallet_ReadableNotes_Save
func SKY_wallet_ReadableNotes_Save(_rns *C.ReadableNotes, _filename string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	rns := (*cipher.ReadableNotes)(unsafe.Pointer(_rns))
	filename := _filename
	____return_err := rns.Save(filename)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_wallet_NewReadableNote
func SKY_wallet_NewReadableNote(_note *C.Note, _arg1 *C.ReadableNote) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	note := *(*cipher.Note)(unsafe.Pointer(_note))
	__arg1 := wallet.NewReadableNote(note)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofReadableNote))
	return
}

// export SKY_wallet_NewReadableNotesFromNotes
func SKY_wallet_NewReadableNotesFromNotes(_w *C.Notes, _arg1 *C.ReadableNotes) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	w := *(*cipher.Notes)(unsafe.Pointer(_w))
	__arg1 := wallet.NewReadableNotesFromNotes(w)
	copyToBuffer(reflect.ValueOf(__arg1[:]), unsafe.Pointer(_arg1), uint(SizeofReadableNotes))
	return
}

// export SKY_wallet_Notes_Save
func SKY_wallet_Notes_Save(_notes *C.Notes, _dir string, _fileName string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	notes := (*cipher.Notes)(unsafe.Pointer(_notes))
	dir := _dir
	fileName := _fileName
	____return_err := notes.Save(dir, fileName)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_wallet_Notes_SaveNote
func SKY_wallet_Notes_SaveNote(_notes *C.Notes, _dir string, _note *C.Note) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	notes := (*cipher.Notes)(unsafe.Pointer(_notes))
	dir := _dir
	note := *(*cipher.Note)(unsafe.Pointer(_note))
	____return_err := notes.SaveNote(dir, note)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_wallet_Notes_ToReadable
func SKY_wallet_Notes_ToReadable(_notes *C.Notes, _arg0 *C.ReadableNotes) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	notes := *(*cipher.Notes)(unsafe.Pointer(_notes))
	__arg0 := notes.ToReadable()
	copyToBuffer(reflect.ValueOf(__arg0[:]), unsafe.Pointer(_arg0), uint(SizeofReadableNotes))
	return
}

// export SKY_wallet_NotesFileExist
func SKY_wallet_NotesFileExist(_dir string, _arg1 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	dir := _dir
	__arg1, ____return_err := wallet.NotesFileExist(dir)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}

// export SKY_wallet_CreateNoteFileIfNotExist
func SKY_wallet_CreateNoteFileIfNotExist(_dir string) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	dir := _dir
	wallet.CreateNoteFileIfNotExist(dir)
	return
}
