package aether

import (
	"fmt"

	"github.com/skycoin/skycoin/src/aether/daemon"
	gnet "github.com/skycoin/skycoin/src/aether/gnet"
	"github.com/skycoin/skycoin/src/cipher"
)

/*
	TODO:
	- finish
*/

type AetherServer struct {
	Service *gnet.Service // Service
}

func NewAetherServer(pubkey cipher.PubKey) *AetherServer {
	var x AetherServer
	return &x
}

func (self *AetherServer) RegisterWithDaemon(daemon *daemon.Daemon) {
	self.Service = daemon.ServiceManager.AddService(
		[]byte("test service"),
		[]byte("{service=\"test service\"}"), 1, self)
}

func (self *AetherServer) OnConnect(c *gnet.Connection) {
	fmt.Printf("AetherServer: OnConnect, addr= %s \n", c.Addr())
}

func (self *AetherServer) OnDisconnect(c *gnet.Connection) {
	fmt.Printf("AetherServer: OnDisconnect, addr= %s \n", c.Addr())
}

func (self *AetherServer) RegisterMessages(d *gnet.Dispatcher) {
	var messageMap map[string](interface{}) = map[string](interface{}){
		//put messages here
		"id01": TestMessage{}, //message id, message type
	}
	d.RegisterMessages(messageMap)
}

//define message we want to be able to handle
type TestMessage struct {
	Text []byte
}

func (self *TestMessage) Handle(context *gnet.MessageContext, state interface{}) error {
	server := state.(*AetherServer) //service server state

	fmt.Printf("TestMessage Handle: ServiceIdLong= %s, Text= %s \n", string(server.Service.IdLong), string(self.Text))
	return nil
}
