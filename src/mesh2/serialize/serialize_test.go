
package serialize

import(
	"testing"
	)

import(
	"github.com/stretchr/testify/assert"
	)

type TestStruct struct {
	Foo string
	Bar uint32
	Gah bool
	Blah uint32
}

var serializer *Serializer

func init() {
	serializer = NewSerializer()
	serializer.RegisterMessageForSerialization(MessagePrefix{44}, TestStruct{})
}

func TestSerializeStruct(t *testing.T) {
	testStruct := TestStruct {
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

