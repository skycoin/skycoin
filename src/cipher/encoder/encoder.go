// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package encoder binary implements translation between struct data and byte sequences
//
// Fields can be ignored with the struct tag `enc:"-"` .
// Unexported struct fields are ignored by default .
//
// Fields can be skipped if empty with the struct tag `enc:",omitempty"`
// Note the comma, which follows package json's conventions.
// Only Slice, Map and String types recognize the omitempty tag.
// When omitempty is set, the no data will be written if the value is empty.
// If the value is empty and omitempty is not set, then a length prefix with value 0 would be written.
// omitempty can only be used for the last field in the struct
//
// Encoding of maps is supported, but note that the use of them results in non-deterministic output.
// If determinism is required, do not use map.
package encoder

import (
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"strings"
)

/*
Todo:
- ensure that invalid input from foreign server cannot crash
- validate packet legnth for incoming
*/

var (
	// ErrBufferUnderflow bytes in input buffer not enough to deserialize expected type
	ErrBufferUnderflow = errors.New("Not enough buffer data to deserialize")
	// ErrInvalidOmitEmpty field tagged with omitempty and it's not last one in struct
	ErrInvalidOmitEmpty = errors.New("omitempty only supported for the final field in the struct")
)

// TODO: constant length byte arrays must not be prefixed

// EncodeInt encodes an Integer type contained in `data`
// into buffer `b`. If `data` is not an Integer type,
// panic message is logged.
func EncodeInt(b []byte, data interface{}) {
	//var b [8]byte
	var bs []byte
	switch v := data.(type) {

	case int8:
		// bs = b[:1]
		b[0] = byte(v)
	case uint8:
		// bs = b[:1]
		b[0] = v
	case int16:
		bs = b[:2]
		lePutUint16(bs, uint16(v))
	case uint16:
		bs = b[:2]
		lePutUint16(bs, v)
	case int32:
		bs = b[:4]
		lePutUint32(bs, uint32(v))
	case uint32:
		bs = b[:4]
		lePutUint32(bs, v)
	case int64:
		bs = b[:8]
		lePutUint64(bs, uint64(v))
	case uint64:
		bs = b[:8]
		lePutUint64(bs, v)
	default:
		log.Panic("PushAtomic, case not handled")
	}
}

// DecodeInt decodes `in` buffer into `data` parameter.
// If `data` is not an Integer type, panic message is logged.
func DecodeInt(in []byte, data interface{}) {

	n := intDestSize(data)
	if len(in) < n {
		log.Panic()
	}
	if n != 0 {
		var b [8]byte
		copy(b[0:n], in[0:n])
		bs := b[:n]

		switch v := data.(type) {
		case *int8:
			*v = int8(b[0])
		case *uint8:
			*v = b[0]
		case *int16:
			*v = int16(leUint16(bs))
		case *uint16:
			*v = leUint16(bs)
		case *int32:
			*v = int32(leUint32(bs))
		case *uint32:
			*v = leUint32(bs)
		case *int64:
			*v = int64(leUint64(bs))
		case *uint64:
			*v = leUint64(bs)
		default:
			//FIX: this does not get triggered on invalid type in
			// pass in struct on unit test
			log.Panic("PopAtomic, case not handled")

		}

	}
}

// DeserializeAtomic deserializes `in` buffer into `data`
// parameter. If `data` is not an atomic type
// (i.e., Integer type or Boolean type), panic message is logged.
func DeserializeAtomic(in []byte, data interface{}) {
	n := intDestSize(data)
	if len(in) < n {
		log.Panic(ErrBufferUnderflow)
	}
	if n != 0 {
		var b [8]byte
		copy(b[0:n], in[0:n])
		bs := b[:n]

		switch v := data.(type) {
		case *bool:
			if b[0] == 1 {
				*v = true
			} else {
				*v = false
			}
		case *int8:
			*v = int8(b[0])
		case *uint8:
			*v = b[0]
		case *int16:
			*v = int16(leUint16(bs))
		case *uint16:
			*v = leUint16(bs)
		case *int32:
			*v = int32(leUint32(bs))
		case *uint32:
			*v = leUint32(bs)
		case *int64:
			*v = int64(leUint64(bs))
		case *uint64:
			*v = leUint64(bs)
		default:
			//FIX: this does not get triggered on invalid type in
			// pass in struct on unit test
			log.Panic("type not atomic")
		}
	}
}

