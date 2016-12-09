package transport

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/domain"
)

func (self *Transport) SendMessage(toPeer cipher.PubKey, contents []byte, _ chan error) error {
	self.status = SENDING
	messageID := self.newMessageID()
	sendMessage := SendMessage{messageID, self.config.MyPeerID, contents}
	sendSerialized := self.serializer.SerializeMessage(sendMessage)
	now := time.Now()
	state := messageSentState{
		toPeer,
		sendSerialized,
		now.Add(self.config.RetransmitDuration),
		false,
	}
	retChan := make(chan error, 0)
	var err error
	go self.physicalTransport.SendMessage(toPeer, sendSerialized, retChan)
	select {
	case err = <-retChan:
		self.status = CONNECTED
		if err == nil {
			self.lock.Lock()
			defer self.lock.Unlock()
			self.messagesSent[messageID] = state
			self.packetIsSent = now
			self.packetsSent++
			self.packetsCount++
		}
	case <-time.After(5 * time.Second):
		self.status = TIMEOUT
		err = errors.New("Sending is timed out")
		fmt.Fprintf(os.Stderr, "Timeout for sending message %s to %s\n", sendSerialized, toPeer)
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
			self.packetsRetransmissions++
			go self.physicalTransport.SendMessage(state.toPeer, state.serialized, nil)
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
	self.status = REPLYING
	reply := ReplyMessage{message.MessageID}
	serialized := self.serializer.SerializeMessage(reply)
	go self.physicalTransport.SendMessage(message.FromPeerID, serialized, nil)
	self.status = CONNECTED
}

func (self *Transport) newMessageID() domain.MessageID {
	return (domain.MessageID)(uuid.NewV4())
}
