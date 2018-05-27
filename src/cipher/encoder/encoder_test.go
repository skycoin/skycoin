package encoder

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
)

func randBytes(n int) []byte { // nolint: unparam
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return bytes
}

/*
* the file name has to end with _test.go to be picked up as a set of tests by go test
* the package name has to be the same as in the source file that has to be tested
* you have to import the package testing
* all test functions should start with Test to be run as a test
* the tests will be executed in the same order that they are appear in the source
* the test function TestXxx functions take a pointer to the type testing.T. You use it to record the test status and also for logging.
* the signature of the test function should always be func TestXxx ( *testing.T). You can have any combination of alphanumeric characters and the hyphen for the Xxx part, the only constraint that it should not begin with a small alphabet, [a-z].
* a call to any of the following functions of testing.T within the test code Error, Errorf, FailNow, Fatal, FatalIf will indicate to go test that the test has failed.
 */

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

func Test_Encode_1(T *testing.T) { //test function starts with "Test" and takes a pointer to type testing.T
	var t TestStruct
	t.X = 345535
	t.Y = 23432435443
	t.Z = 255
	t.K = []byte("TEST6")
	t.W = true
	t.T = "hello"
	t.U = cipher.PubKey{1, 2, 3, 0, 5, 4, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	b := Serialize(t)

	var buf bytes.Buffer
	buf.Write(b)

	var t2 TestStruct
	err := Deserialize(&buf, len(b), &t2)
	if err != nil {
		T.Fatal(err)
	}

	b2 := Serialize(t2)

	if bytes.Compare(b, b2) != 0 {
		T.Fatal()
	}
}

func Test_Encode_2a(T *testing.T) { //test function starts with "Test" and takes a pointer to type testing.T
	var t TestStruct2
	t.X = 345535
	t.Y = 23432435443
	t.Z = 255
	t.W = false
	_tt := []byte("ASDSADFSDFASDFSD")
	for i := 0; i < 8; i++ {
		t.K[i] = _tt[i]
	}

	b := Serialize(t)

	var buf bytes.Buffer
	buf.Write(b)

	var t2 TestStruct2
	err := Deserialize(&buf, len(b), &t2)
	if err != nil {
		T.Fatal(err)
	}

	b2 := Serialize(t2)

	if bytes.Compare(b, b2) != 0 {
		T.Fatal()
	}
}

func Test_Encode_2b(T *testing.T) { //test function starts with "Test" and takes a pointer to type testing.T
	var t TestStruct2
	t.X = 345535
	t.Y = 23432435443
	t.Z = 255
	_tt := []byte("ASDSADFSDFASDFSD")
	for i := 0; i < 8; i++ {
		t.K[i] = _tt[i]
	}

	b := Serialize(t)

	var t2 TestStruct2
	err := DeserializeRaw(b, &t2)
	if err != nil {
		T.Fatal(err)
	}

	b2 := Serialize(t2)

	if bytes.Compare(b, b2) != 0 {
		T.Fatal()
	}
}

type TestStruct3 struct {
	X int32
	K []byte
}

func Test_Encode_3a(T *testing.T) { //test function starts with "Test" and takes a pointer to type testing.T
	var t1 TestStruct3
	t1.X = 345535
	t1.K = randBytes(32)

	b := Serialize(t1)

	var buf bytes.Buffer
	buf.Write(b)

	var t2 TestStruct3
	err := Deserialize(&buf, len(b), &t2)
	if err != nil {
		T.Fatal(err)
	}

	if t1.X != t2.X || len(t1.K) != len(t2.K) || bytes.Compare(t1.K, t2.K) != 0 {
		T.Fatal()
	}

	b2 := Serialize(t2)

	if bytes.Compare(b, b2) != 0 {
		T.Fatal()
	}
}

func Test_Encode_3b(T *testing.T) { //test function starts with "Test" and takes a pointer to type testing.T
	var t1 TestStruct3
	t1.X = 345535
	t1.K = randBytes(32)

	b := Serialize(t1)

	var t2 TestStruct3
	err := DeserializeRaw(b, &t2)
	if err != nil {
		T.Fatal(err)
	}

	if t1.X != t2.X || len(t1.K) != len(t2.K) || bytes.Compare(t1.K, t2.K) != 0 {
		T.Fatal()
	}

	b2 := Serialize(t2)

	if bytes.Compare(b, b2) != 0 {
		T.Fatal()
	}
}

type TestStruct4 struct {
	X int32
	Y int32
}

type TestStruct5 struct {
	X int32
	A []TestStruct4
}

func Test_Encode_4(T *testing.T) {
	var t1 TestStruct5
	t1.X = 345535

	const NUM = 8
	t1.A = make([]TestStruct4, NUM)

	b := Serialize(t1)

	var t2 TestStruct5
	err := DeserializeRaw(b, &t2)
	if err != nil {
		T.Fatal(err)
	}

	if t1.X != t2.X {
		T.Fatal("TestStruct5.X not equal")
	}

	if len(t1.A) != len(t2.A) {
		T.Fatal("Slice lengths not equal")
	}

	for i, ts := range t1.A {
		if ts != t2.A[i] {
			T.Fatal("Slice values not equal")
		}
	}

	b2 := Serialize(t2)

	if bytes.Compare(b, b2) != 0 {
		T.Fatal()
	}
}

// type TestStruct2 struct {
//     X   int32
//     Y   int64
//     Z   uint8
//     K   [8]byte
// }

func Test_Encode_5(T *testing.T) {

	var ts TestStruct2
	ts.X = 345535
	ts.Y = 23432435443
	ts.Z = 255

	b1 := Serialize(ts)

	var t = reflect.TypeOf(ts)
	var v = reflect.New(t) //pointer to type t

	//New returns a Value representing a pointer to a new zero value for the specified type.
	//That is, the returned Value's Type is PtrTo(t).

	_, err := DeserializeRawToValue(b1, v)
	if err != nil {
		T.Fatal(err)
	}

	v = reflect.Indirect(v)
	if v.FieldByName("X").Int() != int64(ts.X) {
		T.Fatalf("X not equal")
	}
	if v.FieldByName("Y").Int() != ts.Y {
		T.Fatalf("Y not equal")
	}
	if v.FieldByName("Z").Uint() != uint64(ts.Z) {
		T.Fatalf("Z not equal")
	}
}

func Test_Encode_IgnoreTagSerialize(T *testing.T) {
	var t TestStructIgnore
	t.X = 345535
	t.Y = 23432435443
	t.Z = 255
	t.K = []byte("TEST6")

	b := Serialize(t)
	var buf bytes.Buffer
	buf.Write(b)

	var t2 TestStructIgnore
	t.X = 0
	t.Y = 0
	t.Z = 0
	t.K = []byte("")
	err := Deserialize(&buf, len(b), &t2)
	if err != nil {
		T.Fatal(err)
	}

	if t2.Z != 0 {
		T.Fatalf("Z should not deserialize. It is %d", t2.Z)
	}

	buf.Reset()
	buf.Write(b)

	var t3 TestStructWithoutIgnore
	err = Deserialize(&buf, len(b), &t3)
	if err != nil {
		T.Fatal(err)
	}

	b2 := Serialize(t2)
	if bytes.Compare(b, b2) != 0 {
		T.Fatal()
	}
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
	if err != nil {
		t.Fatalf("datasizeWrite failed: %v", err)
	}
	if n != size+4 {
		t.Fatal("Wrong data size")
	}
	b := Serialize(c)
	d := Container{}
	err = DeserializeRaw(b, &d)
	if err != nil {
		t.Fatalf("DeserializeRaw failed: %v", err)
	}
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
	if err == nil {
		t.Fatal("Expected error")
	} else if err.Error() != "Deserialization failed" {
		t.Fatalf("Expected different error, but got %s", err.Error())
	}

	// Test with slice
	thing := make([]int, 3)
	err = DeserializeRaw(b, thing)
	if err == nil {
		t.Fatal("Expected error")
	} else if err.Error() != "Deserialization failed" {
		t.Fatal("Expected different error")
	}
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

func TestMultiArrays(T *testing.T) {
	var data [16][16]byte
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			data[i][j] = byte(i * j)
		}
	}

	b := Serialize(data)

	var data2 [16][16]byte

	err := DeserializeRaw(b, &data2)
	if err != nil {
		T.Fatal(err)
	}

	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			if data[i][j] != data2[i][j] {
				T.Fatalf("failed round trip test")
			}
		}
	}

	b2 := Serialize(data2)
	if !bytes.Equal(b, b2) {
		T.Fatalf("Failed round trip test")
	}

	if len(b) != 256 {
		T.Fatalf("decoded to wrong byte length")
	}

}

