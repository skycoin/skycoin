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
func openUDPConn(tr *Transport, peer, pairPeer *messages.Peer) (*UDPConfig, error) {
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
	_, err := self.conn.WriteTo(msg, self.pairAddr)
	return err
}

// receive - listen to port and send to incoming channel
//udp listens to []byte then passes it to incomingChannel, maybe decrypts it first

func (self *UDPConfig) receiveLoop() {
	go_on := true
	go func() {
		for go_on {

			buffer := make([]byte, self.maxPacketSize)

			n, _, err := self.conn.ReadFrom(buffer)

			if err != nil {
				if !go_on && n == 0 {
					break
				} else {
					panic(err)
				}
			} else {
				go self.relatedTransport.getFromUDP(buffer[:n])
			}
		}
	}()
	<-self.closeChannel
	go_on = false
	self.conn.Close()
}
