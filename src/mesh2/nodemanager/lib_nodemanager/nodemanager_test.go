package lib_nodemanager

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh2/domain"
	"github.com/skycoin/skycoin/src/mesh2/node"
	"github.com/stretchr/testify/assert"
)

func TestCreateNodeList(t *testing.T) {
	nodeManager := &NodeManager{}
	nodeManager.CreateNodeConfigList(4)
	assert.Len(t, nodeManager.ConfigList, 4, "Error expected 4 nodes")
	pubKey0 := nodeManager.PubKeyList[0]
	assert.Len(t, nodeManager.ConfigList[pubKey0].PeersToConnect, 0, "Error expected 0 PeersToConnect from Node 1")
}

func TestConnectNodes(t *testing.T) {
	nodeManager := &NodeManager{}
	nodeManager.CreateNodeConfigList(5)
	assert.Len(t, nodeManager.ConfigList, 5, "Error expected 5 nodes")
	pubKey0 := nodeManager.PubKeyList[0]
	assert.Len(t, nodeManager.ConfigList[pubKey0].PeersToConnect, 0, "Error expected 0 PeersToConnect from Node 1")
	nodeManager.ConnectNodes()
	assert.Len(t, nodeManager.ConfigList[pubKey0].PeersToConnect, 1, "Error expected 1 PeersToConnect from Node 1")
	pubKey1 := nodeManager.PubKeyList[1]
	assert.Len(t, nodeManager.ConfigList[pubKey1].PeersToConnect, 2, "Error expected 2 PeersToConnect from Node 2")
	pubKey4 := nodeManager.PubKeyList[4]
	assert.Len(t, nodeManager.ConfigList[pubKey4].PeersToConnect, 1, "Error expected 1 PeersToConnect from Node 5")
}

func TestConnectNodeRandomly(t *testing.T) {
	nodeManager := &NodeManager{Port: 1100}
	index1 := nodeManager.AddNode()
	assert.Len(t, nodeManager.NodesList, 1, "Error expected 1 nodes")
	pubKey1 := nodeManager.PubKeyList[index1]
	assert.Len(t, nodeManager.ConfigList[pubKey1].PeersToConnect, 0, "Error expected 0 PeersToConnect from Node 1")
	nodeManager.ConnectNodeRandomly(index1)
	assert.Len(t, nodeManager.ConfigList[pubKey1].PeersToConnect, 0, "Error expected 0 PeersToConnect from Node 1")

	index2 := nodeManager.AddNode()
	assert.Len(t, nodeManager.NodesList, 2, "Error expected 2 nodes")
	pubKey2 := nodeManager.PubKeyList[index2]
	assert.Len(t, nodeManager.ConfigList[pubKey2].PeersToConnect, 0, "Error expected 0 PeersToConnect from Node 2")
	nodeManager.ConnectNodeRandomly(index2)
	assert.Len(t, nodeManager.ConfigList[pubKey2].PeersToConnect, 1, "Error expected 1 PeersToConnect from Node 2")

	index3 := nodeManager.AddNode()
	assert.Len(t, nodeManager.NodesList, 3, "Error expected 3 nodes")
	pubKey3 := nodeManager.PubKeyList[index3]
	assert.Len(t, nodeManager.ConfigList[pubKey3].PeersToConnect, 0, "Error expected 0 PeersToConnect from Node 3")
	nodeManager.ConnectNodeRandomly(index3)
	assert.Len(t, nodeManager.ConfigList[pubKey3].PeersToConnect, 1, "Error expected 1 PeersToConnect from Node 3")

	index4 := nodeManager.AddNode()
	assert.Len(t, nodeManager.NodesList, 4, "Error expected 4 nodes")
	pubKey4 := nodeManager.PubKeyList[index4]
	assert.Len(t, nodeManager.ConfigList[pubKey4].PeersToConnect, 0, "Error expected 0 PeersToConnect from Node 4")
	nodeManager.ConnectNodeRandomly(index4)
	assert.Len(t, nodeManager.ConfigList[pubKey4].PeersToConnect, 1, "Error expected 1 PeersToConnect from Node 4")

	index5 := nodeManager.AddNode()
	assert.Len(t, nodeManager.NodesList, 5, "Error expected 5 nodes")
	pubKey5 := nodeManager.PubKeyList[index5]
	assert.Len(t, nodeManager.ConfigList[pubKey5].PeersToConnect, 0, "Error expected 0 PeersToConnect from Node 4")
	nodeManager.ConnectNodeRandomly(index5)
	assert.Len(t, nodeManager.ConfigList[pubKey5].PeersToConnect, 1, "Error expected 1 PeersToConnect from Node 5")
}

