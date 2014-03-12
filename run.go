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
	"time"
	//"errors"
	"fmt"
	"log"
)

func blobVerify(data []byte) sync.BlobCallbackResponse {

	fmt.Printf("!!! blob: %v \n", string(data))
	return sync.BlobCallbackResponse{
		Valid:    true,  //is data valid (if false, will discard)
		Ignore:   false, //should be on ignore list?
		KickPeer: false, //kick the peer
	}
}

func daemon_spawn(port int) (*sync.Daemon, *sync.BlobReplicator) {
	cfg := sync.NewConfig()
	cfg.Daemon.Port = port
	cfg.DHT.Disabled = true
	cfg.Peers.AllowLocalhost = true
	cfg.Peers.Ephemerial = true //disable load/save to disable
	daemon := sync.NewDaemon(cfg)

	//the callback response
	callback := func(data []byte) sync.BlobCallbackResponse {
		fmt.Printf("port: %v, callback= %v \n", port, string(data))
		return sync.BlobCallbackResponse{
			Valid:    true,  //is data valid (if false, will discard)
			Ignore:   false, //should be on ignore list?
			KickPeer: false, //kick the peer
		}
	}

	br := daemon.NewBlobReplicator(uint16(1), callback) //channel 0

	daemon.Init() // begins listening here
	return daemon, br
}

func testRep() {

	//quit := make(chan int) //write to this to shutdown

}

func main() {

	d1, br1 := daemon_spawn(5050)
	d2, br2 := daemon_spawn(5051)

	//inject
	blobData := []byte("BLOB DATA") //replicate this to world!
	err := br1.InjectBlob(blobData)
	if err != nil {
		log.Panic(err) //inject will fail if blob data is duplicate
	}

	quit := make(chan int) //write to this to shutdown

	//time.Sleep(1000* time.Millisecond)

	// fmt.Printf("sleep done\n")

	addr1 := "127.0.0.1:5050"
	addr2 := "127.0.0.1:5051"

	d1.Peers.Peers.AddPeer(addr2)
	d2.Peers.Peers.AddPeer(addr1)

	go d1.Start(quit) //goroutine
	go d2.Start(quit) //goroutine

	_ = br1
	_ = br2

	for true {
		time.Sleep(50 * time.Millisecond)
	}
	/*
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
	*/

	// for true {
	//     time.Sleep(50)
	// }
	// quit<-1
	// quit<-1
}
