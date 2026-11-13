// Package gnet src/daemon/gnet/message.go c1-coin-daemon
package gnet

import (
	"errors"
	"fmt"
	"reflect"
)

const messagePrefixLength = 4

// MessagePrefix message prefix identifies a message
type MessagePrefix [messagePrefixLength]byte

// MessagePrefixFromString creates MessagePrefix from string
func MessagePrefixFromString(prefix string) MessagePrefix {
	if len(prefix) == 0 || len(prefix) > 4 {
		logger.Panicf("Invalid prefix %s", prefix)
	}
	p := MessagePrefix{}
	for i, c := range prefix {
		p[i] = byte(c) //nolint:gosec // deliberate truncation: packing a rune into a protocol byte
	}
	for i := len(prefix); i < 4; i++ {
		p[i] = 0x00
	}
	return p
}

// Serializer serialization interface
type Serializer interface {
	EncodeSize() uint64
	Encode([]byte) error
	Decode([]byte) (uint64, error)
}

// Handler message handler interface
type Handler interface {
	// State is user-defined application state that is attached to the Dispatcher.
	// If a non-nil error is returned, the connection will be disconnected.
	Handle(context *MessageContext, state interface{}) error
}

// Message message interface
type Message interface {
	Handler
	Serializer
}

// MessageContext message context
type MessageContext struct {
	ConnID uint64 // connection message was received from
	Addr   string
}

// NewMessageContext creates MessageContext
func NewMessageContext(conn *Connection) *MessageContext {
	if conn.Conn != nil {
		return &MessageContext{ConnID: conn.ID, Addr: conn.Addr()}
	}
	return &MessageContext{ConnID: conn.ID}
}

// MessageIDMap maps message types to their ids
var MessageIDMap = make(map[reflect.Type]MessagePrefix)

// MessageIDReverseMap maps message ids to their types
var MessageIDReverseMap = make(map[MessagePrefix]reflect.Type)

var registeredMsgsCount = 0

// RegisterMessage registers a message struct for recognition by the message handlers.
func RegisterMessage(prefix MessagePrefix, msg interface{}) error {
	t := reflect.TypeOf(msg)
	id := MessagePrefix{}
	copy(id[:], prefix[:])
	_, exists := MessageIDReverseMap[id]
	if exists {
		return fmt.Errorf("attempted to register message prefix %s twice", string(id[:]))
	}
	_, exists = MessageIDMap[t]
	if exists {
		return fmt.Errorf("attempted to register message type %v twice", t)
	}
	MessageIDMap[t] = id
	MessageIDReverseMap[id] = t

	registeredMsgsCount++
	return nil
}

// VerifyMessages returns an error if message registration violates sanity checks
func VerifyMessages() error {
	if registeredMsgsCount != len(MessageIDMap) {
		return errors.New("MessageIDMap was altered without using RegisterMessage")
	}
	if registeredMsgsCount != len(MessageIDReverseMap) {
		return errors.New("MessageIDReverseMap was altered without using RegisterMessage")
	}

	for t, k := range MessageIDMap {
		// No empty prefixes allowed
		if k[0] == 0x00 {
			return errors.New("no empty message prefixes allowed")
		}
		// No non-null bytes allowed after a nul byte
		emptyFound := false
		for _, b := range k {
			if b == 0x00 {
				emptyFound = true
			} else if emptyFound {
				return errors.New("no non-null bytes allowed after a nul byte")
			}
		}
		// All characters must be non-whitespace printable ascii chars/digits
		// No punctation
		for _, b := range k {
			if !((b >= '0' && b <= '9') || (b >= 'A' && b <= 'Z') ||
				(b >= 'a' && b <= 'z') || b == 0x00) {
				return fmt.Errorf("invalid prefix byte %v", b)
			}
		}

		// Confirm that all registered messages support the Message interface
		// This should only be untrue if the user modified the message map
		// directly
		mptr := reflect.PointerTo(t)
		if !mptr.Implements(reflect.TypeOf((*Message)(nil)).Elem()) {
			return fmt.Errorf("invalid message at ID %s: Message must implement the gnet.Message interface", string(k[:]))
		}
	}
	if len(MessageIDMap) != len(MessageIDReverseMap) {
		return errors.New("MessageIdMap mismatch")
	}
	return nil
}

// EraseMessages wipes all recorded message types
func EraseMessages() {
	MessageIDMap = make(map[reflect.Type]MessagePrefix)
	MessageIDReverseMap = make(map[MessagePrefix]reflect.Type)
	registeredMsgsCount = 0
}
