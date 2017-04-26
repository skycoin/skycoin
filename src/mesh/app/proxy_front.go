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
	self.send(requestSerialized)
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

func (self *proxyClient) RegisterAtNode(node messages.NodeInterface) error {
	err := node.RegisterApp(self)
	if err != nil {
		return err
	}
	self.node = node
	return nil
}

func (self *proxyClient) Consume(appMsg *messages.AppMessage) {

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

	//	log.Printf("\nClient accepted %d bytes to %s\n\n", len(data), remoteAddr)

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
