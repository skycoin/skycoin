package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh2/domain"
	"github.com/skycoin/skycoin/src/mesh2/node"
	"github.com/skycoin/skycoin/src/mesh2/nodemanager/lib_nodemanager"
	"github.com/skycoin/skycoin/src/mesh2/transport/reliable"
	"github.com/skycoin/skycoin/src/mesh2/transport/udp"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/stretchr/testify/assert"
)

type RouteConfig struct {
	Id    uuid.UUID
	Peers []cipher.PubKey
}

type MessageToSend struct {
	ThruRoute uuid.UUID
	Contents  []byte
	Reliably  bool
}

type MessageToReceive struct {
	Contents      []byte
	Reply         []byte
	ReplyReliably bool
}

type ToConnect struct {
	Peer cipher.PubKey
	Info string
}

type TestConfig struct {
	Reliable reliable.ReliableTransportConfig
	Udp      udp.UDPConfig
	Node     mesh.NodeConfig

	PeersToConnect    []ToConnect
	RoutesToEstablish []RouteConfig
	MessagesToSend    []MessageToSend
	MessagesToReceive []MessageToReceive
}

var configText1 string = `{
	"Reliable": {
		"MyPeerId": [1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
		"PhysicalReceivedChannelLength": 100,
		"ExpireMessagesInterval": 5000000000,
		"RememberMessageReceivedDuration": 10000000000,
		"RetransmitDuration": 100000000
	},
	"Udp": {
		"SendChannelLength": 100,
		"DatagramLength": 512,
		"LocalAddress": "",
		"NumListenPorts": 1,
		"ListenPortMin": 15000,
		"ExternalAddress": "127.0.0.1",
		"StunEndpoints": []
	},
	"Node": {
		"PubKey": 		[1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
		"ChaCha20Key":	[1,0,0,0,1,0,0,0,1,0,0,0,1,0,0,0,1,11,22,0,1,0,0,0,1,0,0,0,1,0,0,0],
		"MaximumForwardingDuration":	10000000000,
		"RefreshRouteDuration":			5000000000,
		"ExpireMessagesInterval":       5000000000,
		"ExpireRoutesInterval":			5000000000,
		"TimeToAssembleMessage":		10000000000,
		"TransportMessageChannelLength": 100
	},
	"PeersToConnect": [
		{
			"Peer": [3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
			"Info": "7b22446174616772616d4c656e677468223a3531322c2245787465726e616c486f737473223a5b7b224950223a223132372e302e302e31222c22506f7274223a31363030302c225a6f6e65223a22227d5d2c2243727970746f4b6579223a22415463414141454141414142414141414151414141414541414141424141414141514141414145414141413d227d"
		}
	],
	"RoutesToEstablish": [
		{
			"Id": "50000000-0000-0000-0000-000000000001",
			"Peers": [[3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]]
		}
	],
	"MessagesToSend": [
		{
			"ThruRoute": "50000000-0000-0000-0000-000000000001",
			"Contents": [3,4,5,6,7,1,2,3],
			"Reliably": true
		}
	],
	"MessagesToReceive": [
		{
			"Contents": [5,5,5,6],
			"Reply": []
		}
	]
}`

var configText2 string = `{
	"Reliable": {
		"MyPeerId": [3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
		"PhysicalReceivedChannelLength": 100,
		"ExpireMessagesInterval": 5000000000,
		"RememberMessageReceivedDuration": 10000000000,
		"RetransmitDuration": 100000000
	},
	"Udp": {
		"SendChannelLength": 100,
		"DatagramLength": 512,
		"LocalAddress": "",
		"NumListenPorts": 1,
		"ListenPortMin": 17000,
		"ExternalAddress": "127.0.0.1",
		"StunEndpoints": []
	},
	"Node": {
		"PubKey": 		[3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
		"ChaCha20Key":	[1,0,0,0,1,0,44,22,1,0,0,0,1,0,0,0,1,0,0,0,1,0,0,0,1,0,0,0,1,11,0,0],
		"MaximumForwardingDuration":	10000000000,
		"RefreshRouteDuration":			5000000000,
		"ExpireMessagesInterval":       5000000000,
		"ExpireRoutesInterval":			5000000000,
		"TimeToAssembleMessage":		10000000000,
		"TransportMessageChannelLength": 100
	},
	"PeersToConnect": [
		{
			"Peer": [1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
			"Info": "7b22446174616772616d4c656e677468223a3531322c2245787465726e616c486f737473223a5b7b224950223a223132372e302e302e31222c22506f7274223a31363030302c225a6f6e65223a22227d5d2c2243727970746f4b6579223a22415463414141454141414142414141414151414141414541414141424141414141514141414145414141413d227d"
		}
	],
	"MessagesToReceive": [
		{
			"Contents": [3,4,5,6,7,1,2,3],
			"Reply": [5,5,5,6],
			"ReplyReliably": true
		}
	]
}`

