package main

import (
	"github.com/skycoin/skywire/src/lib/gnet"
)

func onConnect(c *gnet.Connection, solicited bool) {

}

func onDisconnect(c *gnet.Connection,
	reason gnet.DisconnectReason) {

}

func onMessage(c *gnet.Connection, channel uint16,
	msg []byte) error {
	return nil
}

//create connection pool and tests

func main() {

	config := gnet.NewConfig()
	config.Port = 6060
	config.DisconnectCallback = onDisconnect
	config.ConnectCallback = onConnect
	config.MessageCallback = onMessage

	_ := gnet.NewConnectionPool(config)
}
