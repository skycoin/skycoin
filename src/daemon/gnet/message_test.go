package gnet

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/cipher/encoder"
)

func TestNewMessageContext(t *testing.T) {
	c := &Connection{}
	mc := NewMessageContext(c)
	require.Equal(t, mc.ConnID, c.ID)
}

func TestRegisterMessage(t *testing.T) {
	EraseMessages()
	RegisterMessage(DummyPrefix, DummyMessage{})
	require.Equal(t, len(MessageIDMap), 1)
	require.Equal(t, len(MessageIDReverseMap), 1)
	require.NotNil(t, MessageIDReverseMap[DummyPrefix])

	RegisterMessage(ErrorPrefix, ErrorMessage{})
	require.Equal(t, len(MessageIDMap), 2)
	require.Equal(t, len(MessageIDReverseMap), 2)
	require.NotNil(t, MessageIDReverseMap[ErrorPrefix])
}

func TestEraseMessages(t *testing.T) {
	EraseMessages()
	RegisterMessage(DummyPrefix, DummyMessage{})
	require.Equal(t, len(MessageIDMap), 1)
	require.Equal(t, len(MessageIDReverseMap), 1)
	EraseMessages()
	require.Equal(t, len(MessageIDMap), 0)
	require.Equal(t, len(MessageIDReverseMap), 0)
}

func TestVerifyMessages(t *testing.T) {
	// VerifyMessages either no-ops or panics. Make sure it doesnt panic
	EraseMessages()
	RegisterMessage(DummyPrefix, DummyMessage{})
	RegisterMessage(ErrorPrefix, ErrorMessage{})
	require.NotPanics(t, VerifyMessages)
}

func TestVerifyMessagesDuplicateRegistered(t *testing.T) {
	EraseMessages()
	RegisterMessage(DummyPrefix, DummyMessage{})
	require.Panics(t, func() { RegisterMessage(DummyPrefix, DummyMessage{}) })
	require.Panics(t, func() { RegisterMessage(BytePrefix, DummyMessage{}) })
}

func TestVerifyMessagesNotAMessage(t *testing.T) {
	EraseMessages()
	RegisterMessage(NothingPrefix, Nothing{})
	require.Panics(t, VerifyMessages)
}

func TestVerifyMessagesBadPrefix(t *testing.T) {
	EraseMessages()
	// Can't be all null
	RegisterMessage(MessagePrefix{0x00, 0x00, 0x00, 0x00}, DummyMessage{})
	require.Panics(t, VerifyMessages)
	EraseMessages()
	// Can't start with nul
	RegisterMessage(MessagePrefix{0x00, 'A', 'A', 'A'}, DummyMessage{})
	require.Panics(t, VerifyMessages)
	EraseMessages()
	// Can't have non nul after nul
	RegisterMessage(MessagePrefix{'A', 0x00, 'A', 'A'}, DummyMessage{})
	require.Panics(t, VerifyMessages)
	EraseMessages()
	// Can't have invalid ascii bytes
	RegisterMessage(MessagePrefix{'A', 'A', 'A', '\n'}, DummyMessage{})
	require.Panics(t, VerifyMessages)
	EraseMessages()
	RegisterMessage(MessagePrefix{'A', 'A', 'A', '@'}, DummyMessage{})
	require.Panics(t, VerifyMessages)
	EraseMessages()
	RegisterMessage(MessagePrefix{'A', 'A', 'A', ' '}, DummyMessage{})
	require.Panics(t, VerifyMessages)
	EraseMessages()
	RegisterMessage(MessagePrefix{'A', 'A', 'A', '{'}, DummyMessage{})
	require.Panics(t, VerifyMessages)
	EraseMessages()
	RegisterMessage(MessagePrefix{'A', 'A', 'A', '`'}, DummyMessage{})
	require.Panics(t, VerifyMessages)
	EraseMessages()
	RegisterMessage(MessagePrefix{'A', 'A', 'A', '['}, DummyMessage{})
	require.Panics(t, VerifyMessages)
	EraseMessages()
	RegisterMessage(MessagePrefix{'A', 'A', 'A', '/'}, DummyMessage{})
	require.Panics(t, VerifyMessages)
	EraseMessages()
	RegisterMessage(MessagePrefix{'A', 'A', 'A', ':'}, DummyMessage{})
	require.Panics(t, VerifyMessages)
	EraseMessages()
	// Some valid messages
	RegisterMessage(MessagePrefix{'1', '9', 'A', 'z'}, DummyMessage{})
	require.NotPanics(t, VerifyMessages)
	EraseMessages()
}

