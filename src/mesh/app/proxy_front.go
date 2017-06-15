package app

import (
	"io"
	"log"
	"net"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

type proxyClient struct {
	app
	connections map[string]*net.Conn
}

func (self *proxyClient) Send(msg []byte) {

	request := &messages.AppMessage{
		0,
		msg,
	}
	requestSerialized := messages.Serialize(messages.MsgAppMessage, request)
	self.sendToMeshnet(requestSerialized)
}

func (self *proxyClient) Shutdown() {
	for _, c := range self.connections {
		(*c).Close()
	}
	self.app.Shutdown()
}

func (self *proxyClient) Listen() {
	proxyAddress := self.ProxyAddress

	l, err := net.Listen("tcp", proxyAddress)
	if err != nil {
		panic(err)
	}

	for {
		userConn, err := l.Accept() // create a connection with the user app (e.g. browser)
		if err != nil {
			log.Println("Cannot accept client's connection")
			return
		}
		defer userConn.Close()

		remoteAddr := userConn.RemoteAddr().String()

		self.lock.Lock()
		self.connections[remoteAddr] = &userConn
		self.lock.Unlock()

		go func() { // run listening the connection for data and sending it through the meshnet to the server
			for {
				message := make([]byte, PROXY_PACKET_SIZE)

				n, err := userConn.Read(message)
				if err != nil {
					return
					if err == io.EOF {
						continue
					} else {
						break
					}
				}

				proxyMessage := messages.ProxyMessage{
					message[:n],
					remoteAddr,
					false,
				}

				proxyMessageS := messages.Serialize(messages.MsgProxyMessage, proxyMessage)

				self.Send(proxyMessageS)
			}
		}()
	}
}

func (self *proxyClient) consume(appMsg *messages.AppMessage) {

	proxyMessageS := appMsg.Payload
	proxyMessage := messages.ProxyMessage{}
	err := messages.Deserialize(proxyMessageS, &proxyMessage)
	if err != nil {
		log.Printf("Cannot deserialize proxy message: %s\n", err.Error())
		return
	}

	remoteAddr := proxyMessage.RemoteAddr

	self.lock.Lock()
	connPointer, ok := self.connections[remoteAddr] // get the connection from existing connections by remote address
	self.lock.Unlock()

	if !ok {
		log.Printf("Cannot find the connection with remote address %s\n", remoteAddr)
		return
	}

	userConn := *connPointer

	if proxyMessage.NeedClose { // if we got a command to close a connection, close it
		log.Printf("Closing connection %s according to a signal from server\n", remoteAddr)
		userConn.Close()

		self.lock.Lock()
		delete(self.connections, remoteAddr)
		self.lock.Unlock()

		return
	}

	data := proxyMessage.Data // otherwise send data to the user app

	_, err = userConn.Write(data)
	if err != nil { // if the write is unsuccessful, close the connection and send closing command to close the corresponding connection on the server
		log.Printf("Cannot write to connection with remote address %s, error is %s\n", proxyMessage.RemoteAddr, err.Error())
		closingMessage := messages.ProxyMessage{
			nil,
			remoteAddr,
			true,
		}
		closingMessageS := messages.Serialize(messages.MsgProxyMessage, closingMessage)
		self.Send(closingMessageS)
		userConn.Close()
	}
}

func (self *proxyClient) RegisterAtNode(nodeAddr string) error {

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

func (self *proxyClient) listenFromNode() {
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

func (self *proxyClient) handleIncomingFromNode(msg []byte) error {
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
