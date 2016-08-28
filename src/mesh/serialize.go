package mesh

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

const messagePrefixLength = 1

// Message prefix identifies a message
type MessagePrefix [messagePrefixLength]byte

type Serializer struct {
	messageIdMap        map[reflect.Type]MessagePrefix
	messageIdReverseMap map[MessagePrefix]reflect.Type
}

func NewSerializer() *Serializer {
	ret := &Serializer{}
	ret.messageIdMap = make(map[reflect.Type]MessagePrefix)
	ret.messageIdReverseMap = make(map[MessagePrefix]reflect.Type)
	return ret
}

// Register a message struct for recognition by the message handlers.
func (self *Serializer) RegisterMessageForSerialization(prefix MessagePrefix, msg interface{}) {
	t := reflect.TypeOf(msg)
	id := MessagePrefix{}
	copy(id[:], prefix[:])
	_, exists := self.messageIdReverseMap[id]
	if exists {
		log.Panicf("Attempted to register message prefix %s twice",
			string(id[:]))
	}
	_, exists = self.messageIdMap[t]
	if exists {
		log.Panicf("Attempts to register message type %v twice", t)
	}
	self.messageIdMap[t] = id
	self.messageIdReverseMap[id] = t
}

func (self *Serializer) UnserializeMessage(msg []byte) (interface{}, error) {
	msgId := [1]byte{}
	if len(msg) < len(msgId) {
		return nil, errors.New("Not enough data to read msg id")
	}
	copy(msgId[:], msg[:len(msgId)])
	msg = msg[len(msgId):]

	t, succ := self.messageIdReverseMap[msgId]
	if !succ {
		logger.Debug("Unknown message id %s msgId %v", string(msgId[:]))
		return nil, fmt.Errorf("Unknown message %s received", string(msgId[:]))
	}

	var m interface{}
	var v reflect.Value = reflect.New(t)
	//logger.Debug("Giving %d bytes to the decoder", len(msg))
	used, err := self.deserializeMessage(msg, v)
	if err != nil {
		return nil, err
	}
	if used != len(msg) {
		return nil, errors.New("Data buffer was not completely decoded")
	}
	m, succ = (v.Elem().Interface()).(interface{})
	if !succ {
		// This occurs only when the user registers an interface that does
		// match the interface{} interface.  They should have known about this
		// earlier via a call to VerifyMessages
		log.Panic("interface{} obtained from map does not match interface{} interface")
		return nil, errors.New("MessageIdMaps contain non-interface{}")
	}
	return m, nil
}

// Wraps encoder.DeserializeRawToValue and traps panics as an error
func (self *Serializer) deserializeMessage(msg []byte, v reflect.Value) (n int, e error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Debug("Recovering from deserializer panic: %v", r)
			switch x := r.(type) {
			case string:
				e = errors.New(x)
			case error:
				e = x
			default:
				e = errors.New("interface{} deserialization failed")
			}
		}
	}()
	n, e = encoder.DeserializeRawToValue(msg, v)
	return
}

// Packgs a interface{} into []byte containing length, id and data
func (self *Serializer) SerializeMessage(msg interface{}) []byte {
	t := reflect.TypeOf(msg)
	msgId, succ := self.messageIdMap[t]
	if !succ {
		txt := "Attempted to serialize message type not in MessageIdMap: %v"
		log.Panicf(txt, t)
	}
	bMsg := encoder.Serialize(msg)

	// message length
	m := make([]byte, 0)
	m = append(m, msgId[:]...) // message id
	m = append(m, bMsg...)     // message bytes
	return m
}
