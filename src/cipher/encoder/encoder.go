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
	// ErrMapDuplicateKeys encountered duplicate map keys while decoding a map
	ErrMapDuplicateKeys = errors.New("Duplicate keys encountered while decoding a map")
	// ErrInvalidBool is returned if the decoder encounters a value other than 0 or 1 for a bool type field
	ErrInvalidBool = errors.New("Invalid value for bool type")
)

// SerializeUint32 serializes a uint32
func SerializeUint32(x uint32) []byte {
	var b [4]byte
	lePutUint32(b[:], x)
	return b[:]
}

// DeserializeUint32 serializes a uint32
func DeserializeUint32(buf []byte) (uint32, uint64, error) {
	if len(buf) < 4 {
		return 0, 0, ErrBufferUnderflow
	}
	return leUint32(buf[:4]), 4, nil
}

// SerializeAtomic encodes an integer or boolean contained in `data` to bytes.
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
func DeserializeAtomic(in []byte, data interface{}) (uint64, error) {
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
	size := datasizeWrite(v)
	buf := make([]byte, size)
	e := &Encoder{
		Buffer: buf,
	}
	e.value(v)
	return buf
}

// DeserializeString deserializes a string from []byte, returning the string and the number of bytes read
func DeserializeString(in []byte, maxlen int) (string, uint64, error) {
	var s string
	v := reflect.ValueOf(&s)
	v = v.Elem()

	inlen := len(in)
	d1 := &Decoder{
		Buffer: make([]byte, inlen),
	}
	copy(d1.Buffer, in)

	err := d1.value(v, maxlen)
	if err != nil {
		return "", 0, err
	}

	return s, uint64(inlen - len(d1.Buffer)), nil
}

// DeserializeRaw deserializes `in` buffer into return
// parameter. If `data` is not a Pointer or Map type an error
// is returned. If `in` buffer can't be deserialized,
// an error message is returned.
// Returns number of bytes read if no error.
func DeserializeRaw(in []byte, data interface{}) (uint64, error) {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
	case reflect.Map:
	default:
		return 0, fmt.Errorf("DeserializeRaw value must be a ptr, is %s", v.Kind().String())
	}

	inlen := len(in)
	d1 := &Decoder{
		Buffer: make([]byte, inlen),
	}
	copy(d1.Buffer, in)

	if err := d1.value(v, 0); err != nil {
		return 0, err
	}

	return uint64(inlen - len(d1.Buffer)), nil
}

// DeserializeRawExact deserializes `in` buffer into return
// parameter. If `data` is not a Pointer or Map type an error
// is returned. If `in` buffer can't be deserialized,
// an error message is returned.
// Returns number of bytes read if no error.
// If the number of bytes read does not equal the length of the input buffer,
// ErrRemainingBytes is returned.
func DeserializeRawExact(in []byte, data interface{}) error {
	n, err := DeserializeRaw(in, data)
	if err != nil {
		return err
	}
	if n != uint64(len(in)) {
		return ErrRemainingBytes
	}
	return nil
}

// DeserializeRawToValue deserializes `in` buffer into
// `dst`'s type and returns the number of bytes used and
// the value of the buffer. If `data` is not either a
// Pointer type an error is returned.
// If `in` buffer can't be deserialized, the number of bytes read and an error message are returned.
func DeserializeRawToValue(in []byte, v reflect.Value) (uint64, error) {
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
	case reflect.Map:
	default:
		return 0, fmt.Errorf("DeserializeRawToValue value must be a ptr, is %s", v.Kind().String())
	}

	inlen := len(in)
	d1 := &Decoder{
		Buffer: make([]byte, inlen),
	}
	copy(d1.Buffer, in)

	err := d1.value(v, 0)
	if err != nil {
		return 0, err
	}

	return uint64(inlen - len(d1.Buffer)), nil
}

// Serialize returns serialized basic type-based `data`
// parameter. Encoding is reflect-based. Panics if `data` is not serializable.
func Serialize(data interface{}) []byte {
	v := reflect.Indirect(reflect.ValueOf(data))
	size := datasizeWrite(v)
	buf := make([]byte, size)
	e := &Encoder{
		Buffer: buf,
	}
	e.value(v)
	return buf
}