// DeserializeRaw deserializes `in` buffer into return
// parameter. If `data` is not either a Pointer type,
// a Slice type or a Struct type, an error message
// is returned. If `in` buffer can't be deserialized,
// an error message is returned.
func DeserializeRaw(in []byte, data interface{}) error {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
	case reflect.Slice:
	case reflect.Struct:
	case reflect.Map:
	default:
		return fmt.Errorf("Invalid type %s", reflect.TypeOf(v).String())
	}

	d1 := &decoder{buf: make([]byte, len(in))}
	copy(d1.buf, in)

	return d1.value(v)
}

// CanDeserialize returns true if `in` buffer can be
// deserialized into `dst`'s type. Returns false in any
// other case.
func CanDeserialize(in []byte, dst reflect.Value) bool {
	d1 := &decoder{buf: make([]byte, len(in))}
	copy(d1.buf, in)
	return d1.dchk(dst) == 0
}

// DeserializeRawToValue deserializes `in` buffer into
// `dst`'s type and returns the number of bytes used and
// the value of the buffer. If `data` is not either a
// Pointer type, a Slice type or a Struct type, 0 and an error
// message are returned. If `in` buffer can't be deserialized, 0 and
// an error message are returned.
func DeserializeRawToValue(in []byte, dst reflect.Value) (int, error) {
	var v reflect.Value
	switch dst.Kind() {
	case reflect.Ptr:
		v = dst.Elem()
	case reflect.Slice:
		v = dst
	case reflect.Struct:
	default:
		return 0, errors.New("binary.Read: invalid type " + reflect.TypeOf(dst).String())
	}

	inlen := len(in)
	d1 := &decoder{buf: make([]byte, inlen)}
	copy(d1.buf, in)

	err := d1.value(v)
	return inlen - len(d1.buf), err
}

// SerializeAtomic returns serialization of `data`
// parameter. If `data` is not an atomic type, panic message is logged.
func SerializeAtomic(data interface{}) []byte {
	var b [8]byte
	var bs []byte
	switch v := data.(type) {
	case *bool:
		bs = b[:1]
		if *v {
			b[0] = 1
		} else {
			b[0] = 0
		}
	case bool:
		bs = b[:1]
		if v {
			b[0] = 1
		} else {
			b[0] = 0
		}
	case *int8:
		bs = b[:1]
		b[0] = byte(*v)
	case int8:
		bs = b[:1]
		b[0] = byte(v)
	case *uint8:
		bs = b[:1]
		b[0] = *v
	case uint8:
		bs = b[:1]
		b[0] = v
	case *int16:
		bs = b[:2]
		lePutUint16(bs, uint16(*v))
	case int16:
		bs = b[:2]
		lePutUint16(bs, uint16(v))
	case *uint16:
		bs = b[:2]
		lePutUint16(bs, *v)
	case uint16:
		bs = b[:2]
		lePutUint16(bs, v)
	case *int32:
		bs = b[:4]
		lePutUint32(bs, uint32(*v))
	case int32:
		bs = b[:4]
		lePutUint32(bs, uint32(v))
	case *uint32:
		bs = b[:4]
		lePutUint32(bs, *v)
	case uint32:
		bs = b[:4]
		lePutUint32(bs, v)
	case *int64:
		bs = b[:8]
		lePutUint64(bs, uint64(*v))
	case int64:
		bs = b[:8]
		lePutUint64(bs, uint64(v))
	case *uint64:
		bs = b[:8]
		lePutUint64(bs, *v)
	case uint64:
		bs = b[:8]
		lePutUint64(bs, v)
	default:
		log.Panic("type not atomic")
	}
	return bs
}

