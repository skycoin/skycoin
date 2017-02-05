package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/skycoin/skycoin/src/mesh/app"
	"github.com/skycoin/skycoin/src/mesh/messages"
	network "github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func main() {
	messages.SetInfoLogLevel()
	pingPong(20, 10)
}

func pingPong(size, pings int) {
	meshnet := network.NewNetwork()
	nodes := meshnet.CreateRandomNetwork(size)
	var clientIndex, serverIndex int
	clientIndex = rand.Intn(size)
	for {
		serverIndex = rand.Intn(size)
		if serverIndex != clientIndex {
			break
		}
	}
	clientAddr, serverAddr := nodes[clientIndex], nodes[serverIndex]

	_, err := pongServer(meshnet, serverAddr)
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

	pingsSum, pongsSum, totalSum := int64(0), int64(0), int64(0)
	receivedPackets, lostPackets := 0, 0

	fmt.Printf("Pinging %s from %s\n\n", serverAddr.Hex(), clientAddr.Hex())
	for i := 0; i < pings; i++ {
		sendTime := time.Now().UnixNano()
		response, err := client.Send([]byte{}) //send a message to the server and wait for a response
		if err != nil {
			fmt.Println("No response")
			lostPackets++
			continue
		}

		returnTime := time.Now().UnixNano()

		receivedPackets++

		result := string(response)
		receivedByServer, err := strconv.ParseInt(result, 10, 64)
		if err != nil {
			fmt.Println("Wrong response:", result)
			panic(err)
		}

		pingTime := receivedByServer - sendTime
		pongTime := returnTime - receivedByServer
		totalTime := pingTime + pongTime

		fmt.Printf("Ping takes %d ns, pong takes %d ns, total %d ns\n", pingTime, pongTime, totalTime)

		pingsSum += pingTime
		pongsSum += pongTime
		totalSum += totalTime
	}

	pings64 := int64(pings)

	pingsAvg := pingsSum / pings64
	pongsAvg := pongsSum / pings64
	totalAvg := totalSum / pings64

	fmt.Printf("\nReceived: %d packets\n", receivedPackets)
	fmt.Printf("Lost: %d packets\n", lostPackets)

	fmt.Printf("\nAverage ping time: %d ns\n", pingsAvg)
	fmt.Printf("Average pong time: %d ns\n", pongsAvg)
	fmt.Printf("Average total time: %d ns\n", totalAvg)
	fmt.Println("")

}
