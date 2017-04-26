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
	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()

	clientConn, serverConn := meshnet.CreateSequenceOfNodes(n) // create sequence and get addresses of the first and the last node in it

	server := app.NewServer(serverConn, func(in []byte) []byte { // register server on last node in meshnet nm
		return append(in, []byte(" OK.")...) // assign callback function which handles incoming messages
	})
	defer server.Shutdown()

	client := app.NewClient(clientConn) // register client on the first node
	defer client.Shutdown()

	err := client.Dial(serverConn.Address()) // client dials to server
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
