package transport

import (
	"fmt"
	"os"
	"reflect"
	"time"
)

func (self *Transport) SetReceiveChannel(received chan []byte) {
	self.output = received
}

func (self *Transport) processReceivedLoop() {
	self.closeWait.Add(1)
	defer self.closeWait.Done()

	for len(self.closing) == 0 {
		select {
		case physicalMsg, ok := <-self.physicalReceived:
			{
				if ok {
					self.processPhysicalMessage(physicalMsg)
					self.packetsReceived++
					self.packetsCount++
				}
			}
		case <-self.closing:
			{
				return
			}
		}
	}
}

func (self *Transport) processSend(message SendMessage) {
	self.status = REPLYING
	self.sendAck(message)

	self.lock.Lock()
	defer self.lock.Unlock()
	_, alreadyReceived := self.messagesReceived[message.MessageID]
	if !alreadyReceived {
		self.output <- message.Contents
		self.messagesReceived[message.MessageID] = time.Now().Add(self.config.RememberMessageReceivedDuration)
		self.status = CONNECTED
	}
}

func (self *Transport) processReply(message ReplyMessage) {
	self.status = ACKWAITING
	self.lock.Lock()
	defer self.lock.Unlock()
	state, exists := self.messagesSent[message.MessageID]
	if !exists {
		fmt.Fprintf(os.Stderr, "Received ack for unknown sent message %v\n", message.MessageID)
		return
	}
	state.receivedAck = true
	self.messagesSent[message.MessageID] = state
	now := time.Now()
	self.latency = (uint64)(now.Unix() - self.packetIsSent.Unix())
	self.status = CONNECTED

	// Test
	if !self.messagesSent[message.MessageID].receivedAck {
		panic("Test error")
	}
}

func (self *Transport) processPhysicalMessage(physicalMessage []byte) {
	message, err := self.serializer.UnserializeMessage(physicalMessage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Deserialization error %v\n", err)
		return
	}
	messageType := reflect.TypeOf(message)

	if messageType == reflect.TypeOf(SendMessage{}) {
		send := message.(SendMessage)
		self.processSend(send)
	} else if messageType == reflect.TypeOf(ReplyMessage{}) {
		reply := message.(ReplyMessage)
		self.processReply(reply)
	} else {
		panic("Internal error: unknown message type")
	}
}
