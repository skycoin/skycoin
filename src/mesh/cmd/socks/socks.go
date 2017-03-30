package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/skycoin/skycoin/src/mesh/app"
	"github.com/skycoin/skycoin/src/mesh/messages"
	network "github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func main() {

	//messages.SetDebugLogLevel()
	messages.SetInfoLogLevel()

	var (
		err error
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

	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()

	clientAddr, serverAddr := meshnet.CreateSequenceOfNodes(hops + 1)

	_, err = app.NewSocksServer(meshnet, serverAddr, "0.0.0.0:8001")
	if err != nil {
		panic(err)
	}

	client, err := app.NewSocksClient(meshnet, clientAddr, "0.0.0.0:8000")
	if err != nil {
		panic(err)
	}

	err = client.Dial(serverAddr)
	if err != nil {
		panic(err)
	}

	client.Listen()

}

func printHelp() {
	fmt.Println("\nFORMAT: go run socks.go n , where n is a number of hops")
	fmt.Println("\nUsage example for 10 meshnet hops:")
	fmt.Println("\ngo run socks.go 10")
	fmt.Println("\nNumber of hops should be more than 0\n")
}
