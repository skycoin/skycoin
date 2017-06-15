package app

import (
	"io"
	"net"
	"sync"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

type Server struct {
	app
}

func NewServer(appId messages.AppId, nodeAddr string, handle func([]byte) []byte) (*Server, error) {

	server := &Server{}
	server.id = appId
	server.lock = &sync.Mutex{}
	server.timeout = APP_TIMEOUT
	server.handle = handle
	server.responseNodeAppChannels = make(map[uint32]chan bool)

	err := server.RegisterAtNode(nodeAddr)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (self *Server) consume(appMsg *messages.AppMessage) {

	sequence := appMsg.Sequence
	go func() {
		responsePayload := self.handle(appMsg.Payload)
		response := &messages.AppMessage{
			sequence,
			responsePayload,
		}
		responseSerialized := messages.Serialize(messages.MsgAppMessage, response)
		self.sendToMeshnet(responseSerialized)
	}()
}

func (self *Server) RegisterAtNode(nodeAddr string) error {

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

func (self *Server) listenFromNode() {
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

func (self *Server) handleIncomingFromNode(msg []byte) error {
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
