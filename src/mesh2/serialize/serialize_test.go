package serialize

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Foo  string
	Bar  uint32
	Gah  bool
	Blah uint32
}

type TestWithBytes struct {
	Blah  uint32
	Slice []byte
}

var serializer *Serializer

func init() {
	serializer = NewSerializer()
	serializer.RegisterMessageForSerialization(MessagePrefix{44}, TestStruct{})
	serializer.RegisterMessageForSerialization(MessagePrefix{88}, TestWithBytes{})
}

func TestSerializeStruct(t *testing.T) {
	testStruct := TestStruct{
		"hello",
		8,
		false,
		11,
	}
	data := serializer.SerializeMessage(testStruct)
	assert.NotEqual(t, 0, len(data))
	result, error := serializer.UnserializeMessage(data)
	assert.Nil(t, error)
	assert.Equal(t, testStruct, result)
}

func TestSerializeNilBytes(t *testing.T) {
	testA := TestWithBytes{55, nil}
	testB := TestWithBytes{55, []byte{44, 55, 1, 2, 3}}
	testC := TestWithBytes{55, []byte{}}
	dataA := serializer.SerializeMessage(testA)
	dataB := serializer.SerializeMessage(testB)
	dataC := serializer.SerializeMessage(testC)
	assert.Equal(t, 5, len(dataB)-len(dataA))
	assert.Equal(t, len(dataA), len(dataC))
}