// Create TestConfig
func createTestConfig(configText string) TestConfig {
	var config TestConfig
	e_parse := json.Unmarshal([]byte(configText), &config)
	if e_parse != nil {
		panic(fmt.Sprintf("Config parse error: %v\n", e_parse))
	}
	return config
}

// Create UDPTransport
func createNewUDPTransport(configUdp udp.UDPConfig) *udp.UDPTransport {
	udpTransport, createUDPError := udp.NewUDPTransport(configUdp)
	if createUDPError != nil {
		panic(createUDPError)
	}
	return udpTransport
}

// Create TestConfig to the test using the functions created in the meshnet library.
func createTestConfig2(port int) TestConfig {
	testConfig := TestConfig{}
	testConfig.Node = lib_nodemanager.NewNodeConfig()
	testConfig.Reliable = reliable.CreateReliable(testConfig.Node.PubKey)
	testConfig.Udp = udp.CreateUdp(port, "127.0.0.1")

	return testConfig
}

// Create two Nodes using the functions created in the meshnet library.
func TestSendMessage(t *testing.T) {
	// Setup for Node 1
	config1 := createTestConfig2(15000)
	// Setup for Node 2
	config2 := createTestConfig2(17000)

	peersToConnect1 := []ToConnect{}
	peerToConnect1 := ToConnect{}
	peerToConnect1.Peer = config2.Node.PubKey
	peerToConnect1.Info = udp.CreateUDPCommConfig("127.0.0.1:17000", config2.Node.ChaCha20Key[:])
	peersToConnect1 = append(peersToConnect1, peerToConnect1)
	config1.PeersToConnect = peersToConnect1

	routesToEstablish1 := []RouteConfig{}
	routeToEstablish1 := RouteConfig{}
	routeToEstablish1.Id = uuid.NewV4()
	routeToEstablish1.Peers = append(routeToEstablish1.Peers, config2.Node.PubKey)
	routesToEstablish1 = append(routesToEstablish1, routeToEstablish1)
	config1.RoutesToEstablish = routesToEstablish1

	messagesToSend1 := []MessageToSend{}
	messageToSend1 := MessageToSend{}
	messageToSend1.ThruRoute = routeToEstablish1.Id
	messageToSend1.Contents = []byte("Message 1")
	messageToSend1.Reliably = true
	messagesToSend1 = append(messagesToSend1, messageToSend1)
	config1.MessagesToSend = messagesToSend1

	messagesToReceive1 := []MessageToReceive{}
	messageToReceive1 := MessageToReceive{}
	messageToReceive1.Contents = []byte("Message 2")
	messagesToReceive1 = append(messagesToReceive1, messageToReceive1)
	config1.MessagesToReceive = messagesToReceive1

	peersToConnect2 := []ToConnect{}
	peerToConnect2 := ToConnect{}
	peerToConnect2.Peer = config1.Node.PubKey
	peerToConnect2.Info = udp.CreateUDPCommConfig("127.0.0.1:15000", config1.Node.ChaCha20Key[:])
	peersToConnect2 = append(peersToConnect2, peerToConnect2)
	config2.PeersToConnect = peersToConnect2

	messagesToReceive2 := []MessageToReceive{}
	messageToReceive2 := MessageToReceive{}
	messageToReceive2.Contents = []byte("Message 1")
	messageToReceive2.Reply = []byte("Message 2")
	messageToReceive2.ReplyReliably = true
	messagesToReceive2 = append(messagesToReceive2, messageToReceive2)
	config2.MessagesToReceive = messagesToReceive2

	var wg sync.WaitGroup
	wg.Add(2)

	statusChannel := make(chan bool, 2)

	// Initialize Node 2
	go InitializeNode(2, config2, &wg, statusChannel)

	// Initialize Node 1
	go InitializeNode(1, config1, &wg, statusChannel)

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
				assert.Fail(t, "Error TimeOut", "")
				break
			}
		}
	}
	wg.Wait()
	fmt.Println("Done")
}

