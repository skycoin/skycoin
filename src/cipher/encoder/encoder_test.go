package encoder

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
)

func randBytes(t *testing.T, n int) []byte { // nolint: unparam
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	_, err := rand.Read(bytes)
	require.NoError(t, err)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return bytes
}

//Size of= 13
type TestStruct struct {
	X int32
	Y int64
	Z uint8
	K []byte
	W bool
	T string
	U cipher.PubKey
}

type TestStruct2 struct {
	X int32
	Y int64
	Z uint8
	K [8]byte
	W bool
}

type TestStructIgnore struct {
	X int32
	Y int64
	Z uint8 `enc:"-"`
	K []byte
}

type TestStructWithoutIgnore struct {
	X int32
	Y int64
	K []byte
}

//func (*B) Fatal

func Test_Encode_1(t *testing.T) {
	var ts TestStruct
	ts.X = 345535
	ts.Y = 23432435443
	ts.Z = 255
	ts.K = []byte("TEST6")
	ts.W = true
	ts.T = "hello"
	ts.U = cipher.PubKey{1, 2, 3, 0, 5, 4, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	b := Serialize(ts)

	var ts2 TestStruct
	err := DeserializeRaw(b, &ts2)
	require.NoError(t, err)

	b2 := Serialize(ts2)

	c := bytes.Compare(b, b2)
	require.Equal(t, c, 0)
}

func Test_Encode_2a(t *testing.T) {
	var ts TestStruct2
	ts.X = 345535
	ts.Y = 23432435443
	ts.Z = 255
	ts.W = false
	_tt := []byte("ASDSADFSDFASDFSD")
	copy(ts.K[:], _tt)

	b := Serialize(ts)

	var ts2 TestStruct2
	err := DeserializeRaw(b, &ts2)
	require.NoError(t, err)

	b2 := Serialize(ts2)

	c := bytes.Compare(b, b2)
	require.Equal(t, c, 0)
}

func Test_Encode_2b(t *testing.T) {
	var ts TestStruct2
	ts.X = 345535
	ts.Y = 23432435443
	ts.Z = 255
	_tt := []byte("ASDSADFSDFASDFSD")
	copy(ts.K[:], _tt)

	b := Serialize(ts)

	var ts2 TestStruct2
	err := DeserializeRaw(b, &ts2)
	require.NoError(t, err)

	b2 := Serialize(ts2)

	c := bytes.Compare(b, b2)
	require.Equal(t, c, 0)
}

type TestStruct3 struct {
	X int32
	K []byte
}

func Test_Encode_3a(t *testing.T) {
	var t1 TestStruct3
	t1.X = 345535
	t1.K = randBytes(t, 32)

	b := Serialize(t1)

	var t2 TestStruct3
	err := DeserializeRaw(b, &t2)
	require.NoError(t, err)

	require.False(t, t1.X != t2.X || len(t1.K) != len(t2.K) || !bytes.Equal(t1.K, t2.K))

	b2 := Serialize(t2)

	c := bytes.Compare(b, b2)
	require.Equal(t, c, 0)
}

func Test_Encode_3b(t *testing.T) {
	var t1 TestStruct3
	t1.X = 345535
	t1.K = randBytes(t, 32)

	b := Serialize(t1)

	var t2 TestStruct3
	err := DeserializeRaw(b, &t2)
	require.NoError(t, err)

	require.False(t, t1.X != t2.X || len(t1.K) != len(t2.K) || !bytes.Equal(t1.K, t2.K))

	b2 := Serialize(t2)

	c := bytes.Compare(b, b2)
	require.Equal(t, c, 0)
}

type TestStruct4 struct {
	X int32
	Y int32
}

type TestStruct5 struct {
	X int32
	A []TestStruct4
}

func Test_Encode_4(t *testing.T) {
	var t1 TestStruct5
	t1.X = 345535

	t1.A = make([]TestStruct4, 8)

	b := Serialize(t1)

	var t2 TestStruct5
	err := DeserializeRaw(b, &t2)
	require.NoError(t, err)

	require.Equal(t, t1.X, t2.X, "TestStruct5.X not equal")

	require.Equal(t, len(t1.A), len(t2.A), "Slice lengths not equal: %d != %d", len(t1.A), len(t2.A))

	for i, ts := range t1.A {
		require.Equal(t, ts, t2.A[i], "Slice values not equal")
	}

	b2 := Serialize(t2)

	require.True(t, bytes.Equal(b, b2))
}

func Test_Encode_5(t *testing.T) {
	var ts TestStruct2
	ts.X = 345535
	ts.Y = 23432435443
	ts.Z = 255

	b1 := Serialize(ts)

	var tts = reflect.TypeOf(ts)
	var v = reflect.New(tts) // pointer to type tts

	_, err := DeserializeRawToValue(b1, v)
	require.NoError(t, err)

	v = reflect.Indirect(v)
	if v.FieldByName("X").Int() != int64(ts.X) {
		t.Fatalf("X not equal")
	}
	if v.FieldByName("Y").Int() != ts.Y {
		t.Fatalf("Y not equal")
	}
	if v.FieldByName("Z").Uint() != uint64(ts.Z) {
		t.Fatalf("Z not equal")
	}
}

func Test_Encode_IgnoreTagSerialize(t *testing.T) {
	var ts TestStructIgnore
	ts.X = 345535
	ts.Y = 23432435443
	ts.Z = 255
	ts.K = []byte("TEST6")

	b := Serialize(ts)

	var ts2 TestStructIgnore
	ts.X = 0
	ts.Y = 0
	ts.Z = 0
	ts.K = []byte("")
	err := DeserializeRaw(b, &ts2)
	require.NoError(t, err)

	if ts2.Z != 0 {
		t.Fatalf("Z should not deserialize. It is %d", ts2.Z)
	}

	var ts3 TestStructWithoutIgnore
	err = DeserializeRaw(b, &ts3)
	require.NoError(t, err)

	b2 := Serialize(ts2)
	c := bytes.Compare(b, b2)
	require.Equal(t, c, 0)
}

type Contained struct {
	X     uint32
	Y     uint64
	Bytes []uint8
	Ints  []uint16
}

type Container struct {
	Elements []Contained
}

func TestEncodeNestedSlice(t *testing.T) {
	size := 0
	elems := make([]Contained, 4)
	for i := range elems {
		elems[i].X = uint32(i)
		size += 4
		elems[i].Y = uint64(i)
		size += 8
		elems[i].Bytes = make([]uint8, i)
		for j := range elems[i].Bytes {
			elems[i].Bytes[j] = uint8(j)
		}
		size += 4 + i*1
		elems[i].Ints = make([]uint16, i)
		for j := range elems[i].Ints {
			elems[i].Ints[j] = uint16(j)
		}
		size += 4 + i*2
	}
	c := Container{elems}
	n, err := datasizeWrite(reflect.ValueOf(c))
	require.NoError(t, err)
	require.False(t, n != size+4, "Wrong data size")

	b := Serialize(c)
	d := Container{}
	err = DeserializeRaw(b, &d)
	require.NoError(t, err)

	for i, e := range d.Elements {
		if c.Elements[i].X != e.X || c.Elements[i].Y != e.Y {
			t.Fatalf("Deserialized x, y to invalid value. "+
				"Expected %d,%d but got %d,%d", c.Elements[i].X,
				c.Elements[i].Y, e.X, e.Y)
		}
		if len(c.Elements[i].Bytes) != len(e.Bytes) {
			t.Fatal("Deserialized Bytes to invalid length")
		}
		for j, b := range c.Elements[i].Bytes {
			if c.Elements[i].Bytes[j] != b {
				t.Fatal("Deserialized to invalid value")
			}
		}
		if len(c.Elements[i].Ints) != len(e.Ints) {
			t.Fatal("Deserialized Ints to invalid length")
		}
		for j, b := range c.Elements[i].Ints {
			if c.Elements[i].Ints[j] != b {
				t.Fatal("Deserialized Ints to invalid value")
			}
		}
	}
}

type Array struct {
	Arr []int
}

func TestDecodeNotEnoughLength(t *testing.T) {
	b := make([]byte, 2)
	var d Array
	err := DeserializeRaw(b, &d)
	require.Error(t, err)
	require.Equal(t, ErrBufferUnderflow, err)

	// Test with slice
	thing := make([]int, 3)
	err = DeserializeRaw(b, thing)
	require.Error(t, err)
	require.Equal(t, ErrBufferUnderflow, err)
}

func TestFlattenMultidimensionalBytes(t *testing.T) {
	var data [16][16]byte
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			data[i][j] = byte(i * j)
		}
	}

	b := Serialize(data)
	expect := 16 * 16
	if len(b) != expect {
		t.Fatalf("Expected %d bytes, decoded to %d bytes", expect, len(b))
	}

}

