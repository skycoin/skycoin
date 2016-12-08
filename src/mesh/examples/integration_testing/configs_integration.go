package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/mesh/domain"
	"github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	statusChannel := make(chan bool, 2)

	configs := nodemanager.TestConfigsFromFile("config")

	config1 := configs[0]
	config2 := configs[1]

	nodemanager.ConnectNodeToNodeNew(config1, config2)

	//config1.AddRouteToEstablishNew(config2)
	config1.AddMessageToSend(config1.RoutesConfigsToEstablish[0].RouteID, "Message 1")
	config1.AddMessageToReceive("Message 2", "")

	config2.AddMessageToReceive("Message 1", "Message 2")

	go sendMessage(2, *config1, &wg, statusChannel)

	go sendMessage(1, *config2, &wg, statusChannel)

	timeout := 15 * time.Second
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
func sendMessage(configID int, config nodemanager.TestConfig, wg *sync.WaitGroup, statusChannel chan bool) {
	fmt.Fprintf(os.Stderr, "Starting Config: %v\n", configID)
	defer wg.Done()

	node := nodemanager.CreateNode(config)
	nodemanager.AddPeersToNodeNew(node, config)

	defer node.Close()

	nodemanager.AddRoutesToEstablish(node, config.RoutesConfigsToEstablish)

	// Send messages
	for _, messageToSend := range config.MessagesToSend {
		sendMsgErr := node.SendMessageThruRoute((domain.RouteID)(messageToSend.ThruRoute), messageToSend.Contents)
		if sendMsgErr != nil {
			panic(sendMsgErr)
		}
		fmt.Fprintf(os.Stdout, "Send message %v from Node: %v to Node: %v\n", messageToSend.Contents, configID, node.GetConnectedPeers()[0].Hex())
	}

	// Receive messages
	received := make(chan domain.MeshMessage, 2*len(config.MessagesToReceive))
	node.SetReceiveChannel(received)

	// Wait for messages to pass through
	recvMap := make(map[string]domain.ReplyTo)
	for timeEnd := time.Now().Add(1 * time.Second); time.Now().Before(timeEnd); {

		if len(received) > 0 {
			fmt.Fprintf(os.Stdout, "Len Receive Channel %v in Node: %v \n", len(received), configID)
			msgRecvd := <-received
			recvMap[fmt.Sprintf("%v", msgRecvd.Contents)] = msgRecvd.ReplyTo

			for _, messageToReceive := range config.MessagesToReceive {
				if fmt.Sprintf("%v", messageToReceive.Contents) == fmt.Sprintf("%v", msgRecvd.Contents) {
					if len(messageToReceive.Reply) > 0 {
						sendBackErr := node.SendMessageBackThruRoute(msgRecvd.ReplyTo, messageToReceive.Reply)
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
			fmt.Fprintf(os.Stdout, "Didn't receive message contents: %v - Node: %v\n", messageToReceive.Contents, configID)
		}
	}
	// Wait for messages to pass back
	time.Sleep(1 * time.Second)

	fmt.Fprintf(os.Stdout, "-- Finished test -- %v\n", time.Now())
	if success {
		fmt.Fprint(os.Stdout, "\t Success!\n")
	} else {
		fmt.Fprint(os.Stderr, "\t Failure. \n")
	}

	statusChannel <- success
}
