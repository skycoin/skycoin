package mesh_rpc

import (
	"encoding/gob"
	"github.com/satori/go.uuid"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/skycoin/skycoin/src/mesh2/messages"
	"github.com/skycoin/skycoin/src/mesh2/node"
)

type RPC struct {
	rpcChan      chan RPCMessage
	backChannels map[uuid.UUID]chan []byte
}

type RPCMessage struct {
	Command    string
	Arguments  Args
	BackChanId uuid.UUID
}

type Args []interface{}

func NewRPC() *RPC {
	newRPC := &RPC{}
	newRPC.rpcChan = make(chan RPCMessage)
	newRPC.backChannels = make(map[uuid.UUID]chan []byte)
	return newRPC
}

func (self *RPC) Start() {
	go self.serveRPC()
	time.Sleep(1 * time.Second)
	go self.clientRPC()
}

func (self *RPC) serveRPC() {
	gob.Register(&node.Node{})
	receiver := new(RPCReceiver)
	err := rpc.Register(receiver)
	if err != nil {
		panic(err)
	}
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	log.Println("Serving RPC")
	err = http.Serve(l, nil)
	if err != nil {
		panic(err)
	}
}

func (self *RPC) clientRPC() {
	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	var result []byte
	for {
		select {
		case msg := <-self.rpcChan:
			{
				err = client.Call("RPCReceiver."+msg.Command, msg.Arguments, &result)
				if err != nil {
					panic(err)
				}
				back := self.backChannels[msg.BackChanId]
				back <- result
			}
		}
	}
}

func (self *RPC) sendToRPC(command string, args []interface{}) []byte {
	backChan := make(chan []byte)
	bcId := uuid.NewV4()
	self.backChannels[bcId] = backChan
	msg := RPCMessage{
		command,
		args,
		bcId,
	}
	self.rpcChan <- msg
	result := <-backChan
	close(backChan)
	delete(self.backChannels, bcId)
	return result
}

func (self *RPC) CreateControlChannel(node *node.Node) uuid.UUID {
	command := "CreateControlChannel"
	args := []interface{}{node}
	res := self.sendToRPC(command, args)

	result := uuid.UUID{}
	err := messages.Deserialize(res, &result)
	if err != nil {
		log.Println(err)
	}
	log.Println("CHANNEL ID:", result)
	return result
}
