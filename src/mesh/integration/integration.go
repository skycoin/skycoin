package main

import (
	"fmt"

	"github.com/skycoin/skycoin/src/mesh/app"
	network "github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func main() {
	testSendAndReceive(20)
}

func testSendAndReceive(n int) {
	meshnet := network.NewNetwork()
	clientAddr, serverAddr := meshnet.CreateSequenceOfNodes(n) // create sequence and get addresses of the first and the last node in it

	_, err := app.NewServer(meshnet, serverAddr, func(in []byte) []byte { // register server on last node in meshnet nm
		return append(in, []byte(" OK.")...) // assign callback function which handles incoming messages
	})
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
