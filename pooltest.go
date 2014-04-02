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

	//process data and connection events
	go func() {

		//required for connection event
		for true {
			time.Sleep(time.Second * 1)
			fmt.Printf("wtf \n")
			cpool1.HandleMessages()
			cpool2.HandleMessages()
		}

	}()

	con, err := cpool1.Connect("127.0.0.1:6061")
	_ = con
	if err != nil {
		log.Panic(err)
	}

	cpool1.SendMessage(con, 0, []byte("test message"))

	time.Sleep(time.Second * 10)
}