func TestMultiArrays(t *testing.T) {
	var data [16][16]byte
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			data[i][j] = byte(i * j)
		}
	}

	b := Serialize(data)

	var data2 [16][16]byte

	err := DeserializeRaw(b, &data2)
	require.NoError(t, err)

	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			if data[i][j] != data2[i][j] {
				t.Fatalf("failed round trip test")
			}
		}
	}

	b2 := Serialize(data2)
	if !bytes.Equal(b, b2) {
		t.Fatalf("Failed round trip test")
	}

	if len(b) != 256 {
		t.Fatalf("decoded to wrong byte length")
	}

}

func TestDeserializeAtomic(t *testing.T) {
	var sp uint64 = 0x000C8A9E1809F720
	b := make([]byte, 8)
	SerializeAtomic(b, sp)

	var i uint64
	DeserializeAtomic(b, &i)

	if i != sp {
		t.Fatal("round trip atomic fail")
	}
}

func TestSerializeDeserializeAtomic(t *testing.T) {
	var di64 int64
	err := DeserializeAtomic(nil, &di64)
	require.Equal(t, ErrBufferUnderflow, err)

	var db bool
	err = DeserializeAtomic([]byte{2}, &db)
	require.Equal(t, ErrBoolDecode, err)

	err = SerializeAtomic(nil, int8(0))
	require.Equal(t, ErrBufferOverflow, err)
	err = SerializeAtomic(nil, int16(0))
	require.Equal(t, ErrBufferOverflow, err)
	err = SerializeAtomic(nil, int32(0))
	require.Equal(t, ErrBufferOverflow, err)
	err = SerializeAtomic(nil, int64(0))
	require.Equal(t, ErrBufferOverflow, err)
	err = SerializeAtomic(nil, uint8(0))
	require.Equal(t, ErrBufferOverflow, err)
	err = SerializeAtomic(nil, uint16(0))
	require.Equal(t, ErrBufferOverflow, err)
	err = SerializeAtomic(nil, uint32(0))
	require.Equal(t, ErrBufferOverflow, err)
	err = SerializeAtomic(nil, uint64(0))
	require.Equal(t, ErrBufferOverflow, err)

	d := make([]byte, 8)

	b := false
	err = SerializeAtomic(d, b)
	require.NoError(t, err)
	var bb bool
	err = DeserializeAtomic(d, &bb)
	require.NoError(t, err)
	require.Equal(t, b, bb)

	b = true
	bb = false
	err = SerializeAtomic(d, b)
	require.NoError(t, err)
	err = DeserializeAtomic(d, &bb)
	require.NoError(t, err)
	require.Equal(t, b, bb)

	var u8 uint8 = 0xF7
	err = SerializeAtomic(d, u8)
	require.NoError(t, err)
	var u8b uint8
	err = DeserializeAtomic(d, &u8b)
	require.NoError(t, err)
	require.Equal(t, u8, u8b)

	var u16 uint16 = 0xF720
	err = SerializeAtomic(d, u16)
	require.NoError(t, err)
	var u16b uint16
	err = DeserializeAtomic(d, &u16b)
	require.NoError(t, err)
	require.Equal(t, u16, u16b)

	var u32 uint32 = 0x1809F720
	err = SerializeAtomic(d, u32)
	require.NoError(t, err)
	var u32b uint32
	err = DeserializeAtomic(d, &u32b)
	require.NoError(t, err)
	require.Equal(t, u32, u32b)

	var u64 uint64 = 0x000C8A9E1809F720
	err = SerializeAtomic(d, u64)
	require.NoError(t, err)
	var u64b uint64
	err = DeserializeAtomic(d, &u64b)
	require.NoError(t, err)
	require.Equal(t, u64, u64b)

	var i8 int8 = 0x69
	err = SerializeAtomic(d, i8)
	require.NoError(t, err)
	var i8b int8
	err = DeserializeAtomic(d, &i8b)
	require.NoError(t, err)
	require.Equal(t, i8, i8b)

	var i16 int16 = 0x6920
	err = SerializeAtomic(d, i16)
	require.NoError(t, err)
	var i16b int16
	err = DeserializeAtomic(d, &i16b)
	require.NoError(t, err)
	require.Equal(t, i16, i16b)

	var i32 int32 = 0x1809F720
	err = SerializeAtomic(d, i32)
	require.NoError(t, err)
	var i32b int32
	err = DeserializeAtomic(d, &i32b)
	require.NoError(t, err)
	require.Equal(t, i32, i32b)

	var i64 int64 = 0x000C8A9E1809F720
	err = SerializeAtomic(d, i64)
	require.NoError(t, err)
	var i64b int64
	err = DeserializeAtomic(d, &i64b)
	require.NoError(t, err)
	require.Equal(t, i64, i64b)
}

