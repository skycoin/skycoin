package app

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type SocksClient struct {
	app
	connections map[string]*net.Conn
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

	socksClient.connections = map[string]*net.Conn{}

	socksClient.socksAddress = proxyAddress

	return socksClient, nil
}

func (self *SocksClient) Send(msg []byte) {

	request := &messages.AppMessage{
		0,
		false,
		msg,
	}
	requestSerialized := messages.Serialize(messages.MsgAppMessage, request)
	self.send(requestSerialized)
}

func (self *SocksClient) Listen() {
	setLimit(16384)
	proxyAddress := self.socksAddress
	l, err := net.Listen("tcp", proxyAddress)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Cannot accept client's connection")
			return
		}
		defer conn.Close()

		remoteAddr := conn.RemoteAddr().String()
		self.connections[remoteAddr] = &conn

		go func() {
			for {
				message := make([]byte, config.SocksPacketSize)

				n, err := conn.Read(message)
				if err != nil {
					return
					if err == io.EOF {
						continue
					} else {
						break
					}
				}

				socksMessage := messages.SocksMessage{
					message[:n],
					remoteAddr,
					false,
				}

				socksMessageS := messages.Serialize(messages.MsgSocksMessage, socksMessage)

				self.Send(socksMessageS)
			}
		}()
	}
}

func (self *SocksClient) Consume(msg []byte) {
	appMsg := messages.AppMessage{}
	err := messages.Deserialize(msg, &appMsg)
	if err != nil {
		fmt.Printf("Cannot deserialize application message: %s\n", err.Error())
		return
	}

	socksMessageS := appMsg.Payload
	socksMessage := messages.SocksMessage{}
	err = messages.Deserialize(socksMessageS, &socksMessage)
	if err != nil {
		fmt.Printf("Cannot deserialize socks message: %s\n", err.Error())
		return
	}

	remoteAddr := socksMessage.RemoteAddr

	connPointer, ok := self.connections[remoteAddr]
	if !ok {
		fmt.Printf("Cannot fint the connection with remote address %s\n", remoteAddr)
		return
	}

	conn := *connPointer

	if socksMessage.NeedClose {
		fmt.Printf("Closing connection %s according to a signal from server\n", remoteAddr)
		conn.Close()
		delete(self.connections, remoteAddr)
		return
	}

	data := socksMessage.Data
	_, err = conn.Write(data)
	if err != nil {
		fmt.Printf("Cannot write to connection with remote address %s, error is %s\n", socksMessage.RemoteAddr, err.Error())
		socksMessage := messages.SocksMessage{
			nil,
			remoteAddr,
			true,
		}
		socksMessageS := messages.Serialize(messages.MsgSocksMessage, socksMessage)
		self.Send(socksMessageS)
	}
}
