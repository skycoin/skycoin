package app

import (
	"fmt"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh2/connection"
	"github.com/skycoin/skycoin/src/mesh2/errors"
	"github.com/skycoin/skycoin/src/mesh2/messages"
	"github.com/skycoin/skycoin/src/mesh2/nodemanager"
)

type Client struct {
	app
	Timeout          int
	connection       *connection.Connection
	responseChannels map[uint32]chan []byte
}

const DEFAULT_TIMEOUT int = 10000

func NewClient(meshnet *nodemanager.NodeManager, address cipher.PubKey) (*Client, error) {
	client := &Client{}
	client.Timeout = DEFAULT_TIMEOUT
	client.responseChannels = make(map[uint32]chan []byte)
	client.Register(meshnet, address)
	err := meshnet.AssignConsumer(address, client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (self *Client) Dial(address cipher.PubKey) error {
	conn, err := connection.NewConnection(self.Meshnet, self.Address, address)
	if err != nil {
		return err
	}
	self.connection = conn
	conn.Tick()
	return nil
}

func (self *Client) DialWithRoutes(route, backRoute messages.RouteId) error {
	conn, err := connection.NewConnectionWithRoutes(self.Meshnet, self.Address, route, backRoute)
	if err != nil {
		return err
	}
	self.connection = conn
	conn.Tick()
	return nil
}

func (self *Client) Consume(sequence uint32, response []byte, _ chan<- []byte) {
	responseChannel, ok := self.responseChannels[sequence]
	if !ok {
		return
	}
	responseChannel <- response
}

func (self *Client) Send(msg []byte) ([]byte, error) {
	sequence, err := self.connection.Send(msg)
	if err != nil {
		return nil, err
	}
	responseChannel := make(chan []byte, 1024)
	self.responseChannels[sequence] = responseChannel
	select {
	case response := <-responseChannel:
		fmt.Println("RESPONSE HAS COME:", response)
		return response, nil
	case <-time.After(time.Duration(self.Timeout) * time.Millisecond):
		fmt.Println("TIMEOUT")
		self.connection.Close()
		return nil, errors.ERR_TIMEOUT
	}
}

/*
Client manager Connection - creates it by Dial() and closes it if Send gets a timeout
Consume accepts message from node and retranslates it to Send
Send sends message and waits for response, retransmits(?), closes connection if timeout
*/
