package main

import (
	"fmt"
	"log"
	"time"
	//"github.com/skycoin/skycoin/src/aether/dht"
)

/*
 Takes a string, hashes it and finds people who are also looking for people
 who are looking for people with the same hash
*/

func PeerCallback(infoHash string, peerAddress string) {
	fmt.Printf("PeerCallback: infoHash= %s, peerAddres= %s \n", infoHash, peerAddress)
}

func main() {

	config := daemon_dht.NewDHTConfig()
	config.AddPeerCallback = PeerCallback

	dht := daemon_dht.NewDHT(config)

	if err := dht.Init(); err != nil {
		log.Panic()
	}

	go dht.Start()
	//go dht.Listen() //flushes

	for i := 0; i < 10; i++ {
		dht.FlushResults()
		dht.RequestPeers("skycoin-skycoin-skycoin-skycoin-skycoin-skycoin-skycoin")
		time.Sleep(time.Second * 1)
	}

	time.Sleep(time.Second * 60)
	dht.Shutdown()
}