type TestStruct5a struct {
	Test uint64
}

func TestSerializeAtomicPanics(t *testing.T) {
	var x float32
	require.PanicsWithValue(t, "SerializeAtomic unhandled type", func() {
		SerializeAtomic(nil, x) // nolint: errcheck
	})

	var tst TestStruct5a
	d := make([]byte, 8)
	require.PanicsWithValue(t, "SerializeAtomic unhandled type", func() {
		SerializeAtomic(d, &tst) // nolint: errcheck
	})
}

func TestDeserializeAtomicPanics(t *testing.T) {
	var y int8
	require.PanicsWithValue(t, "DeserializeAtomic unhandled type", func() {
		DeserializeAtomic(nil, y) // nolint: errcheck
	})

	var x float32
	require.PanicsWithValue(t, "DeserializeAtomic unhandled type", func() {
		DeserializeAtomic(nil, &x) // nolint: errcheck
	})

	var tst TestStruct5a
	d := make([]byte, 8)
	require.PanicsWithValue(t, "DeserializeAtomic unhandled type", func() {
		DeserializeAtomic(d, &tst) // nolint: errcheck
	})
}

func TestByteArray(t *testing.T) {

	tstr := "7105a46cb4c2810f0c916e0bb4b4e4ef834ad42040c471b42c96d356a9fd1b21"

	d, err := hex.DecodeString(tstr)
	if err != nil {
		t.Fail()
	}

	buff := Serialize(d)
	var buff2 [32]byte
	copy(buff2[0:32], buff[0:32])

	if len(buff2) != 32 {
		t.Errorf("incorrect serialization length for fixed sized arrays: %d byte fixed sized array serialized to %d bytes \n", len(d), len(buff2))
	}

}

