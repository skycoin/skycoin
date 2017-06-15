package nodemanager

import (
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type MsgServer struct {
	nm               *NodeManager
	conn             *net.UDPConn
	maxPacketSize    int
	closeChannel     chan bool
	nodeAddrs        map[cipher.PubKey]*net.UDPAddr
	sequence         uint32
	responseChannels map[uint32]chan bool
	timeout          time.Duration
	lock             *sync.Mutex
}

func newMsgServer(nm *NodeManager) (*MsgServer, error) {
	msgSrv := &MsgServer{}
	msgSrv.nm = nm

	msgSrv.responseChannels = make(map[uint32]chan bool)

	msgSrv.maxPacketSize = config.MaxPacketSize
	msgSrv.timeout = time.Duration(config.MsgSrvTimeout) * time.Millisecond

	fullhost := nm.ctrlAddr

	hostdata := strings.Split(fullhost, ":")
	if len(hostdata) != 2 {
		return nil, messages.ERR_INCORRECT_HOST
	}

	host := net.ParseIP(hostdata[0])
	port, err := strconv.Atoi(hostdata[1])
	if err != nil {
		return nil, messages.ERR_INCORRECT_HOST
	}

	addr := &net.UDPAddr{IP: host, Port: port}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	msgSrv.conn = conn
	msgSrv.closeChannel = make(chan bool)

	msgSrv.nodeAddrs = make(map[cipher.PubKey]*net.UDPAddr)

	msgSrv.lock = &sync.Mutex{}
	msgSrv.sequence = uint32(1) // 0 for no-wait sends

	go msgSrv.receiveLoop()

	return msgSrv, nil
}

// close
func (self *MsgServer) shutdown() {
	self.closeChannel <- true
}

func (self *MsgServer) sendMessage(node cipher.PubKey, msg []byte) error {

	responseChannel := make(chan bool)
	sequence := self.sequence
	self.lock.Lock()
	self.responseChannels[sequence] = responseChannel
	self.lock.Unlock()
	self.sequence++

	err := self.send(sequence, node, msg)
	if err != nil {
		return err
	}

	select {
	case <-responseChannel:
		return nil
	case <-time.After(self.timeout * time.Millisecond):
		return messages.ERR_MSG_SRV_TIMEOUT
	}
}

func (self *MsgServer) sendNoWait(node cipher.PubKey, msg []byte) error {
	err := self.send(uint32(0), node, msg)
	return err
}

func (self *MsgServer) sendAck(sequence uint32, node cipher.PubKey, msg []byte) error {
	err := self.send(sequence, node, msg)
	return err
}

func (self *MsgServer) send(sequence uint32, node cipher.PubKey, msg []byte) error {
	self.lock.Lock()
	addr := self.nodeAddrs[node]
	self.lock.Unlock()

	inControlMsg := messages.InControlMessage{
		messages.ChannelId(0),
		sequence,
		msg,
	}
	inControlS := messages.Serialize(messages.MsgInControlMessage, inControlMsg)
	_, err := self.conn.WriteTo(inControlS, addr)
	return err
}

func (self *MsgServer) receiveLoop() {
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
				cm := messages.InControlMessage{}
				err := messages.Deserialize(buffer[:n], &cm)
				if err != nil {
					log.Println("Incorrect InControlMessage:", buffer[:n])
					continue
				}
				go self.nm.handleControlMessage(&cm)
			}
		}
	}()
	<-self.closeChannel
	go_on = false
	self.conn.Close()
}

func (self *MsgServer) getResponse(sequence uint32, response *messages.CommonCMAck) {
	self.lock.Lock()
	responseChannel, ok := self.responseChannels[sequence]
	self.lock.Unlock()
	if !ok {
		return
	}
	responseChannel <- response.Ok
}
