// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package binary implements translation between numbers and byte sequences
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

// Fast path for atomic types.
func DeserializeAtomic(in []byte, data interface{}) error {

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
			*v = int16(le_Uint16(bs))
		case *uint16:
			*v = le_Uint16(bs)
		case *int32:
			*v = int32(le_Uint32(bs))
		case *uint32:
			*v = le_Uint32(bs)
		case *int64:
			*v = int64(le_Uint64(bs))
		case *uint64:
			*v = le_Uint64(bs)
		}
		return nil
	}

	log.Panic()
	return errors.New("type not atomic")
}

//deserialize from a
func DeserializeRaw(in []byte, data interface{}) error {
	var v reflect.Value
	switch d := reflect.ValueOf(data); d.Kind() {
	case reflect.Ptr:
		v = d.Elem()
	case reflect.Slice:
		v = d
	case reflect.Struct:
	default:
		return errors.New("binary.Read: invalid type " + reflect.TypeOf(d).String())
	}

	d := &decoder{buf: make([]byte, len(in))}
	copy(d.buf, in)

	d.value(v)

	return nil
}

//takes reader and number of bytes to read
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

	d := &decoder{buf: make([]byte, dsize)}
	if _, err := io.ReadFull(r, d.buf); err != nil {
		return err
	}
	d.value(v)

	return nil
}

func DeserializeRawToValue(in []byte, dst reflect.Value) error {
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

	d := &decoder{buf: make([]byte, len(in))}
	copy(d.buf, in)

	d.value(v)

	return nil
}

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

	d := &decoder{buf: make([]byte, dsize)}
	if _, err := io.ReadFull(r, d.buf); err != nil {
		return err
	}
	d.value(v)

	return nil
}

//serialize int or other atomic
func SerializeAtomic(data interface{}) ([]byte, error) {
	var b [8]byte
	var bs []byte
	switch v := data.(type) {
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
		le_PutUint16(bs, uint16(*v))
	case int16:
		bs = b[:2]
		le_PutUint16(bs, uint16(v))
	case *uint16:
		bs = b[:2]
		le_PutUint16(bs, *v)
	case uint16:
		bs = b[:2]
		le_PutUint16(bs, v)
	case *int32:
		bs = b[:4]
		le_PutUint32(bs, uint32(*v))
	case int32:
		bs = b[:4]
		le_PutUint32(bs, uint32(v))
	case *uint32:
		bs = b[:4]
		le_PutUint32(bs, *v)
	case uint32:
		bs = b[:4]
		le_PutUint32(bs, v)
	case *int64:
		bs = b[:8]
		le_PutUint64(bs, uint64(*v))
	case int64:
		bs = b[:8]
		le_PutUint64(bs, uint64(v))
	case *uint64:
		bs = b[:8]
		le_PutUint64(bs, *v)
	case uint64:
		bs = b[:8]
		le_PutUint64(bs, v)
	}
	if bs != nil {
		//_, err := w.Write(bs)
		//return err
		return bs, nil
	}

	log.Panic()
	return nil, errors.New("type not atomic")

}

//serialize struct
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
	//_, err = w.Write(buf)
	return buf
}

// Size returns how many bytes Write would generate to encode the value v, which
// must be a fixed-size value or a slice of fixed-size values, or a pointer to such data.
func Size(v interface{}) int {
	n, err := datasizeWrite(reflect.Indirect(reflect.ValueOf(v))) //n, err := datasizeWrite(reflect.Indirect(reflect.ValueOf(v)))
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
	t := v.Type()     //type of the value
	switch t.Kind() { //kind of the type

	case reflect.Slice: //fixed array
		//n, err := datasizeWrite(t.Elem())
		var e int = int(t.Elem().Size()) //should be
		return 4 + e*v.Len(), nil        //return v.Len()*n, nil

	case reflect.Array:
		var n int = int(t.Elem().Size())
		if n != 1 {
			log.Panic("non-byte arrays not supported yet")
		}
		return n * v.Len(), nil

		//return 0, errors.New("invalid type: arrays not supported")

	//this is entry level
	case reflect.Struct:
		sum := 0
		for i, n := 0, t.NumField(); i < n; i++ {
			s, err := datasizeWrite(v.Field(i)) //t.Field(i).Type
			if err != nil {
				log.Panic(err)
				return 0, err
			}
			sum += s
		}
		return sum, nil

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64:
		return int(t.Size()), nil
	}
	return 0, errors.New("invalid type " + t.String())

}

/*
	Internals
*/

func le_Uint16(b []byte) uint16 { return uint16(b[0]) | uint16(b[1])<<8 }

