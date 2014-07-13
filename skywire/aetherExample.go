package main

import (
	//"fmt"
	"github.com/skycoin/skywire/src/aether"
	"github.com/skycoin/skywire/src/cipher"
	"github.com/skycoin/skywire/src/daemon"
	//"github.com/skycoin/skywire/src/lib/gnet"
	//"log"
	//"time"
)

func main() {

	var quit1 chan int

	//create the daemon
	config := daemon.NewConfig()
	//config.Daemon.LocalhostOnly = true
	//config.DHT.Disabled = true
	config.Daemon.Port = 8080
	daemon := daemon.NewDaemon(config)

	//create aether server

	pubkey, seckey := cipher.GenerateDeterministicKeyPair([]byte("seed"))
	_ = seckey
	_ = pubkey

	a := aether.NewAetherServer(pubkey)

	a.RegisterWithDaemon(daemon)

	//start daemon mainloop
	go daemon.Start(quit1)

}
