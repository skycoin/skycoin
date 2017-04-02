package app

import (
	"fmt"
	//	"io"
	"net"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type SocksClient struct {
	app
	socksConn net.Conn
}

func NewSocksClient(meshnet messages.Network, address cipher.PubKey, proxyAddress string) (*SocksClient, error) {
	socksClient := &SocksClient{}
	socksClient.register(meshnet, address)
	socksClient.lock = &sync.Mutex{}
	socksClient.timeout = time.Duration(messages.GetConfig().AppTimeout)

	conn, err := meshnet.NewConnection(address)
	if err != nil {
		return nil, err
	}

	socksClient.connection = conn

	err = meshnet.Register(address, socksClient)
	if err != nil {
		return nil, err
	}
	socksClient.socksAddress = proxyAddress
	go socksClient.socksClient(proxyAddress)

	return socksClient, nil
}

func (self *SocksClient) Send(msg []byte) {

	request := &messages.AppMessage{
		0,
		true,
		msg,
	}
	requestSerialized := messages.Serialize(messages.MsgAppMessage, request)
	self.send(requestSerialized)
	//	fmt.Println("sent", requestSerialized, string(msg))
}

func (self *SocksClient) socksClient(proxyAddress string) {
	l, err := net.Listen("tcp", proxyAddress)
	if err != nil {
		panic(err)
	}

	conn, err := l.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	self.socksConn = conn

	for {
		message := make([]byte, config.SocksPacketSize)
		fmt.Println("ready for the next packet from USER")
		n, err := conn.Read(message)
		if err != nil {
			//			if err != io.EOF {
			panic(err)
			//			}
		}
		fmt.Println("got from USER:", n)

		self.Send(message[:n])
	}
}

func (self *SocksClient) Consume(msg []byte) {
	appMsg := messages.AppMessage{}
	err := messages.Deserialize(msg, &appMsg)
	if err != nil {
		fmt.Printf("Cannot consume a message: %s\n", err.Error())
		return
	}
	conn := self.socksConn
	response := appMsg.Payload
	n, err := conn.Write(response)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nresponse %d is written to USER (local %s, remote %s), n = %d\n\n", response, conn.LocalAddr(), conn.RemoteAddr(), n)
}
