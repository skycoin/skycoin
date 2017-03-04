package transport

import (
	"net"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

type UDPConfig struct {
	relatedTransport *Transport
	addr             net.Addr
	pairAddr         net.Addr
	conn             *net.UDPConn
	maxPacketSize    int
	closeChannel     chan bool
}

// create
func openConn(tr *Transport, peer, pairPeer *messages.Peer) (*UDPConfig, error) {
	maxPacketSize := messages.GetConfig().MaxPacketSize
	host, pairHost := net.ParseIP(peer.Host), net.ParseIP(pairPeer.Host)
	port, pairPort := int(peer.Port), int(pairPeer.Port)
	addr, pairAddr := &net.UDPAddr{IP: host, Port: port}, &net.UDPAddr{IP: pairHost, Port: pairPort}

	udp := &UDPConfig{relatedTransport: tr, addr: addr, pairAddr: pairAddr, maxPacketSize: maxPacketSize}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	udp.conn = conn
	udp.closeChannel = make(chan bool)
	return udp, nil
}

func (self *UDPConfig) Tick() {
	go self.receiveLoop()
}

// close
func (self *UDPConfig) closeConn() {
	self.closeChannel <- true
}

// send - serialize and send to peer

func (self *UDPConfig) send(msg []byte) error {
	//	_, err := self.conn.WriteTo(msg, self.pairAddr)

	//	messages.RegisterEvent("udp.send start")
	/*
		retChan := make(chan error)
		go self.sendGo(msg, retChan)
		err := <-retChan
	*/
	_, err := self.conn.WriteTo(msg, self.pairAddr)

	//	messages.RegisterEvent("udp.send finish")

	return err
}

func (self *UDPConfig) sendGo(msg []byte, retChan chan error) {
	_, err := self.conn.WriteTo(msg, self.pairAddr)
	retChan <- err
}

// receive - listen to port and send to incoming channel
//udp listens to []byte then passes it to incomingChannel, maybe decrypts it first

func (self *UDPConfig) receiveLoop() {
	go_on := true
	incomingChannel := self.relatedTransport.incomingChannel
	go func() {
		for go_on {

			//			messages.RegisterEvent("udp.receiveLoop - making buffer")

			buffer := make([]byte, self.maxPacketSize)

			//			messages.RegisterEvent("udp.receiveLoop - reading from buffer")

			n, _, err := self.conn.ReadFrom(buffer)

			//			messages.RegisterEvent("udp.receiveLoop - read from buffer")

			if err != nil {
				if !go_on && n == 0 {
					break
				} else {
					panic(err)
				}
			} else {
				/*	if addr.String() != self.pairAddr.String() {
					panic("wrong address")
				}*/
				go func() { incomingChannel <- buffer[:n] }()

				//				messages.RegisterEvent("udp.receiveLoop - sent to transport.incomingChannel")

			}
		}
	}()
	<-self.closeChannel
	go_on = false
	self.conn.Close()
}
