package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	mesh "github.com/skycoin/skycoin/src/mesh3/node"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	statusChannel := make(chan bool, 2)

	// Setup for Node 1
	config1 := mesh.CreateTestConfig(15000)
	// Setup for Node 2
	config2 := mesh.CreateTestConfig(17000)

	cryptoKey1 := []byte{1, 0, 0, 0, 1, 0, 44, 22, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 11, 0, 0}
	config1.AddPeerToConnect("127.0.0.1:17000", config2, cryptoKey1)
	config1.AddRouteToEstablish(config2)
	config1.AddMessageToSend(config1.RoutesToEstablish[0].Id, "Message 1", true)
	config1.AddMessageToReceive("Message 2", "", true)

	cryptoKey2 := []byte{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 11, 22, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0}
	config2.AddPeerToConnect("127.0.0.1:15000", config1, cryptoKey2)
	config2.AddMessageToReceive("Message 1", "Message 2", true)

	go sendMessage(2, *config2, &wg, statusChannel)

	go sendMessage(1, *config1, &wg, statusChannel)

	timeout := 30 * time.Second
	for i := 1; i <= 2; i++ {
		select {
		case status, ok := <-statusChannel:
			{
				if ok {
					if !status {
						fmt.Fprintln(os.Stderr, "Error expected Status True")
					}
				}
			}
		case <-time.After(timeout):
			{
				fmt.Fprintln(os.Stderr, "Error TimeOut")
				break
			}
		}
	}
	wg.Wait()
	fmt.Println("Done")
}

// Initialize the Nodes for communication and sending messages
func sendMessage(idConfig int, config mesh.TestConfig, wg *sync.WaitGroup, statusChannel chan bool) {
	fmt.Fprintf(os.Stderr, "Starting Config: %v\n", idConfig)
	defer wg.Done()

	node := mesh.CreateNode(config)
	node.AddTransportToNode(config)

	defer node.Close()

	node.AddRoutesToEstablish(config)

	// Send messages
	for _, messageToSend := range config.MessagesToSend {
		fmt.Fprintf(os.Stdout, "Is Reliably: %v\n", messageToSend.Reliably)
		sendMsgErr := node.SendMessageThruRoute((mesh.RouteId)(messageToSend.ThruRoute), messageToSend.Contents, messageToSend.Reliably)
		if sendMsgErr != nil {
			panic(sendMsgErr)
		}
		fmt.Fprintf(os.Stdout, "Send message %v from Node: %v to Node: %v\n", messageToSend.Contents, idConfig, node.GetConnectedPeers()[0].Hex())
	}

	// Receive messages
	received := make(chan mesh.MeshMessage, 2*len(config.MessagesToReceive))
	node.SetReceiveChannel(received)

	// Wait for messages to pass thru
	recvMap := make(map[string]mesh.ReplyTo)
	for timeEnd := time.Now().Add(5 * time.Second); time.Now().Before(timeEnd); {

		if len(received) > 0 {
			fmt.Fprintf(os.Stdout, "Len Receive Channel %v in Node: %v \n", len(received), idConfig)
			msgRecvd := <-received
			recvMap[fmt.Sprintf("%v", msgRecvd.Contents)] = msgRecvd.ReplyTo

			for _, messageToReceive := range config.MessagesToReceive {
				if fmt.Sprintf("%v", messageToReceive.Contents) == fmt.Sprintf("%v", msgRecvd.Contents) {
					if len(messageToReceive.Reply) > 0 {
						sendBackErr := node.SendMessageBackThruRoute(msgRecvd.ReplyTo, messageToReceive.Reply, messageToReceive.ReplyReliably)
						if sendBackErr != nil {
							panic(sendBackErr)
						}
						fmt.Fprintf(os.Stdout, "=== Send back %v\n", time.Now())
					}
				}
			}
		}
		runtime.Gosched()
	}

	success := true

	for _, messageToReceive := range config.MessagesToReceive {
		_, received := recvMap[fmt.Sprintf("%v", messageToReceive.Contents)]
		if !received {
			success = false
			fmt.Fprintf(os.Stdout, "Didn't receive message contents: %v - Node: %v\n", messageToReceive.Contents, idConfig)
		}
	}
	// Wait for messages to pass back
	time.Sleep(5 * time.Second)

	fmt.Fprintf(os.Stdout, "-- Finished test -- %v\n", time.Now())
	if success {
		fmt.Fprintf(os.Stdout, "\t Success!\n")
	} else {
		fmt.Fprintf(os.Stderr, "\t Failure. \n")
	}

	statusChannel <- success
}
