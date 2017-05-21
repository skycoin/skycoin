package main

import (
	"fmt"
	"log"
	"time"

	gnet "github.com/skycoin/skycoin/src/aether/gnet"
)

//this is called when client connects
func onConnect(c *gnet.Connection, solicited bool) {
	fmt.Printf("Event Callback: connnect event \n")
}

//this is called when client disconnects
func onDisconnect(c *gnet.Connection,
	reason gnet.DisconnectReason) {
	fmt.Printf("Event Callback: disconnect event \n")
}

//this is called when a message is received
func onMessage(c *gnet.Connection, channel uint16,
	msg []byte) error {
	fmt.Printf("Event Callback: message event: addr= %s, channel %v, msg= %s \n", c.Addr(), channel, msg)
	return nil
}

func SpawnConnectionPool(Port int) *gnet.ConnectionPool {
	config := gnet.NewConfig()
	config.Port = uint16(Port)               //set listening port
	config.DisconnectCallback = onDisconnect //disconnect callback
	config.ConnectCallback = onConnect       //connect callback
	config.MessageCallback = onMessage       //message callback

	//create pool
	cpool := gnet.NewConnectionPool(config)
	//open lsitening port
	if err := cpool.StartListen(); err != nil {
		log.Panic(err)
	}

	//listen for connections in new goroutine
	go cpool.AcceptConnections()

	//handle income data in new goroutine
	go func() {
		for true {
			time.Sleep(time.Millisecond * 100)
			cpool.HandleMessages()
		}
	}()

	return cpool
}

//define message we want to be able to handle
type TestMessage struct {
	Text []byte
}

func (self *TestMessage) Handle(context *gnet.MessageContext, state interface{}) error {

	fmt.Printf("TestMessage Handle: Text= %s \n", string(self.Text))
	return nil
}

//array of messages to register
var messageMap map[string](interface{}) = map[string](interface{}){
	"id01": TestMessage{}, //message id, message type
}

//create connection pool and tests
func main() {

	cpool1 := SpawnConnectionPool(6060)
	cpool2 := SpawnConnectionPool(6061)

	//connect to peer
	con, err := cpool1.Connect("127.0.0.1:6061")
	_ = con

	if err != nil {
		log.Panic(err)
	}

	//new dispatch manager
	dm1 := gnet.NewDispatcherManager()
	cpool1.Config.MessageCallback = dm1.OnMessage //set message handler
	//new dispatcher for handling messages on channel 3
	d1 := dm1.NewDispatcher(cpool1, 3, nil) //dispatcher 1
	d1.RegisterMessages(messageMap)

	dm2 := gnet.NewDispatcherManager()
	cpool2.Config.MessageCallback = dm2.OnMessage
	d2 := dm2.NewDispatcher(cpool2, 3, nil) //dispatcher 2
	d2.RegisterMessages(messageMap)

	//create a message to send
	tm := TestMessage{Text: []byte("Message test")}
	d1.SendMessage(con, 3, &tm)

	time.Sleep(time.Second * 10)
}
