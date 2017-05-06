package viscript

import (
	"log"
	"net"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

type ViscriptServer struct {
	conn          net.Conn
	maxPacketSize int
	closeChannel  chan bool
}

const viscriptAddr = "127.0.0.1:7999"

func (self *ViscriptServer) Init(sequence, appId uint32) {
	conn, err := net.Dial("tcp", viscriptAddr)
	if err != nil {
		panic(err)
	}
	self.conn = conn
	log.Println("Waiting for Viscript messages")
	self.sendFirstAck(sequence, appId)
	self.serve()
}

func (self *ViscriptServer) Shutdown() {
	self.conn.Close()
}

func (self *ViscriptServer) serve() {
	go_on := true
	go func() {
		for go_on {

			buffer := make([]byte, self.maxPacketSize)

			n, err := self.conn.Read(buffer)

			if err != nil {
				if !go_on && n == 0 {
					break
				} else {
					panic(err)
				}
			} else {
				if n > 0 {
					log.Printf("connection at %s received %d bytes\n", self.conn.LocalAddr().String(), n)
					uc := messages.UserCommand{}
					err := messages.Deserialize(buffer[:n], &uc)
					if err != nil {
						log.Println("Incorrect UserCommand:", buffer[:n])
						continue
					}
					go self.handleUserCommand(&uc)
				}
			}
		}
	}()
	<-self.closeChannel
	go_on = false
	self.conn.Close()
}

func (self *ViscriptServer) handleUserCommand(uc *messages.UserCommand) {
	log.Println("command received:", uc)
}

func (self *ViscriptServer) SendAck(ackS []byte, sequence, appId uint32) {
	ucAck := &messages.UserCommand{
		sequence,
		appId,
		ackS,
	}
	ucAckS := messages.Serialize(messages.MsgUserCommand, ucAck)
	self.send(ucAckS)
}

func (self *ViscriptServer) send(data []byte) {
	_, err := self.conn.Write(data)
	if err != nil {
		log.Println("Unsuccessful sending to viscript from nodemanager")
	}
}

func (self *ViscriptServer) sendFirstAck(sequence, appId uint32) {
	ack := messages.CreateAck{}
	ackS := messages.Serialize(messages.MsgCreateAck, ack)
	self.SendAck(ackS, sequence, appId)
}
