package app

import (
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type Client struct {
	app
	ResponsesReceived uint32
	ResponsesErrors   uint32
	ResponsesSent     uint32
	Consumed          uint32
	ConsumedSent      uint32
}

func NewClient(meshnet messages.Network, address cipher.PubKey) (*Client, error) {
	client := &Client{}
	client.register(meshnet, address)
	client.lock = &sync.Mutex{}
	client.timeout = time.Duration(messages.GetConfig().AppTimeout)
	client.responseChannels = make(map[uint32]chan messages.AppResponse)

	conn, err := meshnet.NewConnection(address)
	if err != nil {
		return nil, err
	}

	client.connection = conn

	err = meshnet.Register(address, client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (self *Client) Dial(address cipher.PubKey) error {
	err := self.Network.Connect(self.Address, address)
	return err
}

/*
func (self *Client) DialWithRoutes(route, backRoute messages.RouteId) error {
	conn, err := self.Network.NewConnectionWithRoutes(self.Address, route, backRoute)
	if err != nil {
		return err
	}
	self.connection = conn
	return nil
}
*/

func (self *Client) Send(msg []byte) ([]byte, error) {

	responseChannel := make(chan messages.AppResponse)
	sequence := self.setResponseChannel(responseChannel)

	request := &messages.AppMessage{
		sequence,
		false,
		msg,
	}
	requestSerialized := messages.Serialize(messages.MsgAppMessage, request)
	go self.send(requestSerialized)

	select {
	case appResponse := <-responseChannel:
		return appResponse.Response, appResponse.Err
	case <-time.After(self.timeout * time.Millisecond):
		return nil, messages.ERR_APP_TIMEOUT
	}
}

/*
func (self *Client) GetConnection() messages.Connection {
	return self.connection
}
*/