func TestSerializeAtomic(t *testing.T) {

	var sp uint64 = 0x000C8A9E1809F720
	b := SerializeAtomic(sp)

	var i uint64
	DeserializeAtomic(b, &i)

	if i != sp {
		t.Fatal("round trip atomic fail")
	}
}

func TestPushPop(t *testing.T) {
	var sp uint64 = 0x000C8A9E1809F720

	var d [8]byte
	EncodeInt(d[0:8], sp)

	//fmt.Printf("d= %X \n", d[:])

	var ti uint64
	DecodeInt(d[0:8], &ti)

	if ti != sp {
		//fmt.Printf("sp= %X ti= %X \n", sp,ti)
		t.Error("roundtrip failed")
	}
}

type TestStruct5a struct {
	Test uint64
}

func TestPanicTest(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Error("EncodeInt Did not panic")
		}
	}()

	log.Panic()
}

func TestPushPopNegative(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Error("EncodeInt Did not panic on invalid input type")
		}
	}()

	var tst TestStruct5a
	//var sp uint64 = 0x000C8A9E1809F720
	var d [8]byte
	EncodeInt(d[0:8], &tst) //attemp to encode invalid type

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
		0x01234567,
		0x0123456789ABCDEF,
		map[uint8]TestStruct{
			0x01: TestStruct{
				0x01234567,
				0x0123456789ABCDEF,
				0x01,
				[]byte{0, 1, 2},
				true,
				"ab",
				cipher.PubKey{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
					17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
				},
			},
			0x23: TestStruct{
				0x01234567,
				0x0123456789ABCDEF,
				0x01,
				[]byte{0, 1, 2},
				true,
				"cd",
				cipher.PubKey{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
					17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
				},
			},
		},
		[]byte{0, 1, 2, 3, 4},
	}
	buff := Serialize(s1)
	if len(buff) == 0 {
		t.Fail()
	}

	s2 := TestStructWithDict{}
	if DeserializeRaw(buff, &s2) != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(s1, s2) {
		t.Errorf("Expected %v but got %v", s1, s2)
	}
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
