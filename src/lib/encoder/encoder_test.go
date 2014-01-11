package encoder

import (
	"bytes"
	"reflect"
	"testing"
)

import (
	"crypto/rand"
)

func randBytes(n int) []byte {
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
}

type TestStruct2 struct {
	X int32
	Y int64
	Z uint8
	K [8]byte
}

//func (*B) Fatal

func Test_Encode_1(T *testing.T) { //test function starts with "Test" and takes a pointer to type testing.T
	var t TestStruct
	t.X = 345535
	t.Y = 23432435443
	t.Z = 255
	t.K = []byte("TEST6")

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

func Test_Encode_4(T *testing.T) { //test function starts with "Test" and takes a pointer to type testing.T
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
		T.Fatal()
	}

	if len(t1.A) != len(t2.A) {
		T.Fatal("t1.A= %v. t2.A= %v \n", len(t1.A), len(t2.A))
	}

	//serialize A for t1.K and t2.K
	//if bytes.Compare(t1.K, t2.K) != 0 {
	//}

	b2 := Serialize(t2)

	if bytes.Compare(b, b2) != 0 {
		T.Fatal()
	}
}

func Test_Encode_5(T *testing.T) {

	var ts TestStruct2
	ts.X = 345535
	ts.Y = 23432435443
	ts.Z = 255

	b1 := Serialize(ts)

	var t reflect.Type = reflect.TypeOf(ts)
	var v reflect.Value = reflect.New(t) //pointer to type t

	//New returns a Value representing a pointer to a new zero value for the specified type.
	//That is, the returned Value's Type is PtrTo(t).

	err := DeserializeRawToValue(b1, v)
	if err != nil {
		T.Fatal(err)
	}

}
