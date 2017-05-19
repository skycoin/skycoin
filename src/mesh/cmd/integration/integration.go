package main

import (
	"fmt"

	"github.com/skycoin/skycoin/src/mesh/app"
	"github.com/skycoin/skycoin/src/mesh/messages"
	network "github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func main() {
	messages.SetDebugLogLevel()
	testSendAndReceive(20)
}

func testSendAndReceive(n int) {
	meshnet, _ := network.NewNetwork("test.network", "127.0.0.1:5999")
	defer meshnet.Shutdown()

	clientNode, serverNode := meshnet.CreateSequenceOfNodes(n, 14000) // create sequence and get addresses of the first and the last node in it

	serverId := messages.MakeAppId("1")

	server, err := app.NewServer(serverId, serverNode.AppTalkAddr(), func(in []byte) []byte { // register server on last node in meshnet nm
		return append(in, []byte(" OK.")...) // assign callback function which handles incoming messages
	})
	if err != nil {
		panic(err)
	}
	defer server.Shutdown()

	client, err := app.NewClient(messages.MakeAppId("2"), clientNode.AppTalkAddr()) // register client on the first node
	if err != nil {
		panic(err)
	}
	defer client.Shutdown()

	err = client.Connect(serverId, serverNode.Id().Hex()) // client dials to server
	if err != nil {
		panic(err)
	}

	response, err := client.Send([]byte("Integration test")) //send a message to the server and wait for a response
	if err != nil {
		panic(err)
	}

	result := string(response)

	if result == "Integration test OK." {
		fmt.Println("PASSED:", result)
	} else {
		fmt.Println("FAILED, wrong message:", result)
	}
}
