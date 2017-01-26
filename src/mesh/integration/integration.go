package main

import (
	"fmt"

	"github.com/skycoin/skycoin/src/mesh/app"
	"github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func main() {
	testSend(20)
}

func testSend(n int) {
	nm := nodemanager.NewNodeManager()
	nodeList := nm.CreateNodeList(n)
	nm.Tick()

	_, err := nm.ConnectAll() // connect all sequentially
	if err != nil {
		panic(err)
	}
	nm.RebuildRoutes()

	clientNode, serverNode := nodeList[0], nodeList[len(nodeList)-1] // get addresses for server and client

	_, err = app.NewServer(nm, serverNode, func(in []byte) []byte { // register server on last node in meshnet nm
		return append(in, []byte(" OK.")...) // assign callback function which handles incoming messages
	})
	if err != nil {
		panic(err)
	}

	client, err := app.NewClient(nm, clientNode) // register client on the first node
	if err != nil {
		panic(err)
	}

	err = client.Dial(serverNode) // client dials to server
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
