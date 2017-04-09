package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type app struct {
	Address          cipher.PubKey
	Network          messages.Network
	ProxyAddress     string
	handle           func([]byte) []byte
	timeout          time.Duration
	sequence         uint32
	connection       messages.Connection
	responseChannels map[uint32]chan messages.AppResponse
	lock             *sync.Mutex
}

var config = messages.GetConfig()

func (app *app) register(meshnet messages.Network, address cipher.PubKey) {
	app.Network = meshnet
	app.Address = address
}

func (self *app) Dial(address cipher.PubKey) error {
	err := self.Network.Connect(self.Address, address)
	return err
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
			fmt.Println("response:", responseSerialized)
			self.send(responseSerialized)
		}()
	} else {
		responseChannel, err := self.getResponseChannel(sequence)
		if err != nil {
			fmt.Println("error:", err)
			responseChannel <- messages.AppResponse{nil, err}
		} else {
			fmt.Println("response:", appMsg.Payload)
			responseChannel <- messages.AppResponse{appMsg.Payload, nil}
		}
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
