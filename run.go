package main

/*
import (
 "github.com/skycoin/skycoin/src/blockchain"
)

func main() {
	blockchain_server := blockchain.NewServer(blockchain.ServerConfig{})
	blockchain_server.Start()
}
*/

import (
	"time"
	"github.com/skycoin/skycoin/src/sync"
)

func main() {

	//dcfg := sync.NewDaemonConfig()
	cfg := sync.NewConfig()
	daemon := sync.NewDaemon(cfg)
	quit := make(chan int) //write to this to shutdown
	daemon.Start(quit)	

	for true {
		time.Sleep(50)
	}
}
