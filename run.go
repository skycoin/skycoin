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
 "github.com/skycoin/skycoin/src/sync"
)

func main() {

	dcfg := sync.NewDaemonConfig()
	daemon := sync.NewDaemon(dcfg)
}