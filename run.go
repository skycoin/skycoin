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
	//"errors"
	"log"
	"fmt"
)


func blobVerify(data []byte) sync.BlobCallbackResponse {

	fmt.Printf("blob: %v \n", string(data))
	return sync.BlobCallbackResponse {
		Valid : true, //is data valid (if false, will discard)
		Ignore : false, //should be on ignore list?
		KickPeer: false, //kick the peer
	}
}


func main() {

	cfg := sync.NewConfig()
	cfg.Daemon.Port = 5050
	cfg.DHT.Disabled = true
	daemon := sync.NewDaemon(cfg)

	quit := make(chan int) //write to this to shutdown

	br := daemon.NewBlobReplicator(uint16(0), blobVerify) //channel 0

	blobData := []byte ("BLOB DATA") //replicate this to world!
	err := br.InjectBlob(blobData)
	if err != nil {
		log.Panic(err) //inject will fail if blob data is duplicate
	}

	daemon.Start(quit)	

	for true {
		time.Sleep(50)
	}
}
