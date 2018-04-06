package main

import (
	bucket "github.com/skycoin/skycoin/src/bucket"
	cipher "github.com/skycoin/skycoin/src/cipher"
	reflect "reflect"
	unsafe "unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_bucket_New
func SKY_bucket_New(_name *C.GoSlice_, _db *C.DB, _arg2 *C.Bucket) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	name := *(*[]byte)(unsafe.Pointer(_name))
	__arg2, ____return_err := bucket.New(name, db)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
		copyToBuffer(reflect.ValueOf((*__arg2)[:]), unsafe.Pointer(_arg2), uint(SizeofBucket))
	}
	return
}

// export SKY_bucket_Bucket_Reset
func SKY_bucket_Bucket_Reset(_b *C.Bucket) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := (*cipher.Bucket)(unsafe.Pointer(_b))
	____return_err := b.Reset()
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_bucket_Bucket_Get
func SKY_bucket_Bucket_Get(_b *C.Bucket, _key *C.GoSlice_, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Bucket)(unsafe.Pointer(_b))
	key := *(*[]byte)(unsafe.Pointer(_key))
	__arg1 := b.Get(key)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_bucket_Bucket_GetWithTx
func SKY_bucket_Bucket_GetWithTx(_b *C.Bucket, _tx *C.Tx, _key *C.GoSlice_, _arg2 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Bucket)(unsafe.Pointer(_b))
	key := *(*[]byte)(unsafe.Pointer(_key))
	__arg2 := b.GetWithTx(tx, key)
	copyToGoSlice(reflect.ValueOf(__arg2), _arg2)
	return
}

// export SKY_bucket_Bucket_GetAll
func SKY_bucket_Bucket_GetAll(_b *C.Bucket, _arg0 map[interface{}][]byte) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := (*cipher.Bucket)(unsafe.Pointer(_b))
	__arg0 := b.GetAll()
	return
}

// export SKY_bucket_Bucket_GetSlice
func SKY_bucket_Bucket_GetSlice(_b *C.Bucket, _keys *C.GoSlice_, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := (*cipher.Bucket)(unsafe.Pointer(_b))
	__arg1 := b.GetSlice(keys)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_bucket_Bucket_Put
func SKY_bucket_Bucket_Put(_b *C.Bucket, _key *C.GoSlice_, _value *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Bucket)(unsafe.Pointer(_b))
	key := *(*[]byte)(unsafe.Pointer(_key))
	value := *(*[]byte)(unsafe.Pointer(_value))
	____return_err := b.Put(key, value)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_bucket_Bucket_PutWithTx
func SKY_bucket_Bucket_PutWithTx(_b *C.Bucket, _tx *C.Tx, _key *C.GoSlice_, _value *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Bucket)(unsafe.Pointer(_b))
	key := *(*[]byte)(unsafe.Pointer(_key))
	value := *(*[]byte)(unsafe.Pointer(_value))
	____return_err := b.PutWithTx(tx, key, value)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_bucket_Bucket_Find
func SKY_bucket_Bucket_Find(_b *C.Bucket, _filter C.Handle, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := *(*cipher.Bucket)(unsafe.Pointer(_b))
	__arg1 := b.Find(filter)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_bucket_Bucket_Update
func SKY_bucket_Bucket_Update(_b *C.Bucket, _key *C.GoSlice_, _f C.Handle) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := (*cipher.Bucket)(unsafe.Pointer(_b))
	key := *(*[]byte)(unsafe.Pointer(_key))
	____return_err := b.Update(key, f)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_bucket_Bucket_Delete
func SKY_bucket_Bucket_Delete(_b *C.Bucket, _key *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := (*cipher.Bucket)(unsafe.Pointer(_b))
	key := *(*[]byte)(unsafe.Pointer(_key))
	____return_err := b.Delete(key)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_bucket_Bucket_DeleteWithTx
func SKY_bucket_Bucket_DeleteWithTx(_b *C.Bucket, _tx *C.Tx, _key *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := (*cipher.Bucket)(unsafe.Pointer(_b))
	key := *(*[]byte)(unsafe.Pointer(_key))
	____return_err := b.DeleteWithTx(tx, key)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_bucket_Bucket_RangeUpdate
func SKY_bucket_Bucket_RangeUpdate(_b *C.Bucket, _f C.Handle) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := (*cipher.Bucket)(unsafe.Pointer(_b))
	____return_err := b.RangeUpdate(f)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_bucket_Bucket_IsExist
func SKY_bucket_Bucket_IsExist(_b *C.Bucket, _k *C.GoSlice_, _arg1 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := (*cipher.Bucket)(unsafe.Pointer(_b))
	k := *(*[]byte)(unsafe.Pointer(_k))
	__arg1 := b.IsExist(k)
	*_arg1 = __arg1
	return
}

// export SKY_bucket_Bucket_IsEmpty
func SKY_bucket_Bucket_IsEmpty(_b *C.Bucket, _arg0 *bool) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := (*cipher.Bucket)(unsafe.Pointer(_b))
	__arg0 := b.IsEmpty()
	*_arg0 = __arg0
	return
}

// export SKY_bucket_Bucket_ForEach
func SKY_bucket_Bucket_ForEach(_b *C.Bucket, _f C.Handle) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := (*cipher.Bucket)(unsafe.Pointer(_b))
	____return_err := b.ForEach(f)
	____return_var = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}

// export SKY_bucket_Bucket_Len
func SKY_bucket_Bucket_Len(_b *C.Bucket, _arg0 *int) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	b := (*cipher.Bucket)(unsafe.Pointer(_b))
	__arg0 := b.Len()
	*_arg0 = __arg0
	return
}

// export SKY_bucket_Itob
func SKY_bucket_Itob(_v uint64, _arg1 *C.GoSlice_) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	v := _v
	__arg1 := bucket.Itob(v)
	copyToGoSlice(reflect.ValueOf(__arg1), _arg1)
	return
}

// export SKY_bucket_Btoi
func SKY_bucket_Btoi(_v *C.GoSlice_, _arg1 *uint64) (____return_var uint32) {
	____return_var = 0
	defer func() {
		____return_var = catchApiPanic(recover())
	}()
	v := *(*[]byte)(unsafe.Pointer(_v))
	__arg1 := bucket.Btoi(v)
	*_arg1 = __arg1
	return
}
