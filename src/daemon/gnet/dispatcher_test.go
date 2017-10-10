package gnet

import (
	"bytes"
	"errors"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	_sendByteMessage = sendByteMessage
	_encodeMessage   = encodeMessage
)

func resetHandler() {
	sendByteMessage = _sendByteMessage
	encodeMessage = _encodeMessage
}

func TestConvertToMessage(t *testing.T) {
	EraseMessages()
	resetHandler()
	RegisterMessage(BytePrefix, ByteMessage{})
	VerifyMessages()
	c := &Connection{}
	b := make([]byte, 0)
	b = append(b, BytePrefix[:]...)
	b = append(b, byte(7))
	m, err := convertToMessage(c.ID, b, testing.Verbose())
	assert.Nil(t, err)
	assert.NotNil(t, m)
	if m == nil {
		t.Fatalf("ConvertToMessage failed")
	}
	bm := m.(*ByteMessage)
	assert.Equal(t, bm.X, byte(7))
}

func TestConvertToMessageNoMessageId(t *testing.T) {
	EraseMessages()
	resetHandler()
	c := &Connection{}
	b := []byte{}
	m, err := convertToMessage(c.ID, b, testing.Verbose())
	assert.Nil(t, m)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Not enough data to read msg id")
}

func TestConvertToMessageUnknownMessage(t *testing.T) {
	EraseMessages()
	resetHandler()
	c := &Connection{}
	b := MessagePrefix{'C', 'C', 'C', 'C'}
	m, err := convertToMessage(c.ID, b[:], testing.Verbose())
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Unknown message CCCC received")
	assert.Nil(t, m)
}

func TestConvertToMessageBadDeserialize(t *testing.T) {
	EraseMessages()
	resetHandler()
	RegisterMessage(DummyPrefix, DummyMessage{})
	RegisterMessage(BytePrefix, ByteMessage{})
	VerifyMessages()
	c := &Connection{}
	// Test with too many bytes
	b := append(DummyPrefix[:], []byte{0, 1, 1, 1}...)
	m, err := convertToMessage(c.ID, b, testing.Verbose())
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Data buffer was not completely decoded")
	assert.Nil(t, m)

	// Test with not enough bytes
	b = append([]byte{}, BytePrefix[:]...)
	m, err = convertToMessage(c.ID, b, testing.Verbose())
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Deserialization failed")
	assert.Nil(t, m)
}

func TestConvertToMessageNotMessage(t *testing.T) {
	EraseMessages()
	resetHandler()
	RegisterMessage(NothingPrefix, Nothing{})
	// don't verify messages
	c := &Connection{}
	assert.Panics(t, func() {
		convertToMessage(c.ID, NothingPrefix[:], testing.Verbose())
	})
}

func TestDeserializeMessageTrapsPanic(t *testing.T) {
	resetHandler()
	EraseMessages()
	p := 7
	m := PointerMessage{Ptr: &p}
	b := []byte{4, 4, 4, 4, 4, 4, 4, 4}
	_, err := deserializeMessage(b, reflect.ValueOf(m))
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(),
		"Decode error: kind invalid not handled")
}

func TestEncodeMessage(t *testing.T) {
	resetHandler()
	EraseMessages()
	RegisterMessage(BytePrefix, ByteMessage{})
	VerifyMessages()
	m := NewByteMessage(7)
	b := encodeMessage(m)
	assert.True(t, bytes.Equal(b, []byte{5, 0, 0, 0, 'B', 'Y', 'T', 'E', 7}))
}

func TestEncodeMessageUnknownMessage(t *testing.T) {
	resetHandler()
	EraseMessages()
	assert.Panics(t, func() { encodeMessage(&DummyMessage{}) })
}

func TestSendByteMessage(t *testing.T) {
	resetHandler()
	b := []byte{1}
	c := NewCaptureConn()
	err := sendByteMessage(c, b, 0)
	assert.Nil(t, err)
	assert.True(t, bytes.Equal(c.(*CaptureConn).Wrote, b))
	assert.True(t, c.(*CaptureConn).WriteDeadlineSet)
}

func TestSendByteMessageWithTimeout(t *testing.T) {
	resetHandler()
	b := []byte{1}
	c := NewCaptureConn()
	err := sendByteMessage(c, b, time.Minute)
	assert.Nil(t, err)
	assert.True(t, bytes.Equal(c.(*CaptureConn).Wrote, b))
	assert.True(t, c.(*CaptureConn).WriteDeadlineSet)
}

func TestSendByteMessageWriteFailed(t *testing.T) {
	resetHandler()
	c := &FailingWriteConn{}
	err := sendByteMessage(c, nil, 0)
	assert.NotNil(t, err)
}

func TestSendByteMessageWriteDeadlineFailed(t *testing.T) {
	resetHandler()
	c := &FailingWriteDeadlineConn{}
	err := sendByteMessage(c, nil, 0)
	assert.NotNil(t, err)
}

func TestSendMessage(t *testing.T) {
	resetHandler()
	EraseMessages()
	RegisterMessage(BytePrefix, ByteMessage{})
	VerifyMessages()
	m := NewByteMessage(7)
	sendByteMessage = func(conn net.Conn, msg []byte, tm time.Duration) error {
		expect := []byte{5, 0, 0, 0, 'B', 'Y', 'T', 'E', 7}
		assert.True(t, bytes.Equal(msg, expect))
		return nil
	}
	err := sendMessage(nil, m, 0)
	assert.Nil(t, err)
}

/* Helpers */

func noopSendByteMessage(conn net.Conn, m []byte, tm time.Duration) error {
	return nil
}

func failingSendByteMessage(conn net.Conn, m []byte, tm time.Duration) error {
	return errors.New("failed")
}

func noopEncodeMessage(msg Message) []byte {
	return []byte{}
}

type CaptureConn struct {
	Wrote            []byte
	WriteDeadlineSet bool
	net.Conn
}

func NewCaptureConn() net.Conn {
	return &CaptureConn{Wrote: nil, WriteDeadlineSet: false}
}

func (cc *CaptureConn) Write(b []byte) (int, error) {
	cc.Wrote = b
	return len(b), nil
}

func (cc *CaptureConn) SetWriteDeadline(t time.Time) error {
	cc.WriteDeadlineSet = true
	return nil
}

type FailingWriteDeadlineConn struct {
	net.Conn
}

func (c *FailingWriteDeadlineConn) SetWriteDeadline(t time.Time) error {
	return errors.New("failed")
}

type FailingWriteConn struct {
	net.Conn
}

func (c *FailingWriteConn) Write(b []byte) (int, error) {
	return 0, errors.New("failed")
}

func (c *FailingWriteConn) SetWriteDeadline(t time.Time) error {
	return nil
}