// Serialize returns serialized basic type-based `data`
// parameter. Encoding is reflect-based.
func Serialize(data interface{}) []byte {
	// Fast path for basic types.
	// Fallback to reflect-based encoding.
	v := reflect.Indirect(reflect.ValueOf(data))
	size, err := datasizeWrite(v)
	if err != nil {
		log.Panic(err)
	}
	buf := make([]byte, size)
	e := &encoder{buf: buf}
	e.value(v)
	return buf
}

// Size returns how many bytes would it take to encode the
// value v, which must be a fixed-size value (struct) or a
// slice of fixed-size values, or a pointer to such data.
// Reflect-based encoding is used.
func Size(v interface{}) int {
	n, err := datasizeWrite(reflect.Indirect(reflect.ValueOf(v)))
	if err != nil {
		return -1
	}
	return n
}

// isEmpty returns true if a value is "empty".
// Only supports Slice, Map and String.
// All other values are never considered empty.
func isEmpty(v reflect.Value) bool {
	t := v.Type()
	switch t.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Map:
		return v.IsNil() || v.Len() == 0
	case reflect.Slice:
		return v.IsNil() || v.Len() == 0
	default:
		return false
	}
}

// datasizeWrite returns the number of bytes the actual data represented by v occupies in memory.
// For compound structures, it sums the sizes of the elements. Thus, for instance, for a slice
// it returns the length of the slice times the element size and does not count the memory
// occupied by the header.
func datasizeWrite(v reflect.Value) (int, error) {
	t := v.Type()
	switch t.Kind() {
	case reflect.Interface:
		return datasizeWrite(v.Elem())

	case reflect.Array:
		// Arrays are a fixed size, so the length is not written
		t := v.Type()
		elem := t.Elem()
		switch elem.Kind() {
		case reflect.Uint8, reflect.Int8:
			return v.Len(), nil
		case reflect.Uint16, reflect.Int16:
			return v.Len() * 2, nil
		case reflect.Uint32, reflect.Int32, reflect.Float32:
			return v.Len() * 4, nil
		case reflect.Uint64, reflect.Int64, reflect.Float64:
			return v.Len() * 8, nil
		default:
			size := 0
			for i := 0; i < v.Len(); i++ {
				elem := v.Index(i)
				s, err := datasizeWrite(elem)
				if err != nil {
					return 0, err
				}
				size += s
			}
			return size, nil
		}

	case reflect.Slice:
		t := v.Type()
		elem := t.Elem()
		switch elem.Kind() {
		case reflect.Uint8, reflect.Int8:
			return 4 + v.Len(), nil
		case reflect.Uint16, reflect.Int16:
			return 4 + v.Len()*2, nil
		case reflect.Uint32, reflect.Int32, reflect.Float32:
			return 4 + v.Len()*4, nil
		case reflect.Uint64, reflect.Int64, reflect.Float64:
			return 4 + v.Len()*8, nil
		default:
			size := 0
			for i := 0; i < v.Len(); i++ {
				elem := v.Index(i)
				s, err := datasizeWrite(elem)
				if err != nil {
					return 0, err
				}
				size += s
			}
			return 4 + size, nil
		}

	case reflect.Map:
		// length prefix
		size := 4
		for _, key := range v.MapKeys() {
			s, err := datasizeWrite(key)
			if err != nil {
				return 0, err
			}
			size += s
			elem := v.MapIndex(key)
			s, err = datasizeWrite(elem)
			if err != nil {
				return 0, err
			}
			size += s
		}
		return size, nil

	case reflect.Struct:
		sum := 0
		nFields := t.NumField()
		for i, n := 0, nFields; i < n; i++ {
			ff := t.Field(i)
			// Skip unexported fields
			if ff.PkgPath != "" {
				continue
			}

			tag, omitempty := ParseTag(ff.Tag.Get("enc"))

			if omitempty && i != nFields-1 {
				log.Panic(ErrInvalidOmitEmpty)
			}

			if tag != "-" {
				fv := v.Field(i)
				if !omitempty || !isEmpty(fv) {
					s, err := datasizeWrite(fv)
					if err != nil {
						return 0, err
					}
					sum += s
				}
			}
		}
		return sum, nil

	case reflect.Bool:
		return 1, nil

	case reflect.String:
		return v.Len() + 4, nil

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64:
		return int(t.Size()), nil

	default:
		return 0, fmt.Errorf("invalid type %s", t.String())
	}
}

