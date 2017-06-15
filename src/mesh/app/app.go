package app

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

type app struct {
	ProxyAddress            string
	id                      messages.AppId
	nodeConn                net.Conn
	handle                  func([]byte) []byte
	timeout                 time.Duration
	meshConnId              messages.ConnectionId
	nodeAppSequence         uint32
	responseNodeAppChannels map[uint32]chan bool
	lock                    *sync.Mutex
	viscriptServer          *AppViscriptServer
}

const PACKET_SIZE = 1024

var APP_TIMEOUT = 10000 * time.Duration(time.Millisecond)

func (self *app) Id() messages.AppId {
	return self.id
}

func (self *app) Connect(appId messages.AppId, address string) error {

	msg := messages.ConnectToAppMessage{
		address,
		self.id,
		appId,
	}

	msgS := messages.Serialize(messages.MsgConnectToAppMessage, msg)

	err := self.sendToNode(msgS)

	return err
}

func (self *app) Shutdown() {
	self.nodeConn.Close()
}

func (self *app) consume(_ *messages.AppMessage) {
	panic("STUB-CONSUMING, no consume method in app implementation")
	//stub
}

func (self *app) sendToMeshnet(payload []byte) error {

	msg := messages.SendFromAppMessage{
		self.meshConnId,
		payload,
	}

	msgS := messages.Serialize(messages.MsgSendFromAppMessage, msg)

	err := self.sendToNode(msgS)
	return err
}

func (self *app) sendToNode(payload []byte) error {
	if self.nodeConn == nil {
		return nil // return error
	}

	respChan := make(chan bool)

	sequence := self.setResponseNodeAppChannel(respChan)
	nodeAppMessage := messages.NodeAppMessage{
		sequence,
		self.id,
		payload,
	}

	msgS := messages.Serialize(messages.MsgNodeAppMessage, nodeAppMessage)
	sizeMessage := messages.NumToBytes(len(msgS), 8)

	_, err := self.nodeConn.Write(sizeMessage)
	if err != nil {
		return err
	}

	_, err = self.nodeConn.Write(msgS)
	if err != nil {
		return err
	}

	select {
	case <-respChan:
		return nil
	case <-time.After(self.timeout * time.Millisecond):
		return messages.ERR_APP_TIMEOUT
	}
}

func (self *app) getResponseNodeAppChannel(sequence uint32) (chan bool, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	ch, ok := self.responseNodeAppChannels[sequence]
	if !ok {
		return nil, messages.ERR_NO_APP_RESPONSE_CHANNEL
	}
	return ch, nil
}

func (self *app) setResponseNodeAppChannel(responseChannel chan bool) uint32 {
	self.lock.Lock()
	defer self.lock.Unlock()

	sequence := self.nodeAppSequence
	self.nodeAppSequence++
	self.responseNodeAppChannels[sequence] = responseChannel
	return sequence
}

func (self *app) RegisterAtNode(nodeAddr string) error {

	nodeConn, err := net.Dial("tcp", nodeAddr)
	if err != nil {
		panic(err)
		return err
	}

	self.nodeConn = nodeConn

	go self.listenFromNode()

	registerMessage := messages.RegisterAppMessage{}

	rmS := messages.Serialize(messages.MsgRegisterAppMessage, registerMessage)

	err = self.sendToNode(rmS)
	return err
}

func (self *app) listenFromNode() {
	conn := self.nodeConn
	for {
		message := make([]byte, PACKET_SIZE)

		n, err := conn.Read(message)
		if err != nil {
			return
			if err == io.EOF {
				continue
			} else {
				break
			}
		}

		self.handleIncomingFromNode(message[:n])
	}
}

func (self *app) handleIncomingFromNode(msg []byte) error {
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
