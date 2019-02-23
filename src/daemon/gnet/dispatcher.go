package gnet

import (
	"errors"
	"fmt"
	"math"
	"net"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/util/mathutil"
)

var (
	// ErrMsgExceedsMaxLen is returned if trying to send a message that exceeds the configured max length
	ErrMsgExceedsMaxLen = errors.New("Message exceeds max message length")
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
func sendMessage(conn net.Conn, msg Message, timeout time.Duration, maxMsgLength int) error {
	m, err := EncodeMessage(msg)
	if err != nil {
		return err
	}
	if len(m) > maxMsgLength {
		return ErrMsgExceedsMaxLen
	}
	return sendByteMessage(conn, m, timeout)
}

// msgIDStringSafe formats msgID bytes to a string that is safe for logging (e.g. not impacted by ascii control chars)
func msgIDStringSafe(msgID [4]byte) string {
	x := fmt.Sprintf("%q", msgID)
	return x[1 : len(x)-1] // trim quotes that are added by %q formatting
}

// Event handler that is called after a Connection sends a complete message
func convertToMessage(id uint64, msg []byte, debugPrint bool) (Message, error) {
	msgID := [4]byte{}
	if len(msg) < len(msgID) {
		logger.WithError(ErrDisconnectTruncatedMessageID).WithField("connID", id).Warning()
		return nil, ErrDisconnectTruncatedMessageID
	}

	copy(msgID[:], msg[:len(msgID)])

	if debugPrint {
		logger.WithField("msgID", msgIDStringSafe(msgID)).Debug("Received message")
	}

	msg = msg[len(msgID):]
	t, ok := MessageIDReverseMap[msgID]
	if !ok {
		logger.WithError(ErrDisconnectUnknownMessage).WithFields(logrus.Fields{
			"msgID":  msgIDStringSafe(msgID),
			"connID": id,
		}).Warning()
		return nil, ErrDisconnectUnknownMessage
	}

	if debugPrint {
		logger.WithFields(logrus.Fields{
			"connID":      id,
			"messageType": fmt.Sprintf("%v", t),
		}).Debugf("convertToMessage")
	}

	v := reflect.New(t)
	m, ok := (v.Interface()).(Message)
	if !ok {
		// This occurs only when the user registers an interface that does not
		// match the Message interface.  They should have known about this
		// earlier via a call to VerifyMessages
		logger.Panic("Message obtained from map does not match Message interface")
		return nil, errors.New("MessageIdMaps contain non-Message")
	}

	used, err := deserializeMessage(msg, v)
	if err != nil {
		logger.Critical().WithError(err).WithFields(logrus.Fields{
			"connID":      id,
			"messageType": fmt.Sprintf("%v", t),
		}).Warning("deserializeMessage failed")
		return nil, ErrDisconnectMalformedMessage
	}

	if used != uint64(len(msg)) {
		logger.WithError(ErrDisconnectMessageDecodeUnderflow).WithFields(logrus.Fields{
			"connID":      id,
			"messageType": fmt.Sprintf("%v", t),
		}).Warning()
		return nil, ErrDisconnectMessageDecodeUnderflow
	}

	return m, nil
}

// Wraps Serializer.Decode and traps panics as an error
func deserializeMessage(msg []byte, v reflect.Value) (n uint64, err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Critical().Warningf("Recovering from deserializer panic: %v", r)
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Message deserialization failed")
			}
		}
	}()

	iface := v.Interface()
	x, ok := iface.(Serializer)
	if !ok {
		return 0, errors.New("deserializeMessage object does not have Serializer interface")
	}

	return x.Decode(msg)
}

// EncodeMessage packs a Message into []byte containing length, id and data
func EncodeMessage(msg Serializer) ([]byte, error) {
	t := reflect.ValueOf(msg).Elem().Type()

	// Lookup message ID
	msgID, succ := MessageIDMap[t]
	if !succ {
		logger.Panicf("Attempted to serialize message struct not in MessageIDMap: %v", msg)
	}
	if uint64(len(msgID)) > math.MaxUint32 {
		return nil, errors.New("Message ID length exceeds math.MaxUint32")
	}

	// Compute size of encoded Message object
	bMsgLen := msg.EncodeSize()
	if bMsgLen > math.MaxUint32 {
		return nil, errors.New("Message length exceeds math.MaxUint32")
	}

	// Compute message + message ID length
	bLen, err := mathutil.AddUint32(uint32(bMsgLen), uint32(len(msgID)))
	if err != nil {
		return nil, err
	}

	// Serialize total message length
	bLenPrefix := encoder.SerializeUint32(bLen)
	if uint64(len(bLenPrefix)) > math.MaxUint32 {
		return nil, errors.New("Message length prefix length exceeds math.MaxUint32")
	}

	mLen, err := mathutil.AddUint32(bLen, uint32(len(bLenPrefix)))
	if err != nil {
		return nil, err
	}

	// Allocate message bytes
	m := make([]byte, mLen)

	// Write the total message length to the buffer
	copy(m[:], bLenPrefix[:])

	// Write the message ID to the buffer
	copy(m[len(bLenPrefix):], msgID[:])

	// Encode the message into the message buffer
	if err := msg.Encode(m[len(bLenPrefix)+len(msgID):]); err != nil {
		return nil, err
	}

	return m, nil
}

// Sends []byte over a net.Conn
var sendByteMessage = func(conn net.Conn, msg []byte, timeout time.Duration) error {
	deadline := time.Time{}
	if timeout != 0 {
		deadline = time.Now().Add(timeout)
	}
	if err := conn.SetWriteDeadline(deadline); err != nil {
		return err
	}
	if _, err := conn.Write(msg); err != nil {
		return &WriteError{
			Err: err,
		}
	}
	return nil
}