// Validates that public keys generated to random default are different.
func TestPubKey(t *testing.T) {
	b1 := cipher.RandByte(33)
	pubKeyRandom1 := cipher.NewPubKey(b1)
	fmt.Fprintf(os.Stdout, "Public Key Random 1: %v \n", pubKeyRandom1)
	assert.True(t, bytes.Equal(pubKeyRandom1[:], b1))

	b2 := cipher.RandByte(33)
	pubKeyRandom2 := cipher.NewPubKey(b2)
	fmt.Fprintf(os.Stdout, "Public Key Random 2: %v \n", pubKeyRandom2)
	assert.True(t, bytes.Equal(pubKeyRandom2[:], b2))

	assert.False(t, bytes.Equal(pubKeyRandom1[:], pubKeyRandom2[:]))
}

// Validates that info to peer connect is equal.
func TestUDPCommConfig(t *testing.T) {
	cryptoKey := []byte{1, 55, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0}
	enc := udp.CreateUDPCommConfig("127.0.0.1:16000", cryptoKey)

	expected := "7b22446174616772616d4c656e677468223a3531322c2245787465726e616c486f737473223a5b7b224950223a223132372e302e302e31222c22506f7274223a31363030302c225a6f6e65223a22227d5d2c2243727970746f4b6579223a22415463414141454141414142414141414151414141414541414141424141414141514141414145414141413d227d"
	assert.Equal(t, expected, enc, "Error in encoding")

}

// Integration Test that:
// 1. create two nodes (Nodo 1 y Nodo 2), these task are execute in Goroutines separated.
// 2. Assign transport the nodes (Nodo 1 y Nodo 2).
// 3. Connect the nodes (Nodo 1 y Nodo 2) together.
// 4. Create a route, send data over the route, confirm receipt of data
func TestNodeCase1(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	statusChannel := make(chan bool, 2)

	// Initialize Node 2
	config2 := createTestConfig(configText2)
	cryptoKey2 := []byte{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 11, 22, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0}
	config2.PeersToConnect[0].Info = udp.CreateUDPCommConfig("127.0.0.1:15000", cryptoKey2)
	go InitializeNode(2, config2, &wg, statusChannel)

	// Initialize Node 1
	config1 := createTestConfig(configText1)
	cryptoKey1 := []byte{1, 0, 0, 0, 1, 0, 44, 22, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 11, 0, 0}
	config1.PeersToConnect[0].Info = udp.CreateUDPCommConfig("127.0.0.1:17000", cryptoKey1)
	go InitializeNode(1, config1, &wg, statusChannel)

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
				assert.Fail(t, "Error TimeOut", "")
				break
			}
		}
	}
	wg.Wait()
	fmt.Println("Done")
}

