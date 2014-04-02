package main

import (
	"github.com/skycoin/skywire/src/lib/gnet"
	"log"
	"time"
)

func onConnect(c *gnet.Connection, solicited bool) {
	fmt.Printf("connnect event")
}

func onDisconnect(c *gnet.Connection,
	reason gnet.DisconnectReason) {
	fmt.Printf("disconnect event")
}

func onMessage(c *gnet.Connection, channel uint16,
	msg []byte) error {

	fmt.Printf("message event: channel %v", channel)
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

	err := cpool2.StartListen()
	if err = nil {
		log.Panic()
	}

	//blocks, run in goroutine
	go cpool1.AcceptConnections()
	go cpool2.AcceptConnections()

	go {
		con, err := cpoo1.Connect("127.0.0.1:6061")
		_ = con
		if err != nil {
			log.Panic(err)
		}
	}
	cpool1 
	for true {
		time.Sleep(time.Second * 1)
		cpool1.HandleMessages()
		cpool2.HandleMessages()
	}
	
	//HandleMessages()

	//go cpool1.Accept()

/*
	_ = cpool

	//pool tickers
	clearStaleConnectionsTicker := time.Tick(self.Pool.Config.ClearStaleRate)
	idleCheckTicker := time.Tick(self.Pool.Config.IdleCheckRate)
	messageHandlingTicker := time.Tick(self.Pool.Config.MessageHandlingRate)
*/
}
