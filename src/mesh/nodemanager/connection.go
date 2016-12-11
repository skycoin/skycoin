package nodemanager

import (
	"fmt"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/domain"
	"github.com/skycoin/skycoin/src/mesh/node"
	"github.com/skycoin/skycoin/src/mesh/serialize"
)

type IConnection interface {
	SendData([]byte) error
	SetOutputChannel(chan []byte)
}

type Connection struct {
	Node                   *mesh.Node
	ConnectedPeerID        cipher.PubKey
	Route                  *domain.Route
	Output                 chan []byte
	Serializer             *serialize.Serializer
	MessagesBeingAssembled map[domain.MessageID]*domain.MessageUnderAssembly
	TimeToAssembleMessage  time.Duration
	ExpireMessagesInterval time.Duration

	closing chan bool
}

func NewConnection(node *mesh.Node, peerID cipher.PubKey) IConnection {
	connection := &Connection{
		Node:                   node,
		ConnectedPeerID:        peerID,
		Serializer:             serialize.NewSerializer(),
		MessagesBeingAssembled: make(map[domain.MessageID]*domain.MessageUnderAssembly),
		TimeToAssembleMessage:  2 * time.Second,
		closing:                make(chan bool, 10),
	}

	connection.Serializer.RegisterMessageForSerialization(serialize.MessagePrefix{1}, domain.UserMessage{})
	connection.Serializer.RegisterMessageForSerialization(serialize.MessagePrefix{2}, domain.SetRouteMessage{})
	connection.Serializer.RegisterMessageForSerialization(serialize.MessagePrefix{3}, domain.RefreshRouteMessage{})
	connection.Serializer.RegisterMessageForSerialization(serialize.MessagePrefix{4}, domain.DeleteRouteMessage{})
	connection.Serializer.RegisterMessageForSerialization(serialize.MessagePrefix{5}, domain.SetRouteReply{})
	connection.Serializer.RegisterMessageForSerialization(serialize.MessagePrefix{6}, domain.AddNodeMessage{})

	go connection.expireOldMessagesLoop()

	return connection
}

func (self *Connection) SendData(data []byte) error {
	userMessages := self.FragmentMessage(data)
	for _, userMessage := range userMessages {
		serializedMessage := self.SerializeMessage(userMessage)
		err := self.Node.SendMessageThruRoute(self.Route.ForwardRewriteSendRouteID, serializedMessage)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Connection) FragmentMessage(contents []byte) []domain.UserMessage {
	messageBase := domain.MessageBase{
		SendRouteID: self.Route.ForwardRewriteSendRouteID,
		SendBack:    false,
		FromPeerID:  self.Node.Config.PubKey,
		Nonce:       mesh.GenerateNonce(),
	}

	fragmentsBuf := make([]domain.UserMessage, 0)
	maxContentLength := self.GetMaximumContentLength()
	fmt.Fprintf(os.Stdout, "MaxContentLength: %v\n", maxContentLength)
	remainingBytes := contents[:]
	messageID := (domain.MessageID)(uuid.NewV4())
	for len(remainingBytes) > 0 {
		nBytesThisMessage := min(maxContentLength, (uint64)(len(remainingBytes)))
		bytesThisMessage := remainingBytes[:nBytesThisMessage]
		remainingBytes = remainingBytes[nBytesThisMessage:]
		message := domain.UserMessage{
			MessageBase: messageBase,
			MessageID:   messageID,
			Index:       (uint64)(len(fragmentsBuf)),
			Count:       0,
			Contents:    bytesThisMessage,
		}
		fragmentsBuf = append(fragmentsBuf, message)
	}
	fragments := make([]domain.UserMessage, 0)
	for _, message := range fragmentsBuf {
		message.Count = (uint64)(len(fragmentsBuf))
		fragments = append(fragments, message)
	}
	fmt.Fprintf(os.Stdout, "Message fragmented in %v packets.\n", len(fragments))
	return fragments
}

// Returns nil if reassembly didn't happen (incomplete message)
func (self *Connection) ReassembleMessage(msgIn domain.UserMessage) []byte {

	_, assembledExists := self.MessagesBeingAssembled[msgIn.MessageID]
	if !assembledExists {
		beingAssembled := &domain.MessageUnderAssembly{
			Fragments:   make(map[uint64]domain.UserMessage),
			SendRouteID: msgIn.SendRouteID,
			SendBack:    msgIn.SendBack,
			Count:       msgIn.Count,
			Dropped:     false,
			ExpiryTime:  time.Now().Add(self.TimeToAssembleMessage),
		}
		self.MessagesBeingAssembled[msgIn.MessageID] = beingAssembled
	}

	beingAssembled, _ := self.MessagesBeingAssembled[msgIn.MessageID]

	if beingAssembled.Dropped {
		return nil
	}

	if beingAssembled.Count != msgIn.Count {
		fmt.Fprintf(os.Stderr, "Fragments of message %v have different total counts!\n", msgIn.MessageID)
		beingAssembled.Dropped = true
		return nil
	}

	if beingAssembled.SendRouteID != msgIn.SendRouteID {
		fmt.Fprintf(os.Stderr, "Fragments of message %v have different send ids!\n", msgIn.SendRouteID)
		beingAssembled.Dropped = true
		return nil
	}

	if beingAssembled.SendBack != msgIn.SendBack {
		fmt.Fprintf(os.Stderr, "Fragments of message %v have different send directions!\n", msgIn.SendRouteID)
		beingAssembled.Dropped = true
		return nil
	}

	_, messageExists := beingAssembled.Fragments[msgIn.Index]
	if messageExists {
		fmt.Fprintf(os.Stderr, "Fragment %v of message %v is duplicated, dropping message\n", msgIn.Index, msgIn.MessageID)
		return nil
	}

	beingAssembled.Fragments[msgIn.Index] = msgIn
	if (uint64)(len(beingAssembled.Fragments)) == beingAssembled.Count {
		delete(self.MessagesBeingAssembled, msgIn.MessageID)
		reassembled := []byte{}
		for i := (uint64)(0); i < beingAssembled.Count; i++ {
			reassembled = append(reassembled, beingAssembled.Fragments[i].Contents...)
		}
		return reassembled
	}

	return nil
}

func (self *Connection) SetOutputChannel(output chan []byte) {
	self.Output = output
}

func (self *Connection) DeserializeMessage(msg []byte) (interface{}, error) {
	return self.Serializer.UnserializeMessage(msg)
}

func (self *Connection) SerializeMessage(msg interface{}) []byte {
	return self.Serializer.SerializeMessage(msg)
}

func (self *Connection) expireOldMessages() {
	timeNow := time.Now()

	lastMessages := self.MessagesBeingAssembled
	self.MessagesBeingAssembled = make(map[domain.MessageID]*domain.MessageUnderAssembly)
	for messageID, message := range lastMessages {
		if timeNow.Before(message.ExpiryTime) {
			self.MessagesBeingAssembled[messageID] = message
		}
	}
}

func (self *Connection) expireOldMessagesLoop() {
	for len(self.closing) == 0 {
		select {
		case <-time.After(self.ExpireMessagesInterval):
			{
				self.expireOldMessages()
			}
		case <-self.closing:
			{
				return
			}
		}
	}
}

func (self *Connection) Close() error {
	for i := 0; i < 10; i++ {
		self.closing <- true
	}
	return nil
}

// WTF
func (self *Connection) GetMaximumContentLength() uint64 {
	empty := domain.UserMessage{}
	emptySerialized := self.Serializer.SerializeMessage(empty)
	return self.Node.GetMaximumContentLength(self.ConnectedPeerID, emptySerialized)
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
