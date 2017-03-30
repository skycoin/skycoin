package app

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"golang.org/x/net/proxy"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/go-socks5"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type SocksServer struct {
	app
	socksAddress string
	Socks        *socks5.Server
	dialer       proxy.Dialer
	connections  map[string]net.Conn
}

func NewSocksServer(meshnet messages.Network, address cipher.PubKey, socksAddress string) (*SocksServer, error) {
	socksServer := &SocksServer{}
	socksServer.register(meshnet, address)
	socksServer.lock = &sync.Mutex{}
	socksServer.timeout = time.Duration(messages.GetConfig().AppTimeout)
	socksServer.responseChannels = make(map[uint32]chan messages.AppResponse)
	socksServer.socksAddress = socksAddress
	socksServer.connections = map[string]net.Conn{}

	conn, err := meshnet.NewConnection(address)
	if err != nil {
		return nil, err
	}

	socksServer.connection = conn

	err = meshnet.Register(address, socksServer)
	if err != nil {
		return nil, err
	}

	go socksServer.Serve()
	fmt.Println("ready to accept requests")

	return socksServer, nil
}

func (self *SocksServer) Consume(msg []byte) {
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
	/*
		dataReader, dataWriter := io.Pipe()
		replyReader, replyWriter := io.Pipe()

		go self.Socks.ServeConn(dataReader, replyWriter)
		go self.getFromSocks(replyReader, socksMessage.RemoteAddr)

		n, err := dataWriter.Write(socksMessage.Data)
		if err != nil {
			panic(err)
		}
	*/

	data := socksMessage.Data
	remoteAddr := socksMessage.RemoteAddr
	socksConn, ok := self.connections[remoteAddr]
	if !ok {
		socksConn, err = self.dialer.Dial("tcp", self.socksAddress)
		if err != nil {
			panic(err)
		}
		self.connections[remoteAddr] = socksConn
		go self.getFromSocks(socksConn, remoteAddr)
	}

	_, err = socksConn.Write(data)
	if err != nil {
		panic(err)
	}

	//	fmt.Println("sent to socks server:", n)
}

func (self *SocksServer) Serve() {

	go func() {
		conf := &socks5.Config{}
		server, err := socks5.New(conf)
		if err != nil {
			panic(err)
		}

		self.Socks = server
		if err := server.ListenAndServe("tcp", self.socksAddress); err != nil {
			panic(err)
		}
	}()
	dialer, err := proxy.SOCKS5("tcp", self.socksAddress, nil, proxy.Direct)
	if err != nil {
		panic(err)
	}
	self.dialer = dialer
}

func (self *SocksServer) getFromSocks(conn net.Conn, remoteAddr string) {

	for {
		//		fmt.Println("waiting from socks:", conn.RemoteAddr(), conn.LocalAddr())
		buffer := make([]byte, config.SocksPacketSize)
		n, err := conn.Read(buffer)
		if err != nil {
			return
			if err == io.EOF {
				continue
			} else {
				panic(err)
			}
		}
		//		fmt.Printf("got from %s to %s %d bytes\n", conn.RemoteAddr(), conn.LocalAddr(), n)
		socksMessage := messages.SocksMessage{
			buffer[:n],
			remoteAddr,
		}
		socksMessageS := messages.Serialize(messages.MsgSocksMessage, socksMessage)
		response := &messages.AppMessage{
			self.sequence,
			false,
			socksMessageS,
		}
		responseSerialized := messages.Serialize(messages.MsgAppMessage, response)
		self.send(responseSerialized)
		self.sequence++
	}
}
