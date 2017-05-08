package app

import (
	"io"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/net/proxy"

	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/proxy/go-socks5"
)

type SocksServer struct {
	proxyServer
	dialer proxy.Dialer
}

func NewSocksServer(appId messages.AppId, nodeAddr string, proxyAddress string) (*SocksServer, error) {
	socksServer := &SocksServer{}
	socksServer.id = appId
	socksServer.lock = &sync.Mutex{}
	socksServer.timeout = time.Duration(messages.GetConfig().AppTimeout)
	socksServer.responseNodeAppChannels = make(map[uint32]chan bool)
	socksServer.ProxyAddress = proxyAddress
	socksServer.targetConns = map[string]net.Conn{}

	err := socksServer.RegisterAtNode(nodeAddr)
	if err != nil {
		return nil, err
	}

	go socksServer.serveSocks()
	log.Println("ready to accept requests")

	return socksServer, nil
}

func (self *SocksServer) consume(msg *messages.AppMessage) {

	proxyMessage := getProxyMessage(msg)
	if proxyMessage == nil {
		return
	}

	remoteAddr := proxyMessage.RemoteAddr // user address
	needClose := proxyMessage.NeedClose   // the message can be a comand to close the coresponding connection

	self.lock.Lock()
	targetConn, ok := self.targetConns[remoteAddr] // find the existing connection
	self.lock.Unlock()

	if needClose { // if we got a command to close a connection but there is no such one (already closed), just return
		if ok {
			log.Printf("Closing connection %s according to a signal from client\n", remoteAddr)
			self.closeConns(remoteAddr)
		}
		return
	}

	if !ok && !needClose { // otherwise if there is no such connection create one
		var err error
		targetConn, err = self.dialer.Dial("tcp", self.ProxyAddress)
		if err != nil {
			log.Println("Cannot dial to proxy server on ", self.ProxyAddress)
			return
		}

		self.lock.Lock()
		self.targetConns[remoteAddr] = targetConn
		self.lock.Unlock()

		go self.getFromConn(targetConn, remoteAddr)
	}

	// if we haven't returned from the procedure yet, write the data to the connection

	data := proxyMessage.Data
	_, err := targetConn.Write(data)
	if err != nil { // if write is unsuccessful, close this connection and send the closing command to the corresponding client connection
		log.Println("Cannot write to connection:", targetConn.LocalAddr().String(), targetConn.RemoteAddr().String())
		self.sendClose(remoteAddr)
		self.closeConns(remoteAddr)
	}
}

func (self *SocksServer) serveSocks() { //run socks server and get dialer
	go func() {
		conf := &socks5.Config{}
		server, err := socks5.New(conf)
		if err != nil {
			panic(err)
		}

		if err := server.ListenAndServe("tcp", self.ProxyAddress); err != nil {
			panic(err)
		}
	}()
	dialer, err := proxy.SOCKS5("tcp", self.ProxyAddress, nil, proxy.Direct)
	if err != nil {
		panic(err)
	}
	self.dialer = dialer
}

func (self *SocksServer) RegisterAtNode(nodeAddr string) error {

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

func (self *SocksServer) listenFromNode() {
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

func (self *SocksServer) handleIncomingFromNode(msg []byte) error {
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
