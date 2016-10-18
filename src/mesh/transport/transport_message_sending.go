package transport

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/domain"
)

func (self *Transport) SendMessage(toPeer cipher.PubKey, contents []byte) error {
	messageID := self.newMessageID()
	sendMessage := SendMessage{messageID, self.config.MyPeerID, contents}
	sendSerialized := self.serializer.SerializeMessage(sendMessage)
	state := messageSentState{toPeer,
		sendSerialized,
		time.Now().Add(self.config.RetransmitDuration),
		false}
	err := self.physicalTransport.SendMessage(toPeer, sendSerialized)
	if err == nil {
		self.lock.Lock()
		defer self.lock.Unlock()
		self.messagesSent[messageID] = state
	}
	return err
}

func (self *Transport) SetCrypto(crypto ITransportCrypto) {
	self.physicalTransport.SetCrypto(crypto)
}

func (self *Transport) doRetransmits() {
	self.lock.Lock()
	defer self.lock.Unlock()
	for _, state := range self.messagesSent {
		if !state.receivedAck {
			go self.physicalTransport.SendMessage(state.toPeer, state.serialized)
		}
	}
}

func (self *Transport) retransmitLoop() {
	self.closeWait.Add(1)
	defer self.closeWait.Done()
	for len(self.closing) == 0 {
		select {
		case <-time.After(self.config.RetransmitDuration):
			{
				self.doRetransmits()
			}
		case <-self.closing:
			{
				return
			}
		}
	}
}

func (self *Transport) sendAck(message SendMessage) {
	reply := ReplyMessage{message.MessageID}
	serialized := self.serializer.SerializeMessage(reply)
	go self.physicalTransport.SendMessage(message.FromPeerID, serialized)
}

func (self *Transport) newMessageID() domain.MessageID {
	return (domain.MessageID)(uuid.NewV4())
}
