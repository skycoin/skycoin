package app

import (
	"errors"
	"io"
	"log"
	"net"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

type proxyServer struct {
	app
	targetConns map[string]net.Conn
}

const PROXY_PACKET_SIZE uint32 = 16384

func (self *proxyServer) Send(data []byte) {
	message := &messages.AppMessage{
		0,
		data,
	}
	messageS := messages.Serialize(messages.MsgAppMessage, message)
	self.sendToMeshnet(messageS)
}

func (self *proxyServer) Shutdown() {
	for _, c := range self.targetConns {
		c.Close()
	}
	self.app.Shutdown()
}

func getProxyMessage(appMsg *messages.AppMessage) *messages.ProxyMessage {

	proxyMessageS := appMsg.Payload
	proxyMessage := messages.ProxyMessage{}

	err := messages.Deserialize(proxyMessageS, &proxyMessage)
	if err != nil {
		log.Printf("Cannot deserialize proxy message: %s\n", err.Error())
		return nil
	}

	return &proxyMessage
}

func (self *proxyServer) writeToTarget(request []byte, remoteAddr string) error {
	self.lock.Lock()
	targetConn, ok := self.targetConns[remoteAddr]
	self.lock.Unlock()

	if !ok {
		return errors.New("Target connection not found for address:" + remoteAddr)
	}
	_, err := targetConn.Write(request)
	if err != nil {
		self.sendClose(remoteAddr)
		self.closeConns(remoteAddr)
	}
	return err
}

func (self *proxyServer) sendClose(remoteAddr string) {
	closingMessage := messages.ProxyMessage{
		nil,
		remoteAddr,
		true,
	}
	closingMessageS := messages.Serialize(messages.MsgProxyMessage, closingMessage)
	self.Send(closingMessageS)
}

func (self *proxyServer) getFromConn(conn net.Conn, remoteAddr string) { // permanently read data from proxy server and send it through the meshnet until an error

	for {
		packet, err := getPacketFromConn(conn)
		if err != nil {
			self.closeConns(remoteAddr)
			return
		}

		proxyMessage := messages.ProxyMessage{
			packet,
			remoteAddr,
			false,
		}
		proxyMessageS := messages.Serialize(messages.MsgProxyMessage, proxyMessage)
		self.Send(proxyMessageS)
	}
}

func getPacketFromConn(conn io.Reader) ([]byte, error) {
	buffer := make([]byte, PROXY_PACKET_SIZE)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:n], nil
}

func (self *proxyServer) closeConns(remoteAddr string) {
	self.lock.Lock()
	defer self.lock.Unlock()
	targetConn, ok := self.targetConns[remoteAddr]
	if ok {
		targetConn.Close()
	}
	delete(self.targetConns, remoteAddr)
}
