package gnet

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessageContext(t *testing.T) {
	c := &Connection{}
	mc := NewMessageContext(c)
	assert.Equal(t, mc.ConnID, c.ID)
}

func TestRegisterMessage(t *testing.T) {
	EraseMessages()
	RegisterMessage(DummyPrefix, DummyMessage{})
	assert.Equal(t, len(MessageIDMap), 1)
	assert.Equal(t, len(MessageIDReverseMap), 1)
	assert.NotNil(t, MessageIDReverseMap[DummyPrefix])

	RegisterMessage(ErrorPrefix, ErrorMessage{})
	assert.Equal(t, len(MessageIDMap), 2)
	assert.Equal(t, len(MessageIDReverseMap), 2)
	assert.NotNil(t, MessageIDReverseMap[ErrorPrefix])
}

func TestEraseMessages(t *testing.T) {
	EraseMessages()
	RegisterMessage(DummyPrefix, DummyMessage{})
	assert.Equal(t, len(MessageIDMap), 1)
	assert.Equal(t, len(MessageIDReverseMap), 1)
	EraseMessages()
	assert.Equal(t, len(MessageIDMap), 0)
	assert.Equal(t, len(MessageIDReverseMap), 0)
}

func TestVerifyMessages(t *testing.T) {
	// VerifyMessages either no-ops or panics. Make sure it doesnt panic
	EraseMessages()
	RegisterMessage(DummyPrefix, DummyMessage{})
	RegisterMessage(ErrorPrefix, ErrorMessage{})
	assert.NotPanics(t, VerifyMessages)
}

func TestVerifyMessagesDuplicateRegistered(t *testing.T) {
	EraseMessages()
	RegisterMessage(DummyPrefix, DummyMessage{})
	assert.Panics(t, func() { RegisterMessage(DummyPrefix, DummyMessage{}) })
	assert.Panics(t, func() { RegisterMessage(BytePrefix, DummyMessage{}) })
}

func TestVerifyMessagesNotAMessage(t *testing.T) {
	EraseMessages()
	RegisterMessage(NothingPrefix, Nothing{})
	assert.Panics(t, VerifyMessages)
}

func TestVerifyMessagesBadPrefix(t *testing.T) {
	EraseMessages()
	// Can't be all null
	RegisterMessage(MessagePrefix{0x00, 0x00, 0x00, 0x00}, DummyMessage{})
	assert.Panics(t, VerifyMessages)
	EraseMessages()
	// Can't start with nul
	RegisterMessage(MessagePrefix{0x00, 'A', 'A', 'A'}, DummyMessage{})
	assert.Panics(t, VerifyMessages)
	EraseMessages()
	// Can't have non nul after nul
	RegisterMessage(MessagePrefix{'A', 0x00, 'A', 'A'}, DummyMessage{})
	assert.Panics(t, VerifyMessages)
	EraseMessages()
	// Can't have invalid ascii bytes
	RegisterMessage(MessagePrefix{'A', 'A', 'A', '\n'}, DummyMessage{})
	assert.Panics(t, VerifyMessages)
	EraseMessages()
	RegisterMessage(MessagePrefix{'A', 'A', 'A', '@'}, DummyMessage{})
	assert.Panics(t, VerifyMessages)
	EraseMessages()
	RegisterMessage(MessagePrefix{'A', 'A', 'A', ' '}, DummyMessage{})
	assert.Panics(t, VerifyMessages)
	EraseMessages()
	RegisterMessage(MessagePrefix{'A', 'A', 'A', '{'}, DummyMessage{})
	assert.Panics(t, VerifyMessages)
	EraseMessages()
	RegisterMessage(MessagePrefix{'A', 'A', 'A', '`'}, DummyMessage{})
	assert.Panics(t, VerifyMessages)
	EraseMessages()
	RegisterMessage(MessagePrefix{'A', 'A', 'A', '['}, DummyMessage{})
	assert.Panics(t, VerifyMessages)
	EraseMessages()
	RegisterMessage(MessagePrefix{'A', 'A', 'A', '/'}, DummyMessage{})
	assert.Panics(t, VerifyMessages)
	EraseMessages()
	RegisterMessage(MessagePrefix{'A', 'A', 'A', ':'}, DummyMessage{})
	assert.Panics(t, VerifyMessages)
	EraseMessages()
	// Some valid messages
	RegisterMessage(MessagePrefix{'1', '9', 'A', 'z'}, DummyMessage{})
	assert.NotPanics(t, VerifyMessages)
	EraseMessages()
}

func TestVerifyMessagesCorruptMap(t *testing.T) {
	// MessageIdMap circumvented
	EraseMessages()
	mtype := reflect.TypeOf(DummyMessage{})
	MessageIDMap[mtype] = DummyPrefix
	assert.Panics(t, VerifyMessages)
	delete(MessageIDMap, mtype)
	// MessageIdReverseMap circumvented
	EraseMessages()
	MessageIDReverseMap[DummyPrefix] = mtype
	assert.Panics(t, VerifyMessages)
	delete(MessageIDReverseMap, DummyPrefix)
}

func TestMessagePrefixFromString(t *testing.T) {
	EraseMessages()
	assert.Panics(t, func() { MessagePrefixFromString("") })
	assert.Panics(t, func() { MessagePrefixFromString("xxxxx") })
	assert.Equal(t, MessagePrefixFromString("abcd"),
		MessagePrefix{'a', 'b', 'c', 'd'})
	assert.Equal(t, MessagePrefixFromString("abc"),
		MessagePrefix{'a', 'b', 'c', 0x00})
	assert.Equal(t, MessagePrefixFromString("ab"),
		MessagePrefix{'a', 'b', 0x00, 0x00})
	assert.Equal(t, MessagePrefixFromString("a"),
		MessagePrefix{'a', 0x00, 0x00, 0x00})
}

/* Helpers */

type Nothing struct{}

var NothingPrefix = MessagePrefix{'N', 'O', 'T', 'H'}

type DummyMessage struct{}

var DummyPrefix = MessagePrefix{'D', 'U', 'M', 'Y'}

func (dm *DummyMessage) Handle(context *MessageContext, x interface{}) error {
	return nil
}

func NewDummyMessage() Message {
	return &DummyMessage{}
}

type ErrorMessage struct{}

var ErrorPrefix = MessagePrefix{'E', 'R', 'R', 0x00}

func (em *ErrorMessage) Handle(context *MessageContext, x interface{}) error {
	return errors.New("Bad")
}

func NewErrorMessage() Message {
	return &ErrorMessage{}
}

type ByteMessage struct {
	X byte
}

var BytePrefix = MessagePrefix{'B', 'Y', 'T', 'E'}

func (bm *ByteMessage) Handle(c *MessageContext, x interface{}) error {
	return nil
}

func NewByteMessage(x byte) Message {
	return &ByteMessage{X: x}
}

type PointerMessage struct {
	Ptr *int
}