// ParseTag to extract encoder args from raw string. Returns the tag name and if omitempty was specified
func ParseTag(tag string) (string, bool) {
	commaIndex := strings.Index(tag, ",")
	if commaIndex == -1 {
		return tag, false
	}

	if tag[commaIndex+1:] == "omitempty" {
		return tag[:commaIndex], true
	}

	return tag[:commaIndex], false
}

/*
	Internals
*/

func leUint16(b []byte) uint16 { return uint16(b[0]) | uint16(b[1])<<8 }

func lePutUint16(b []byte, v uint16) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}

func leUint32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func lePutUint32(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

func leUint64(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func lePutUint64(b []byte, v uint64) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
}

type coder struct {
	buf []byte // nolint: structcheck
}

type decoder coder
type encoder coder

func (d *decoder) bool() (bool, error) {
	if len(d.buf) < 1 {
		return false, ErrBufferUnderflow
	}
	x := d.buf[0]
	d.buf = d.buf[1:] // advance slice
	return x != 0, nil
}

func (e *encoder) bool(x bool) {
	if x {
		e.buf[0] = 1
	} else {
		e.buf[0] = 0
	}
	e.buf = e.buf[1:]
}

func (d decoder) string() (string, error) { // nolint: unused,megacheck
	ul, err := d.uint32() // pop length
	if err != nil {
		return "", err
	}
	l := int(ul)
	t := d.buf[:l]
	d.buf = d.buf[l:]
	return string(t), nil
}

func (e encoder) string(xs string) { // nolint: unused,megacheck
	x := []byte(xs)
	l := len(x)
	for i := 0; i < l; i++ {
		e.buf[i] = x[i]
	} // memcpy
	e.buf = e.buf[l:] // advance slice l bytes
}

func (d *decoder) uint8() (uint8, error) {
	if len(d.buf) < 1 {
		return 0, ErrBufferUnderflow
	}

	x := d.buf[0]
	d.buf = d.buf[1:] // advance slice
	return x, nil
}

func (e *encoder) uint8(x uint8) {
	e.buf[0] = x
	e.buf = e.buf[1:]
}

func (d *decoder) uint16() (uint16, error) {
	if len(d.buf) < 2 {
		return 0, ErrBufferUnderflow
	}

	x := leUint16(d.buf[0:2])
	d.buf = d.buf[2:]
	return x, nil
}

func (e *encoder) uint16(x uint16) {
	lePutUint16(e.buf[0:2], x)
	e.buf = e.buf[2:]
}

func (d *decoder) uint32() (uint32, error) {
	if len(d.buf) < 4 {
		return 0, ErrBufferUnderflow
	}

	x := leUint32(d.buf[0:4])
	d.buf = d.buf[4:]
	return x, nil
}

func (e *encoder) uint32(x uint32) {
	lePutUint32(e.buf[0:4], x)
	e.buf = e.buf[4:]
}

func (d *decoder) uint64() (uint64, error) {
	if len(d.buf) < 8 {
		return 0, ErrBufferUnderflow
	}

	x := leUint64(d.buf[0:8])
	d.buf = d.buf[8:]
	return x, nil
}

func (e *encoder) uint64(x uint64) {
	lePutUint64(e.buf[0:8], x)
	e.buf = e.buf[8:]
}

func (e *encoder) bytes(x []byte) {
	e.uint32(uint32(len(x)))
	copy(e.buf, x)
	e.buf = e.buf[len(x):]
}

func (d *decoder) int8() (int8, error) {
	u, err := d.uint8()
	if err != nil {
		return 0, err
	}

	return int8(u), nil
}

func (e *encoder) int8(x int8) { e.uint8(uint8(x)) }

func (d *decoder) int16() (int16, error) {
	u, err := d.uint16()
	if err != nil {
		return 0, err
	}

	return int16(u), nil
}

func (e *encoder) int16(x int16) { e.uint16(uint16(x)) }

func (d *decoder) int32() (int32, error) {
	u, err := d.uint32()
	if err != nil {
		return 0, err
	}

	return int32(u), nil
}

func (e *encoder) int32(x int32) { e.uint32(uint32(x)) }

func (d *decoder) int64() (int64, error) {
	u, err := d.uint64()
	if err != nil {
		return 0, err
	}

	return int64(u), nil
}

func (e *encoder) int64(x int64) { e.uint64(uint64(x)) }

func (d *decoder) value(v reflect.Value) error {
	kind := v.Kind()
	switch kind {
	case reflect.Array:

		t := v.Type()
		elem := t.Elem()

		// Arrays are a fixed size, so the length is not written
		length := v.Len()

		switch elem.Kind() {
		case reflect.Uint8:
			if length > len(d.buf) {
				return ErrBufferUnderflow
			}

			reflect.Copy(v, reflect.ValueOf(d.buf[:length]))
			d.buf = d.buf[length:]
		default:
			for i := 0; i < length; i++ {
				if err := d.value(v.Index(i)); err != nil {
					return err
				}
			}
		}

	case reflect.Map:
		if len(d.buf) < 4 {
			return ErrBufferUnderflow
		}

		ul, err := d.uint32()
		if err != nil {
			return err
		}

		length := int(ul)
		if length < 0 || length > len(d.buf) {
			return ErrBufferUnderflow
		}

		t := v.Type()
		key := t.Key()
		elem := t.Elem()

		if v.IsNil() {
			v.Set(reflect.Indirect(reflect.MakeMap(t)))
		}

		for i := 0; i < length; i++ {
			keyv := reflect.Indirect(reflect.New(key))
			elemv := reflect.Indirect(reflect.New(elem))
			if err := d.value(keyv); err != nil {
				return err
			}
			if err := d.value(elemv); err != nil {
				return err
			}
			v.SetMapIndex(keyv, elemv)
		}

	case reflect.Slice:
		if len(d.buf) < 4 {
			return ErrBufferUnderflow
		}

		ul, err := d.uint32()
		if err != nil {
			return err
		}

		length := int(ul)
		if length < 0 || length > len(d.buf) {
			return ErrBufferUnderflow
		}

		if length == 0 {
			return nil
		}

		t := v.Type()
		elem := t.Elem()

		switch elem.Kind() {
		case reflect.Uint8:
			v.SetBytes(d.buf[:length])
			d.buf = d.buf[length:]
		default:
			elemvs := reflect.MakeSlice(t, length, length)
			for i := 0; i < length; i++ {
				elemv := reflect.Indirect(elemvs.Index(i))
				if err := d.value(elemv); err != nil {
					return err
				}
			}
			v.Set(elemvs)
		}

	case reflect.Struct:
		t := v.Type()
		nFields := v.NumField()
		for i := 0; i < nFields; i++ {
			ff := t.Field(i)
			// Skip unexported fields
			if ff.PkgPath != "" {
				continue
			}

			tag, omitempty := ParseTag(ff.Tag.Get("enc"))

			if omitempty && i != nFields-1 {
				log.Panic(ErrInvalidOmitEmpty)
			}

			if tag != "-" {
				fv := v.Field(i)
				if fv.CanSet() && ff.Name != "_" {
					if err := d.value(fv); err != nil {
						// omitempty fields at the end of the buffer are ignored
						if !(omitempty && len(d.buf) == 0) {
							return err
						}
					}
				}
			}
		}

	case reflect.String:
		if len(d.buf) < 4 {
			return ErrBufferUnderflow
		}

		ul, err := d.uint32()
		if err != nil {
			return err
		}

		length := int(ul)
		if length < 0 || length > len(d.buf) {
			return ErrBufferUnderflow
		}

		v.SetString(string(d.buf[:length]))
		d.buf = d.buf[length:]

	case reflect.Bool:
		b, err := d.bool()
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Int8:
		i, err := d.int8()
		if err != nil {
			return err
		}
		v.SetInt(int64(i))
	case reflect.Int16:
		i, err := d.int16()
		if err != nil {
			return err
		}
		v.SetInt(int64(i))
	case reflect.Int32:
		i, err := d.int32()
		if err != nil {
			return err
		}
		v.SetInt(int64(i))
	case reflect.Int64:
		i, err := d.int64()
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Uint8:
		u, err := d.uint8()
		if err != nil {
			return err
		}
		v.SetUint(uint64(u))
	case reflect.Uint16:
		u, err := d.uint16()
		if err != nil {
			return err
		}
		v.SetUint(uint64(u))
	case reflect.Uint32:
		u, err := d.uint32()
		if err != nil {
			return err
		}
		v.SetUint(uint64(u))
	case reflect.Uint64:
		u, err := d.uint64()
		if err != nil {
			return err
		}
		v.SetUint(u)

	case reflect.Float32:
		u, err := d.uint32()
		if err != nil {
			return err
		}
		v.SetFloat(float64(math.Float32frombits(u)))
	case reflect.Float64:
		u, err := d.uint64()
		if err != nil {
			return err
		}
		v.SetFloat(math.Float64frombits(u))

	default:
		log.Panicf("Decode error: kind %s not handled", v.Kind().String())
	}

	return nil
}

// advance, returns -1 on failure
//returns 0 on success
func (d *decoder) adv(n int) int {
	if n > len(d.buf) {
		n = len(d.buf)
		d.buf = d.buf[n:]
		return -1
	}
	d.buf = d.buf[n:]
	return 0
}

//recursive size
func (d *decoder) dchk(v reflect.Value) int {
	kind := v.Kind()
	switch kind {
	case reflect.Array:
		// Arrays are a fixed size, so the length is not written
		for i := 0; i < v.Len(); i++ {
			if d.dchk(v.Index(i)) < 0 {
				return -1
			}
		}
		return 0

	case reflect.Map:
		if len(d.buf) < 4 {
			return -1 //error
		}

		length := int(leUint32(d.buf[0:4]))
		if d.adv(4) < 0 {
			return -1
		}

		key := v.Type().Key()
		elem := v.Type().Elem()

		for i := 0; i < length; i++ {
			keyv := reflect.Indirect(reflect.New(key))
			elemv := reflect.Indirect(reflect.New(elem))

			if d.dchk(keyv) < 0 {
				return -1
			}

			if d.dchk(elemv) < 0 {
				return -1
			}
		}
		return 0

	case reflect.Slice:
		if len(d.buf) < 4 {
			return -1
		}

		length := int(leUint32(d.buf[0:4]))
		if d.adv(4) < 0 {
			return -1
		}

		if length < 0 || length > len(d.buf) {
			return -1
		}

		elem := v.Type().Elem()
		if elem.Kind() == reflect.Uint8 {
			return d.adv(length)
		}

		for i := 0; i < length; i++ {
			elemv := reflect.Indirect(reflect.New(elem))

			if d.dchk(elemv) < 0 {
				return -1
			}
		}
		return 0

	case reflect.Struct:
		t := v.Type()
		nFields := v.NumField()
		for i := 0; i < nFields; i++ {
			ff := t.Field(i)
			// Skip unexported fields
			if ff.PkgPath != "" {
				continue
			}

			tag, omitempty := ParseTag(ff.Tag.Get("enc"))

			if omitempty && i != nFields-1 {
				log.Panic(ErrInvalidOmitEmpty)
			}

			if tag != "-" {
				fv := v.Field(i)
				if !omitempty && fv.CanSet() && ff.Name != "_" {
					if d.dchk(fv) < 0 {
						return -1
					}
				}
			}
		}
		return 0

	case reflect.String:
		if len(d.buf) < 4 {
			return -1
		}

		length := int(leUint32(d.buf[0:4]))
		if d.adv(4) < 0 {
			return -1
		}

		return d.adv(length)

	case reflect.Bool:
		return d.adv(1)
	case reflect.Int8:
		return d.adv(1)
	case reflect.Int16:
		return d.adv(2)
	case reflect.Int32:
		return d.adv(4)
	case reflect.Int64:
		return d.adv(8)

	case reflect.Uint8:
		return d.adv(1)
	case reflect.Uint16:
		return d.adv(2)
	case reflect.Uint32:
		return d.adv(4)
	case reflect.Uint64:
		return d.adv(8)

	case reflect.Float32:
		return d.adv(4)
	case reflect.Float64:
		return d.adv(8)

	default:
		log.Panicf("Decode error: kind %s not handled", v.Kind().String())
	}

	log.Panic()
	return 0
}

func (e *encoder) value(v reflect.Value) {
	switch v.Kind() {
	case reflect.Interface:
		e.value(v.Elem())

	case reflect.Array:
		// Arrays are a fixed size, so the length is not written
		t := v.Type()
		elem := t.Elem()
		switch elem.Kind() {
		case reflect.Uint8:
			reflect.Copy(reflect.ValueOf(e.buf), v)
			e.buf = e.buf[v.Len():]
		default:
			for i := 0; i < v.Len(); i++ {
				e.value(v.Index(i))
			}
		}

	case reflect.Slice:
		t := v.Type()
		elem := t.Elem()
		switch elem.Kind() {
		case reflect.Uint8:
			e.bytes(v.Bytes())
		default:
			e.uint32(uint32(v.Len()))
			for i := 0; i < v.Len(); i++ {
				e.value(v.Index(i))
			}
		}

	case reflect.Map:
		e.uint32(uint32(v.Len()))
		for _, key := range v.MapKeys() {
			e.value(key)
			e.value(v.MapIndex(key))
		}

	case reflect.Struct:
		t := v.Type()
		nFields := v.NumField()
		for i := 0; i < nFields; i++ {
			// see comment for corresponding code in decoder.value()
			ff := t.Field(i)
			// Skip unexported fields
			if ff.PkgPath != "" {
				continue
			}

			tag, omitempty := ParseTag(ff.Tag.Get("enc"))

			if omitempty && i != nFields-1 {
				log.Panic(ErrInvalidOmitEmpty)
			}

			if tag != "-" {
				fv := v.Field(i)
				if !(omitempty && isEmpty(fv)) && (fv.CanSet() || ff.Name != "_") {
					e.value(fv)
				}
			}
		}

	case reflect.Bool:
		e.bool(v.Bool())

	case reflect.String:
		e.bytes([]byte(v.String()))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v.Type().Kind() {
		case reflect.Int8:
			e.int8(int8(v.Int()))
		case reflect.Int16:
			e.int16(int16(v.Int()))
		case reflect.Int32:
			e.int32(int32(v.Int()))
		case reflect.Int64:
			e.int64(v.Int())
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch v.Type().Kind() {
		case reflect.Uint8:
			e.uint8(uint8(v.Uint()))
		case reflect.Uint16:
			e.uint16(uint16(v.Uint()))
		case reflect.Uint32:
			e.uint32(uint32(v.Uint()))
		case reflect.Uint64:
			e.uint64(v.Uint())
		}

	case reflect.Float32, reflect.Float64:
		switch v.Type().Kind() {
		case reflect.Float32:
			e.uint32(math.Float32bits(float32(v.Float())))
		case reflect.Float64:
			e.uint64(math.Float64bits(v.Float()))
		}

	default:
		log.Panic("Encoding unhandled type " + v.Type().Name())

	}

}

func (d *decoder) skip(v reflect.Value) { // nolint: unused,megacheck
	n, err := datasizeWrite(v)
	if err != nil {
		panic(err)
	}
	d.buf = d.buf[n:]
}

//skip with byte size return
/*
func (d *decoder) skipn(v reflect.Value) int {
    n := intDestSize(&v)
    if n == 0 {
        log.Panic()
    }
    d.buf = d.buf[n:]
    return n
}
*/
func (e *encoder) skip(v reflect.Value) { // nolint: unused,megacheck
	n, err := datasizeWrite(v)
	if err != nil {
		panic(err)
	}
	for i := range e.buf[0:n] {
		e.buf[i] = 0
	}
	e.buf = e.buf[n:]
}

// intDestSize returns the size of the integer that ptrType points to,
// or 0 if the type is not supported.
func intDestSize(ptrType interface{}) int {
	switch ptrType.(type) {
	case *bool:
		return 1
	case *int8, *uint8:
		return 1
	case *int16, *uint16:
		return 2
	case *int32, *uint32:
		return 4
	case *int64, *uint64:
		return 8
	}
	return 0
}
