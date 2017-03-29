package app

import (
	"fmt"
	//	"io"
	"net"
	"sync"
	"time"

	"github.com/armon/go-socks5"
	"golang.org/x/net/proxy"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type SocksServer struct {
	app
	socksConn net.Conn
}

func NewSocksServer(meshnet messages.Network, address cipher.PubKey, socksAddress string) (*SocksServer, error) {
	socksServer := &SocksServer{}
	socksServer.register(meshnet, address)
	socksServer.lock = &sync.Mutex{}
	socksServer.timeout = time.Duration(messages.GetConfig().AppTimeout)
	socksServer.socksAddress = socksAddress

	conn, err := meshnet.NewConnection(address)
	if err != nil {
		return nil, err
	}

	socksServer.connection = conn

	err = meshnet.Register(address, socksServer)
	if err != nil {
		return nil, err
	}

	return socksServer, nil
}

func (self *SocksServer) Consume(msg []byte) {
	appMsg := messages.AppMessage{}
	err := messages.Deserialize(msg, &appMsg)
	if err != nil {
		fmt.Printf("Cannot consume a message: %s\n", err.Error())
		return
	}

	originalRequest := appMsg.Payload

	n, err := self.socksConn.Write(originalRequest)
	fmt.Println("sent to SOCKS:", n)
}

func (self *SocksServer) Serve() {

	socksAddress := self.socksAddress

	go func() {
		conf := &socks5.Config{}
		server, err := socks5.New(conf)
		if err != nil {
			panic(err)
		}

		if err := server.ListenAndServe("tcp", socksAddress); err != nil {
			panic(err)
		}
	}()

	dialer, err := proxy.SOCKS5("tcp", socksAddress, nil, proxy.Direct)
	if err != nil {
		panic(err)
	}

	for {
		socksConn, err := dialer.Dial("tcp", socksAddress)
		if err == nil {
			self.socksConn = socksConn
			break
		}
	}

	fmt.Println("ready to accept requests")
	self.getFromSocks()
}

func (self *SocksServer) getFromSocks() {

	for {
		fmt.Println("Waiting from SOCKS")
		buffer := make([]byte, config.SocksPacketSize)
		n, err := self.socksConn.Read(buffer)
		if err != nil {
			//			if err != io.EOF {
			panic(err)
			//			}
		}
		fmt.Println("got from SOCKS:", n)
		response := &messages.AppMessage{
			0,
			false,
			buffer[:n],
		}
		responseSerialized := messages.Serialize(messages.MsgAppMessage, response)
		self.send(responseSerialized)
	}
}
