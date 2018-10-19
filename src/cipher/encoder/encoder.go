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
//
// A length restriction to certain fields can be applied when decoding.
// Use the tag `,maxlen=` on a struct field to apply this restriction.
// `maxlen` works for string and slice types. The length is interpreted as the length
// of the string or the number of elements in the slice.
// Note that maxlen does not affect serialization; it may serialize objects which could fail deserialization.
// Callers should check their length restricted values manually prior to serialization.
package encoder

import (
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
)

var (
	// ErrBufferUnderflow bytes in input buffer not enough to deserialize expected type
	ErrBufferUnderflow = errors.New("Not enough buffer data to deserialize")
	// ErrBufferOverflow bytes in output buffer not enough to serialize expected type
	ErrBufferOverflow = errors.New("Not enough buffer data to serialize")
	// ErrInvalidOmitEmpty field tagged with omitempty and it's not last one in struct
	ErrInvalidOmitEmpty = errors.New("omitempty only supported for the final field in the struct")
	// ErrRemainingBytes bytes remain in buffer after deserializing object
	ErrRemainingBytes = errors.New("Bytes remain in buffer after deserializing object")
	// ErrMaxLenExceeded a specified maximum length was exceeded when serializing or deserializing a variable length field
	ErrMaxLenExceeded = errors.New("Maximum length exceeded for variable length field")
)

// SerializeAtomic encoder an integer or boolean contained in `data` to bytes.
// Panics if `data` is not an integer or boolean type.
func SerializeAtomic(data interface{}) []byte {
	var b [8]byte

	switch v := data.(type) {
	case bool:
		if v {
			b[0] = 1
		} else {
			b[0] = 0
		}
		return b[:1]
	case int8:
		b[0] = byte(v)
		return b[:1]
	case uint8:
		b[0] = v
		return b[:1]
	case int16:
		lePutUint16(b[:2], uint16(v))
		return b[:2]
	case uint16:
		lePutUint16(b[:2], v)
		return b[:2]
	case int32:
		lePutUint32(b[:4], uint32(v))
		return b[:4]
	case uint32:
		lePutUint32(b[:4], v)
		return b[:4]
	case int64:
		lePutUint64(b[:8], uint64(v))
		return b[:8]
	case uint64:
		lePutUint64(b[:8], v)
		return b[:8]
	default:
		log.Panic("SerializeAtomic unhandled type")
		return nil
	}
}

// DeserializeAtomic deserializes `in` buffer into `data`
// parameter. Panics if `data` is not an integer or boolean type.
// Returns the number of bytes read.
func DeserializeAtomic(in []byte, data interface{}) (int, error) {
	switch v := data.(type) {
	case *bool:
		if len(in) < 1 {
			return 0, ErrBufferUnderflow
		}
		if in[0] == 0 {
			*v = false
		} else {
			*v = true
		}
		return 1, nil
	case *int8:
		if len(in) < 1 {
			return 0, ErrBufferUnderflow
		}
		*v = int8(in[0])
		return 1, nil
	case *uint8:
		if len(in) < 1 {
			return 0, ErrBufferUnderflow
		}
		*v = in[0]
		return 1, nil
	case *int16:
		if len(in) < 2 {
			return 0, ErrBufferUnderflow
		}
		*v = int16(leUint16(in[:2]))
		return 2, nil
	case *uint16:
		if len(in) < 2 {
			return 0, ErrBufferUnderflow
		}
		*v = leUint16(in[:2])
		return 2, nil
	case *int32:
		if len(in) < 4 {
			return 0, ErrBufferUnderflow
		}
		*v = int32(leUint32(in[:4]))
		return 4, nil
	case *uint32:
		if len(in) < 4 {
			return 0, ErrBufferUnderflow
		}
		*v = leUint32(in[:4])
		return 4, nil
	case *int64:
		if len(in) < 8 {
			return 0, ErrBufferUnderflow
		}
		*v = int64(leUint64(in[:8]))
		return 8, nil
	case *uint64:
		if len(in) < 8 {
			return 0, ErrBufferUnderflow
		}
		*v = leUint64(in[:8])
		return 8, nil
	default:
		log.Panic("DeserializeAtomic unhandled type")
		return 0, nil
	}
}

// SerializeString serializes a string to []byte
func SerializeString(s string) []byte {
	v := reflect.ValueOf(s)
	size, err := datasizeWrite(v)
	if err != nil {
		log.Panic(err)
	}
	buf := make([]byte, size)
	e := &encoder{buf: buf}
	e.value(v)
	return buf
}

// DeserializeString deserializes a string from []byte, returning the string and the number of bytes read
func DeserializeString(in []byte, maxlen int) (string, int, error) {
	var s string
	v := reflect.ValueOf(&s)
	v = v.Elem()

	inlen := len(in)
	d1 := &decoder{buf: make([]byte, inlen)}
	copy(d1.buf, in)

	err := d1.value(v, maxlen)
	if err != nil {
		return "", 0, err
	}

	return s, inlen - len(d1.buf), nil
}