// Recover flow control in the tests
func recoverFlowControl(t *testing.T, index1, index2 int) {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "Error: Recovered in TestConnectTwoNodes: %v.\nIt can't connect Node %v to Node %v.\n", r, index1, index2)
	}
}

// Initialize the Node for communication and sending messages
func sendMessage(idConfig int, config mesh.TestConfig, node *mesh.Node, wg *sync.WaitGroup, statusChannel chan bool, t *testing.T, index1, index2 int) {
	defer recoverFlowControl(t, index1, index2)

	fmt.Fprintf(os.Stderr, "Starting Config: %v\n", idConfig)
	defer wg.Done()

	node.AddTransportToNode(config)
	node.AddRoutesToEstablish(config)

	defer node.Close()

	// Send messages
	for _, messageToSend := range config.MessagesToSend {
		fmt.Fprintf(os.Stdout, "Is Reliably: %v\n", messageToSend.Reliably)
		sendMsgErr := node.SendMessageThruRoute((domain.RouteId)(messageToSend.ThruRoute), messageToSend.Contents, messageToSend.Reliably)
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

func random(min, max int) int {
	time.Local = time.UTC
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

func TestRandomNumber(t *testing.T) {
	myrand := random(0, 8)
	fmt.Println(myrand)

	myrand = random(0, 9)
	fmt.Println(myrand)

	myrand = random(0, 7)
	fmt.Println(myrand)

	myrand = random(0, 15)
	fmt.Println(myrand)
}

// Connect two nodes and send one message between them with success
func TestConnectTwoNodesSuccess(t *testing.T) {
	var index1, index2 int

	nodeManager := &NodeManager{Port: 2100}
	// Connect 20 nodes randomly
	for a := 1; a <= 20; a++ {
		if a <= 10 {
			nodeManager.ConnectNodeToNetwork()
		} else {
			if index1 != index2 && index2 >= 0 {
				nodeManager.ConnectNodeToNetwork()
			} else {
				index1, index2 = nodeManager.ConnectNodeToNetwork()
			}
		}
	}

	pubKey1 := nodeManager.PubKeyList[index1]
	config1 := nodeManager.ConfigList[pubKey1]
	node1 := nodeManager.NodesList[pubKey1]

	pubKey2 := nodeManager.PubKeyList[index2]
	config2 := nodeManager.ConfigList[pubKey2]
	node2 := nodeManager.NodesList[pubKey2]

	message1 := "Message to send from Node1 to Node2"
	message2 := "Message to receive from Node2 to Node1"

	// Add route from node1 to node2
	config1.AddRouteToEstablish(config2)

	config1.AddMessageToSend(config1.RoutesToEstablish[0].Id, message1, true)
	config1.AddMessageToReceive(message2, "", true)

	config2.AddMessageToReceive(message1, message2, true)

	var wg sync.WaitGroup
	wg.Add(2)

	statusChannel := make(chan bool, 2)

	fmt.Fprintf(os.Stdout, "Send message from node %v to node %v\n", index1, index2)

	go sendMessage(index2, *config2, node2, &wg, statusChannel, t, index1, index2)

	time.Sleep(1 * time.Second)

	go sendMessage(index1, *config1, node1, &wg, statusChannel, t, index1, index2)

	timeout := 30 * time.Second
	for i := 1; i <= 2; i++ {
		select {
		case status, ok := <-statusChannel:
			{
				if ok {
					assert.True(t, status, "Error expected Status True")
				}
			}
		case <-time.After(timeout):
			{
				t.Error("Error TimeOut")
				break
			}
		}
	}
	wg.Wait()
	fmt.Println("Done")

}

// Connect two nodes and send one message between them with fail
func TestConnectTwoNodesFail(t *testing.T) {
	var index1, index2 int

	nodeManager := &NodeManager{Port: 3100}
	// Connect 20 nodes randomly
	for a := 1; a <= 20; a++ {
		nodeManager.ConnectNodeToNetwork()
	}
	rang := len(nodeManager.ConfigList)
	index1 = rand.Intn(rang)
	pubKey1 := nodeManager.PubKeyList[index1]
	config1 := nodeManager.ConfigList[pubKey1]
	node1 := nodeManager.NodesList[pubKey1]
	index2 = rand.Intn(rang)
	pubKey2 := nodeManager.PubKeyList[index2]
	config2 := nodeManager.ConfigList[pubKey2]
	node2 := nodeManager.NodesList[pubKey2]

	message1 := "Message to send from Node1 to Node2"
	message2 := "Message to receive from Node2 to Node1"

	// Add route from node1 to node2
	config1.AddRouteToEstablish(config2)

	config1.AddMessageToSend(config1.RoutesToEstablish[0].Id, message1, true)
	config1.AddMessageToReceive(message2, "", true)

	config2.AddMessageToReceive(message1, message2, true)

	var wg sync.WaitGroup
	wg.Add(2)

	statusChannel := make(chan bool, 2)

	go sendMessage(index2, *config2, node2, &wg, statusChannel, t, index1, index2)

	go sendMessage(index1, *config1, node1, &wg, statusChannel, t, index1, index2)

	timeout := 30 * time.Second
	for i := 1; i <= 2; i++ {
		select {
		case status, ok := <-statusChannel:
			{
				if ok {
					assert.False(t, status, "Error expected Status False")
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

// Connect two Nodes (Node A - Node B) through one route with various nodes.
func _TestBuildRouteWithSuccess(t *testing.T) {
	nodeManager := &NodeManager{Port: 3100}
	// Connect 200 nodes randomly
	for a := 1; a <= 20; a++ {
		nodeManager.ConnectNodeToNetwork()
	}

	var index1, index2 int

	rang := len(nodeManager.ConfigList)
	index1 = rand.Intn(rang)
	pubKey1 := nodeManager.PubKeyList[index1]
	config1 := nodeManager.ConfigList[pubKey1]
	index2 = rand.Intn(rang)
	pubKey2 := nodeManager.PubKeyList[index2]
	config2 := nodeManager.ConfigList[pubKey2]

	assert.False(t, bytes.Equal(pubKey1[:], pubKey2[:]), "Error expected that pubKey1 and pubKey2 were different")

	existConn := false
	for _, v := range config1.PeersToConnect {
		if bytes.Equal(v.Peer[:], pubKey2[:]) {
			existConn = true
		}
	}

	configList1 := []*mesh.TestConfig{}
	routeList := []cipher.PubKey{}

	if !existConn {

		for _, v := range config1.PeersToConnect {
			configN := nodeManager.ConfigList[v.Peer]
			if len(configN.PeersToConnect) > 1 {
				configList1 = append(configList1, configN)
			}
		}

		configList2 := []*mesh.TestConfig{}
		for _, v := range config2.PeersToConnect {
			configN := nodeManager.ConfigList[v.Peer]
			if len(configN.PeersToConnect) > 1 {
				configList2 = append(configList2, configN)
			}
		}

		for _, c := range configList1 {
			for _, p := range c.PeersToConnect {
				if bytes.Equal(p.Peer[:], pubKey2[:]) {
					existConn = true
					routeList = append(routeList, p.Peer)
					routeList = append(routeList, pubKey2)
					break
				}
				for _, v := range configList2 {
					for _, p2 := range v.PeersToConnect {
						if bytes.Equal(p2.Peer[:], p.Peer[:]) {
							existConn = true
							routeList = append(routeList, p.Peer)
							routeList = append(routeList, p2.Peer)
							routeList = append(routeList, pubKey2)
						}
					}
				}
			}
			if existConn {
				break
			}
		}
	}
	if assert.True(t, existConn, "Error route not found from Node A to Node B") {
		t.Log(routeList)
	}
}

func FindRoute(config *mesh.TestConfig, pubKey cipher.PubKey, routeList *[]cipher.PubKey) {
	for _, p := range config.PeersToConnect {
		if bytes.Equal(p.Peer[:], pubKey[:]) {
			*routeList = append(*routeList, pubKey)
			break
		}
	}
}

func TestBuildRoutes(t *testing.T) {
	nodeManager := &NodeManager{Port: 3100}
	// Connect 200 nodes randomly
	for a := 1; a <= 10; a++ {
		nodeManager.ConnectNodeToNetwork()
	}

	nodeManager.BuildRoutes()

	rang := len(nodeManager.ConfigList)
	index1 := rand.Intn(rang)
	pubKey1 := nodeManager.PubKeyList[index1]

	index2 := rand.Intn(rang)
	pubKey2 := nodeManager.PubKeyList[index2]

	routeKey := RouteKey{SourceNode: pubKey1, TargetNode: pubKey2}

	t.Logf("Find a route between Node %v and Node %v", index1, index2)
	route, ok := nodeManager.Routes[routeKey]

	if assert.True(t, ok, "Error expected find a route") {
		t.Log("Route:", route.RoutesToEstablish)
	}
}

func TestAddTransportsToNode(t *testing.T) {
	nodeManager := &NodeManager{Port: 5100}
	nodeManager.CreateNodeConfigList(10)
	nodeManager.ConnectNodes()

	config := CreateTestConfig(nodeManager.Port)
	nodeManager.Port++

	pubKey := nodeManager.PubKeyList[1]
	node := nodeManager.NodesList[pubKey]

	assert.Len(t, node.GetTransports(), 1, "Error expected 1 transport in the node")

	nodeManager.AddTransportsToNode(*config, 1)

	assert.Len(t, node.GetTransports(), 2, "Error expected 2 transport in the node")

	config2 := CreateTestConfig(nodeManager.Port)
	nodeManager.Port++

	pubKey2 := nodeManager.PubKeyList[3]
	node2 := nodeManager.NodesList[pubKey2]

	assert.Len(t, node2.GetTransports(), 1, "Error expected 1 transport in the node2")

	nodeManager.AddTransportsToNode(*config2, 3)

	assert.Len(t, node2.GetTransports(), 2, "Error expected 2 transport in the node2")
}

func TestGetTransportsFromNode(t *testing.T) {
	nodeManager := &NodeManager{Port: 6100}
	nodeManager.CreateNodeConfigList(10)
	nodeManager.ConnectNodes()

	pubKey := nodeManager.PubKeyList[2]
	node := nodeManager.NodesList[pubKey]

	assert.Len(t, node.GetTransports(), 1, "Error expected 1 transport in the node")
}

func TestRemoveTransportsFromNode(t *testing.T) {
	nodeManager := &NodeManager{Port: 7100}
	nodeManager.CreateNodeConfigList(10)
	nodeManager.ConnectNodes()

	pubKey := nodeManager.PubKeyList[4]
	node := nodeManager.NodesList[pubKey]

	assert.Len(t, node.GetTransports(), 1, "Error expected 1 transport in the node")

	config := CreateTestConfig(nodeManager.Port)
	nodeManager.Port++
	nodeManager.AddTransportsToNode(*config, 4)

	assert.Len(t, node.GetTransports(), 2, "Error expected 2 transport in the node")

	config2 := CreateTestConfig(nodeManager.Port)
	nodeManager.Port++
	nodeManager.AddTransportsToNode(*config2, 4)

	assert.Len(t, node.GetTransports(), 3, "Error expected 3 transport in the node")

	transport := node.GetTransports()[0]

	nodeManager.RemoveTransportsFromNode(4, transport)

	assert.Len(t, node.GetTransports(), 2, "Error expected 2 transport in the node")
}
