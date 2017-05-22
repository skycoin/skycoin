// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package encoder binary implements translation between numbers and byte sequences
// and encoding and decoding of varints.
//
// Numbers are translated by reading and writing fixed-size values.
// A fixed-size value is either a fixed-size arithmetic
// type (int8, uint8, int16, float32, complex64, ...)
// or an array or struct containing only fixed-size values.
//
// Varints are a method of encoding integers using one or more bytes;
// numbers with smaller absolute value take a smaller number of bytes.
// For a specification, see http://code.google.com/apis/protocolbuffers/docs/encoding.html.
package encoder

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"reflect"
)

/*
Todo:
- ensure that invalid input from foreign server cannot crash
- validate packet legnth for incoming
*/

// TODO: constant length byte arrays must not be prefixed

// EncodeInt encodes int
func EncodeInt(b []byte, data interface{}) {
	//var b [8]byte
	var bs []byte
	switch v := data.(type) {

	case int8:
		bs = b[:1]
		b[0] = byte(v)
	case uint8:
		bs = b[:1]
		b[0] = byte(v)
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

// DecodeInt decodes int
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

// DeserializeAtomic fast path for atomic types.
func DeserializeAtomic(in []byte, data interface{}) {
	n := intDestSize(data)
	if len(in) < n {
		log.Panic("Not enough data to deserialize")
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

// DeserializeRaw deserialize raw
func DeserializeRaw(in []byte, data interface{}) error {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
	case reflect.Slice:
	case reflect.Struct:
	default:
		return fmt.Errorf("Invalid type %s", reflect.TypeOf(v).String())
	}

	d1 := &decoder{buf: make([]byte, len(in))}
	copy(d1.buf, in)

	//check if can deserialize
	d2 := &decoder{buf: make([]byte, len(in))}
	copy(d2.buf, in)
	if d2.dchk(v) != 0 {
		return errors.New("Deserialization failed")
	}

	return d1.value(v)
}

// Deserialize takes reader and number of bytes to read
func Deserialize(r io.Reader, dsize int, data interface{}) error {
	// Fallback to reflect-based decoding.
	//fmt.Printf("A1 v is type %s \n", reflect.TypeOf(data).String() )
	//fmt.Printf("A2 v is value/type %s \n", reflect.ValueOf(data).Type().String() )
	//fmt.Printf("A2 v is value,kind %s \n", reflect.ValueOf(data).Kind().String() )

	var v reflect.Value
	switch d := reflect.ValueOf(data); d.Kind() {
	//case reflect.

	case reflect.Ptr:
		v = d.Elem()
	case reflect.Slice:
		v = d
	case reflect.Struct:

	default:
		return errors.New("binary.Read: invalid type " + reflect.TypeOf(d).String())
	}
	//size, err := datasizeWrite(v)
	//if err != nil {
	//	return errors.New("binary.Read: " + err.Error())
	//}

	//fmt.Printf("B v is type %s \n", v.Type().String() )
	//fmt.Printf("C v is type %s \n", reflect.TypeOf(v).String() )
	//fmt.Printf("D v is type %s \n", reflect.TypeOf(reflect.TypeOf(v)).String() )

	d1 := &decoder{buf: make([]byte, dsize)}
	if _, err := io.ReadFull(r, d1.buf); err != nil {
		return err
	}

	//check if can deserialize
	d2 := &decoder{buf: make([]byte, dsize)}
	copy(d2.buf, d1.buf)
	if d2.dchk(v) != 0 {
		return errors.New("Deserialization failed")
	}

	return d1.value(v)
}

// CanDeserialize does a check to see if serialization would be successful
func CanDeserialize(in []byte, dst reflect.Value) bool {
	d1 := &decoder{buf: make([]byte, len(in))}
	copy(d1.buf, in)
	if d1.dchk(dst) != 0 {
		return false
	}
	return true
}

// DeserializeRawToValue returns number of bytes used and an error if deserialization failed
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

	//check if can deserialize
	d2 := &decoder{buf: make([]byte, inlen)}
	copy(d2.buf, d1.buf)
	if d2.dchk(v) != 0 {
		return 0, errors.New("Deserialization failed")
	}

	err := d1.value(v)
	return inlen - len(d1.buf), err
}

// DeserializeToValue deserialize to value
func DeserializeToValue(r io.Reader, dsize int, dst reflect.Value) error {

	//fmt.Printf("*A1 v is type %s \n", data.Type().String() )		//this is the type of the value

	var v reflect.Value
	switch dst.Kind() {
	case reflect.Ptr:
		v = dst.Elem()
	case reflect.Slice:
		v = dst
	case reflect.Struct:

	default:
		return errors.New("binary.Read: invalid type " + reflect.TypeOf(dst).String())
	}

	//fmt.Printf("*A2 v is type %s \n", v.Type().String() )		//this is the type of the value

	d1 := &decoder{buf: make([]byte, dsize)}
	if _, err := io.ReadFull(r, d1.buf); err != nil {
		return err
	}

	return d1.value(v)
}

// SerializeAtomic serializes int or other atomic
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
		b[0] = byte(v)
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

// Serialize serialize struct
func Serialize(data interface{}) []byte {
	// Fast path for basic types.
	// Fallback to reflect-based encoding.
	v := reflect.Indirect(reflect.ValueOf(data))
	size, err := datasizeWrite(v)
	if err != nil {
		//return nil, errors.New("binary.Write: " + err.Error())
		log.Panic(err)
	}
	buf := make([]byte, size)
	e := &encoder{buf: buf}
	e.value(v)
	return buf
}

// Size returns how many bytes Write would generate to encode the value v, which
// must be a fixed-size value or a slice of fixed-size values, or a pointer to such data.
func Size(v interface{}) int {
	n, err := datasizeWrite(reflect.Indirect(reflect.ValueOf(v)))
	if err != nil {
		return -1
	}
	return n
}

// dataSize returns the number of bytes the actual data represented by v occupies in memory.
// For compound structures, it sums the sizes of the elements. Thus, for instance, for a slice
// it returns the length of the slice times the element size and does not count the memory
// occupied by the header.

/* Datasize needs to write variable length slice fields */
/* Datasize for serialization is different than for serialization */
func datasizeWrite(v reflect.Value) (int, error) {
	t := v.Type()
	switch t.Kind() {
	case reflect.Interface:
		//fmt.Println(v.Elem())
		return datasizeWrite(v.Elem())
	case reflect.Array:
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

	case reflect.Slice:
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

	case reflect.Map:
		size := 0
		for _, key := range v.MapKeys() {
			elem := v.MapIndex(key)
			s, err := datasizeWrite(elem)
			if err != nil {
				return 0, err
			}
			size += s
		}
		return 4 + size, nil

	case reflect.Struct:
		sum := 0
		for i, n := 0, t.NumField(); i < n; i++ {
			f := t.Field(i)
			if f.Tag.Get("enc") != "-" {
				s, err := datasizeWrite(v.Field(i))
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
		return len(v.String()) + 4, nil

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64:
		return int(t.Size()), nil

	default:
		return 0, errors.New("invalid type " + t.String())
	}
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
	buf []byte
}

type decoder coder
type encoder coder

func (d *decoder) bool() bool {
	x := d.buf[0]
	d.buf = d.buf[1:] //advance slice
	if x == 0 {
		return false
	}
	return true
}

func (e *encoder) bool(x bool) {
	if x {
		e.buf[0] = 1
	} else {
		e.buf[0] = 0
	}
	e.buf = e.buf[1:]
}

func (d decoder) string() string {
	l := int(d.uint32()) //pop length
	t := d.buf[:l]
	d.buf = d.buf[l:]
	return string(t)
}

func (e encoder) string(xs string) {
	x := []byte(xs)
	l := len(x)
	for i := 0; i < l; i++ {
		e.buf[i] = x[i]
	} //memcpy
	e.buf = e.buf[l:] //advance slice l bytes

} //write this

func (d *decoder) uint8() uint8 {
	x := d.buf[0]
	d.buf = d.buf[1:] //advance slice
	return x
}

func (e *encoder) uint8(x uint8) {
	e.buf[0] = x
	e.buf = e.buf[1:]
}

func (d *decoder) uint16() uint16 {
	x := leUint16(d.buf[0:2])
	d.buf = d.buf[2:]
	return x
}

func (e *encoder) uint16(x uint16) {
	lePutUint16(e.buf[0:2], x)
	e.buf = e.buf[2:]
}

func (d *decoder) uint32() uint32 {
	x := leUint32(d.buf[0:4])
	d.buf = d.buf[4:]
	return x
}

func (e *encoder) uint32(x uint32) {
	lePutUint32(e.buf[0:4], x)
	e.buf = e.buf[4:]
}

func (d *decoder) uint64() uint64 {
	x := leUint64(d.buf[0:8])
	d.buf = d.buf[8:]
	return x
}

func (e *encoder) uint64(x uint64) {
	lePutUint64(e.buf[0:8], x)
	e.buf = e.buf[8:]
}

//v.SetBytes(d.bytes())
func (d decoder) bytes() []byte {
	l := int(d.uint32()) //pop length
	t := d.buf[:l]
	d.buf = d.buf[l:]
	return t
}

func (e encoder) bytes(x []byte) {
	l := len(x)
	for i := 0; i < l; i++ {
		e.buf[i] = x[i]
	} //memcpy
	e.buf = e.buf[l:] //advance slice l bytes

} //write this

func (d *decoder) int8() int8 { return int8(d.uint8()) }

func (e *encoder) int8(x int8) { e.uint8(uint8(x)) }

func (d *decoder) int16() int16 { return int16(d.uint16()) }

func (e *encoder) int16(x int16) { e.uint16(uint16(x)) }

func (d *decoder) int32() int32 { return int32(d.uint32()) }

func (e *encoder) int32(x int32) { e.uint32(uint32(x)) }

func (d *decoder) int64() int64 { return int64(d.uint64()) }

func (e *encoder) int64(x int64) { e.uint64(uint64(x)) }

func (d *decoder) value(v reflect.Value) error {
	kind := v.Kind()
	switch kind {

	case reflect.Array:
		//if len(d.buf) < 4 {
		//    return errors.New("Not enough buffer data to deserialize length")
		//}
		//length := int(d.uint32())
		//if length < 0 || length > len(d.buf) {
		//    return fmt.Errorf("Invalid length: %d", length)
		//}
		//if length != v.Len() {
		//    return errors.New("Incomplete fixed length array received")
		//}

		for i := 0; i < v.Len(); i++ {
			if err := d.value(v.Index(i)); err != nil {
				return err
			}
		}

	case reflect.Slice:
		if len(d.buf) < 4 {
			return errors.New("Not enough buffer data to deserialize length")
		}
		length := int(d.uint32())
		if length < 0 || length > len(d.buf) {
			return fmt.Errorf("Invalid length: %d", length)
		}
		elem := v.Type().Elem()
		if elem.Kind() == reflect.Uint8 {
			v.SetBytes(d.buf[:length])
			d.buf = d.buf[length:]
		} else {
			for i := 0; i < length; i++ {
				elemv := reflect.Indirect(reflect.New(elem))
				if err := d.value(elemv); err != nil {
					return err
				}
				v.Set(reflect.Append(v, elemv))
			}
		}

	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			ff := t.Field(i)
			if ff.Tag.Get("enc") != "-" {
				if fv.CanSet() && ff.Name != "_" {
					if err := d.value(fv); err != nil {
						return err
					}
				} else {
					//dont decode anything
					//d.skip(fv) //BUG!?
				}
			}
		}

	case reflect.String:
		if len(d.buf) < 4 {
			return errors.New("Not enough buffer data to deserialize length")
		}
		length := int(d.uint32())
		if length < 0 || length > len(d.buf) {
			return fmt.Errorf("Invalid length: %d", length)
		}
		v.SetString(string(d.buf[:length]))
		d.buf = d.buf[length:]

	case reflect.Bool:
		v.SetBool(d.bool())
	case reflect.Int8:
		v.SetInt(int64(d.int8()))
	case reflect.Int16:
		v.SetInt(int64(d.int16()))
	case reflect.Int32:
		v.SetInt(int64(d.int32()))
	case reflect.Int64:
		v.SetInt(d.int64())

	case reflect.Uint8:
		v.SetUint(uint64(d.uint8()))
	case reflect.Uint16:
		v.SetUint(uint64(d.uint16()))
	case reflect.Uint32:
		v.SetUint(uint64(d.uint32()))
	case reflect.Uint64:
		v.SetUint(d.uint64())

	case reflect.Float32:
		v.SetFloat(float64(math.Float32frombits(d.uint32())))
	case reflect.Float64:
		v.SetFloat(math.Float64frombits(d.uint64()))

	default:
		log.Panicf("Decode error: kind %s not handled", v.Kind().String())
	}

	return nil
}

func (d *decoder) cmp(n int, m int) int {
	if n != 0 {
		return -1
	}
	if m != 0 {
		return -1
	}
	return 0
}

//advance, returns -1 on failure
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
		c := 0
		for i := 0; i < v.Len(); i++ {
			//t := d.dchk(v.Index(i))
			//c += t
			c = d.cmp(c, d.dchk(v.Index(i)))
		}
		return c

	case reflect.Slice:
		if len(d.buf) < 4 {
			return -1 //error
		}

		length := int(leUint32(d.buf[0:4]))
		d.adv(4) //must succeed

		if length < 0 || length > len(d.buf) {
			return -1 //error
		}

		elem := v.Type().Elem()
		if elem.Kind() == reflect.Uint8 {
			return d.cmp(0, d.adv(length)) //already advanced 4
		}

		c := 0
		for i := 0; i < length; i++ {
			elemv := reflect.Indirect(reflect.New(elem))

			c = d.cmp(c, d.dchk(elemv))
			//c += d.adv(d.dchk(elemv))

			//t := d.dchk(elemv)
			//d.buf = d.buf[t:]
			//c += t

			//v.Set(reflect.Append(v, elemv))
		}
		return c

	case reflect.Struct:
		t := v.Type()
		c := 0
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			ff := t.Field(i)
			if ff.Tag.Get("enc") != "-" {
				if fv.CanSet() && ff.Name != "_" {
					//c += d.adv(d.dchk(fv))
					//c += d.dchk(fv)
					c = d.cmp(c, d.dchk(fv))
				} else {
					//dont try to decode anything
					//d.skip(fv) //BUG!?
				}
			}
		}
		return c

	case reflect.Bool:
		return d.adv(1)
	case reflect.String:
		length := int(leUint32(d.buf[0:4]))
		d.adv(4) //must succeed
		return d.cmp(0, d.adv(length))
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

	case reflect.Array: //fixed size
		//e.uint32(uint32(v.Len()))
		for i := 0; i < v.Len(); i++ {
			e.value(v.Index(i))
		}

	case reflect.Slice:
		e.uint32(uint32(v.Len()))
		for i := 0; i < v.Len(); i++ {
			e.value(v.Index(i))
		}

	case reflect.Map:
		e.uint32(uint32(v.Len()))
		for _, key := range v.MapKeys() {
			e.value(v.MapIndex(key))
		}

	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			// see comment for corresponding code in decoder.value()
			v := v.Field(i)
			f := t.Field(i)
			if f.Tag.Get("enc") != "-" {
				if v.CanSet() || f.Name != "_" {
					e.value(v)
				} else {
					//dont write anything
					//e.skip(v)
				}
			}
		}

	// case reflect.Slice:
	//     t := v.Type() //type of the value

	//     //handle byte array
	//     if t.Elem().Kind() == reflect.Uint8 {
	//         b := v.Bytes()
	//         n := len(b)
	//         e.uint32(uint32(n))
	//         for i := 0; i < n; i++ {
	//             e.buf[i] = b[i]
	//         }   //memcpy
	//         e.buf = e.buf[n:] //advance slice n bytes
	//     } else { //handle struct array
	//         s := int(t.Elem().Size())
	//         if s <= 1 {
	//             log.Panic()
	//         }
	//         n := v.Len()            //const
	//         e.uint32(uint32(n * s)) //push number of bytes
	//         for i := 0; i < n; i++ {
	//             e.value(v.Index(i))
	//         }
	//     }

	case reflect.Bool:
		e.bool(v.Bool())

	case reflect.String:
		vb := []byte(v.String())
		e.uint32(uint32(len(vb)))
		for i := 0; i < len(vb); i++ {
			e.uint8(vb[i])
		}

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

func (d *decoder) skip(v reflect.Value) {
	n, _ := datasizeWrite(v)
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
func (e *encoder) skip(v reflect.Value) {
	n, _ := datasizeWrite(v)
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
