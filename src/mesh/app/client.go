package app

import (
	"sync"
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
	lock             *sync.Mutex
}

type ConnResponse struct {
	Response []byte
	Err      error
}

const DEFAULT_TIMEOUT int = 50000

var packetSize = messages.GetConfig().MaxPacketSize

func NewClient(meshnet messages.Network, address cipher.PubKey) (*Client, error) {
	client := &Client{}
	client.Timeout = DEFAULT_TIMEOUT
	client.responseChannels = make(map[uint32]chan []byte)
	client.lock = &sync.Mutex{}
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

/*
func (self *Client) Send(msg []byte) ([]byte, error) {
	responseChannel := make(chan []byte)

	conn := self.connection

	sequence, err := conn.Send(msg)
	if err != nil {
		conn.Close()
		return nil, err
	}

	self.responseChannels[sequence] = responseChannel

	select {
	case response := <-responseChannel:
		return response, nil
	case <-time.After(time.Duration(self.Timeout) * time.Millisecond):
		conn.Close()
		return nil, errors.ERR_TIMEOUT
	}
}
*/

func (self *Client) getResponseChannel(sequence uint32) (chan []byte, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	ch, ok := self.responseChannels[sequence]
	if !ok {
		return nil, errors.ERR_NO_CLIENT_RESPONSE_CHANNEL
	}
	return ch, nil
}

func (self *Client) setResponseChannel(sequence uint32, responseChannel chan []byte) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.responseChannels[sequence] = responseChannel
}

func (self *Client) Send(msg []byte) chan *ConnResponse {

	retChan := make(chan *ConnResponse, 1024)
	responseChannel := make(chan []byte, 1024)

	conn := self.connection

	sequence, err := conn.Send(msg)
	if err != nil {
		conn.Close()
		retChan <- &ConnResponse{nil, err}
		return retChan
	}

	self.setResponseChannel(sequence, responseChannel)

	select {
	case response := <-responseChannel:

		retChan <- &ConnResponse{response, nil}
		return retChan

	case <-time.After(time.Duration(self.Timeout) * time.Millisecond):
		conn.Close()
		retChan <- &ConnResponse{nil, errors.ERR_TIMEOUT}
		return retChan
	}
}

/*
Client manager Connection - creates it by Dial() and closes it if Send gets a timeout
Consume accepts message from node and retranslates it to Send
Send sends message and waits for response, retransmits(?), closes connection if timeout
*/

/*
Client manager Connection - creates it by Dial() and closes it if Send gets a timeout
Consume accepts message from node and retranslates it to Send
Send sends message and waits for response, retransmits(?), closes connection if timeout
*/
