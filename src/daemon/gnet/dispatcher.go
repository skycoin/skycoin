package gnet

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// SendResult result of a single message send
type SendResult struct {
	Addr    string
	Message Message
	Error   error
}

func newSendResult(addr string, m Message, err error) SendResult {
	return SendResult{
		Addr:    addr,
		Message: m,
		Error:   err,
	}
}

// Serializes a Message over a net.Conn
func sendMessage(conn net.Conn, msg Message, timeout time.Duration) error {
	m := encodeMessage(msg)
	return sendByteMessage(conn, m, timeout)
}

// Event handler that is called after a Connection sends a complete message
func convertToMessage(id int, msg []byte, debugPrint bool) (Message, error) {
	msgID := [4]byte{}
	if len(msg) < len(msgID) {
		return nil, errors.New("Not enough data to read msg id")
	}
	copy(msgID[:], msg[:len(msgID)])
	msg = msg[len(msgID):]
	t, succ := MessageIDReverseMap[msgID]
	if !succ {
		return nil, fmt.Errorf("Unknown message %s received", string(msgID[:]))
	}

	if debugPrint {
		logger.Debug("Convert, Message type %v", t)
	}

	var m Message
	v := reflect.New(t)
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
		logger.Panic("Message obtained from map does not match Message interface")
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
	msgID, succ := MessageIDMap[t]
	if !succ {
		txt := "Attempted to serialize message struct not in MessageIdMap: %v"
		logger.Panicf(txt, msg)
	}
	bMsg := encoder.Serialize(msg)

	// message length
	bLen := encoder.SerializeAtomic(uint32(len(bMsg) + len(msgID)))
	m := make([]byte, 0)
	m = append(m, bLen...)     // length prefix
	m = append(m, msgID[:]...) // message id
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
