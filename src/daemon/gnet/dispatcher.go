package gnet

import (
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"time"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// Result of a single message send
type SendResult struct {
	Connection *Connection
	Message    Message
	Error      error
}

func newSendResult(c *Connection, m Message, err error) SendResult {
	return SendResult{
		Connection: c,
		Message:    m,
		Error:      err,
	}
}

// Serializes a Message over a net.Conn
func sendMessage(conn net.Conn, msg Message, timeout time.Duration) error {
	m := encodeMessage(msg)
	return sendByteMessage(conn, m, timeout)
}

// Event handler that is called after a Connection sends a complete message
func convertToMessage(c *Connection, msg []byte) (Message, error) {
	msgId := [4]byte{}
	if len(msg) < len(msgId) {
		return nil, errors.New("Not enough data to read msg id")
	}
	copy(msgId[:], msg[:len(msgId)])
	msg = msg[len(msgId):]

	t, succ := MessageIdReverseMap[msgId]
	if !succ {
		logger.Debug("Connection %d sent unknown message id %s",
			c.Id, string(msgId[:]))
		return nil, fmt.Errorf("Unknown message %s received", string(msgId[:]))
	}

	if DebugPrint {
		logger.Debug("Convert, Message type %v", t)
	}

	var m Message
	var v reflect.Value = reflect.New(t)
	//logger.Debug("Giving %d bytes to the decoder", len(msg))
	used, err := deserializeMessage(msg, v)
	if err != nil {
		return nil, err
	}
	if used != len(msg) {
		return nil, errors.New("Data buffer was not completely decoded")
	}

	m, succ = (v.Interface()).(Message)
	if !succ {
		// This occurs only when the user registers an interface that does
		// match the Message interface.  They should have known about this
		// earlier via a call to VerifyMessages
		log.Panic("Message obtained from map does not match Message interface")
		return nil, errors.New("MessageIdMaps contain non-Message")
	}
	return m, nil
}

// Wraps encoder.DeserializeRawToValue and traps panics as an error
func deserializeMessage(msg []byte, v reflect.Value) (n int, e error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Debug("Recovering from deserializer panic: %v", r)
			switch x := r.(type) {
			case string:
				e = errors.New(x)
			case error:
				e = x
			default:
				e = errors.New("Message deserialization failed")
			}
		}
	}()
	n, e = encoder.DeserializeRawToValue(msg, v)
	return
}

// Packgs a Message into []byte containing length, id and data
var encodeMessage = func(msg Message) []byte {
	t := reflect.ValueOf(msg).Elem().Type()
	msgId, succ := MessageIdMap[t]
	if !succ {
		txt := "Attempted to serialize message struct not in MessageIdMap: %v"
		log.Panicf(txt, msg)
	}
	bMsg := encoder.Serialize(msg)

	// message length
	bLen := encoder.SerializeAtomic(uint32(len(bMsg) + len(msgId)))
	m := make([]byte, 0)
	m = append(m, bLen...)     // length prefix
	m = append(m, msgId[:]...) // message id
	m = append(m, bMsg...)     // message bytes
	return m
}

// Sends []byte over a net.Conn
var sendByteMessage = func(conn net.Conn, msg []byte,
	timeout time.Duration) error {
	deadline := time.Time{}
	if timeout != 0 {
		deadline = time.Now().Add(timeout)
	}
	if err := conn.SetWriteDeadline(deadline); err != nil {
		return err
	}
	if _, err := conn.Write(msg); err != nil {
		return err
	}
	return nil
}
