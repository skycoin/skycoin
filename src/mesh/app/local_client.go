package app

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

type Client struct {
	app
	sequence         uint32
	responseChannels map[uint32]chan messages.AppResponse
}

func NewClient(appId messages.AppId, nodeAddr string) (*Client, error) {

	client := &Client{}
	client.id = appId
	client.lock = &sync.Mutex{}
	client.timeout = APP_TIMEOUT
	client.responseChannels = make(map[uint32]chan messages.AppResponse)
	client.responseNodeAppChannels = make(map[uint32]chan bool)

	err := client.RegisterAtNode(nodeAddr)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (self *Client) Send(msg []byte) ([]byte, error) {

	responseChannel := make(chan messages.AppResponse)
	sequence := self.setResponseChannel(responseChannel)

	request := &messages.AppMessage{
		sequence,
		msg,
	}
	requestSerialized := messages.Serialize(messages.MsgAppMessage, request)
	go self.sendToMeshnet(requestSerialized)

	select {
	case appResponse := <-responseChannel:
		return appResponse.Response, appResponse.Err
	case <-time.After(self.timeout * time.Millisecond):
		return nil, messages.ERR_APP_TIMEOUT
	}
}

func (self *Client) getResponseChannel(sequence uint32) (chan messages.AppResponse, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	ch, ok := self.responseChannels[sequence]
	if !ok {
		return nil, messages.ERR_NO_CLIENT_RESPONSE_CHANNEL
	}
	return ch, nil
}

func (self *Client) setResponseChannel(responseChannel chan messages.AppResponse) uint32 {
	self.lock.Lock()
	defer self.lock.Unlock()

	sequence := self.sequence
	self.sequence++
	self.responseChannels[sequence] = responseChannel
	return sequence
}

func (self *Client) consume(appMsg *messages.AppMessage) {

	sequence := appMsg.Sequence
	responseChannel, err := self.getResponseChannel(sequence)
	if err != nil {
		responseChannel <- messages.AppResponse{nil, err}
	} else {
		responseChannel <- messages.AppResponse{appMsg.Payload, nil}
	}
}

func (self *Client) RegisterAtNode(nodeAddr string) error {

	nodeConn, err := net.Dial("tcp", nodeAddr)
	if err != nil {
		panic(nodeAddr)
		return err
	}

	self.nodeConn = nodeConn

	go self.listenFromNode()

	registerMessage := messages.RegisterAppMessage{}

	rmS := messages.Serialize(messages.MsgRegisterAppMessage, registerMessage)

	err = self.sendToNode(rmS)
	return err
}

func (self *Client) listenFromNode() {
	conn := self.nodeConn
	for {
		message, err := getFullMessage(conn)
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				break
			}
		} else {
			go self.handleIncomingFromNode(message)
		}
	}
}

func (self *Client) handleIncomingFromNode(msg []byte) error {
	switch messages.GetMessageType(msg) {

	case messages.MsgAssignConnectionNAM:
		m1 := &messages.AssignConnectionNAM{}
		err := messages.Deserialize(msg, m1)
		if err != nil {
			return err
		}
		self.meshConnId = m1.ConnectionId
		return nil

	case messages.MsgAppMessage:
		appMsg := &messages.AppMessage{}
		err := messages.Deserialize(msg, appMsg)
		if err != nil {
			return err
		}
		go self.consume(appMsg)
		return nil

	case messages.MsgNodeAppResponse:
		nar := &messages.NodeAppResponse{}
		err := messages.Deserialize(msg, nar)
		if err != nil {
			return err
		}

		sequence := nar.Sequence
		respChan, err := self.getResponseNodeAppChannel(sequence)
		if err != nil {
			panic(err)
			return err
		} else {
			respChan <- true
			return nil
		}

	default:
		return messages.ERR_INCORRECT_MESSAGE_TYPE
	}
}