// Size returns how many bytes would it take to encode the
// value v, which must be a fixed-size value (struct) or a
// slice of fixed-size values, or a pointer to such data.
// Reflect-based encoding is used.
func Size(v interface{}) uint64 {
	return datasizeWrite(reflect.Indirect(reflect.ValueOf(v)))
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
func datasizeWrite(v reflect.Value) uint64 {
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
			return uint64(v.Len())
		case reflect.Uint16, reflect.Int16:
			return uint64(v.Len()) * 2
		case reflect.Uint32, reflect.Int32, reflect.Float32:
			return uint64(v.Len()) * 4
		case reflect.Uint64, reflect.Int64, reflect.Float64:
			return uint64(v.Len()) * 8
		default:
			size := uint64(0)
			for i := 0; i < v.Len(); i++ {
				elem := v.Index(i)
				s := datasizeWrite(elem)
				size += s
			}
			return size
		}

	case reflect.Slice:
		t := v.Type()
		elem := t.Elem()
		switch elem.Kind() {
		case reflect.Uint8, reflect.Int8:
			return 4 + uint64(v.Len())
		case reflect.Uint16, reflect.Int16:
			return 4 + uint64(v.Len())*2
		case reflect.Uint32, reflect.Int32, reflect.Float32:
			return 4 + uint64(v.Len())*4
		case reflect.Uint64, reflect.Int64, reflect.Float64:
			return 4 + uint64(v.Len())*8
		default:
			size := uint64(0)
			for i := 0; i < v.Len(); i++ {
				elem := v.Index(i)
				s := datasizeWrite(elem)
				size += s
			}
			return 4 + size
		}

	case reflect.Map:
		// length prefix
		size := uint64(4)
		for _, key := range v.MapKeys() {
			s := datasizeWrite(key)
			size += s
			elem := v.MapIndex(key)
			s = datasizeWrite(elem)
			size += s
		}
		return size

	case reflect.Struct:
		sum := uint64(0)
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
				s := datasizeWrite(fv)
				sum += s
			}
		}
		return sum

	case reflect.Bool:
		return 1

	case reflect.String:
		return 4 + uint64(v.Len())

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64:
		return uint64(t.Size())

	default:
		log.Panicf("invalid type %s", t.String())
		return 0
	}
}

// TagOmitempty returns true if the tag specifies omitempty
func TagOmitempty(tag string) bool {
	return strings.Contains(tag, ",omitempty")
}

func tagName(tag string) string { //nolint:deadcode,megacheck
	commaIndex := strings.Index(tag, ",")
	if commaIndex == -1 {
		return tag
	}

	return tag[:commaIndex]
}

