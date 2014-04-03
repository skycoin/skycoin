package main

import (
	"fmt"
	"github.com/skycoin/skywire/src/lib/gnet"
	"log"
	"time"
)

func onConnect(c *gnet.Connection, solicited bool) {
	fmt.Printf("connnect event \n")
}

func onDisconnect(c *gnet.Connection,
	reason gnet.DisconnectReason) {
	fmt.Printf("disconnect event \n")
}

func onMessage(c *gnet.Connection, channel uint16,
	msg []byte) error {

	fmt.Printf("message event: channel %v, msg= %s \n", channel, msg)
	return nil
}

type TestMessage struct {
	Text []byte
}

func (self *TestMessage) Handle(context *gnet.MessageContext, state interface{}) error {
	fmt.Printf("Handle Test Message: Text= %s ", string(self.Text))
}

//create connection pool and tests
func main() {

	config := gnet.NewConfig()
	config.Port = 6060
	config.DisconnectCallback = onDisconnect
	config.ConnectCallback = onConnect
	config.MessageCallback = onMessage

	cpool1 := gnet.NewConnectionPool(config)

	config.Port = 6061
	cpool2 := gnet.NewConnectionPool(config)

	err := cpool1.StartListen()
	if err != nil {
		log.Panic()
	}

	err = cpool2.StartListen()
	if err != nil {
		log.Panic()
	}

	//blocks, run in goroutine
	go cpool1.AcceptConnections()
	go cpool2.AcceptConnections()

	//process data and connection events in another goroutine
	go func() {
		//required for connection event
		for true {
			time.Sleep(time.Second * 1)
			cpool1.HandleMessages()
			cpool2.HandleMessages()
		}

	}()

	//connect to remote peer
	con, err := cpool1.Connect("127.0.0.1:6061")
	_ = con
	if err != nil {
		log.Panic(err)
	}


	//array of messages to register
	var map[string](interface{}) =  {
		"test" : TestMessage{},
	}

	//cpool1.SendMessage(con, 0, []byte("test message"))

	dm1 := gnet.NewDispatcherManager()
	dm1.NewDispatcher(3)                         //channel 3
	cpoo1.Config.MessageCallback = dm1.OnMessage //set message handler

	dm2 := gnet.NewDispatcherManager()
	dm2.NewDispatcher(3)
	cool2.Config.MessageCallback = dm2.OnMessage



	time.Sleep(time.Second * 10)
}
