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
	socksClient.responseChannels = make(map[uint32]chan messages.AppResponse)

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

	responseChannel := make(chan messages.AppResponse)
	sequence := self.setResponseChannel(responseChannel)

	request := &messages.AppMessage{
		sequence,
		true,
		msg,
	}
	requestSerialized := messages.Serialize(messages.MsgAppMessage, request)
	self.send(requestSerialized)
	//	fmt.Println("sent", requestSerialized)
}

func (self *SocksClient) Listen() {
	proxyAddress := self.socksAddress
	l, err := net.Listen("tcp", proxyAddress)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		//		fmt.Println(conn.LocalAddr(), conn.RemoteAddr())
		remoteAddr := conn.RemoteAddr().String()
		self.connections[remoteAddr] = &conn

		go func() {
			for {
				message := make([]byte, config.SocksPacketSize)
				//	fmt.Println("ready for the next packet")
				n, err := conn.Read(message)
				if err != nil {
					return
					if err == io.EOF {
						continue
					} else {
						panic(err)
					}
				}
				//				fmt.Println("got from client:", n)

				socksMessage := messages.SocksMessage{
					message[:n],
					remoteAddr,
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
		fmt.Printf("Cannot consume a message: %s\n", err.Error())
		return
	}

	socksMessageS := appMsg.Payload
	socksMessage := messages.SocksMessage{}
	err = messages.Deserialize(socksMessageS, &socksMessage)
	if err != nil {
		panic(err)
	}
	connPointer, ok := self.connections[socksMessage.RemoteAddr]
	if !ok {
		panic(socksMessage.RemoteAddr)
	}
	conn := *connPointer
	data := socksMessage.Data
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	//	fmt.Printf("\nresponse is written to local %s, remote %s, n = %d\n\n", conn.LocalAddr(), conn.RemoteAddr(), n)
}