func TestEncodeDictInt2Int(t *testing.T) {
	m1 := map[uint8]uint64{0: 0, 1: 1, 2: 2}
	buff := Serialize(m1)
	if len(buff) != 4 /* Length */ +(1+8)*len(m1) /* 1b key + 8b value per entry */ {
		t.Fail()
	}
	m2 := make(map[uint8]uint64)
	if DeserializeRaw(buff, m2) != nil {
		t.Fail()
	}
	if len(m1) != len(m2) {
		t.Errorf("Expected length %d but got %d", len(m1), len(m2))
	}
	for key := range m1 {
		if m1[key] != m2[key] {
			t.Errorf("Expected value %d for key %d but got %d", m1[key], key, m2[key])
		}
	}
}

type TestStructWithDict struct {
	X int32
	Y int64
	M map[uint8]TestStruct
	K []byte
}

func TestEncodeDictNested(t *testing.T) {
	s1 := TestStructWithDict{
		X: 0x01234567,
		Y: 0x0123456789ABCDEF,
		M: map[uint8]TestStruct{
			0x01: TestStruct{
				X: 0x01234567,
				Y: 0x0123456789ABCDEF,
				Z: 0x01,
				K: []byte{0, 1, 2},
				W: true,
				T: "ab",
				U: cipher.PubKey{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
					17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
				},
			},
			0x23: TestStruct{
				X: 0x01234567,
				Y: 0x0123456789ABCDEF,
				Z: 0x01,
				K: []byte{0, 1, 2},
				W: true,
				T: "cd",
				U: cipher.PubKey{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
					17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
				},
			},
		},
		K: []byte{0, 1, 2, 3, 4},
	}
	buff := Serialize(s1)
	require.NotEmpty(t, buff)

	s2 := TestStructWithDict{}
	err := DeserializeRaw(buff, &s2)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(s1, s2), "Expected %v but got %v", s1, s2)
}

