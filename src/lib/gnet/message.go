package gnet

import (
	"log"
	"reflect"
)

const messagePrefixLength = 4

// Message prefix identifies a message
type MessagePrefix [messagePrefixLength]byte

func MessagePrefixFromString(prefix string) MessagePrefix {
	if len(prefix) == 0 || len(prefix) > 4 {
		log.Panicf("Invalid prefix %s", prefix)
	}
	p := MessagePrefix{}
	for i, c := range prefix {
		p[i] = byte(c)
	}
	for i := len(prefix); i < 4; i++ {
		p[i] = 0x00
	}
	return p
}

type Message interface {
	// State is user-defined application state that is attached to the
	// Dispatcher.
	// Return a non-nil error from handle only if you've disconnected the
	// client.  You don't have to return the DisconnectReason but that may
	// be the most convenient.  If error is not nil, event buffer processing
	// is aborted.
	Handle(context *MessageContext, state interface{}) error
}

type MessageContext struct {
	Conn *Connection // connection message was received from
}

func NewMessageContext(conn *Connection) *MessageContext {
	return &MessageContext{Conn: conn}
}

// Maps message types to their ids
//var MessageIdMap = make(map[reflect.Type]MessagePrefix)

// Maps message ids to their types
//var MessageIdReverseMap = make(map[MessagePrefix]reflect.Type)

// Register a message struct for recognition by the message handlers.
// A maximum of 255 messages are allowed
func (self *Dispatcher) RegisterMessage(prefix MessagePrefix, msg interface{}) {
	t := reflect.TypeOf(msg)
	id := MessagePrefix{}
	copy(id[:], prefix[:])
	_, exists := self.MessageIdReverseMap[id]
	if exists {
		log.Panicf("Attempted to register message prefix %s twice",
			string(id[:]))
	}
	_, exists = self.MessageIdMap[t]
	if exists {
		log.Panicf("Attempts to register message type %v twice", t)
	}
	self.MessageIdMap[t] = id
	self.MessageIdReverseMap[id] = t
}

//Takes in map with string to message and registers them
func (self *Dispatcher) RegisterMessages(msgMap map[string](interface{})) {
	for k, msg := range msgMap {
		prefix := MessagePrefixFromString(k)
		self.RegisterMessage(prefix, msg)
	}
	self.VerifyMessages()

	//log.Printf("message_map: %v \n", self.MessageIdReverseMap)
}

// Calls log.Panic if message registration violates sanity checks
func (self *Dispatcher) VerifyMessages() {
	for t, k := range self.MessageIdMap {
		// No empty prefixes allowed
		if k[0] == 0x00 {
			log.Panic("No empty message prefixes allowed")
		}
		// No non-null bytes allowed after a nul byte
		hasEmpty := false
		for _, b := range k {
			if b == 0x00 {
				hasEmpty = true
			} else if hasEmpty {
				log.Panic("No non-null bytes allowed after a nul byte")
			}
		}
		// All characters must be non-whitespace printable ascii chars/digits
		// No punctation
		for _, b := range k {
			if !((b >= '0' && b <= '9') || (b >= 'A' && b <= 'Z') ||
				(b >= 'a' && b <= 'z') || b == 0x00) {
				log.Panicf("Invalid prefix byte %v", b)
			}
		}

		// Confirm that all registered messages support the Message interface
		// This should only be untrue if the user modified the message map
		// directly
		mptr := reflect.PtrTo(t)
		if !mptr.Implements(reflect.TypeOf((*Message)(nil)).Elem()) {
			m := "Message must implement the gnet.Message interface"
			log.Panicf("Invalid message at id %d: %s", k, m)
		}
	}
	if len(self.MessageIdMap) != len(self.MessageIdReverseMap) {
		log.Panic("MessageIdMap mismatch")
	}
	// No empty prefixes
	// All prefixes must be 0 padded
}

// Wipes all recorded messages
func (self *Dispatcher) EraseMessages() {
	self.MessageIdMap = make(map[reflect.Type]MessagePrefix)
	self.MessageIdReverseMap = make(map[MessagePrefix]reflect.Type)
}