// TagMaxLen returns the maxlen value tagged on a struct. Returns 0 if no tag is present.
func TagMaxLen(tag string) int {
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

// Decoder decodes an object from the skycoin binary encoding format
type Decoder struct {
	Buffer []byte
}

// Encoder encodes an object to the skycoin binary encoding format
type Encoder struct {
	Buffer []byte
}

// Bool decodes bool
func (d *Decoder) Bool() (bool, error) {
	if len(d.Buffer) < 1 {
		return false, ErrBufferUnderflow
	}
	x := d.Buffer[0]
	d.Buffer = d.Buffer[1:] // advance slice

	switch x {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, ErrInvalidBool
	}
}

// Bool encodes bool
func (e *Encoder) Bool(x bool) {
	if x {
		e.Buffer[0] = 1
	} else {
		e.Buffer[0] = 0
	}
	e.Buffer = e.Buffer[1:]
}

// Uint8 decodes uint8
func (d *Decoder) Uint8() (uint8, error) {
	if len(d.Buffer) < 1 {
		return 0, ErrBufferUnderflow
	}

	x := d.Buffer[0]
	d.Buffer = d.Buffer[1:] // advance slice
	return x, nil
}

// Uint8 encodes uint8
func (e *Encoder) Uint8(x uint8) {
	e.Buffer[0] = x
	e.Buffer = e.Buffer[1:]
}

// Uint16 decodes uint16
func (d *Decoder) Uint16() (uint16, error) {
	if len(d.Buffer) < 2 {
		return 0, ErrBufferUnderflow
	}

	x := leUint16(d.Buffer[0:2])
	d.Buffer = d.Buffer[2:]
	return x, nil
}

// Uint16 encodes uint16
func (e *Encoder) Uint16(x uint16) {
	lePutUint16(e.Buffer[0:2], x)
	e.Buffer = e.Buffer[2:]
}

// Uint32 decodes a Uint32
func (d *Decoder) Uint32() (uint32, error) {
	if len(d.Buffer) < 4 {
		return 0, ErrBufferUnderflow
	}

	x := leUint32(d.Buffer[0:4])
	d.Buffer = d.Buffer[4:]
	return x, nil
}

// Uint32 encodes a Uint32
func (e *Encoder) Uint32(x uint32) {
	lePutUint32(e.Buffer[0:4], x)
	e.Buffer = e.Buffer[4:]
}

// Uint64 decodes uint64
func (d *Decoder) Uint64() (uint64, error) {
	if len(d.Buffer) < 8 {
		return 0, ErrBufferUnderflow
	}

	x := leUint64(d.Buffer[0:8])
	d.Buffer = d.Buffer[8:]
	return x, nil
}

// Uint64 encodes uint64
func (e *Encoder) Uint64(x uint64) {
	lePutUint64(e.Buffer[0:8], x)
	e.Buffer = e.Buffer[8:]
}

// ByteSlice encodes []byte
func (e *Encoder) ByteSlice(x []byte) {
	e.Uint32(uint32(len(x)))
	e.CopyBytes(x)
}

// CopyBytes copies bytes to the buffer, without a length prefix
func (e *Encoder) CopyBytes(x []byte) {
	if len(x) == 0 {
		return
	}
	copy(e.Buffer, x)
	e.Buffer = e.Buffer[len(x):]
}

// Int8 decodes int8
func (d *Decoder) Int8() (int8, error) {
	u, err := d.Uint8()
	if err != nil {
		return 0, err
	}

	return int8(u), nil
}

// Int8 encodes int8
func (e *Encoder) Int8(x int8) {
	e.Uint8(uint8(x))
}

// Int16 decodes int16
func (d *Decoder) Int16() (int16, error) {
	u, err := d.Uint16()
	if err != nil {
		return 0, err
	}

	return int16(u), nil
}

// Int16 encodes int16
func (e *Encoder) Int16(x int16) {
	e.Uint16(uint16(x))
}

// Int32 decodes int32
func (d *Decoder) Int32() (int32, error) {
	u, err := d.Uint32()
	if err != nil {
		return 0, err
	}

	return int32(u), nil
}

// Int32 encodes int32
func (e *Encoder) Int32(x int32) {
	e.Uint32(uint32(x))
}

// Int64 decodes int64
func (d *Decoder) Int64() (int64, error) {
	u, err := d.Uint64()
	if err != nil {
		return 0, err
	}

	return int64(u), nil
}

// Int64 encodes int64
func (e *Encoder) Int64(x int64) {
	e.Uint64(uint64(x))
}

func (d *Decoder) value(v reflect.Value, maxlen int) error {
	kind := v.Kind()
	switch kind {
	case reflect.Array:

		t := v.Type()
		elem := t.Elem()

		// Arrays are a fixed size, so the length is not written
		length := v.Len()

		switch elem.Kind() {
		case reflect.Uint8:
			if length > len(d.Buffer) {
				return ErrBufferUnderflow
			}

			reflect.Copy(v, reflect.ValueOf(d.Buffer[:length]))
			d.Buffer = d.Buffer[length:]
		default:
			for i := 0; i < length; i++ {
				if err := d.value(v.Index(i), 0); err != nil {
					return err
				}
			}
		}

	case reflect.Map:
		ul, err := d.Uint32()
		if err != nil {
			return err
		}

		length := int(ul)
		if length < 0 || length > len(d.Buffer) {
			return ErrBufferUnderflow
		}

		if length == 0 {
			return nil
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

		if v.Len() != length {
			return ErrMapDuplicateKeys
		}

	case reflect.Slice:
		ul, err := d.Uint32()
		if err != nil {
			return err
		}

		length := int(ul)
		if length < 0 || length > len(d.Buffer) {
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
			v.SetBytes(d.Buffer[:length])
			d.Buffer = d.Buffer[length:]
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
				maxlen := TagMaxLen(tag)

				if err := d.value(fv, maxlen); err != nil {
					if err == ErrMaxLenExceeded {
						return err
					}

					// omitempty fields at the end of the buffer are ignored if missing
					if !omitempty || len(d.Buffer) != 0 {
						return err
					}
				}
			}
		}

	case reflect.String:
		ul, err := d.Uint32()
		if err != nil {
			return err
		}

		length := int(ul)
		if length < 0 || length > len(d.Buffer) {
			return ErrBufferUnderflow
		}

		if maxlen > 0 && length > maxlen {
			return ErrMaxLenExceeded
		}

		v.SetString(string(d.Buffer[:length]))
		d.Buffer = d.Buffer[length:]

	case reflect.Bool:
		b, err := d.Bool()
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Int8:
		i, err := d.Int8()
		if err != nil {
			return err
		}
		v.SetInt(int64(i))
	case reflect.Int16:
		i, err := d.Int16()
		if err != nil {
			return err
		}
		v.SetInt(int64(i))
	case reflect.Int32:
		i, err := d.Int32()
		if err != nil {
			return err
		}
		v.SetInt(int64(i))
	case reflect.Int64:
		i, err := d.Int64()
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Uint8:
		u, err := d.Uint8()
		if err != nil {
			return err
		}
		v.SetUint(uint64(u))
	case reflect.Uint16:
		u, err := d.Uint16()
		if err != nil {
			return err
		}
		v.SetUint(uint64(u))
	case reflect.Uint32:
		u, err := d.Uint32()
		if err != nil {
			return err
		}
		v.SetUint(uint64(u))
	case reflect.Uint64:
		u, err := d.Uint64()
		if err != nil {
			return err
		}
		v.SetUint(u)

	case reflect.Float32:
		u, err := d.Uint32()
		if err != nil {
			return err
		}
		v.SetFloat(float64(math.Float32frombits(u)))
	case reflect.Float64:
		u, err := d.Uint64()
		if err != nil {
			return err
		}
		v.SetFloat(math.Float64frombits(u))

	default:
		log.Panicf("Decode error: kind %s not handled", v.Kind().String())
	}

	return nil
}

func (e *Encoder) value(v reflect.Value) {
	switch v.Kind() {
	case reflect.Interface:
		e.value(v.Elem())

	case reflect.Array:
		// Arrays are a fixed size, so the length is not written
		t := v.Type()
		elem := t.Elem()
		switch elem.Kind() {
		case reflect.Uint8:
			reflect.Copy(reflect.ValueOf(e.Buffer), v)
			e.Buffer = e.Buffer[v.Len():]
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
			e.ByteSlice(v.Bytes())
		default:
			e.Uint32(uint32(v.Len()))
			for i := 0; i < v.Len(); i++ {
				e.value(v.Index(i))
			}
		}

	case reflect.Map:
		e.Uint32(uint32(v.Len()))
		for _, key := range v.MapKeys() {
			e.value(key)
			e.value(v.MapIndex(key))
		}

	case reflect.Struct:
		t := v.Type()
		nFields := v.NumField()
		for i := 0; i < nFields; i++ {
			// see comment for corresponding code in Decoder.value()
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
		e.Bool(v.Bool())

	case reflect.String:
		e.ByteSlice([]byte(v.String()))

	case reflect.Int8:
		e.Int8(int8(v.Int()))
	case reflect.Int16:
		e.Int16(int16(v.Int()))
	case reflect.Int32:
		e.Int32(int32(v.Int()))
	case reflect.Int64:
		e.Int64(v.Int())

	case reflect.Uint8:
		e.Uint8(uint8(v.Uint()))
	case reflect.Uint16:
		e.Uint16(uint16(v.Uint()))
	case reflect.Uint32:
		e.Uint32(uint32(v.Uint()))
	case reflect.Uint64:
		e.Uint64(v.Uint())

	case reflect.Float32:
		e.Uint32(math.Float32bits(float32(v.Float())))
	case reflect.Float64:
		e.Uint64(math.Float64bits(v.Float()))

	default:
		log.Panicf("Encoding unhandled type %s", v.Type().Name())
	}
}