func TestEncodeDictString2Int64(t *testing.T) {
	v := map[string]int64{
		"foo": 1,
		"bar": 2,
	}

	b := Serialize(v)

	v2 := make(map[string]int64)
	err := DeserializeRaw(b, &v2)
	require.NoError(t, err)

	require.Equal(t, v, v2)
}

func TestOmitEmptyString(t *testing.T) {

	type omitString struct {
		A string `enc:"a,omitempty"`
	}

	cases := []struct {
		name                string
		input               omitString
		outputShouldBeEmpty bool
	}{
		{
			name: "string not empty",
			input: omitString{
				A: "foo",
			},
		},

		{
			name:                "string empty",
			input:               omitString{},
			outputShouldBeEmpty: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := Serialize(tc.input)

			if tc.outputShouldBeEmpty {
				require.Empty(t, b)
			} else {
				require.NotEmpty(t, b)
			}

			var y omitString
			err := DeserializeRaw(b, &y)
			require.NoError(t, err)

			require.Equal(t, tc.input, y)
		})
	}

}

func TestOmitEmptySlice(t *testing.T) {
	type omitSlice struct {
		B []byte `enc:"b,omitempty"`
	}

	cases := []struct {
		name                string
		input               omitSlice
		expect              *omitSlice
		outputShouldBeEmpty bool
	}{
		{
			name: "slice not empty",
			input: omitSlice{
				B: []byte("foo"),
			},
		},

		{
			name:                "slice nil",
			input:               omitSlice{},
			outputShouldBeEmpty: true,
		},

		{
			name: "slice empty but not nil",
			input: omitSlice{
				B: []byte{},
			},
			expect:              &omitSlice{},
			outputShouldBeEmpty: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := Serialize(tc.input)

			if tc.outputShouldBeEmpty {
				require.Empty(t, b)
			} else {
				require.NotEmpty(t, b)
			}

			var y omitSlice
			err := DeserializeRaw(b, &y)
			require.NoError(t, err)

			expect := tc.expect
			if expect == nil {
				expect = &tc.input
			}

			require.Equal(t, *expect, y)
		})
	}
}

func TestOmitEmptyMap(t *testing.T) {

	type omitMap struct {
		C map[string]int64 `enc:"d,omitempty"`
	}

	cases := []struct {
		name                string
		input               omitMap
		expect              *omitMap
		outputShouldBeEmpty bool
	}{
		{
			name: "map not empty",
			input: omitMap{
				C: map[string]int64{"foo": 1},
			},
		},

		{
			name:                "map nil",
			input:               omitMap{},
			outputShouldBeEmpty: true,
		},

		{
			name: "map empty but not nil",
			input: omitMap{
				C: map[string]int64{},
			},
			expect:              &omitMap{},
			outputShouldBeEmpty: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := Serialize(tc.input)

			if tc.outputShouldBeEmpty {
				require.Empty(t, b)
			} else {
				require.NotEmpty(t, b)
			}

			var y omitMap
			err := DeserializeRaw(b, &y)
			require.NoError(t, err)

			expect := tc.expect
			if expect == nil {
				expect = &tc.input
			}

			require.Equal(t, *expect, y)
		})
	}
}

