package app

import (
	"fmt"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/errors"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type Client struct {
	app
	Timeout          int
	connection       messages.Connection
	responseChannels map[uint32]chan []byte
}

const DEFAULT_TIMEOUT int = 10000

func NewClient(meshnet messages.Network, address cipher.PubKey) (*Client, error) {
	client := &Client{}
	client.Timeout = DEFAULT_TIMEOUT
	client.responseChannels = make(map[uint32]chan []byte)
	client.register(meshnet, address)
	err := meshnet.Register(address, client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (self *Client) Dial(address cipher.PubKey) error {
	conn, err := self.Network.NewConnection(self.Address, address)
	if err != nil {
		return err
	}
	self.connection = conn
	return nil
}

func (self *Client) DialWithRoutes(route, backRoute messages.RouteId) error {
	conn, err := self.Network.NewConnectionWithRoutes(self.Address, route, backRoute)
	if err != nil {
		return err
	}
	self.connection = conn
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
	conn := self.connection
	sequence, err := conn.Send(msg)
	if err != nil {
		conn.Close()
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
		conn.Close()
		return nil, errors.ERR_TIMEOUT
	}
}

/*
Client manager Connection - creates it by Dial() and closes it if Send gets a timeout
Consume accepts message from node and retranslates it to Send
Send sends message and waits for response, retransmits(?), closes connection if timeout
*/
