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
	dialer       proxy.Dialer
	connections  map[string]net.Conn
}

func NewSocksServer(meshnet messages.Network, address cipher.PubKey, socksAddress string) (*SocksServer, error) {
	socksServer := &SocksServer{}
	socksServer.register(meshnet, address)
	socksServer.lock = &sync.Mutex{}
	socksServer.timeout = time.Duration(messages.GetConfig().AppTimeout)
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
	needClose := socksMessage.NeedClose

	socksConn, ok := self.connections[remoteAddr]

	if !ok && needClose {
		return
	}

	if !ok && !needClose {
		socksConn, err = self.dialer.Dial("tcp", self.socksAddress)
		if err != nil {
			fmt.Println("Cannot dial to socks server on ", self.socksAddress)
			return
		}
		self.connections[remoteAddr] = socksConn
		go self.getFromSocks(socksConn, remoteAddr)
	}

	if needClose {
		fmt.Printf("Closing connection %s according to a signal from client\n", remoteAddr)
		socksConn.Close()
		delete(self.connections, remoteAddr)
		return
	}

	data := socksMessage.Data
	_, err = socksConn.Write(data)
	if err != nil {
		fmt.Println("Cannot write to connection:", socksConn.LocalAddr().String(), socksConn.RemoteAddr().String())
		closingMessage := messages.SocksMessage{
			nil,
			remoteAddr,
			true,
		}
		closingMessageS := messages.Serialize(messages.MsgSocksMessage, closingMessage)
		self.Send(closingMessageS)
	}
}

func (self *SocksServer) Serve() {

	go func() {
		conf := &socks5.Config{}
		server, err := socks5.New(conf)
		if err != nil {
			panic(err)
		}

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
		buffer := make([]byte, config.SocksPacketSize)
		n, err := conn.Read(buffer)
		if err != nil {
			return
			if err == io.EOF {
				continue
			} else {
				conn.Close()
				return
			}
		}

		socksMessage := messages.SocksMessage{
			buffer[:n],
			remoteAddr,
			false,
		}
		socksMessageS := messages.Serialize(messages.MsgSocksMessage, socksMessage)
		self.Send(socksMessageS)
	}
}

func (self *SocksServer) Send(data []byte) {
	message := &messages.AppMessage{
		0,
		false,
		data,
	}
	messageS := messages.Serialize(messages.MsgAppMessage, message)
	self.send(messageS)
}