func TestOmitEmptyMixedFinalByte(t *testing.T) {
	type omitMixed struct {
		A string
		B []byte `enc:",omitempty"`
	}

	cases := []struct {
		name   string
		input  omitMixed
		expect omitMixed
	}{
		{
			name: "none empty",
			input: omitMixed{
				A: "foo",
				B: []byte("foo"),
			},
			expect: omitMixed{
				A: "foo",
				B: []byte("foo"),
			},
		},

		{
			name: "byte nil",
			input: omitMixed{
				A: "foo",
			},
			expect: omitMixed{
				A: "foo",
			},
		},

		{
			name: "byte empty but not nil",
			input: omitMixed{
				A: "foo",
				B: []byte{},
			},
			expect: omitMixed{
				A: "foo",
			},
		},

		{
			name: "first string empty but not omitted",
			input: omitMixed{
				B: []byte("foo"),
			},
			expect: omitMixed{
				B: []byte("foo"),
			},
		},

		{
			name:   "all empty",
			input:  omitMixed{},
			expect: omitMixed{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := Serialize(tc.input)
			require.NotEmpty(t, b)

			var y omitMixed
			err := DeserializeRaw(b, &y)
			require.NoError(t, err)

			require.Equal(t, tc.expect, y)
		})
	}
}

func TestOmitEmptyFinalFieldOnly(t *testing.T) {
	type bad struct {
		A string
		B string `enc:",omitempty"`
		C string
	}

	require.Panics(t, func() {
		var b bad
		Serialize(b)
	})
}

func TestParseTag(t *testing.T) {
	cases := []struct {
		tag       string
		name      string
		omitempty bool
	}{
		{
			tag:  "foo",
			name: "foo",
		},
		{
			tag:  "foo,",
			name: "foo",
		},
		{
			tag:  "foo,asdasd",
			name: "foo",
		},
		{
			tag:       "foo,omitempty",
			name:      "foo",
			omitempty: true,
		},
		{
			tag:  "omitempty",
			name: "omitempty",
		},
		{
			tag:       ",omitempty",
			omitempty: true,
		},
		{
			tag: "",
		},
		{
			tag:  "-",
			name: "-",
		},
	}

	for _, tc := range cases {
		t.Run(tc.tag, func(t *testing.T) {
			name, omitempty := ParseTag(tc.tag)
			require.Equal(t, tc.name, name)
			require.Equal(t, tc.omitempty, omitempty)
		})
	}
}

type primitiveInts struct {
	A int8
	B uint8
	C int16
	D uint16
	E int32
	F uint32
	G int64
	H uint64
}

func TestPrimitiveInts(t *testing.T) {
	cases := []struct {
		name string
		c    primitiveInts
	}{
		{
			name: "all maximums",
			c: primitiveInts{
				A: math.MaxInt8,
				B: math.MaxUint8,
				C: math.MaxInt16,
				D: math.MaxUint16,
				E: math.MaxInt32,
				F: math.MaxUint32,
				G: math.MaxInt64,
				H: math.MaxUint64,
			},
		},
		{
			name: "negative integers",
			c: primitiveInts{
				A: -99,
				C: -99,
				E: -99,
				G: -99,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bytes := Serialize(tc.c)
			require.NotEmpty(t, bytes)

			var obj primitiveInts
			err := DeserializeRaw(bytes, &obj)
			require.NoError(t, err)
			require.Equal(t, tc.c, obj)
		})
	}
}

type hasEveryType struct {
	A int8
	B int16
	C int32
	D int64
	E uint8
	F uint16
	G uint32
	H uint64
	I bool
	J byte
	K string
	L []byte   // slice, byte type
	M []int64  // slice, non-byte type
	N [3]byte  // array, byte type
	O [3]int64 // array, non-byte type
	P struct {
		A int8
		B uint16
	} // struct
	Q map[string]string // map
	R float32
	S float64
}

