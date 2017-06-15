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

	meshnet, _ := network.NewNetwork("mesh.network", "127.0.0.1:5999")
	defer meshnet.Shutdown()

	clientNode, serverNode := meshnet.CreateSequenceOfNodes(hops+1, 15000)

	serverId := messages.MakeAppId("vpn_server")
	clientId := messages.MakeAppId("vpn_client")

	server, err := app.NewVPNServer(serverId, serverNode.AppTalkAddr())
	if err != nil {
		panic(err)
	}
	defer server.Shutdown()

	client, err := app.NewVPNClient(clientId, clientNode.AppTalkAddr(), "0.0.0.0:4321")
	if err != nil {
		panic(err)
	}
	defer client.Shutdown()

	err = client.Connect(serverId, serverNode.Id().Hex())
	if err != nil {
		panic(err)
	}

	client.Listen()

}

func printHelp() {
	fmt.Println("\nFORMAT: go run vpn.go n , where n is a number of hops")
	fmt.Println("\nUsage example for 10 meshnet hops:")
	fmt.Println("\ngo run vpn.go 10")
	fmt.Println("\nNumber of hops should be more than 0\n")
}
