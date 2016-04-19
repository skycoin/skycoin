
package mesh

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

func init() {
	RegisterMessageForSerialization(messagePrefix{44}, TestStruct{})
}

func TestSerializeStruct(t *testing.T) {
	testStruct := TestStruct {
		"hello",
		8,
		false,
		11,
	}
	data := SerializeMessage(testStruct)
	assert.NotEqual(t, 0, len(data))
	result, error := UnserializeMessage(data)
	assert.Nil(t, error)
	assert.Equal(t, testStruct, result)
}

