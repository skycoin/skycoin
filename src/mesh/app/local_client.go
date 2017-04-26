package app

import (
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/node"
)

type Client struct {
	app
}

func BrandNewClient(appId messages.AppId, host, meshnet string) (*Client, error) {

	client := newClient(appId)

	node, err := node.CreateAndConnectNode(host, meshnet)
	if err != nil {
		return nil, err
	}

	err = client.RegisterAtNode(node)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewClient(appId messages.AppId, node messages.NodeInterface) (*Client, error) {

	client := newClient(appId)

	err := client.RegisterAtNode(node)
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
		true,
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

func newClient(appId messages.AppId) *Client {
	client := &Client{}
	client.id = appId
	client.lock = &sync.Mutex{}
	client.timeout = APP_TIMEOUT
	client.responseChannels = make(map[uint32]chan messages.AppResponse)
	return client
}
