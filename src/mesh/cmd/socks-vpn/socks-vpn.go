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

	socksServerId := messages.MakeAppId("socksServer0")
	socksClientId := messages.MakeAppId("socksClient0")

	socksServer, err := app.NewSocksServer(socksServerId, serverNode.AppTalkAddr(), "0.0.0.0:8001")
	if err != nil {
		panic(err)
	}
	defer socksServer.Shutdown()

	socksClient, err := app.NewSocksClient(socksClientId, clientNode.AppTalkAddr(), "0.0.0.0:8000")
	if err != nil {
		panic(err)
	}
	defer socksClient.Shutdown()

	err = socksClient.Connect(socksServerId, serverNode.Id().Hex())
	if err != nil {
		panic(err)
	}

	go socksClient.Listen()

	vpnServerId := messages.MakeAppId("vpn_server")
	vpnClientId := messages.MakeAppId("vpn_client")

	vpnServer, err := app.NewVPNServer(vpnServerId, serverNode.AppTalkAddr())
	if err != nil {
		panic(err)
	}
	defer vpnServer.Shutdown()

	vpnClient, err := app.NewVPNClient(vpnClientId, clientNode.AppTalkAddr(), "0.0.0.0:4321")
	if err != nil {
		panic(err)
	}
	defer vpnClient.Shutdown()

	err = vpnClient.Connect(vpnServerId, serverNode.Id().Hex())
	if err != nil {
		panic(err)
	}

	vpnClient.Listen()

}

func printHelp() {
	fmt.Println("\nFORMAT: go run socks-vpn.go n , where n is a number of hops")
	fmt.Println("\nUsage example for 10 meshnet hops:")
	fmt.Println("\ngo run socks-vpn.go 10")
	fmt.Println("\nNumber of hops should be more than 0\n")
}
