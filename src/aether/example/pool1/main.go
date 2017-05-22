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

//create connection pool and tests
func main() {

	cpool1 := SpawnConnectionPool(6060)
	cpool2 := SpawnConnectionPool(6061)

	//connect to peer
	con, err := cpool1.Connect("127.0.0.1:6061")

	if err != nil {
		log.Panic(err)
	}

	//send a raw binary message on channel 0
	cpool1.SendMessage(con, 0, []byte("test message"))
	_ = cpool2

	time.Sleep(time.Second * 10)
}