func TestVerifyMessagesCorruptMap(t *testing.T) {
	// MessageIdMap circumvented
	EraseMessages()
	mtype := reflect.TypeOf(DummyMessage{})
	MessageIDMap[mtype] = DummyPrefix
	require.Panics(t, VerifyMessages)
	delete(MessageIDMap, mtype)
	// MessageIdReverseMap circumvented
	EraseMessages()
	MessageIDReverseMap[DummyPrefix] = mtype
	require.Panics(t, VerifyMessages)
	delete(MessageIDReverseMap, DummyPrefix)
}

func TestMessagePrefixFromString(t *testing.T) {
	EraseMessages()
	require.Panics(t, func() { MessagePrefixFromString("") })
	require.Panics(t, func() { MessagePrefixFromString("xxxxx") })
	require.Equal(t, MessagePrefixFromString("abcd"), MessagePrefix{'a', 'b', 'c', 'd'})
	require.Equal(t, MessagePrefixFromString("abc"), MessagePrefix{'a', 'b', 'c', 0x00})
	require.Equal(t, MessagePrefixFromString("ab"), MessagePrefix{'a', 'b', 0x00, 0x00})
	require.Equal(t, MessagePrefixFromString("a"), MessagePrefix{'a', 0x00, 0x00, 0x00})
}

/* Helpers */

type Nothing struct{}

var NothingPrefix = MessagePrefix{'N', 'O', 'T', 'H'}

type DummyMessage struct{}

var DummyPrefix = MessagePrefix{'D', 'U', 'M', 'Y'}

// EncodeSize implements gnet.Serializer
func (dm *DummyMessage) EncodeSize() uint64 {
	return uint64(encoder.Size(dm))
}

// Encode implements gnet.Serializer
func (dm *DummyMessage) Encode(buf []byte) error {
	buf2 := encoder.Serialize(dm)
	if len(buf) < len(buf2) {
		return errors.New("Not enough buffer data to encode")
	}
	copy(buf[:], buf2[:])
	return nil
}

// Decode implements gnet.Serializer
func (dm *DummyMessage) Decode(buf []byte) (uint64, error) {
	return encoder.DeserializeRaw(buf, dm)
}

func (dm *DummyMessage) Handle(context *MessageContext, x interface{}) error {
	return nil
}

type ErrorMessage struct{}

var (
	ErrorPrefix            = MessagePrefix{'E', 'R', 'R', 0x00}
	ErrErrorMessageHandler = errors.New("Bad")
)

// EncodeSize implements gnet.Serializer
func (em *ErrorMessage) EncodeSize() uint64 {
	return uint64(encoder.Size(em))
}

// Encode implements gnet.Serializer
func (em *ErrorMessage) Encode(buf []byte) error {
	buf2 := encoder.Serialize(em)
	if len(buf) < len(buf2) {
		return errors.New("Not enough buffer data to encode")
	}
	copy(buf[:], buf2[:])
	return nil
}

// Decode implements gnet.Serializer
func (em *ErrorMessage) Decode(buf []byte) (uint64, error) {
	return encoder.DeserializeRaw(buf, em)
}

func (em *ErrorMessage) Handle(context *MessageContext, x interface{}) error {
	return ErrErrorMessageHandler
}

type ByteMessage struct {
	X byte
}

var BytePrefix = MessagePrefix{'B', 'Y', 'T', 'E'}

// EncodeSize implements gnet.Serializer
func (bm *ByteMessage) EncodeSize() uint64 {
	return uint64(encoder.Size(bm))
}

// Encode implements gnet.Serializer
func (bm *ByteMessage) Encode(buf []byte) error {
	buf2 := encoder.Serialize(bm)
	if len(buf) < len(buf2) {
		return errors.New("Not enough buffer data to encode")
	}
	copy(buf[:], buf2[:])
	return nil
}

// Decode implements gnet.Serializer
func (bm *ByteMessage) Decode(buf []byte) (uint64, error) {
	return encoder.DeserializeRaw(buf, bm)
}

func (bm *ByteMessage) Handle(c *MessageContext, x interface{}) error {
	return nil
}

func NewByteMessage(x byte) Message {
	return &ByteMessage{X: x}
}

type PointerMessage struct {
	Ptr *int
}

// EncodeSize implements gnet.Serializer
func (m *PointerMessage) EncodeSize() uint64 {
	return uint64(encoder.Size(m))
}

// Encode implements gnet.Serializer
func (m *PointerMessage) Encode(buf []byte) error {
	buf2 := encoder.Serialize(m)
	if len(buf) < len(buf2) {
		return errors.New("Not enough buffer data to encode")
	}
	copy(buf[:], buf2[:])
	return nil
}

// Decode implements gnet.Serializer
func (m *PointerMessage) Decode(buf []byte) (uint64, error) {
	return encoder.DeserializeRaw(buf, m)
}
