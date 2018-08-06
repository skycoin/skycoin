package main

import (
	"unsafe"

	wallet "github.com/skycoin/skycoin/src/wallet"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_wallet_NewNotesFilename
func SKY_wallet_NewNotesFilename(_arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := wallet.NewNotesFilename()
	copyString(__arg0, _arg0)
	return
}

//export SKY_wallet_LoadNotes
func SKY_wallet_LoadNotes(_dir string, _arg1 *C.WalletNotes_Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	dir := _dir
	__arg1, ____return_err := wallet.LoadNotes(dir)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerWalletNotesHandle(&__arg1)
	}
	return
}

//export SKY_wallet_LoadReadableNotes
func SKY_wallet_LoadReadableNotes(_filename string, _arg1 *C.WalletReadableNotes_Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	filename := _filename
	__arg1, ____return_err := wallet.LoadReadableNotes(filename)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = registerWalletReadableNotesHandle(__arg1)
	}
	return
}

//export SKY_wallet_ReadableNotes_Load
func SKY_wallet_ReadableNotes_Load(_rns C.WalletReadableNotes_Handle, _filename string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	rns, ok := lookupWalletReadableNotesHandle(_rns)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	filename := _filename
	____return_err := rns.Load(filename)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_wallet_ReadableNotes_ToNotes
func SKY_wallet_ReadableNotes_ToNotes(_rns C.WalletReadableNotes_Handle, _arg0 *C.WalletNotes_Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	rns, ok := lookupWalletReadableNotesHandle(_rns)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	__arg0, ____return_err := rns.ToNotes()
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		notes := wallet.Notes(__arg0)
		*_arg0 = registerWalletNotesHandle(&notes)
	}
	return
}

//export SKY_wallet_ReadableNotes_Save
func SKY_wallet_ReadableNotes_Save(_rns C.WalletReadableNotes_Handle, _filename string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	rns, ok := lookupWalletReadableNotesHandle(_rns)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	filename := _filename
	____return_err := rns.Save(filename)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_wallet_NewReadableNote
func SKY_wallet_NewReadableNote(_note *C.wallet__Note, _arg1 *C.wallet__ReadableNote) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	note := *(*wallet.Note)(unsafe.Pointer(_note))
	__arg1 := wallet.NewReadableNote(note)
	copyString(__arg1.TransactionID, &_arg1.TransactionID)
	copyString(__arg1.ActualNote, &_arg1.ActualNote)
	return
}

//export SKY_wallet_NewReadableNotesFromNotes
func SKY_wallet_NewReadableNotesFromNotes(_w C.WalletNotes_Handle, _arg1 *C.WalletReadableNotes_Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	w, ok := lookupWalletNotesHandle(_w)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	__arg1 := wallet.NewReadableNotesFromNotes(*w)
	*_arg1 = registerWalletReadableNotesHandle(&__arg1)
	return
}

//export SKY_wallet_Notes_Save
func SKY_wallet_Notes_Save(_notes C.WalletNotes_Handle, _dir string, _fileName string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	notes, ok := lookupWalletNotesHandle(_notes)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	dir := _dir
	fileName := _fileName
	____return_err := notes.Save(dir, fileName)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_wallet_Notes_SaveNote
func SKY_wallet_Notes_SaveNote(_notes C.WalletNotes_Handle, _dir string, _note *C.wallet__Note) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	notes, ok := lookupWalletNotesHandle(_notes)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	dir := _dir
	note := *(*wallet.Note)(unsafe.Pointer(_note))
	____return_err := notes.SaveNote(dir, note)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

//export SKY_wallet_Notes_ToReadable
func SKY_wallet_Notes_ToReadable(_notes C.WalletNotes_Handle, _arg0 *C.WalletReadableNotes_Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	notes, ok := lookupWalletNotesHandle(_notes)
	if !ok {
		____error_code = SKY_BAD_HANDLE
		return
	}
	__arg0 := notes.ToReadable()
	*_arg0 = registerWalletReadableNotesHandle(&__arg0)
	return
}

//export SKY_wallet_NotesFileExist
func SKY_wallet_NotesFileExist(_dir string, _arg1 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	dir := _dir
	__arg1, ____return_err := wallet.NotesFileExist(dir)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = __arg1
	}
	return
}

//export SKY_wallet_CreateNoteFileIfNotExist
func SKY_wallet_CreateNoteFileIfNotExist(_dir string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	dir := _dir
	wallet.CreateNoteFileIfNotExist(dir)
	return
}
