package main

import (
	"fmt"
	"github.com/skycoin/skywire/src/aether"
	"github.com/skycoin/skywire/src/cipher"
	"github.com/skycoin/skywire/src/daemon"
	"github.com/skycoin/skywire/src/lib/gnet"
	//"log"
	//"time"
)

func main() {

	//create the daemon
	config := daemon.NewConfig()
	//config.Daemon.LocalhostOnly = true
	//config.DHT.Disabled = true
	config.Daemon.Port = 8080
	daemon := daemon.NewDaemon(config)

	//create aether server

	pubkey, seckey := hashchain.GenerateDeterministicKeyPair([]byte("seed"))
	_ = seckey
	_ = pubkey

	a := aether.NewAetherServer(pubkey)

	d1.ServiceManager.AddService(
		[]byte("test service"),
		[]byte("{service=\"test service\"}"), 1, tss1)

	//start daemon mainloop
	go d1.Start(quit1)

}
