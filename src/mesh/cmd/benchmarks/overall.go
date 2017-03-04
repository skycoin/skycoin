package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/app"
	"github.com/skycoin/skycoin/src/mesh/messages"
	network "github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func main() {

	//messages.SetDebugLogLevel()
	messages.SetInfoLogLevel()

	var (
		size int
		err  error
	)

	sizeStr := strings.ToLower(os.Args[1])

	kb := 1024
	mb := kb * kb
	gb := mb * kb

	if strings.HasSuffix(sizeStr, "mb") {
		sizemb := strings.TrimSuffix(sizeStr, "mb")
		size, err = strconv.Atoi(sizemb)
		if err != nil {
			panic(err)
		}
		size *= mb
	} else if strings.HasSuffix(sizeStr, "gb") {
		sizegb := strings.TrimSuffix(sizeStr, "gb")
		size, err = strconv.Atoi(sizegb)
		if err != nil {
			panic(err)
		}
		size *= gb
	} else if strings.HasSuffix(sizeStr, "kb") {
		sizekb := strings.TrimSuffix(sizeStr, "kb")
		size, err = strconv.Atoi(sizekb)
		if err != nil {
			panic(err)
		}
		size *= kb
	} else {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			panic(err)
		}
		sizeStr += "b"
	}

	networkSize := 3

	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()

	clientAddr, serverAddr := meshnet.CreateSequenceOfNodes(networkSize)

	_, err = echoServer(meshnet, serverAddr)
	if err != nil {
		panic(err)
	}

	client, err := app.NewClient(meshnet, clientAddr) // register client on the first node
	if err != nil {
		panic(err)
	}

	err = client.Dial(serverAddr) // client dials to server
	if err != nil {
		panic(err)
	}

	duration := benchmark(client, size)
	log.Println(sizeStr+" duration:", duration)
}

func benchmark(client *app.Client, msgSize int) time.Duration {

	if msgSize < 1 {
		panic("message should be at least 1 byte")
	}
	packetSize := messages.GetConfig().MaxPacketSize / 2
	packets := (msgSize-1)/packetSize + 1

	msg := make([]byte, packetSize)

	wg := &sync.WaitGroup{}
	wg.Add(packets)

	start := time.Now()
	for p := 0; p < packets; p++ {
		go func() {
			//retChan := make(chan *app.ConnResponse, 1024)
			retChan := client.Send(msg)
			response := <-retChan
			wg.Done()
			if response.Err != nil {
				panic(response.Err)
			}
		}()
	}

	wg.Wait()

	duration := time.Now().Sub(start)

	return duration
}

func echoServer(meshnet *network.NodeManager, serverAddr cipher.PubKey) (*app.Server, error) {

	srv, err := app.NewServer(meshnet, serverAddr, func(in []byte) []byte {
		return in
	})
	return srv, err
}
