package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/mesh/app"
	"github.com/skycoin/skycoin/src/mesh/messages"
	network "github.com/skycoin/skycoin/src/mesh/nodemanager"
)

var config = messages.GetConfig()

func main() {

	//messages.SetDebugLogLevel()
	messages.SetInfoLogLevel()

	var (
		size int
		err  error
	)

	args := os.Args
	if len(args) < 2 {
		printHelp()
		return
	}

	hopsStr := os.Args[1]

	if hopsStr == "--help" {
		printHelp()
		return
	}

	hops, err := strconv.Atoi(hopsStr)
	if err != nil {
		fmt.Println("\nThe first argument should be a number of hops\n")
		return
	}

	if hops < 1 {
		fmt.Println("\nThe number of hops should be a positive number > 0\n")
		return
	}

	sizeStr := strings.ToLower(os.Args[2])

	kb := 1024
	mb := kb * kb
	gb := mb * kb

	if strings.HasSuffix(sizeStr, "mb") {
		sizemb := strings.TrimSuffix(sizeStr, "mb")
		size, err = strconv.Atoi(sizemb)
		if err != nil {
			fmt.Println("Incorrect number of megabytes:", sizemb)
			return
		}
		size *= mb
	} else if strings.HasSuffix(sizeStr, "gb") {
		sizegb := strings.TrimSuffix(sizeStr, "gb")
		size, err = strconv.Atoi(sizegb)
		if err != nil {
			fmt.Println("Incorrect number of gigabytes:", sizegb)
			return
		}
		size *= gb
	} else if strings.HasSuffix(sizeStr, "kb") {
		sizekb := strings.TrimSuffix(sizeStr, "kb")
		size, err = strconv.Atoi(sizekb)
		if err != nil {
			fmt.Println("Incorrect number of kilobytes:", sizekb)
			return
		}
		size *= kb
	} else if strings.HasSuffix(sizeStr, "b") {
		sizeb := strings.TrimSuffix(sizeStr, "b")
		size, err = strconv.Atoi(sizeb)
		if err != nil {
			fmt.Println("Incorrect number of bytes:", sizeb)
			return
		}
	} else {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			fmt.Println("Incorrect number of bytes:", size)
			return
		}
		sizeStr += "b"
	}

	meshnet, _ := network.NewNetwork("test.network", "127.0.0.1:5999")
	defer meshnet.Shutdown()

	clientNode, serverNode := meshnet.CreateSequenceOfNodes(hops+1, 14000)
	clientAddr, serverAddr := clientNode.Id(), serverNode.Id()

	server := echoServer(serverNode)

	client, err := app.NewClient(messages.MakeAppId("echoClient"), clientNode.AppTalkAddr()) // register client on the first node
	if err != nil {
		panic(err)
	}

	err = client.Connect(messages.MakeAppId("echoServer"), serverAddr.Hex()) // client dials to server
	if err != nil {
		panic(err)
	}

	duration := benchmark(client, server, size)

	fmt.Println("server:", serverAddr.Hex())
	fmt.Println("client:", clientAddr.Hex())
	log.Println(sizeStr+" duration:", duration)
	log.Println("Ticks:", meshnet.GetTicks())
}

func benchmark(client *app.Client, server *app.Server, msgSize int) time.Duration {

	if msgSize < 1 {
		panic("message should be at least 1 byte")
	}

	msg := make([]byte, msgSize)

	start := time.Now()

	_, err := client.Send(msg)

	if err != nil {
		panic(err)
	}

	duration := time.Now().Sub(start)

	return duration
}

func echoServer(serverNode messages.NodeInterface) *app.Server {

	srv, err := app.NewServer(messages.MakeAppId("echoServer"), serverNode.AppTalkAddr(), func(in []byte) []byte {
		return in
	})
	if err != nil {
		panic(err)
	}
	return srv
}

func printHelp() {
	fmt.Println("")
	fmt.Println("Usage: go run overall.go hops_number data_size\n")
	fmt.Println("Usage example:")
	fmt.Println("go run overall.go 40 100\t- 40 hops 100 bytes")
	fmt.Println("go run overall.go 200 100b\t- 200 hops 100 bytes")
	fmt.Println("go run overall.go 2 10kb\t- 2 hops 10 kilobytes")
	fmt.Println("go run overall.go 10 10mb\t- 10 hops 10 megabytes")
	fmt.Println("go run overall.go 50 1gb\t- 50 hops 1 gigabyte")
	fmt.Println("")
}
