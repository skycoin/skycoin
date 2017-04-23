package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type app struct {
	//	Address          cipher.PubKey
	ProxyAddress     string
	handle           func([]byte) []byte
	timeout          time.Duration
	sequence         uint32
	connection       messages.Connection
	responseChannels map[uint32]chan messages.AppResponse
	lock             *sync.Mutex
}

var APP_TIMEOUT = 100000 * time.Duration(time.Millisecond)

func (self *app) Dial(address cipher.PubKey) error {
	return self.connection.Dial(address)
}

func (self *app) Consume(msg []byte) {
	appMsg := messages.AppMessage{}
	err := messages.Deserialize(msg, &appMsg)
	if err != nil {
		fmt.Printf("Cannot consume a message: %s\n", err.Error())
		return
	}

	sequence := appMsg.Sequence
	if appMsg.ResponseRequired {
		go func() {
			responsePayload := self.handle(appMsg.Payload)
			response := &messages.AppMessage{
				sequence,
				false,
				responsePayload,
			}
			responseSerialized := messages.Serialize(messages.MsgAppMessage, response)
			self.send(responseSerialized)
		}()
	} else {
		responseChannel, err := self.getResponseChannel(sequence)
		if err != nil {
			fmt.Println("error:", err)
			responseChannel <- messages.AppResponse{nil, err}
		} else {
			responseChannel <- messages.AppResponse{appMsg.Payload, nil}
		}
	}
}

func (self *app) Shutdown() {
	if self.connection != nil {
		self.connection.Shutdown()
	}
}

func (self *app) send(msg []byte) {

	conn := self.connection
	conn.Send(msg)
}

func (self *app) getResponseChannel(sequence uint32) (chan messages.AppResponse, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	ch, ok := self.responseChannels[sequence]
	if !ok {
		return nil, messages.ERR_NO_CLIENT_RESPONSE_CHANNEL
	}
	return ch, nil
}

func (self *app) setResponseChannel(responseChannel chan messages.AppResponse) uint32 {
	self.lock.Lock()
	defer self.lock.Unlock()

	sequence := self.sequence
	self.sequence++
	self.responseChannels[sequence] = responseChannel
	return sequence
}