func TestEncodeStable(t *testing.T) {
	// Tests encoding against previously encoded data on disk to verify
	// that encoding results have not changed
	update := false

	x := hasEveryType{
		A: -127,
		B: math.MaxInt16,
		C: math.MaxInt32,
		D: math.MaxInt64,
		E: math.MaxInt8 + 1,
		F: math.MaxInt16 + 1,
		G: math.MaxInt32 + 1,
		H: math.MaxInt64 + 1,
		I: true,
		J: byte(128),
		K: "foo",
		L: []byte("bar"),
		M: []int64{math.MaxInt64, math.MaxInt64 / 2, -10000},
		N: [3]byte{'b', 'a', 'z'},
		O: [3]int64{math.MaxInt64, math.MaxInt64 / 2, -10000},
		P: struct {
			A int8
			B uint16
		}{
			A: -127,
			B: math.MaxUint16,
		},
		Q: map[string]string{"foo": "bar"},
		R: float32(123.45),
		S: float64(123.45),
	}

	goldenFile := "testdata/encode-every-type.golden"

	if update {
		f, err := os.Create(goldenFile)
		require.NoError(t, err)
		defer f.Close()

		b := Serialize(x)
		_, err = f.Write(b)
		require.NoError(t, err)
		return
	}

	f, err := os.Open(goldenFile)
	require.NoError(t, err)
	defer f.Close()

	d, err := ioutil.ReadAll(f)
	require.NoError(t, err)

	var y hasEveryType
	err = DeserializeRaw(d, &y)
	require.NoError(t, err)
	require.Equal(t, x, y)

	b := Serialize(x)
	require.Equal(t, len(d), len(b))
	require.Equal(t, d, b)
}

func TestEncodeByteSlice(t *testing.T) {
	type foo struct {
		W int8
		X []byte
		Y int8 // these are added to make sure extra fields don't interact with the byte encoding
	}

	f := foo{
		W: 1,
		X: []byte("abc"),
		Y: 2,
	}

	expect := []byte{1, 3, 0, 0, 0, 97, 98, 99, 2}

	b := Serialize(f)
	require.Equal(t, expect, b)
}

func TestEncodeByteArray(t *testing.T) {
	type foo struct {
		W int8
		X [3]byte
		Y int8 // these are added to make sure extra fields don't interact with the byte encoding
	}

	f := foo{
		W: 1,
		X: [3]byte{'a', 'b', 'c'},
		Y: 2,
	}

	expect := []byte{1, 97, 98, 99, 2}

	b := Serialize(f)
	require.Equal(t, expect, b)
}

func TestEncodeEmptySlice(t *testing.T) {
	// Decoding an empty slice should not allocate
	type foo struct {
		X []byte
		Y []int64
	}

	f := &foo{}
	b := Serialize(f)

	var g foo
	err := DeserializeRaw(b, &g)
	require.NoError(t, err)
	require.Nil(t, g.X)
	require.Nil(t, g.Y)
}

func TestRandomGarbage(t *testing.T) {
	// Basic fuzz test to check for panics, deserializes random data

	// initialize the struct with data in the variable sized fields
	x := hasEveryType{
		K: "string",
		L: []byte("bar"),
		M: []int64{math.MaxInt64, math.MaxInt64 / 2, -10000},
		Q: map[string]string{"foo": "bar", "cat": "dog"},
	}

	size, err := datasizeWrite(reflect.ValueOf(x))
	require.NoError(t, err)

	var y hasEveryType
	for j := 0; j < 100; j++ {
		for i := 0; i < size*2; i++ {
			b := randBytes(t, i)
			DeserializeRaw(b, &y) // nolint: errcheck
		}
	}

	for i := 0; i < 10000; i++ {
		b := randBytes(t, size)
		DeserializeRaw(b, &y) // nolint: errcheck
	}
}