func le_PutUint16(b []byte, v uint16) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}

func le_Uint32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func le_PutUint32(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

func le_Uint64(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func le_PutUint64(b []byte, v uint64) {
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
	x := le_Uint16(d.buf[0:2])
	d.buf = d.buf[2:]
	return x
}

func (e *encoder) uint16(x uint16) {
	le_PutUint16(e.buf[0:2], x)
	e.buf = e.buf[2:]
}

func (d *decoder) uint32() uint32 {
	x := le_Uint32(d.buf[0:4])
	d.buf = d.buf[4:]
	return x
}

func (e *encoder) uint32(x uint32) {
	le_PutUint32(e.buf[0:4], x)
	e.buf = e.buf[4:]
}

func (d *decoder) uint64() uint64 {
	x := le_Uint64(d.buf[0:8])
	d.buf = d.buf[8:]
	return x
}

func (e *encoder) uint64(x uint64) {
	le_PutUint64(e.buf[0:8], x)
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

func (d *decoder) value(v reflect.Value) {

	//fmt.Printf("len(d.buf)= %v \n", len(d.buf))

	//fmt.Printf("v.Kind()= %v \n", v.Kind().String())

	switch v.Kind() {

	case reflect.Array:
		l := v.Len() //length is hardcoded in structs
		for i := 0; i < l; i++ {
			vtmp := v.Index(i)
			d.value(vtmp)
		}

	//this is entry level
	case reflect.Struct:
		t := v.Type()
		l := v.NumField()
		for i := 0; i < l; i++ {
			if v := v.Field(i); v.CanSet() || t.Field(i).Name != "_" {
				d.value(v)
			} else {
				d.skip(v)
			}
		}

	/*
		Pop length
	*/
	case reflect.Slice:
		l := int(d.uint32()) //pop length

		if l < 0 || l > len(d.buf) {
			fmt.Printf("ERROR, binary decoder: l= %v, len(d.buf)= %v \n", l, len(d.buf))
			break
		}

		t := v.Type() //type of the value

		if t.Elem().Kind() == reflect.Uint8 { //handle byte stream
			v.SetBytes(d.buf[:l])
			d.buf = d.buf[l:]
		} else {
			s := int(t.Elem().Size())
			if l%s != 0 {
				log.Panic("ERROR, binary decoder, array size is not multiple of struct size\n")
			}

			n := l / s
			v.Set(reflect.MakeSlice(v.Type(), n, n))
			for i := 0; i < n; i++ {
				d.value(v.Index(i))
			}
		}
	/*
		*v = reflect.MakeSlice(reflect.TypeOf([]byte{}), l, l)
		for i := 0; i < l; i++ {
			vtmp := v.Index(i)
			vtmp.SetInt(int64(d.uint8())) //set each index
		}
	*/

	/*
		l := v.Len()	//
		e.uint32(uint32(l))
		for i := 0; i < l; i++ {
			e.uint8(uint8(v.Index(i)) //calls self
	*/

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

		fmt.Printf("Binary, Encoder Error, type %s not handled \n", v.Kind().String())
		log.Panic()
	}
}

func (e *encoder) value(v reflect.Value) {

	switch v.Kind() {
	case reflect.Array:
		l := v.Len() //length is hard coded in struct
		for i := 0; i < l; i++ {
			e.value(v.Index(i))
		}

	case reflect.Struct:
		t := v.Type()
		l := v.NumField()
		for i := 0; i < l; i++ {
			// see comment for corresponding code in decoder.value()
			if v := v.Field(i); v.CanSet() || t.Field(i).Name != "_" {
				e.value(v)
			} else {
				e.skip(v)
			}
		}

	/*
		Push length
	*/
	case reflect.Slice:
		t := v.Type() //type of the value

		//handle byte array
		if t.Elem().Kind() == reflect.Uint8 {
			b := v.Bytes()
			n := len(b)
			e.uint32(uint32(n))
			for i := 0; i < n; i++ {
				e.buf[i] = b[i]
			} //memcpy
			e.buf = e.buf[n:] //advance slice n bytes
		} else { //handle struct array
			s := int(t.Elem().Size())
			if s <= 1 {
				log.Panic()
			}
			n := v.Len()            //const
			e.uint32(uint32(n * s)) //push number of bytes
			for i := 0; i < n; i++ {
				e.value(v.Index(i))
			}
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
		log.Panic()

	}

}

func (d *decoder) skip(v reflect.Value) {
	n, _ := datasizeWrite(v)
	d.buf = d.buf[n:]
}

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