// DeserializeRaw deserializes `in` buffer into return
// parameter. If `data` is not a Pointer or Map type an error
// is returned. If `in` buffer can't be deserialized,
// an error message is returned. If there are remaining
// bytes in `in` after decoding to data, ErrRemainingBytes is returned.
func DeserializeRaw(in []byte, data interface{}) error {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
	case reflect.Map:
	default:
		return fmt.Errorf("DeserializeRaw value must be a ptr, is %s", v.Kind().String())
	}

	d1 := &decoder{buf: make([]byte, len(in))}
	copy(d1.buf, in)

	if err := d1.value(v, 0); err != nil {
		return err
	}

	if len(d1.buf) != 0 {
		return ErrRemainingBytes
	}

	return nil
}

// DeserializeRawToValue deserializes `in` buffer into
// `dst`'s type and returns the number of bytes used and
// the value of the buffer. If `data` is not either a
// Pointer type an error is returned.
// If `in` buffer can't be deserialized, the number of bytes read and an error message are returned.
func DeserializeRawToValue(in []byte, v reflect.Value) (int, error) {
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
	case reflect.Map:
	default:
		return 0, fmt.Errorf("DeserializeRawToValue value must be a ptr, is %s", v.Kind().String())
	}

	inlen := len(in)
	d1 := &decoder{buf: make([]byte, inlen)}
	copy(d1.buf, in)

	err := d1.value(v, 0)
	if err != nil {
		return 0, err
	}

	return inlen - len(d1.buf), nil
}

// Serialize returns serialized basic type-based `data`
// parameter. Encoding is reflect-based. Panics if `data` is not serializable.
func Serialize(data interface{}) []byte {
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
func Size(v interface{}) (int, error) {
	n, err := datasizeWrite(reflect.Indirect(reflect.ValueOf(v)))
	if err != nil {
		return 0, err
	}
	return n, nil
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

			tag := ff.Tag.Get("enc")
			omitempty := TagOmitempty(tag)

			if omitempty && i != nFields-1 {
				log.Panic(ErrInvalidOmitEmpty)
			}

			if len(tag) > 0 && tag[0] == '-' {
				continue
			}

			fv := v.Field(i)
			if !omitempty || !isEmpty(fv) {
				s, err := datasizeWrite(fv)
				if err != nil {
					return 0, err
				}
				sum += s
			}
		}
		return sum, nil

	case reflect.Bool:
		return 1, nil

	case reflect.String:
		return 4 + v.Len(), nil

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64:
		return int(t.Size()), nil

	default:
		return 0, fmt.Errorf("invalid type %s", t.String())
	}
}

// TagOmitempty returns true if the tag specifies omitempty
func TagOmitempty(tag string) bool {
	return strings.Contains(tag, ",omitempty")
}

func tagName(tag string) string { // nolint: deadcode,megacheck
	commaIndex := strings.Index(tag, ",")
	if commaIndex == -1 {
		return tag
	}

	return tag[:commaIndex]
}

func tagMaxLen(tag string) int {
	maxlenIndex := strings.Index(tag, ",maxlen=")
	if maxlenIndex == -1 {
		return 0
	}

	maxlenRem := tag[maxlenIndex+len(",maxlen="):]
	commaIndex := strings.Index(maxlenRem, ",")
	if commaIndex != -1 {
		maxlenRem = maxlenRem[:commaIndex]
	}

	maxlen, err := strconv.Atoi(maxlenRem)
	if err != nil {
		panic("maxlen must be a number")
	}

	return maxlen
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

type decoder struct {
	buf []byte
}

type encoder struct {
	buf []byte
}

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

func (d *decoder) value(v reflect.Value, maxlen int) error {
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
				if err := d.value(v.Index(i), 0); err != nil {
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
			if err := d.value(keyv, 0); err != nil {
				return err
			}
			if err := d.value(elemv, 0); err != nil {
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

		if maxlen > 0 && length > maxlen {
			return ErrMaxLenExceeded
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
				if err := d.value(elemv, 0); err != nil {
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

			tag := ff.Tag.Get("enc")
			omitempty := TagOmitempty(tag)

			if omitempty && i != nFields-1 {
				log.Panic(ErrInvalidOmitEmpty)
			}

			if len(tag) > 0 && tag[0] == '-' {
				continue
			}

			fv := v.Field(i)
			if fv.CanSet() && ff.Name != "_" {
				maxlen := tagMaxLen(tag)

				if err := d.value(fv, maxlen); err != nil {
					if err == ErrMaxLenExceeded {
						return err
					}

					// omitempty fields at the end of the buffer are ignored if missing
					if !omitempty || len(d.buf) != 0 {
						return err
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

		if maxlen > 0 && length > maxlen {
			return ErrMaxLenExceeded
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

			tag := ff.Tag.Get("enc")
			omitempty := TagOmitempty(tag)

			if omitempty && i != nFields-1 {
				log.Panic(ErrInvalidOmitEmpty)
			}

			if len(tag) > 0 && tag[0] == '-' {
				continue
			}

			fv := v.Field(i)
			if !(omitempty && isEmpty(fv)) && (fv.CanSet() || ff.Name != "_") {
				e.value(fv)
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

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
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
		log.Panicf("Encoding unhandled type %s", v.Type().Name())
	}
}
