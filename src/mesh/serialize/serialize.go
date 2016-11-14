package serialize

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

const messagePrefixLength = 1

// Message prefix identifies a message
type MessagePrefix [messagePrefixLength]byte

type Serializer struct {
	messageTypeToPrefix map[reflect.Type]MessagePrefix
	messagePrefixToType map[MessagePrefix]reflect.Type
}

func NewSerializer() *Serializer {
	serializer := &Serializer{}
	serializer.messageTypeToPrefix = make(map[reflect.Type]MessagePrefix)
	serializer.messagePrefixToType = make(map[MessagePrefix]reflect.Type)
	return serializer
}

// Register a message struct for recognition by the message handlers.
func (self *Serializer) RegisterMessageForSerialization(messagePrefix MessagePrefix, message interface{}) {
	messageType := reflect.TypeOf(message)
	messagePrefixBuf := MessagePrefix{}
	copy(messagePrefixBuf[:], messagePrefix[:])
	_, exists := self.messagePrefixToType[messagePrefixBuf]
	if exists {
		log.Panicf("Attempted to register message prefix %s twice",
			string(messagePrefixBuf[:]))
	}
	_, exists = self.messageTypeToPrefix[messageType]
	if exists {
		log.Panicf("Attempts to register message type %v twice", messageType)
	}
	self.messageTypeToPrefix[messageType] = messagePrefixBuf
	self.messagePrefixToType[messagePrefixBuf] = messageType
}

func (self *Serializer) UnserializeMessage(messageData []byte) (interface{}, error) {
	messagePrefix := [1]byte{}
	if len(messageData) < len(messagePrefix) {
		return nil, errors.New("Not enough data to read message prefix")
	}
	copy(messagePrefix[:], messageData[:len(messagePrefix)])
	messageData = messageData[len(messagePrefix):]

	messageType, exists := self.messagePrefixToType[messagePrefix]
	if !exists {
		return nil, fmt.Errorf("Unknown message %s received", string(messagePrefix[:]))
	}
	var reflectValue reflect.Value = reflect.New(messageType)
	//logger.Debug("Giving %d bytes to the decoder", len(msg))
	usedBytesCount, err := self.deserializeMessage(messageData, reflectValue)
	if err != nil {
		return nil, err
	}
	if usedBytesCount != len(messageData) {
		return nil, errors.New("Data buffer was not completely decoded")
	}
	message, exists := (reflectValue.Elem().Interface()).(interface{})
	if !exists {
		// This occurs only when the user registers an interface that does
		// match the interface{} interface.  They should have known about this
		// earlier via a call to VerifyMessages
		log.Panic("interface{} obtained from map does not match interface{} interface")
		return nil, errors.New("messagePrefixToType contain non-interface{}")
	}
	return message, nil
}

// Wraps encoder.DeserializeRawToValue and traps panics as an error
func (self *Serializer) deserializeMessage(messageData []byte, reflectValue reflect.Value) (usedBytesCount int, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("interface{} deserialization failed")
			}
		}
	}()
	usedBytesCount, err = encoder.DeserializeRawToValue(messageData, reflectValue)
	return
}

// Packs a interface{} into []byte containing length, id and data
func (self *Serializer) SerializeMessage(message interface{}) []byte {
	messageType := reflect.TypeOf(message)
	messagePrefix, exists := self.messageTypeToPrefix[messageType]
	if !exists {
		log.Panicf("Attempted to serialize message with unknown type: %v", messageType)
	}
	serializedMessageData := encoder.Serialize(message)

	// message length
	finalMessage := make([]byte, 0)
	finalMessage = append(finalMessage, messagePrefix[:]...)      // message prefix
	finalMessage = append(finalMessage, serializedMessageData...) // message bytes
	return finalMessage
}