// Initialize the Nodes for communication and sending messages
func InitializeNode(idConfig int, config TestConfig, wg *sync.WaitGroup, statusChannel chan bool) {
	fmt.Fprintf(os.Stderr, "Starting Config: %v\n", idConfig)
	defer wg.Done()

	udpTransport := createNewUDPTransport(config.Udp)

	// Connect
	for _, connectTo := range config.PeersToConnect {
		connectError := udpTransport.ConnectToPeer(connectTo.Peer, connectTo.Info)
		if connectError != nil {
			panic(connectError)
		}
	}

	// Reliable transport closes UDPTransport
	reliableTransport := reliable.NewReliableTransport(udpTransport, config.Reliable)
	defer reliableTransport.Close()

	node, createNodeError := mesh.NewNode(config.Node)
	if createNodeError != nil {
		panic(createNodeError)
	}
	defer node.Close()
	node.AddTransport(reliableTransport, config.Node.ChaCha20Key)

	fmt.Fprintf(os.Stdout, "UDP connect info: %v\n", udpTransport.GetTransportConnectInfo())

	// Setup route
	for _, routeConfig := range config.RoutesToEstablish {
		if len(routeConfig.Peers) == 0 {
			continue
		}
		addRouteErr := node.AddRoute((domain.RouteId)(routeConfig.Id), routeConfig.Peers[0])
		if addRouteErr != nil {
			panic(addRouteErr)
		}
		for peer := 1; peer < len(routeConfig.Peers); peer++ {
			extendErr := node.ExtendRoute((domain.RouteId)(routeConfig.Id), routeConfig.Peers[peer], 5*time.Second)
			if extendErr != nil {
				panic(extendErr)
			}
		}
	}

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

/*
type ReadableTransactionOutput struct {
	Hash    string `json:"uxid"`
	Address string `json:"dst"`
	Coins   string `json:"coins"`
	Hours   uint64 `json:"hours"`
}

type ReadableTransaction struct {
	Length    uint32 `json:"length"`
	Type      uint8  `json:"type"`
	Hash      string `json:"txid"`
	InnerHash string `json:"inner_hash"`

	Sigs []string                    `json:"sigs"`
	In   []string                    `json:"inputs"`
	Out  []ReadableTransactionOutput `json:"outputs"`
}

type ReadableUnconfirmedTxn struct {
	Txn       ReadableTransaction `json:"transaction"`
	Received  time.Time           `json:"received"`
	Checked   time.Time           `json:"checked"`
	Announced time.Time           `json:"announced"`
}
*/

func TestMarshalReadableStruc(t *testing.T) {

	readableTransactionOutput := visor.ReadableTransactionOutput{}
	readableTransactionOutput.Hash = "7b22446174616772616d4c656e677"
	readableTransactionOutput.Address = "127.0.0.1:5470"
	readableTransactionOutput.Coins = visor.StrBalance(1000)
	readableTransactionOutput.Hours = uint64(time.Now().UnixNano())

	readableTransaction := visor.ReadableTransaction{}
	readableTransaction.Length = uint32(3200)
	readableTransaction.Type = uint8(47)
	readableTransaction.Hash = "7b22446174616772616d4c656e677"
	readableTransaction.InnerHash = "2616d4c656e6777b2244617461677"
	readableTransaction.Sigs = []string{"a", "b", "c"}
	readableTransaction.In = []string{"dd", "ee", "ff"}
	readableTransaction.Out = append(readableTransaction.Out, readableTransactionOutput)

	readableUnconfirmedTxn := visor.ReadableUnconfirmedTxn{}
	readableUnconfirmedTxn.Txn = readableTransaction
	readableUnconfirmedTxn.Received, _ = time.Parse("2006-01-02", "2016-09-07")
	readableUnconfirmedTxn.Checked, _ = time.Parse("2006-01-02", "2016-09-08")
	readableUnconfirmedTxn.Announced, _ = time.Parse("2006-01-02", "2016-09-09")

	value, _ := json.Marshal(readableUnconfirmedTxn)
	fmt.Fprintln(os.Stdout, string(value))

	messageJSON := `{
	"transaction": {
		"length": 3200,
		"type": 47,
		"txid": "7b22446174616772616d4c656e677",
		"inner_hash": "2616d4c656e6777b2244617461677",
		"sigs": [
			"a",
			"b",
			"c"
		],
		"inputs": [
			"dd",
			"ee",
			"ff"
		],
		"outputs": [
			{
				"uxid": "7b22446174616772616d4c656e677",
				"dst": "127.0.0.1:5470",
				"coins": "0.1000",
				"hours": 1473390392083830819
			}
		]
	},
	"received": "2016-09-07T00:00:00Z",
	"checked": "2016-09-08T00:00:00Z",
	"announced": "2016-09-09T00:00:00Z"
}`

	readableUnconfirmedTxn2 := visor.ReadableUnconfirmedTxn{}
	json.Unmarshal([]byte(messageJSON), &readableUnconfirmedTxn2)

	value2, _ := json.Marshal(readableUnconfirmedTxn2)

	assert.Equal(t, len(value), len(value2), "Error expected equal length between value and value2")
}
